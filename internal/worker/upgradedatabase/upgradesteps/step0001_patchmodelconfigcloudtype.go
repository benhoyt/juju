// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package upgradesteps

import (
	"context"

	"github.com/canonical/sqlair"

	"github.com/juju/juju/core/database"
	"github.com/juju/juju/core/model"
)

var (
	getCloudTypeQuery = sqlair.MustPrepare(`
SELECT ct.type AS &cloudType.type
FROM   model AS m
JOIN   cloud AS c ON m.cloud_uuid = c.uuid
JOIN   cloud_type AS ct ON c.cloud_type_id = ct.id
WHERE  m.uuid = $uuid.uuid;
`, cloudType{}, uuid{})
	updateModelDBCloudType = sqlair.MustPrepare(`
INSERT INTO model_config (key, value) VALUES ('type', $cloudType.type)
ON CONFLICT(key) DO UPDATE SET value = $cloudType.type
  `, cloudType{})
)

// Step0001_PatchModelConfigCloudType assigns the cloud type to existing
// model-config entries based on the cloud the model is deployed to.
func Step0001_PatchModelConfigCloudType(ctx context.Context, controllerDB, modelDB database.TxnRunner, modelUUID model.UUID) error {
	var ct cloudType
	if err := controllerDB.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		if err := tx.Query(ctx, getCloudTypeQuery, uuid{UUID: modelUUID}).Get(&ct); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return modelDB.Txn(ctx, func(ctx context.Context, tx *sqlair.TX) error {
		return tx.Query(ctx, updateModelDBCloudType, ct).Run()
	})
}

type cloudType struct {
	Type string `db:"type"`
}

type uuid struct {
	UUID model.UUID `db:"uuid"`
}
