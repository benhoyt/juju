// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/api/provisioner (interfaces: MachineProvisioner)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	provisioner "github.com/juju/juju/api/provisioner"
	params "github.com/juju/juju/apiserver/params"
	status "github.com/juju/juju/core/status"
	watcher "github.com/juju/juju/core/watcher"
	instance "github.com/juju/juju/instance"
	version "github.com/juju/version"
	names_v2 "gopkg.in/juju/names.v2"
	reflect "reflect"
)

// MockMachineProvisioner is a mock of MachineProvisioner interface
type MockMachineProvisioner struct {
	ctrl     *gomock.Controller
	recorder *MockMachineProvisionerMockRecorder
}

// MockMachineProvisionerMockRecorder is the mock recorder for MockMachineProvisioner
type MockMachineProvisionerMockRecorder struct {
	mock *MockMachineProvisioner
}

// NewMockMachineProvisioner creates a new mock instance
func NewMockMachineProvisioner(ctrl *gomock.Controller) *MockMachineProvisioner {
	mock := &MockMachineProvisioner{ctrl: ctrl}
	mock.recorder = &MockMachineProvisionerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMachineProvisioner) EXPECT() *MockMachineProvisionerMockRecorder {
	return m.recorder
}

// AvailabilityZone mocks base method
func (m *MockMachineProvisioner) AvailabilityZone() (string, error) {
	ret := m.ctrl.Call(m, "AvailabilityZone")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AvailabilityZone indicates an expected call of AvailabilityZone
func (mr *MockMachineProvisionerMockRecorder) AvailabilityZone() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AvailabilityZone", reflect.TypeOf((*MockMachineProvisioner)(nil).AvailabilityZone))
}

// CharmProfileChangeInfo mocks base method
func (m *MockMachineProvisioner) CharmProfileChangeInfo() (provisioner.CharmProfileChangeInfo, error) {
	ret := m.ctrl.Call(m, "CharmProfileChangeInfo")
	ret0, _ := ret[0].(provisioner.CharmProfileChangeInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CharmProfileChangeInfo indicates an expected call of CharmProfileChangeInfo
func (mr *MockMachineProvisionerMockRecorder) CharmProfileChangeInfo() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CharmProfileChangeInfo", reflect.TypeOf((*MockMachineProvisioner)(nil).CharmProfileChangeInfo))
}

// DistributionGroup mocks base method
func (m *MockMachineProvisioner) DistributionGroup() ([]instance.Id, error) {
	ret := m.ctrl.Call(m, "DistributionGroup")
	ret0, _ := ret[0].([]instance.Id)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DistributionGroup indicates an expected call of DistributionGroup
func (mr *MockMachineProvisionerMockRecorder) DistributionGroup() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DistributionGroup", reflect.TypeOf((*MockMachineProvisioner)(nil).DistributionGroup))
}

// EnsureDead mocks base method
func (m *MockMachineProvisioner) EnsureDead() error {
	ret := m.ctrl.Call(m, "EnsureDead")
	ret0, _ := ret[0].(error)
	return ret0
}

// EnsureDead indicates an expected call of EnsureDead
func (mr *MockMachineProvisionerMockRecorder) EnsureDead() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnsureDead", reflect.TypeOf((*MockMachineProvisioner)(nil).EnsureDead))
}

// Id mocks base method
func (m *MockMachineProvisioner) Id() string {
	ret := m.ctrl.Call(m, "Id")
	ret0, _ := ret[0].(string)
	return ret0
}

// Id indicates an expected call of Id
func (mr *MockMachineProvisionerMockRecorder) Id() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Id", reflect.TypeOf((*MockMachineProvisioner)(nil).Id))
}

// InstanceId mocks base method
func (m *MockMachineProvisioner) InstanceId() (instance.Id, error) {
	ret := m.ctrl.Call(m, "InstanceId")
	ret0, _ := ret[0].(instance.Id)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InstanceId indicates an expected call of InstanceId
func (mr *MockMachineProvisionerMockRecorder) InstanceId() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceId", reflect.TypeOf((*MockMachineProvisioner)(nil).InstanceId))
}

