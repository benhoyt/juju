// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package gce

import (
	"github.com/juju/errors"

	"github.com/juju/juju/instance"
	"github.com/juju/juju/network"
	"github.com/juju/juju/provider/common"
	"github.com/juju/juju/provider/gce/gceapi"
)

type environInstance struct {
	base *gceapi.Instance
	env  *environ
}

var _ instance.Instance = (*environInstance)(nil)

func newInstance(base *gceapi.Instance, env *environ) *environInstance {
	return &environInstance{
		base: base,
		env:  env,
	}
}

func (inst *environInstance) Id() instance.Id {
	return instance.Id(inst.base.ID)
}

func (inst *environInstance) Status() string {
	return inst.base.Status()
}

func (inst *environInstance) Refresh() error {
	env := inst.env.getSnapshot()
	err := inst.base.Refresh(env.gce)
	return errors.Trace(err)
}

func (inst *environInstance) Addresses() ([]network.Address, error) {
	result, err := inst.base.Addresses()
	return result, errors.Trace(err)
}

func findInst(id instance.Id, instances []instance.Instance) instance.Instance {
	for _, inst := range instances {
		if id == inst.Id() {
			return inst
		}
	}
	return nil
}

// firewall stuff

// OpenPorts opens the given ports on the instance, which
// should have been started with the given machine id.
func (inst *environInstance) OpenPorts(machineId string, ports []network.PortRange) error {
	// TODO(ericsnow) Make sure machineId matches inst.Id()?
	name := common.MachineFullName(inst.env, machineId)
	env := inst.env.getSnapshot()
	err := env.gce.OpenPorts(name, ports)
	return errors.Trace(err)
}

// ClosePorts closes the given ports on the instance, which
// should have been started with the given machine id.
func (inst *environInstance) ClosePorts(machineId string, ports []network.PortRange) error {
	name := common.MachineFullName(inst.env, machineId)
	env := inst.env.getSnapshot()
	err := env.gce.ClosePorts(name, ports)
	return errors.Trace(err)
}

// Ports returns the set of ports open on the instance, which
// should have been started with the given machine id.
// The ports are returned as sorted by SortPorts.
func (inst *environInstance) Ports(machineId string) ([]network.PortRange, error) {
	name := common.MachineFullName(inst.env, machineId)
	env := inst.env.getSnapshot()
	ports, err := env.gce.Ports(name)
	return ports, errors.Trace(err)
}
