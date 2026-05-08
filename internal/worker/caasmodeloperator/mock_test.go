// Copyright 2020 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package caasmodeloperator_test

import (
	"time"

	"github.com/juju/names/v6"

	"github.com/juju/juju/agent"
	coretesting "github.com/juju/juju/internal/testing"
)

type mockAgentConfig struct {
	agent.Config
}

func (m *mockAgentConfig) Controller() names.ControllerTag {
	return coretesting.ControllerTag
}

func (m *mockAgentConfig) DataDir() string {
	return "/var/lib/juju"
}

func (m *mockAgentConfig) LogDir() string {
	return "/var/log/juju"
}

func (m *mockAgentConfig) OldPassword() string {
	return "old password"
}

func (m *mockAgentConfig) CACert() string {
	return coretesting.CACert
}

func (m *mockAgentConfig) OpenTelemetryEnabled() bool {
	return false
}

func (m *mockAgentConfig) OpenTelemetryEndpoint() string {
	return ""
}

func (m *mockAgentConfig) OpenTelemetryInsecure() bool {
	return false
}

func (m *mockAgentConfig) OpenTelemetryStackTraces() bool {
	return false
}

func (m *mockAgentConfig) OpenTelemetrySampleRatio() float64 {
	return 0.1000
}

func (m *mockAgentConfig) OpenTelemetryTailSamplingThreshold() time.Duration {
	return time.Millisecond
}