// InstanceStatus mocks base method
func (m *MockMachineProvisioner) InstanceStatus() (status.Status, string, error) {
	ret := m.ctrl.Call(m, "InstanceStatus")
	ret0, _ := ret[0].(status.Status)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// InstanceStatus indicates an expected call of InstanceStatus
func (mr *MockMachineProvisionerMockRecorder) InstanceStatus() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InstanceStatus", reflect.TypeOf((*MockMachineProvisioner)(nil).InstanceStatus))
}

// KeepInstance mocks base method
func (m *MockMachineProvisioner) KeepInstance() (bool, error) {
	ret := m.ctrl.Call(m, "KeepInstance")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// KeepInstance indicates an expected call of KeepInstance
func (mr *MockMachineProvisionerMockRecorder) KeepInstance() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeepInstance", reflect.TypeOf((*MockMachineProvisioner)(nil).KeepInstance))
}

// Life mocks base method
func (m *MockMachineProvisioner) Life() params.Life {
	ret := m.ctrl.Call(m, "Life")
	ret0, _ := ret[0].(params.Life)
	return ret0
}

// Life indicates an expected call of Life
func (mr *MockMachineProvisionerMockRecorder) Life() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Life", reflect.TypeOf((*MockMachineProvisioner)(nil).Life))
}

// MachineTag mocks base method
func (m *MockMachineProvisioner) MachineTag() names_v2.MachineTag {
	ret := m.ctrl.Call(m, "MachineTag")
	ret0, _ := ret[0].(names_v2.MachineTag)
	return ret0
}

// MachineTag indicates an expected call of MachineTag
func (mr *MockMachineProvisionerMockRecorder) MachineTag() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MachineTag", reflect.TypeOf((*MockMachineProvisioner)(nil).MachineTag))
}

// MarkForRemoval mocks base method
func (m *MockMachineProvisioner) MarkForRemoval() error {
	ret := m.ctrl.Call(m, "MarkForRemoval")
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkForRemoval indicates an expected call of MarkForRemoval
func (mr *MockMachineProvisionerMockRecorder) MarkForRemoval() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkForRemoval", reflect.TypeOf((*MockMachineProvisioner)(nil).MarkForRemoval))
}

// ModelAgentVersion mocks base method
func (m *MockMachineProvisioner) ModelAgentVersion() (*version.Number, error) {
	ret := m.ctrl.Call(m, "ModelAgentVersion")
	ret0, _ := ret[0].(*version.Number)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ModelAgentVersion indicates an expected call of ModelAgentVersion
func (mr *MockMachineProvisionerMockRecorder) ModelAgentVersion() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModelAgentVersion", reflect.TypeOf((*MockMachineProvisioner)(nil).ModelAgentVersion))
}

// ProvisioningInfo mocks base method
func (m *MockMachineProvisioner) ProvisioningInfo() (*params.ProvisioningInfo, error) {
	ret := m.ctrl.Call(m, "ProvisioningInfo")
	ret0, _ := ret[0].(*params.ProvisioningInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProvisioningInfo indicates an expected call of ProvisioningInfo
func (mr *MockMachineProvisionerMockRecorder) ProvisioningInfo() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProvisioningInfo", reflect.TypeOf((*MockMachineProvisioner)(nil).ProvisioningInfo))
}

// Refresh mocks base method
func (m *MockMachineProvisioner) Refresh() error {
	ret := m.ctrl.Call(m, "Refresh")
	ret0, _ := ret[0].(error)
	return ret0
}

// Refresh indicates an expected call of Refresh
func (mr *MockMachineProvisionerMockRecorder) Refresh() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockMachineProvisioner)(nil).Refresh))
}

// Remove mocks base method
func (m *MockMachineProvisioner) Remove() error {
	ret := m.ctrl.Call(m, "Remove")
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove
func (mr *MockMachineProvisionerMockRecorder) Remove() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockMachineProvisioner)(nil).Remove))
}

