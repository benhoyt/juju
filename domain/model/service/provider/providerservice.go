// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package provider

import (
	"context"

	"github.com/juju/juju/core/changestream"
	coremodel "github.com/juju/juju/core/model"
	"github.com/juju/juju/core/trace"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/core/watcher/eventsource"
	"github.com/juju/juju/domain/model/service"
)

// WatcherFactory describes methods for creating watchers.
type WatcherFactory interface {
	// NewNotifyWatcher returns a new watcher that filters changes from the input
	// base watcher's db/queue. A single filter option is required, though
	// additional filter options can be provided.
	NewNotifyWatcher(
		ctx context.Context,
		summary string,
		filter eventsource.FilterOption,
		filterOpts ...eventsource.FilterOption,
	) (watcher.NotifyWatcher, error)

	// NewNotifyMapperWatcher returns a new watcher that receives changes from the
	// input base watcher's db/queue. A single filter option is required, though
	// additional filter options can be provided. Filtering of values is done first
	// by the filter, and then subsequently by the mapper. Based on the mapper's
	// logic a subset of them (or none) may be emitted.
	NewNotifyMapperWatcher(
		ctx context.Context,
		summary string,
		mapper eventsource.Mapper,
		filter eventsource.FilterOption,
		filterOpts ...eventsource.FilterOption,
	) (watcher.NotifyWatcher, error)
}

// ModelState is the model state required by the provide service.
type ModelState interface {
	// GetModel returns the model info.
	GetModel(context.Context) (coremodel.ModelInfo, error)
}

// ControllerState is the controller state required by the provide service.
type ControllerState interface {
	service.ProviderControllerState
}

// ProviderService defines a service for interacting with the underlying model
// state, as opposed to the controller state.
type ProviderService struct {
	controllerSt   ControllerState
	modelSt        ModelState
	watcherFactory WatcherFactory
}

// NewProviderService returns a new Service for interacting with a model's state.
func NewProviderService(
	controllerSt ControllerState, modelSt ModelState, watcherFactory WatcherFactory,
) *ProviderService {
	return &ProviderService{
		controllerSt:   controllerSt,
		modelSt:        modelSt,
		watcherFactory: watcherFactory,
	}
}

// Model returns model info for the current service.
//
// The following error types can be expected to be returned:
// - [modelerrors.NotFound]: When the model is not found for a given uuid.
func (s *ProviderService) Model(ctx context.Context) (coremodel.ModelInfo, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	return s.modelSt.GetModel(ctx)
}

// WatchModelCloudCredential returns a new NotifyWatcher watching for changes that
// result in the cloud spec for a model changing. The changes watched for are:
// - updates to model cloud.
// - updates to model credential.
// - changes to the credential set on a model.
// The following errors can be expected:
// - [modelerrors.NotFound] when the model is not found.
func (s *ProviderService) WatchModelCloudCredential(ctx context.Context, modelUUID coremodel.UUID) (watcher.NotifyWatcher, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	return service.WatchModelCloudCredential(ctx, s.controllerSt, s.watcherFactory, modelUUID)
}

// WatchModel returns a watcher that emits an event if the model changes.
func (s ProviderService) WatchModel(ctx context.Context) (watcher.NotifyWatcher, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	return s.watcherFactory.NewNotifyWatcher(
		ctx,
		"provider model watcher",
		eventsource.NamespaceFilter("model", changestream.All),
	)
}
