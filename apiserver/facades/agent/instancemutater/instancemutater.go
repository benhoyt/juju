// Copyright 2019 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package instancemutater

import (
	"context"

	"github.com/juju/names/v6"
	"gopkg.in/tomb.v2"

	"github.com/juju/juju/apiserver/facade"
	"github.com/juju/juju/apiserver/internal"
	"github.com/juju/juju/core/life"
	"github.com/juju/juju/internal/errors"
	"github.com/juju/juju/rpc/params"
)

// InstanceMutaterAPI exists only to allow migrations from Juju 3.6.
// It effectively does nothing.
type InstanceMutaterAPI struct {
	watcherRegistry facade.WatcherRegistry
}

// NewInstanceMutaterAPI creates a new API server endpoint for managing
// charm profiles on juju lxd machines and containers.
func NewInstanceMutaterAPI(watcherRegistry facade.WatcherRegistry) *InstanceMutaterAPI {
	return &InstanceMutaterAPI{
		watcherRegistry: watcherRegistry,
	}
}

// CharmProfilingInfo returns info to update lxd profiles on the machine. If
// the machine is not provisioned, no profile change info will be returned,
// nor will an error.
func (api *InstanceMutaterAPI) CharmProfilingInfo(ctx context.Context, arg params.Entity) (params.CharmProfilingInfoResult, error) {
	return params.CharmProfilingInfoResult{}, nil
}

// ContainerType returns the container type of a machine.
func (api *InstanceMutaterAPI) ContainerType(ctx context.Context, arg params.Entity) (params.ContainerTypeResult, error) {
	return params.ContainerTypeResult{}, nil
}

// SetModificationStatus updates the instance whilst changes are occurring. This
// is different from SetStatus and SetInstanceStatus, by the fact this holds
// information about the ongoing changes that are happening to instances.
// Consider LXD Profile updates that can modify a instance, but may not cause
// the instance to be placed into a error state. This modification status
// serves the purpose of highlighting that to the operator.
// Only machine tags are accepted.
func (api *InstanceMutaterAPI) SetModificationStatus(ctx context.Context, args params.SetStatus) (params.ErrorResults, error) {
	return params.ErrorResults{
		Results: make([]params.ErrorResult, len(args.Entities)),
	}, nil
}

// SetCharmProfiles records the given slice of charm profile names.
func (api *InstanceMutaterAPI) SetCharmProfiles(ctx context.Context, args params.SetProfileArgs) (params.ErrorResults, error) {
	return params.ErrorResults{
		Results: make([]params.ErrorResult, len(args.Args)),
	}, nil
}

// WatchMachines starts a watcher to track machines.
// WatchMachines does not consume the initial event of the watch response, as
// that returns the initial set of machines that are currently available.
func (api *InstanceMutaterAPI) WatchMachines(ctx context.Context) (params.StringsWatchResult, error) {
	// Create a simple watcher that sends the empty string as initial event.
	w := newEmptyStringWatcher()

	watcherID, changes, err := internal.EnsureRegisterWatcher[[]string](ctx, api.watcherRegistry, w)
	if err != nil {
		return params.StringsWatchResult{}, errors.Capture(err)
	}
	return params.StringsWatchResult{
		StringsWatcherId: watcherID,
		Changes:          changes,
	}, nil
}

// WatchModelMachines starts a watcher to track machines, but not containers.
// WatchModelMachines does not consume the initial event of the watch response, as
// that returns the initial set of machines that are currently available.
func (api *InstanceMutaterAPI) WatchModelMachines(ctx context.Context) (params.StringsWatchResult, error) {
	// Create a simple watcher that sends the empty string as initial event.
	w := newEmptyStringWatcher()

	watcherID, changes, err := internal.EnsureRegisterWatcher[[]string](ctx, api.watcherRegistry, w)
	if err != nil {
		return params.StringsWatchResult{}, errors.Capture(err)
	}
	return params.StringsWatchResult{
		StringsWatcherId: watcherID,
		Changes:          changes,
	}, nil
}

// WatchContainers starts a watcher to track Containers on a given
// machine.
func (api *InstanceMutaterAPI) WatchContainers(ctx context.Context, arg params.Entity) (params.StringsWatchResult, error) {
	// Create a simple watcher that sends the empty string as initial event.
	w := newEmptyStringWatcher()

	watcherID, changes, err := internal.EnsureRegisterWatcher[[]string](ctx, api.watcherRegistry, w)
	if err != nil {
		return params.StringsWatchResult{}, errors.Capture(err)
	}
	return params.StringsWatchResult{
		StringsWatcherId: watcherID,
		Changes:          changes,
	}, nil
}

