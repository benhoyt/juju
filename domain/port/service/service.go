// Copyright 2024 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package service

import (
	"context"

	coreapplication "github.com/juju/juju/core/application"
	"github.com/juju/juju/core/logger"
	"github.com/juju/juju/core/machine"
	"github.com/juju/juju/core/network"
	"github.com/juju/juju/core/trace"
	coreunit "github.com/juju/juju/core/unit"
	"github.com/juju/juju/domain/port"
	"github.com/juju/juju/internal/errors"
)

// State describes the methods that a state implementation must provide to
// manage opened ports for units.
type State interface {
	WatcherState

	// GetUnitOpenedPorts returns the opened ports for a given unit uuid,
	// grouped by endpoint.
	GetUnitOpenedPorts(context.Context, coreunit.UUID) (network.GroupedPortRanges, error)

	// GetAllOpenedPorts returns the opened ports in the model, grouped by unit name.
	GetAllOpenedPorts(context.Context) (port.UnitGroupedPortRanges, error)

	// GetMachineOpenedPorts returns the opened ports for all the units on the
	// given machine. Opened ports are grouped first by unit name and then by endpoint.
	GetMachineOpenedPorts(ctx context.Context, machineUUID string) (map[coreunit.Name]network.GroupedPortRanges, error)

	// GetApplicationOpenedPorts returns the opened ports for all the units of the
	// given application. We return opened ports paired with the unit UUIDs, grouped
	// by endpoint.
	GetApplicationOpenedPorts(ctx context.Context, applicationUUID coreapplication.UUID) (port.UnitEndpointPortRanges, error)

	// GetUnitUUID returns the UUID of the unit with the given name.
	GetUnitUUID(ctx context.Context, unitName coreunit.Name) (coreunit.UUID, error)

	// ImportOpenUnitPorts opens ports for the endpoints of a given unit during
	// migration. There can be no conflicts as no other ports for this give
	// unit exist.
	ImportOpenUnitPorts(ctx context.Context, unit coreunit.UUID, openPorts network.GroupedPortRanges) error
}

// Service provides the API for managing the opened ports for units.
type Service struct {
	st     State
	logger logger.Logger
}

// NewService returns a new Service for managing opened ports for units.
func NewService(st State, logger logger.Logger) *Service {
	return &Service{
		st:     st,
		logger: logger,
	}
}

// GetUnitOpenedPorts returns the opened ports for a given unit uuid, grouped by
// endpoint.
func (s *Service) GetUnitOpenedPorts(ctx context.Context, unitUUID coreunit.UUID) (network.GroupedPortRanges, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	return s.st.GetUnitOpenedPorts(ctx, unitUUID)
}

// GetAllOpenedPorts returns the opened ports in the model, grouped by unit name.
//
// NOTE: We do not group by endpoint here. It is not needed. Instead, we just
// group by unit name
func (s *Service) GetAllOpenedPorts(ctx context.Context) (port.UnitGroupedPortRanges, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	return s.st.GetAllOpenedPorts(ctx)
}

// GetMachineOpenedPorts returns the opened ports for all endpoints, for all the
// units on the machine. Opened ports are grouped first by unit name and then by
// endpoint.
func (s *Service) GetMachineOpenedPorts(ctx context.Context, machineUUID machine.UUID) (map[coreunit.Name]network.GroupedPortRanges, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	return s.st.GetMachineOpenedPorts(ctx, machineUUID.String())
}

// GetApplicationOpenedPorts returns the opened ports for all the units of the
// application. Opened ports are grouped first by unit name and then by endpoint.
func (s *Service) GetApplicationOpenedPorts(ctx context.Context, applicationUUID coreapplication.UUID) (map[coreunit.Name]network.GroupedPortRanges, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	openedPorts, err := s.st.GetApplicationOpenedPorts(ctx, applicationUUID)
	if err != nil {
		return nil, errors.Errorf("failed to get opened ports for application %s: %w", applicationUUID, err)
	}
	return openedPorts.ByUnitByEndpoint(), nil
}

// GetApplicationOpenedPortsByEndpoint returns all the opened ports for the given
// application, across all units, grouped by endpoint.
//
// NOTE: The returned port ranges are atomised, meaning we guarantee that each
// port range is of unit length. This is useful for down-stream consumers such
// as k8s, which can only reason with unit-length port ranges.
func (s *Service) GetApplicationOpenedPortsByEndpoint(ctx context.Context, applicationUUID coreapplication.UUID) (network.GroupedPortRanges, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	openedPorts, err := s.st.GetApplicationOpenedPorts(ctx, applicationUUID)
	if err != nil {
		return nil, errors.Errorf("failed to get opened ports for application %s: %w", applicationUUID, err)
	}
	ret := network.GroupedPortRanges{}

	// group port ranges by endpoint across all units and atomise them.
	for _, openedPort := range openedPorts {
		endpoint := openedPort.Endpoint
		ret[endpoint] = append(ret[endpoint], atomisePortRange(openedPort.PortRange)...)
	}

	// de-dupe our port ranges
	for endpoint, portRanges := range ret {
		ret[endpoint] = network.UniquePortRanges(portRanges)
	}

	return ret, nil
}

// atomisePortRange breaks down the input port range into a slice of unit-length
// port ranges.
func atomisePortRange(portRange network.PortRange) []network.PortRange {
	ret := make([]network.PortRange, portRange.Length())
	for i := 0; i < portRange.Length(); i++ {
		ret[i] = network.PortRange{
			Protocol: portRange.Protocol,
			FromPort: portRange.FromPort + i,
			ToPort:   portRange.FromPort + i,
		}
	}
	return ret
}

// ImportOpenUnitPorts opens ports for the endpoints of a given unit during
// migration.
func (s *Service) ImportOpenUnitPorts(ctx context.Context, unitUUID coreunit.UUID, openPorts network.GroupedPortRanges) error {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	if len(openPorts.UniquePortRanges()) == 0 {
		return nil
	}

	if err := s.st.ImportOpenUnitPorts(ctx, unitUUID, openPorts); err != nil {
		return errors.Errorf("failed to update unit ports: %w", err)
	}
	return nil
}

// GetUnitUUID returns the UUID of the unit with the given name.
func (s *Service) GetUnitUUID(ctx context.Context, unitName coreunit.Name) (coreunit.UUID, error) {
	ctx, span := trace.Start(ctx, trace.NameFromFunc())
	defer span.End()

	if err := unitName.Validate(); err != nil {
		return "", errors.Capture(err)
	}
	return s.st.GetUnitUUID(ctx, unitName)
}
