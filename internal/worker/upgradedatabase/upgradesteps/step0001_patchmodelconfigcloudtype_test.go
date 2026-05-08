// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package upgradesteps

import (
	"context"
	"database/sql"
	stdtesting "testing"

	"github.com/juju/tc"

	"github.com/juju/juju/core/model"
	schematesting "github.com/juju/juju/domain/schema/testing"
)

type step0001Suite struct {
	schematesting.ControllerModelSuite
}

func TestStep0001Suite(t *stdtesting.T) {
	tc.Run(t, &step0001Suite{})
}

func (s *step0001Suite) setupTestModel(c *tc.C, cloudTypeID int, modelUUID model.UUID) error {
	controllerDB := s.ControllerTxnRunner()
	return controllerDB.StdTxn(c.Context(), func(ctx context.Context, tx *sql.Tx) error {
		// Insert cloud.
		cloudUUID := modelUUID.String() + "-cloud"
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO cloud (uuid, cloud_type_id, name, endpoint, skip_tls_verify) VALUES (?, ?, ?, ?, ?)",
			cloudUUID, cloudTypeID, "test-cloud", "", 0); err != nil {
			return err
		}

		// Insert a model with all required fields.
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO model (uuid, cloud_uuid, model_type_id, life_id, name, qualifier) VALUES (?, ?, ?, ?, ?, ?)",
			modelUUID.String(), cloudUUID, 0, 0, "test-model", "test"); err != nil {
			return err
		}

		return nil
	})
}

func (s *step0001Suite) TestStep0001_PatchModelConfigCloudType_Success(c *tc.C) {
	controllerDB := s.ControllerTxnRunner()
	modelUUID := tc.Must(c, model.NewUUID)
	modelDB := s.ModelTxnRunner(c, modelUUID.String())

	// Use ec2 (ID 5) for this test.
	cloudTypeID := 5
	expectedCloudType := "ec2"

	// Set up test data in controller database.
	err := s.setupTestModel(c, cloudTypeID, modelUUID)
	c.Assert(err, tc.ErrorIsNil)

	// Execute the upgrade step.
	err = Step0001_PatchModelConfigCloudType(c.Context(), controllerDB, modelDB, modelUUID)
	c.Assert(err, tc.ErrorIsNil)

	// Verify the cloud type was written to model config.
	var resultValue string
	err = modelDB.StdTxn(c.Context(), func(ctx context.Context, tx *sql.Tx) error {
		return tx.QueryRowContext(ctx,
			"SELECT value FROM model_config WHERE key = 'type'").Scan(&resultValue)
	})
	c.Assert(err, tc.ErrorIsNil)
	c.Check(resultValue, tc.Equals, expectedCloudType)
}

func (s *step0001Suite) TestStep0001_PatchModelConfigCloudType_UpdatesExisting(c *tc.C) {
	controllerDB := s.ControllerTxnRunner()
	modelUUID := tc.Must(c, model.NewUUID)
	modelDB := s.ModelTxnRunner(c, modelUUID.String())

	cloudTypeID := 4
	expectedCloudType := "azure"

	// Pre-populate model config with an old value (wrong) value, which used to
	// be the case for the controller model prior to the fix
	// https://github.com/juju/juju/pull/21469.
	oldCloudType := "iaas"
	err := modelDB.StdTxn(c.Context(), func(ctx context.Context, tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx,
			"INSERT INTO model_config (key, value) VALUES (?, ?)",
			"type", oldCloudType)
		return err
	})
	c.Assert(err, tc.ErrorIsNil)

	// Set up test data in controller database with new cloud type.
	err = s.setupTestModel(c, cloudTypeID, modelUUID)
	c.Assert(err, tc.ErrorIsNil)

	err = Step0001_PatchModelConfigCloudType(c.Context(), controllerDB, modelDB, modelUUID)
	c.Assert(err, tc.ErrorIsNil)

	// Verify the cloud type was updated in model config.
	var resultValue string
	err = modelDB.StdTxn(c.Context(), func(ctx context.Context, tx *sql.Tx) error {
		return tx.QueryRowContext(ctx,
			"SELECT value FROM model_config WHERE key = 'type'").Scan(&resultValue)
	})
	c.Assert(err, tc.ErrorIsNil)
	c.Check(resultValue, tc.Equals, expectedCloudType)
}

func (s *step0001Suite) TestStep0001_PatchModelConfigCloudType_ModelNotFound(c *tc.C) {
	controllerDB := s.ControllerTxnRunner()
	nonExistentModelUUID := tc.Must(c, model.NewUUID)
	modelDB := s.ModelTxnRunner(c, nonExistentModelUUID.String())

	// Execute the upgrade step - should fail because model doesn't exist.
	err := Step0001_PatchModelConfigCloudType(c.Context(), controllerDB, modelDB, nonExistentModelUUID)
	c.Assert(err, tc.ErrorMatches, "*no rows in result set")
}

func (s *step0001Suite) TestStep0001_PatchModelConfigCloudType_Idempotent(c *tc.C) {
	controllerDB := s.ControllerTxnRunner()
	modelUUID := tc.Must(c, model.NewUUID)
	modelDB := s.ModelTxnRunner(c, modelUUID.String())

	cloudTypeID := 6
	expectedCloudType := "gce"

	// Set up test data in controller database.
	err := s.setupTestModel(c, cloudTypeID, modelUUID)
	c.Assert(err, tc.ErrorIsNil)

	// Execute the upgrade step first time.
	err = Step0001_PatchModelConfigCloudType(c.Context(), controllerDB, modelDB, modelUUID)
	c.Assert(err, tc.ErrorIsNil)

	// Execute the upgrade step second time - should succeed (idempotent).
	err = Step0001_PatchModelConfigCloudType(c.Context(), controllerDB, modelDB, modelUUID)
	c.Assert(err, tc.ErrorIsNil)

	// Verify the cloud type is still correct.
	var resultValue string
	err = modelDB.StdTxn(c.Context(), func(ctx context.Context, tx *sql.Tx) error {
		return tx.QueryRowContext(ctx,
			"SELECT value FROM model_config WHERE key = 'type'").Scan(&resultValue)
	})
	c.Assert(err, tc.ErrorIsNil)
	c.Check(resultValue, tc.Equals, expectedCloudType)
}
