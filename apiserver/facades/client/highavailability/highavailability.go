// Copyright 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package highavailability

import (
	"context"
	"sort"

	"github.com/juju/errors"
	"github.com/juju/names/v6"

	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/apiserver/facade"
	corelogger "github.com/juju/juju/core/logger"
	"github.com/juju/juju/core/permission"
	controllernodeerrors "github.com/juju/juju/domain/controllernode/errors"
	"github.com/juju/juju/rpc/params"
)

// ControllerNodeService describes the maintenance of controller entries.
type ControllerNodeService interface {
	// GetControllerAPIAddresses returns the list of API addresses for all
	// controllers.
	GetAPIAddressesByControllerIDForClients(ctx context.Context) (map[string][]string, error)
}

// HighAvailabilityAPI implements the HighAvailability interface and is the concrete
// implementation of the api end point.
type HighAvailabilityAPI struct {
	controllerTag         names.ControllerTag
	controllerNodeService ControllerNodeService
	authorizer            facade.Authorizer
	logger                corelogger.Logger
}

// HighAvailabilityAPIV2 implements v2 of the high availability facade.
type HighAvailabilityAPIV2 struct {
	HighAvailabilityAPI
}

// EnableHA adds controller machines as necessary to ensure the
// controller has the number of machines specified.
func (api *HighAvailabilityAPI) EnableHA(
	ctx context.Context, args params.ControllersSpecs,
) (params.ControllersChangeResults, error) {
	return params.ControllersChangeResults{}, apiservererrors.ServerError(errors.NotSupportedf("enable HA"))
}

// ControllerDetails is only available on V3 or later.
func (api *HighAvailabilityAPIV2) ControllerDetails(_ struct{}) {}

// ControllerDetails returns details about each controller node.
func (api *HighAvailabilityAPI) ControllerDetails(
	ctx context.Context,
) (params.ControllerDetailsResults, error) {
	results := params.ControllerDetailsResults{}

	err := api.authorizer.HasPermission(ctx, permission.LoginAccess, api.controllerTag)
	if err != nil {
		return results, apiservererrors.ServerError(apiservererrors.ErrPerm)
	}

	controllerAddresses, err := api.controllerNodeService.GetAPIAddressesByControllerIDForClients(ctx)
	if errors.Is(err, controllernodeerrors.EmptyAPIAddresses) {
		// If there are no API addresses, we return an empty result.
		return results, nil
	} else if err != nil {
		return results, apiservererrors.ServerError(errors.Trace(err))
	}

	details := make([]params.ControllerDetails, 0, len(controllerAddresses))
	for id, addresses := range controllerAddresses {
		details = append(details, params.ControllerDetails{
			ControllerId: id,
			APIAddresses: addresses,
		})
	}

	sort.Slice(details, func(i, j int) bool {
		return details[i].ControllerId < details[j].ControllerId
	})

	results.Results = details

	return results, nil
}
