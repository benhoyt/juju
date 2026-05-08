// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package testing

import (
	"bufio"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/canonical/sqlair"
	"github.com/juju/collections/transform"
	"github.com/juju/tc"
)

// Queryable is an interface that can be used to query a database.
type Queryable interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

// DumpTable dumps the contents of the given table to stdout.
// This is useful for debugging tests. It is not intended for use
// in production code.
func DumpTable(c *tc.C, queryable Queryable, table string, extraTables ...string) {
	for _, t := range append([]string{table}, extraTables...) {
		rows, err := queryable.Query(fmt.Sprintf("SELECT * FROM %q", t))
		c.Assert(err, tc.ErrorIsNil)
		defer func() { _ = rows.Close() }()

		cols, err := rows.Columns()
		c.Assert(err, tc.ErrorIsNil)

		buffer := new(bytes.Buffer)
		writer := tabwriter.NewWriter(buffer, 0, 8, 4, ' ', 0)
		for _, col := range cols {
			_, _ = fmt.Fprintf(writer, "%s\t", col)
		}

		_, _ = fmt.Fprintln(writer)

		vals := make([]any, len(cols))
		for i := range vals {
			vals[i] = new(any)
		}

		for rows.Next() {
			err = rows.Scan(vals...)
			c.Assert(err, tc.ErrorIsNil)

			for _, val := range vals {
				_, _ = fmt.Fprintf(writer, "%v\t", *val.(*any))
			}
			_, _ = fmt.Fprintln(writer)
		}
		err = rows.Err()
		c.Assert(err, tc.ErrorIsNil)
		_ = writer.Flush()

		_, _ = fmt.Fprintf(os.Stdout, "Table - %s:\n", t)

		var width int
		scanner := bufio.NewScanner(bytes.NewBuffer(buffer.Bytes()))
		for scanner.Scan() {
			if num := len(scanner.Text()); num > width {
				width = num
			}
		}

		_, _ = fmt.Fprintln(os.Stdout, strings.Repeat("-", width-4))
		_, _ = fmt.Fprintln(os.Stdout, buffer.String())
		_, _ = fmt.Fprintln(os.Stdout, strings.Repeat("-", width-4))
		_, _ = fmt.Fprintln(os.Stdout)
	}
}

type columnName struct {
	Name string `db:"name"`
}

type tableName struct {
	Name string `db:"name"`
}

// DumpTableSqlair dumps the contents of the given table to stdout using an
// sqlair.TX. This is useful for debugging tests. It is not intended for use
// in production code.
func DumpTableSqlair(c *tc.C, preparer preparer, tx *sqlair.TX, table string, extraTables ...string) {
	getColumnsNamesStmt, err := preparer.Prepare("SELECT &columnName.* FROM pragma_table_info($tableName.name)", columnName{}, tableName{})
	c.Assert(err, tc.ErrorIsNil)

	for _, t := range append([]string{table}, extraTables...) {
		var columnNames []columnName
		err = tx.Query(c.Context(), getColumnsNamesStmt, tableName{Name: t}).GetAll(&columnNames)
		c.Assert(err, tc.ErrorIsNil)

		getAllQuery := fmt.Sprintf("SELECT %s FROM %q",
			strings.Join(transform.Slice(columnNames, func(cn columnName) string { return fmt.Sprintf("&M.%s", cn.Name) }), ", "),
			t)
		getAllStmt, err := preparer.Prepare(getAllQuery, sqlair.M{})
		c.Assert(err, tc.ErrorIsNil)

		var rows []sqlair.M
		err = tx.Query(c.Context(), getAllStmt).GetAll(&rows)
		if !errors.Is(err, sqlair.ErrNoRows) {
			c.Assert(err, tc.ErrorIsNil)
		}

		buffer := new(bytes.Buffer)
		writer := tabwriter.NewWriter(buffer, 0, 8, 4, ' ', 0)
		for _, col := range columnNames {
			_, _ = fmt.Fprintf(writer, "%s\t", col.Name)
		}
		_, _ = fmt.Fprintln(writer)

		for _, row := range rows {
			for _, col := range columnNames {
				_, _ = fmt.Fprintf(writer, "%v\t", row[col.Name])
			}
			_, _ = fmt.Fprintln(writer)
		}
		err = writer.Flush()
		c.Assert(err, tc.ErrorIsNil)

		_, _ = fmt.Fprintf(os.Stdout, "Table - %s:\n", t)

		var width int
		scanner := bufio.NewScanner(bytes.NewBuffer(buffer.Bytes()))
		for scanner.Scan() {
			if num := len(scanner.Text()); num > width {
				width = num
			}
		}

		_, _ = fmt.Fprintln(os.Stdout, strings.Repeat("-", width-4))
		_, _ = fmt.Fprintln(os.Stdout, buffer.String())
		_, _ = fmt.Fprintln(os.Stdout, strings.Repeat("-", width-4))
		_, _ = fmt.Fprintln(os.Stdout)
	}
}

// DumpForeignKeysForTableSqlair dumps the contents of all the tables that have
// a foreign key constraint targeting the given table, using an sqlair.TX. This
// is useful for debugging tests. It is not intended for use in production code.
func DumpForeignKeysForTableSqlair(c *tc.C, preparer preparer, tx *sqlair.TX, table string) {
	getForeignKeysStmt, err := preparer.Prepare(`
WITH tbls(name) AS ( 
	SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' 
) 
SELECT DISTINCT
	t.name    AS &tableName.name 
FROM tbls t, pragma_foreign_key_list(t.name) AS fk 
WHERE fk."table" = $tableName.name
`, tableName{})
	c.Assert(err, tc.ErrorIsNil)

	var sourceTables []tableName
	err = tx.Query(c.Context(), getForeignKeysStmt, tableName{Name: table}).GetAll(&sourceTables)
	if !errors.Is(err, sqlair.ErrNoRows) {
		c.Assert(err, tc.ErrorIsNil)
	}

	if len(sourceTables) == 0 {
		_, _ = fmt.Fprintf(os.Stdout, "Table %q has no foreign key links.\n", table)
		return
	}
	DumpTableSqlair(c, preparer, tx, sourceTables[0].Name, transform.Slice(sourceTables[1:], func(t tableName) string { return t.Name })...)
}

// Preparer is an interface that prepares SQL statements for sqlair.
type preparer interface {
	Prepare(query string, typeSamples ...any) (*sqlair.Statement, error)
}
