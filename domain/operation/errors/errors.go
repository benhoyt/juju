// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package errors

import (
	"fmt"

	"github.com/juju/juju/internal/errors"
)

const (
	// OperationNotFound describes an error that occurs when the given operation does not exist.
	OperationNotFound = errors.ConstError("operation not found")

	// TaskNotFound describes an error that occurs when the task being
	// operated on does not exist.
	TaskNotFound = errors.ConstError("task not found")

	// TaskNotPending describes an error that occurs when a pending task
	// is queried and does not have a pending status.
	TaskNotPending = errors.ConstError("task not pending")
)

// ActionNotDefined describes an error that occurs when the given charm does
// not define the given action.
type ActionNotDefined struct {
	// CharmName is the name of the charm missing the action.
	CharmName string
	// UnitName is the name of the unit where the action has been requested.
	UnitName string
	// HasActions is true if the charm defines some actions.
	HasActions bool
}

// Error implements builtin.error
func (a ActionNotDefined) Error() string {
	return fmt.Sprintf("action not defined for charm %q (target unit: %q)", a.CharmName, a.UnitName)
}