// RemoveUpgradeCharmProfileData mocks base method
func (m *MockMachineProvisioner) RemoveUpgradeCharmProfileData() error {
	ret := m.ctrl.Call(m, "RemoveUpgradeCharmProfileData")
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUpgradeCharmProfileData indicates an expected call of RemoveUpgradeCharmProfileData
func (mr *MockMachineProvisionerMockRecorder) RemoveUpgradeCharmProfileData() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUpgradeCharmProfileData", reflect.TypeOf((*MockMachineProvisioner)(nil).RemoveUpgradeCharmProfileData))
}

// Series mocks base method
func (m *MockMachineProvisioner) Series() (string, error) {
	ret := m.ctrl.Call(m, "Series")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Series indicates an expected call of Series
func (mr *MockMachineProvisionerMockRecorder) Series() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Series", reflect.TypeOf((*MockMachineProvisioner)(nil).Series))
}

// SetCharmProfiles mocks base method
func (m *MockMachineProvisioner) SetCharmProfiles(arg0 []string) error {
	ret := m.ctrl.Call(m, "SetCharmProfiles", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCharmProfiles indicates an expected call of SetCharmProfiles
func (mr *MockMachineProvisionerMockRecorder) SetCharmProfiles(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCharmProfiles", reflect.TypeOf((*MockMachineProvisioner)(nil).SetCharmProfiles), arg0)
}

// SetInstanceInfo mocks base method
func (m *MockMachineProvisioner) SetInstanceInfo(arg0 instance.Id, arg1 string, arg2 *instance.HardwareCharacteristics, arg3 []params.NetworkConfig, arg4 []params.Volume, arg5 map[string]params.VolumeAttachmentInfo, arg6 []string) error {
	ret := m.ctrl.Call(m, "SetInstanceInfo", arg0, arg1, arg2, arg3, arg4, arg5, arg6)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetInstanceInfo indicates an expected call of SetInstanceInfo
func (mr *MockMachineProvisionerMockRecorder) SetInstanceInfo(arg0, arg1, arg2, arg3, arg4, arg5, arg6 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInstanceInfo", reflect.TypeOf((*MockMachineProvisioner)(nil).SetInstanceInfo), arg0, arg1, arg2, arg3, arg4, arg5, arg6)
}

// SetInstanceStatus mocks base method
func (m *MockMachineProvisioner) SetInstanceStatus(arg0 status.Status, arg1 string, arg2 map[string]interface{}) error {
	ret := m.ctrl.Call(m, "SetInstanceStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetInstanceStatus indicates an expected call of SetInstanceStatus
func (mr *MockMachineProvisionerMockRecorder) SetInstanceStatus(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetInstanceStatus", reflect.TypeOf((*MockMachineProvisioner)(nil).SetInstanceStatus), arg0, arg1, arg2)
}

// SetPassword mocks base method
func (m *MockMachineProvisioner) SetPassword(arg0 string) error {
	ret := m.ctrl.Call(m, "SetPassword", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPassword indicates an expected call of SetPassword
func (mr *MockMachineProvisionerMockRecorder) SetPassword(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPassword", reflect.TypeOf((*MockMachineProvisioner)(nil).SetPassword), arg0)
}

// SetStatus mocks base method
func (m *MockMachineProvisioner) SetStatus(arg0 status.Status, arg1 string, arg2 map[string]interface{}) error {
	ret := m.ctrl.Call(m, "SetStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetStatus indicates an expected call of SetStatus
func (mr *MockMachineProvisionerMockRecorder) SetStatus(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetStatus", reflect.TypeOf((*MockMachineProvisioner)(nil).SetStatus), arg0, arg1, arg2)
}

// SetSupportedContainers mocks base method
func (m *MockMachineProvisioner) SetSupportedContainers(arg0 ...instance.ContainerType) error {
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SetSupportedContainers", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetSupportedContainers indicates an expected call of SetSupportedContainers
func (mr *MockMachineProvisionerMockRecorder) SetSupportedContainers(arg0 ...interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetSupportedContainers", reflect.TypeOf((*MockMachineProvisioner)(nil).SetSupportedContainers), arg0...)
}

// SetUpgradeCharmProfileComplete mocks base method
func (m *MockMachineProvisioner) SetUpgradeCharmProfileComplete(arg0 string) error {
	ret := m.ctrl.Call(m, "SetUpgradeCharmProfileComplete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetUpgradeCharmProfileComplete indicates an expected call of SetUpgradeCharmProfileComplete
func (mr *MockMachineProvisionerMockRecorder) SetUpgradeCharmProfileComplete(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetUpgradeCharmProfileComplete", reflect.TypeOf((*MockMachineProvisioner)(nil).SetUpgradeCharmProfileComplete), arg0)
}

// Status mocks base method
func (m *MockMachineProvisioner) Status() (status.Status, string, error) {
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(status.Status)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Status indicates an expected call of Status
func (mr *MockMachineProvisionerMockRecorder) Status() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockMachineProvisioner)(nil).Status))
}

// String mocks base method
func (m *MockMachineProvisioner) String() string {
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String
func (mr *MockMachineProvisionerMockRecorder) String() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockMachineProvisioner)(nil).String))
}

// SupportsNoContainers mocks base method
func (m *MockMachineProvisioner) SupportsNoContainers() error {
	ret := m.ctrl.Call(m, "SupportsNoContainers")
	ret0, _ := ret[0].(error)
	return ret0
}

// SupportsNoContainers indicates an expected call of SupportsNoContainers
func (mr *MockMachineProvisionerMockRecorder) SupportsNoContainers() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SupportsNoContainers", reflect.TypeOf((*MockMachineProvisioner)(nil).SupportsNoContainers))
}

// Tag mocks base method
func (m *MockMachineProvisioner) Tag() names_v2.Tag {
	ret := m.ctrl.Call(m, "Tag")
	ret0, _ := ret[0].(names_v2.Tag)
	return ret0
}

// Tag indicates an expected call of Tag
func (mr *MockMachineProvisionerMockRecorder) Tag() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Tag", reflect.TypeOf((*MockMachineProvisioner)(nil).Tag))
}

// WatchAllContainers mocks base method
func (m *MockMachineProvisioner) WatchAllContainers() (watcher.StringsWatcher, error) {
	ret := m.ctrl.Call(m, "WatchAllContainers")
	ret0, _ := ret[0].(watcher.StringsWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchAllContainers indicates an expected call of WatchAllContainers
func (mr *MockMachineProvisionerMockRecorder) WatchAllContainers() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchAllContainers", reflect.TypeOf((*MockMachineProvisioner)(nil).WatchAllContainers))
}

// WatchContainers mocks base method
func (m *MockMachineProvisioner) WatchContainers(arg0 instance.ContainerType) (watcher.StringsWatcher, error) {
	ret := m.ctrl.Call(m, "WatchContainers", arg0)
	ret0, _ := ret[0].(watcher.StringsWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchContainers indicates an expected call of WatchContainers
func (mr *MockMachineProvisionerMockRecorder) WatchContainers(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchContainers", reflect.TypeOf((*MockMachineProvisioner)(nil).WatchContainers), arg0)
}

// WatchContainersCharmProfiles mocks base method
func (m *MockMachineProvisioner) WatchContainersCharmProfiles(arg0 instance.ContainerType) (watcher.StringsWatcher, error) {
	ret := m.ctrl.Call(m, "WatchContainersCharmProfiles", arg0)
	ret0, _ := ret[0].(watcher.StringsWatcher)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// WatchContainersCharmProfiles indicates an expected call of WatchContainersCharmProfiles
func (mr *MockMachineProvisionerMockRecorder) WatchContainersCharmProfiles(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchContainersCharmProfiles", reflect.TypeOf((*MockMachineProvisioner)(nil).WatchContainersCharmProfiles), arg0)
}
