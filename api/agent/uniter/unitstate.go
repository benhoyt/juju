// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package uniter

import (
	"context"

	"github.com/juju/errors"

	apiservererrors "github.com/juju/juju/apiserver/errors"
	"github.com/juju/juju/rpc/params"
)

// State returns the state persisted by the charm running in this unit
// and the state internal to the uniter for this unit.
func (client *Client) State(ctx context.Context) (params.UnitStateResult, error) {
	var results params.UnitStateResults
	args := params.Entities{
		Entities: []params.Entity{{Tag: client.unitTag.String()}},
	}
	err := client.facade.FacadeCall(ctx, "State", args, &results)
	if err != nil {
		return params.UnitStateResult{}, errors.Trace(apiservererrors.RestoreError(err))
	}
	if len(results.Results) != 1 {
		return params.UnitStateResult{}, errors.Errorf("expected 1 result, got %d", len(results.Results))
	}
	result := results.Results[0]
	if result.Error != nil {
		return params.UnitStateResult{}, result.Error
	}
	return result, nil
}

// SetState sets the state persisted by the charm running in this unit
// and the state internal to the uniter for this unit.
func (client *Client) SetState(ctx context.Context, unitState params.SetUnitStateArg) error {
	unitState.Tag = client.unitTag.String()
	var results params.ErrorResults
	args := params.SetUnitStateArgs{
		Args: []params.SetUnitStateArg{unitState},
	}
	err := client.facade.FacadeCall(ctx, "SetState", args, &results)
	if err != nil {
		return errors.Trace(apiservererrors.RestoreError(err))
	}
	// Make sure we correctly decode quota-related errors.
	return maybeRestoreQuotaLimitError(results.OneError())
}

// maybeRestoreQuotaLimitError checks if the server emitted a quota limit
// exceeded error and restores it back to a typed error from juju/errors.
// Ideally, we would use apiserver/common.RestoreError but apparently, that
// package imports worker/uniter/{operation, remotestate} causing an import
// cycle when api/common is imported by api/uniter.
func maybeRestoreQuotaLimitError(err error) error {
	if params.IsCodeQuotaLimitExceeded(err) {
		return errors.NewQuotaLimitExceeded(nil, err.Error())
	}
	return err
}
