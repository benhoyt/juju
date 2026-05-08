// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"
	"fmt"

	"github.com/juju/collections/set"
	"github.com/juju/collections/transform"

	coreapplication "github.com/juju/juju/core/application"
	"github.com/juju/juju/core/changestream"
	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/core/trace"
	"github.com/juju/juju/core/unit"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/core/watcher/eventsource"
)

// WatchableService provides the API for managing the opened ports for units, as
// well as the ability to watch for changes to opened ports.
type WatchableService struct {
	*Service
	watcherFactory WatcherFactory
}

// NewWatchableService returns a new Service providing an API to manage the opened
// ports for units.
func NewWatchableService(st State, watcherFactory WatcherFactory, logger logger.Logger) *WatchableService {
	return &WatchableService{
		Service:        &Service{st: st, logger: logger},
		watcherFactory: watcherFactory,
	}
}

// WatcherFactory describes methods for creating watchers.
type WatcherFactory interface {
	// NewNamespaceMapperWatcher returns a new watcher that receives changes
	// from the input base watcher's db/queue. Change-log events will be emitted
	// only if the filter accepts them, and dispatching the notifications via
	// the Changes channel, once the mapper has processed them. Filtering of
	// values is done first by the filter, and then by the mapper. Based on the
	// mapper's logic a subset of them (or none) may be emitted. A filter option
	// is required, though additional filter options can be provided.
	NewNamespaceMapperWatcher(
		ctx context.Context,
		initialQuery eventsource.NamespaceQuery,
		summary string,
		mapper eventsource.Mapper,
		filterOption eventsource.FilterOption, filterOptions ...eventsource.FilterOption,
	) (watcher.StringsWatcher, error)

	// NewNamespaceWatcher returns a new watcher that filters changes from the input
	// base watcher's db/queue. Change-log events will be emitted only if the filter
	// accepts them, and dispatching the notifications via the Changes channel. A
	// filter option is required, though additional filter options can be provided.
	NewNamespaceWatcher(
		ctx context.Context,
		initialQuery eventsource.NamespaceQuery,
		summary string,
		filterOption eventsource.FilterOption, filterOptions ...eventsource.FilterOption,
	) (watcher.StringsWatcher, error)

	// NewNotifyMapperWatcher returns a new watcher that receives changes from
	// the input base watcher's db/queue. A single filter option is required,
	// though additional filter options can be provided. Filtering of values is
	// done first by the filter, and then subsequently by the mapper. Based on
	// the mapper's logic a subset of them (or none) may be emitted.
	NewNotifyMapperWatcher(
		ctx context.Context,
		summary string,
		mapper eventsource.Mapper,
		filter eventsource.FilterOption,
		filterOpts ...eventsource.FilterOption,
	) (watcher.NotifyWatcher, error)
}

// WatcherState describes the methods that the service needs for its watchers.
type WatcherState interface {
	// NamespaceForWatchOpenedPort returns the name of the table that should be
	// watched
	NamespaceForWatchOpenedPort() string

	// InitialWatchOpenedPortsStatement returns the name of the table
	// that should be watched and the query to load the
	// initial event for the WatchOpenedPorts watcher
	InitialWatchOpenedPortsStatement() (string, string)

	// FilterUnitUUIDsForApplication returns the subset of provided endpoint
	// uuids that are associated with the provided application.
	FilterUnitUUIDsForApplication(context.Context, []unit.UUID, coreapplication.UUID) (set.Strings, error)
}

// WatchOpenedPorts returns a strings watcher for opened ports. This watcher
// emits events for changes to the opened ports table. Each emitted event
// contains the unit uuids which have seen changes to their opened ports.
func (s *WatchableService) WatchOpenedPorts(ctx context.Context) (watcher.StringsWatcher, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	table, statement := s.st.InitialWatchOpenedPortsStatement()
	return s.watcherFactory.NewNamespaceWatcher(
		ctx,
		eventsource.InitialNamespaceChanges(statement),
		"opened ports watcher",
		eventsource.NamespaceFilter(table, changestream.All),
	)
}

// WatchOpenedPortsForApplication returns a notify watcher for opened ports. This
// watcher emits events for changes to the opened ports table that are associated
// with the given application
func (s *WatchableService) WatchOpenedPortsForApplication(ctx context.Context, applicationUUID coreapplication.UUID) (watcher.NotifyWatcher, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	return s.watcherFactory.NewNotifyMapperWatcher(
		ctx,
		fmt.Sprintf("opened ports watcher for %q", applicationUUID),
		s.filterForApplication(applicationUUID),
		eventsource.NamespaceFilter(s.st.NamespaceForWatchOpenedPort(), changestream.All),
	)
}

// filterForApplication returns an eventsource.Mapper that filters events
// emitted by port range changes to only include events for port range changes
// corresponding to the given application
func (s *WatchableService) filterForApplication(applicationUUID coreapplication.UUID) eventsource.Mapper {
	return func(
		ctx context.Context, events []changestream.ChangeEvent,
	) ([]string, error) {
		unitUUIDs, err := transform.SliceOrErr(events, func(e changestream.ChangeEvent) (unit.UUID, error) {
			return unit.ParseID(e.Changed())
		})
		if err != nil {
			return nil, err
		}

		unitUUIDsForApplication, err := s.st.FilterUnitUUIDsForApplication(ctx, unitUUIDs, applicationUUID)
		if err != nil {
			return nil, err
		}
		var results []string
		for _, event := range events {
			if unitUUIDsForApplication.Contains(event.Changed()) {
				results = append(results, event.Changed())
			}
		}
		return results, nil
	}
}
