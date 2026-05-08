// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package migration

import (
	"context"

	"github.com/juju/juju/core/logger"
	coremodel "github.com/juju/juju/core/model"
	"github.com/juju/juju/core/trace"
	"github.com/juju/juju/domain/model"
	"github.com/juju/juju/domain/model/service"
	secretbackenderrors "github.com/juju/juju/domain/secretbackend/errors"
	"github.com/juju/juju/internal/errors"
	jujusecrets "github.com/juju/juju/internal/secrets/provider/juju"
	kubernetessecrets "github.com/juju/juju/internal/secrets/provider/kubernetes"
)

// State is the combined state required by the migration service.
type State interface {
	service.CreateModelState
}

// MigrationService defines a service for interacting with the underlying state based
// information of a model.
type MigrationService struct {
	st     State
	logger logger.Logger
}

// NewMigrationService returns a new MigrationService for interacting with a models state.
func NewMigrationService(
	st State,
	logger logger.Logger,
) *MigrationService {
	return &MigrationService{
		st:     st,
		logger: logger,
	}
}

// ImportModel is responsible for importing an existing model into this Juju
// controller by creating the model record in the controller database and marking
// it as importing.
//
// The following error types can be expected to be returned:
// - [modelerrors.AlreadyExists]: When the model uuid is already in use or a
// model with the same name and owner already exists.
// - [errors.NotFound]: When the cloud, cloud region, or credential do not
// exist.
// - [github.com/juju/juju/domain/access/errors.NotFound]: When the owner of the
// model can not be found.
// - [secretbackenderrors.NotFound] When the secret backend for the model
// cannot be found.
func (s *MigrationService) ImportModel(
	ctx context.Context,
	args model.ModelImportArgs,
) error {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	if err := args.Validate(); err != nil {
		return errors.Errorf(
			"cannot validate model import args: %w", err,
		)
	}

	modelType, err := service.ModelTypeForCloud(ctx, s.st, args.Cloud)
	if err != nil {
		return errors.Errorf(
			"determining model type when importing model %q: %w",
			args.Name, err,
		)
	}

	if args.SecretBackend == "" {
		switch modelType {
		case coremodel.CAAS:
			args.SecretBackend = kubernetessecrets.BackendName
		case coremodel.IAAS:
			args.SecretBackend = jujusecrets.BackendName
		default:
			return errors.Errorf(
				"%w for model type %q when creating model with name %q",
				secretbackenderrors.NotFound,
				modelType,
				args.Name,
			)
		}
	}

	if err := s.st.ImportModel(ctx, args.UUID, modelType, args.GlobalModelCreationArgs); err != nil {
		return err
	}

	return nil
}

// ActivateModel marks the model as active after a successful import or
// migration.
func (s *MigrationService) ActivateModel(ctx context.Context, modelUUID coremodel.UUID) error {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	if err := modelUUID.Validate(); err != nil {
		return errors.Errorf("invalid model UUID %q: %w", modelUUID, err)
	}
	return s.st.Activate(ctx, modelUUID)
}