// WatchLXDProfileVerificationNeeded starts a watcher to track Applications with
// LXD Profiles.
func (api *InstanceMutaterAPI) WatchLXDProfileVerificationNeeded(ctx context.Context, args params.Entities) (params.NotifyWatchResults, error) {
	result := params.NotifyWatchResults{
		Results: make([]params.NotifyWatchResult, len(args.Entities)),
	}
	for i := range args.Entities {
		// Create a simple notify watcher that only sends one initial event.
		w := newEmptyNotifyWatcher()

		watcherID, _, err := internal.EnsureRegisterWatcher[struct{}](ctx, api.watcherRegistry, w)
		if err != nil {
			return params.NotifyWatchResults{}, errors.Capture(err)
		}
		result.Results[i] = params.NotifyWatchResult{
			NotifyWatcherId: watcherID,
		}
	}

	return result, nil
}

// OneLife returns the life of the specified entity.
func (api *InstanceMutaterAPI) OneLife(tag names.Tag) (life.Value, error) {
	return life.Alive, nil
}

// Life returns the life status of every supplied entity, where available.
func (api *InstanceMutaterAPI) Life(ctx context.Context, args params.Entities) (params.LifeResults, error) {
	result := params.LifeResults{
		Results: make([]params.LifeResult, len(args.Entities)),
	}
	for i := range args.Entities {
		result.Results[i].Life = life.Alive
	}
	return result, nil
}

// newEmptyStringWatcher returns starts and returns a new empty string watcher,
// with an empty string as initial event.
func newEmptyStringWatcher() *emptyStringWatcher {
	changes := make(chan []string)

	w := &emptyStringWatcher{
		changes: changes,
	}
	w.tomb.Go(func() error {
		changes <- []string{""}
		defer close(changes)
		return w.loop()
	})

	return w
}

// emptyStringWatcher implements watcher.StringsWatcher.
type emptyStringWatcher struct {
	changes chan []string
	tomb    tomb.Tomb
}

// Changes returns the event channel for the empty string watcher.
func (w *emptyStringWatcher) Changes() <-chan []string {
	return w.changes
}

// Kill asks the watcher to stop without waiting for it do so.
func (w *emptyStringWatcher) Kill() {
	w.tomb.Kill(nil)
}

// Wait waits for the watcher to die and returns any
// error encountered when it was running.
func (w *emptyStringWatcher) Wait() error {
	return w.tomb.Wait()
}

// Err returns any error encountered while the watcher
// has been running.
func (w *emptyStringWatcher) Err() error {
	return w.tomb.Err()
}

func (w *emptyStringWatcher) loop() error {
	for {
		select {
		case <-w.tomb.Dying():
			return tomb.ErrDying
		}
	}
}

// newEmptyNotifyWatcher returns starts and returns a new empty notify watcher,
// with only the initial event.
func newEmptyNotifyWatcher() *emptyNotifyWatcher {
	changes := make(chan struct{})

	w := &emptyNotifyWatcher{
		changes: changes,
	}
	w.tomb.Go(func() error {
		changes <- struct{}{}
		defer close(changes)
		return w.loop()
	})

	return w
}

// emptyNotifyWatcher implements watcher.NotifyWatcher.
type emptyNotifyWatcher struct {
	changes chan struct{}
	tomb    tomb.Tomb
}

// Changes returns the event channel for the empty notify watcher.
func (w *emptyNotifyWatcher) Changes() <-chan struct{} {
	return w.changes
}

// Kill asks the watcher to stop without waiting for it do so.
func (w *emptyNotifyWatcher) Kill() {
	w.tomb.Kill(nil)
}

// Wait waits for the watcher to die and returns any
// error encountered when it was running.
func (w *emptyNotifyWatcher) Wait() error {
	return w.tomb.Wait()
}

// Err returns any error encountered while the watcher
// has been running.
func (w *emptyNotifyWatcher) Err() error {
	return w.tomb.Err()
}

func (w *emptyNotifyWatcher) loop() error {
	for {
		select {
		case <-w.tomb.Dying():
			return tomb.ErrDying
		}
	}
}
