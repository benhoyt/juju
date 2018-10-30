// Copyright 2018 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package upgrades_test

import (
	"io"
	"log"
	"path/filepath"
	"reflect"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
	"github.com/juju/replicaset"
	"github.com/juju/testing"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/juju/names.v2"

	"github.com/juju/juju/controller"
	"github.com/juju/juju/core/lease"
	"github.com/juju/juju/core/raftlease"
	"github.com/juju/juju/state"
	raftleasestore "github.com/juju/juju/state/raftlease"
	coretesting "github.com/juju/juju/testing"
	"github.com/juju/juju/upgrades"
	raftworker "github.com/juju/juju/worker/raft"
	"github.com/juju/juju/worker/raft/rafttest"
)

type raftSuite struct {
	coretesting.BaseSuite
}

var _ = gc.Suite(&raftSuite{})

func (s *raftSuite) TestBootstrapRaft(c *gc.C) {
	votes := 1
	noVotes := 0
	dataDir := c.MkDir()
	context := &mockContext{
		agentConfig: &mockAgentConfig{
			tag:     names.NewMachineTag("23"),
			dataDir: dataDir,
		},
		state: &mockState{
			members: []replicaset.Member{{
				Address: "somewhere.else:37012",
				Tags:    map[string]string{"juju-machine-id": "42"},
			}, {
				Address: "nowhere.else:37012",
				Tags:    map[string]string{"juju-machine-id": "23"},
				Votes:   &votes,
			}, {
				Address: "everywhere.else:37012",
				Tags:    map[string]string{"juju-machine-id": "7"},
				Votes:   &noVotes,
			}},
			info: state.StateServingInfo{APIPort: 1234},
		},
	}
	err := upgrades.BootstrapRaft(context)
	c.Assert(err, jc.ErrorIsNil)

	// Now make the raft node and check that the configuration is as
	// we expect.
	checkRaftConfiguration(c, dataDir)

	// Check the upgrade is idempotent.
	err = upgrades.BootstrapRaft(context)
	c.Assert(err, jc.ErrorIsNil)

	checkRaftConfiguration(c, dataDir)
}

func withRaft(c *gc.C, dataDir string, fsm raft.FSM, checkFunc func(*raft.Raft)) {
	// Capture logging to include in test output.
	output := captureWriter{c}
	config := raft.DefaultConfig()
	config.LocalID = "23"
	config.Logger = log.New(output, "", 0)
	c.Assert(raft.ValidateConfig(config), jc.ErrorIsNil)

	raftDir := filepath.Join(dataDir, "raft")
	logStore, err := raftboltdb.New(raftboltdb.Options{
		Path: filepath.Join(raftDir, "logs"),
	})
	c.Assert(err, jc.ErrorIsNil)
	defer logStore.Close()

	snapshotStore, err := raft.NewFileSnapshotStore(raftDir, 1, output)
	c.Assert(err, jc.ErrorIsNil)
	_, transport := raft.NewInmemTransport(raft.ServerAddress("nowhere.else"))

	r, err := raft.NewRaft(
		config,
		fsm,
		logStore,
		logStore,
		snapshotStore,
		transport,
	)
	c.Assert(err, jc.ErrorIsNil)
	defer func() {
		c.Assert(r.Shutdown().Error(), jc.ErrorIsNil)
	}()
	checkFunc(r)
}

func checkRaftConfiguration(c *gc.C, dataDir string) {
	withRaft(c, dataDir, &raftworker.SimpleFSM{},
		func(r *raft.Raft) {
			rafttest.CheckConfiguration(c, r, []raft.Server{{
				ID:       "42",
				Address:  "somewhere.else:1234",
				Suffrage: raft.Voter,
			}, {
				ID:       "23",
				Address:  "nowhere.else:1234",
				Suffrage: raft.Voter,
			}, {
				ID:       "7",
				Address:  "everywhere.else:1234",
				Suffrage: raft.Nonvoter,
			}})
		},
	)
}

func (s *raftSuite) TestMigrateLegacyLeases(c *gc.C) {
	dataDir := c.MkDir()
	context := &mockContext{
		agentConfig: &mockAgentConfig{
			tag:     names.NewMachineTag("23"),
			dataDir: dataDir,
		},
		state: &mockState{
			members: []replicaset.Member{{
				Address: "somewhere.else:37012",
				Tags:    map[string]string{"juju-machine-id": "42"},
			}},
			info: state.StateServingInfo{APIPort: 1234},
		},
	}
	err := upgrades.BootstrapRaft(context)
	c.Assert(err, jc.ErrorIsNil)

	var zero time.Time

	// Now we migrate some leases...
	config, err := controller.NewConfig(
		coretesting.ControllerTag.Id(),
		coretesting.CACert,
		nil,
	)
	c.Assert(err, jc.ErrorIsNil)
	var target mockTarget
	context.state = &mockState{
		config: config,
		leases: map[lease.Key]lease.Info{
			{"nonagon", "m1", "gamma"}: {
				Holder: "knife",
				Expiry: zero.Add(30 * time.Second),
			},
			{"reyne", "m2", "keening"}: {
				Holder: "jordan",
				Expiry: zero.Add(45 * time.Second),
			},
		},
		target: &target,
	}

	err = upgrades.MigrateLegacyLeases(context)
	c.Assert(err, jc.ErrorIsNil)

	target.assertClaimed(c, map[lease.Key]string{
		{"nonagon", "m1", "gamma"}: "knife",
		{"reyne", "m2", "keening"}: "jordan",
	})

	expectedLeases := context.state.(*mockState).leases

	// Start up raft with the leases in the snapshot.
	fsm := raftlease.NewFSM()
	withRaft(c, dataDir, fsm, func(r *raft.Raft) {
		// Once the snapshot is loaded the leases should be in the
		// FSM.
		var leases map[lease.Key]lease.Info
		for a := coretesting.LongAttempt.Start(); a.Next(); {
			leases = fsm.Leases(zero)
			if reflect.DeepEqual(leases, expectedLeases) {
				return
			}
		}
		c.Assert(leases, gc.DeepEquals, expectedLeases,
			gc.Commentf("waited %s but didn't see expected leases",
				coretesting.LongAttempt.Total))
	})

	// Check that running the step again works, but doesn't migrate
	// more leases.
	context.state = &mockState{
		config: config,
		leases: map[lease.Key]lease.Info{
			{"songs", "m3", "ferryman"}: {
				Holder: "jordan",
				Expiry: zero.Add(30 * time.Second),
			},
		},
		target: &target,
	}
	target.stub.ResetCalls()

	err = upgrades.MigrateLegacyLeases(context)
	c.Assert(err, jc.ErrorIsNil)

	c.Assert(target.stub.Calls(), gc.HasLen, 0)

	fsm = raftlease.NewFSM()
	withRaft(c, dataDir, fsm, func(r *raft.Raft) {
		var leases map[lease.Key]lease.Info
		for a := coretesting.LongAttempt.Start(); a.Next(); {
			leases = fsm.Leases(zero)
			// The leases haven't been updated with the changed one.
			if reflect.DeepEqual(leases, expectedLeases) {
				return
			}
		}
		c.Assert(leases, gc.DeepEquals, expectedLeases,
			gc.Commentf("waited %s but didn't see expected leases",
				coretesting.LongAttempt.Total))
	})

}

type mockState struct {
	upgrades.StateBackend
	stub    testing.Stub
	members []replicaset.Member
	info    state.StateServingInfo
	config  controller.Config
	leases  map[lease.Key]lease.Info
	target  *mockTarget
}

func (s *mockState) ReplicaSetMembers() ([]replicaset.Member, error) {
	return s.members, s.stub.NextErr()
}

func (s *mockState) StateServingInfo() (state.StateServingInfo, error) {
	return s.info, s.stub.NextErr()
}

func (s *mockState) ControllerConfig() (controller.Config, error) {
	return s.config, s.stub.NextErr()
}

func (s *mockState) LegacyLeases(localTime time.Time) (map[lease.Key]lease.Info, error) {
	s.stub.AddCall("LegacyLeases", localTime)
	return s.leases, s.stub.NextErr()
}

func (s *mockState) LeaseNotifyTarget(log io.Writer, errorLog raftleasestore.Logger) raftlease.NotifyTarget {
	s.stub.AddCall("LeaseNotifyTarget", log, errorLog)
	return s.target
}

type mockTarget struct {
	raftlease.NotifyTarget
	stub testing.Stub
}

func (t *mockTarget) Claimed(key lease.Key, holder string) {
	t.stub.AddCall("Claimed", key, holder)
}

func (t *mockTarget) assertClaimed(c *gc.C, claims map[lease.Key]string) {
	c.Assert(t.stub.Calls(), gc.HasLen, len(claims))
	for _, call := range t.stub.Calls() {
		c.Assert(call.FuncName, gc.Equals, "Claimed")
		c.Assert(call.Args, gc.HasLen, 2)
		key, ok := call.Args[0].(lease.Key)
		c.Assert(ok, gc.Equals, true)
		holder, ok := call.Args[1].(string)
		c.Assert(ok, gc.Equals, true)
		expectedHolder, found := claims[key]
		c.Assert(found, gc.Equals, true)
		c.Assert(holder, gc.Equals, expectedHolder)
		delete(claims, key)
	}
}

type captureWriter struct {
	c *gc.C
}

func (w captureWriter) Write(p []byte) (int, error) {
	w.c.Logf("%s", p[:len(p)-1]) // omit trailling newline
	return len(p), nil
}