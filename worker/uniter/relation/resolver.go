// Copyright 2012-2015 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package relation

import (
	"github.com/juju/collections/set"
	"github.com/juju/errors"
	"github.com/juju/loggo"
	"gopkg.in/juju/charm.v6/hooks"
	"gopkg.in/juju/names.v3"

	"github.com/juju/juju/core/life"
	"github.com/juju/juju/worker/uniter/hook"
	"github.com/juju/juju/worker/uniter/operation"
	"github.com/juju/juju/worker/uniter/remotestate"
	"github.com/juju/juju/worker/uniter/resolver"
)

var logger = loggo.GetLogger("juju.worker.uniter.relation")

//go:generate mockgen -package mocks -destination mocks/mock_subordinate_destroyer.go github.com/juju/juju/worker/uniter/relation SubordinateDestroyer

// SubordinateDestroyer destroys all subordinates of a unit.
type SubordinateDestroyer interface {
	DestroyAllSubordinates() error
}

// NewRelationResolver returns a resolver that handles all relation-related
// hooks (except relation-created) and is wired to the provided RelationStateTracker
// instance.
func NewRelationResolver(stateTracker RelationStateTracker, subordinateDestroyer SubordinateDestroyer) resolver.Resolver {
	return &relationsResolver{
		stateTracker:         stateTracker,
		subordinateDestroyer: subordinateDestroyer,
	}
}

type relationsResolver struct {
	stateTracker         RelationStateTracker
	subordinateDestroyer SubordinateDestroyer
}

// NextOp implements resolver.Resolver.
func (r *relationsResolver) NextOp(localState resolver.LocalState, remoteState remotestate.Snapshot, opFactory operation.Factory) (operation.Operation, error) {
	if err := r.maybeDestroySubordinates(remoteState); err != nil {
		return nil, errors.Trace(err)
	}

	if err := r.stateTracker.SynchronizeScopes(remoteState); err != nil {
		return nil, errors.Trace(err)
	}

	if localState.Kind != operation.Continue {
		return nil, resolver.ErrNoOperation
	}

	// Check whether we need to fire a hook for any of the relations
	for relationId, relationSnapshot := range remoteState.Relations {
		if !r.stateTracker.IsKnown(relationId) {
			continue
		} else if isImplicit, _ := r.stateTracker.IsImplicit(relationId); isImplicit {
			continue
		}

		// If either the unit or the relation are Dying, or the
		// relation becomes suspended, then the relation should be
		// broken.
		var remoteBroken bool
		if remoteState.Life == life.Dying || relationSnapshot.Life == life.Dying || relationSnapshot.Suspended {
			relationSnapshot = remotestate.RelationSnapshot{}
			remoteBroken = true
			// TODO(axw) if relation is implicit, leave scope & remove.
		}

		// Examine local/remote states and figure out if a hook needs
		// to be fired for this relation.
		stateDir, err := r.stateTracker.StateDir(relationId)
		if err != nil {
			return nil, errors.Trace(err)
		}
		hook, err := r.nextHookForRelation(stateDir, relationSnapshot, remoteBroken)
		if err == resolver.ErrNoOperation {
			continue
		}
		return opFactory.NewRunHook(hook)
	}

	return nil, resolver.ErrNoOperation
}

// maybeDestroySubordinates checks whether the remote state indicates that the
// unit is dying and ensures that any related subordinates are properly
// destroyed.
func (r *relationsResolver) maybeDestroySubordinates(remoteState remotestate.Snapshot) error {
	if remoteState.Life != life.Dying {
		return nil
	}

	var destroyAllSubordinates bool
	for relationId, relationSnapshot := range remoteState.Relations {
		if relationSnapshot.Life != life.Alive {
			continue
		} else if hasContainerScope, err := r.stateTracker.HasContainerScope(relationId); err != nil || !hasContainerScope {
			continue
		}

		// Found alive relation to a subordinate
		relationSnapshot.Life = life.Dying
		remoteState.Relations[relationId] = relationSnapshot
		destroyAllSubordinates = true
	}

	if destroyAllSubordinates {
		return r.subordinateDestroyer.DestroyAllSubordinates()
	}

	return nil
}

func (r *relationsResolver) nextHookForRelation(localStateDir *StateDir, remote remotestate.RelationSnapshot, remoteBroken bool) (hook.Info, error) {
	// If there's a guaranteed next hook, return that.
	local := localStateDir.State()
	relationId := local.RelationId
	if local.ChangedPending != "" {
		// ChangedPending should only happen for a unit (not an app). It is a side effect that if we call 'relation-joined'
		// for a unit, we immediately queue up relation-changed for that unit, before we run any other hooks
		// Applications never see "relation-joined".
		unitName := local.ChangedPending
		appName, err := names.UnitApplication(unitName)
		if err != nil {
			return hook.Info{}, errors.Annotate(err, "changed pending held an invalid unit name")
		}
		return hook.Info{
			Kind:              hooks.RelationChanged,
			RelationId:        relationId,
			RemoteUnit:        unitName,
			RemoteApplication: appName,
			ChangeVersion:     remote.Members[unitName],
		}, nil
	}

	// Get related app names, trigger all app hooks first
	allAppNames := set.NewStrings()
	for appName := range local.ApplicationMembers {
		allAppNames.Add(appName)
	}
	for app := range remote.ApplicationMembers {
		allAppNames.Add(app)
	}
	sortedAppNames := allAppNames.SortedValues()

	// Get the union of all relevant units, and sort them, so we produce events
	// in a consistent order (largely for the convenience of the tests).
	allUnitNames := set.NewStrings()
	for unitName := range local.Members {
		allUnitNames.Add(unitName)
	}
	for unitName := range remote.Members {
		allUnitNames.Add(unitName)
	}
	sortedUnitNames := allUnitNames.SortedValues()
	if allUnitNames.Contains("") {
		return hook.Info{}, errors.Errorf("somehow we got the empty unit. local: %v, remote: %v", local.Members, remote.Members)
	}

	// If there are any locally known units that are no longer reflected in
	// remote state, depart them.
	for _, unitName := range sortedUnitNames {
		changeVersion, found := local.Members[unitName]
		if !found {
			continue
		}
		if _, found := remote.Members[unitName]; !found {
			appName, err := names.UnitApplication(unitName)
			if err != nil {
				return hook.Info{}, errors.Trace(err)
			}

			// Consult the life of the local unit and/or app to
			// figure out if its the local or the remote unit going
			// away. Note that if the app is removed, the unit will
			// still be alive but its parent app will by dying.
			localUnitLife, localAppLife, err := r.stateTracker.LocalUnitAndApplicationLife()
			if err != nil {
				return hook.Info{}, errors.Trace(err)
			}

			var departee = unitName
			if localUnitLife != life.Alive || localAppLife != life.Alive {
				departee = r.stateTracker.LocalUnitName()
			}

			return hook.Info{
				Kind:              hooks.RelationDeparted,
				RelationId:        relationId,
				RemoteUnit:        unitName,
				RemoteApplication: appName,
				ChangeVersion:     changeVersion,
				DepartingUnit:     departee,
			}, nil
		}
	}

	// If the relation's meant to be broken, break it. A side-effect of
	// the logic that generates the relation-created hooks is that we may
	// end up in this block for a peer relation.  Since you cannot depart
	// peer relations we can safely ignore this hook.
	isPeer, _ := r.stateTracker.IsPeerRelation(relationId)
	if remoteBroken && !isPeer {
		if !localStateDir.Exists() {
			// The relation may have been suspended and then
			// removed, so we don't want to run the hook twice.
			return hook.Info{}, resolver.ErrNoOperation
		}

		return hook.Info{
			Kind:              hooks.RelationBroken,
			RelationId:        relationId,
			RemoteApplication: r.stateTracker.RemoteApplication(relationId),
		}, nil
	}

	for _, appName := range sortedAppNames {
		changeVersion, found := remote.ApplicationMembers[appName]
		if !found {
			// ?
			continue
		}
		// Note(jam): 2019-10-23 For compatibility purposes, we don't trigger a hook if
		//  local.ApplicationMembers doesn't contain the app and the changeVersion == 0.
		//  This is because otherwise all charms always get a hook with the app
		//  as the context, and that is likely to expose them to something they
		//  may not be ready for. Also, since no app content has been set, there
		//  is nothing for them to respond to.
		if oldVersion := local.ApplicationMembers[appName]; oldVersion != changeVersion {
			return hook.Info{
				Kind:              hooks.RelationChanged,
				RelationId:        relationId,
				RemoteUnit:        "",
				RemoteApplication: appName,
				ChangeVersion:     changeVersion,
			}, nil
		}
	}

	// If there are any remote units not locally known, join them.
	for _, unitName := range sortedUnitNames {
		changeVersion, found := remote.Members[unitName]
		if !found {
			continue
		}
		if _, found := local.Members[unitName]; !found {
			appName, err := names.UnitApplication(unitName)
			if err != nil {
				return hook.Info{}, errors.Trace(err)
			}
			return hook.Info{
				Kind:              hooks.RelationJoined,
				RelationId:        relationId,
				RemoteUnit:        unitName,
				RemoteApplication: appName,
				ChangeVersion:     changeVersion,
			}, nil
		}
	}

	// Finally scan for remote units whose latest version is not reflected
	// in local state.
	for _, unitName := range sortedUnitNames {
		remoteChangeVersion, found := remote.Members[unitName]
		if !found {
			continue
		}
		localChangeVersion, found := local.Members[unitName]
		if !found {
			continue
		}
		appName, err := names.UnitApplication(unitName)
		if err != nil {
			return hook.Info{}, errors.Trace(err)
		}
		// NOTE(axw) we use != and not > to cater due to the
		// use of the relation settings document's txn-revno
		// as the version. When model-uuid migration occurs, the
		// document is recreated, resetting txn-revno.
		if remoteChangeVersion != localChangeVersion {
			return hook.Info{
				Kind:              hooks.RelationChanged,
				RelationId:        relationId,
				RemoteUnit:        unitName,
				RemoteApplication: appName,
				ChangeVersion:     remoteChangeVersion,
			}, nil
		}
	}

	// Nothing left to do for this relation.
	return hook.Info{}, resolver.ErrNoOperation
}

// NewCreatedRelationResolver returns a resolver that handles relation-created
// hooks and is wired to the provided RelationStateTracker instance.
func NewCreatedRelationResolver(stateTracker RelationStateTracker) resolver.Resolver {
	return &createdRelationsResolver{
		stateTracker: stateTracker,
	}
}

type createdRelationsResolver struct {
	stateTracker RelationStateTracker
}

// NextOp implements resolver.Resolver.
func (r *createdRelationsResolver) NextOp(localState resolver.LocalState, remoteState remotestate.Snapshot, opFactory operation.Factory) (operation.Operation, error) {
	// Nothing to do if not yet installed or if the unit is dying.
	if !localState.Installed || remoteState.Life == life.Dying {
		return nil, resolver.ErrNoOperation
	}

	if err := r.stateTracker.SynchronizeScopes(remoteState); err != nil {
		return nil, errors.Trace(err)
	}

	for relationId, relationSnapshot := range remoteState.Relations {
		if relationSnapshot.Life != life.Alive {
			continue
		}

		hook, err := r.nextHookForRelation(relationId, relationSnapshot)
		if err != nil {
			if err == resolver.ErrNoOperation {
				continue
			}

			return nil, errors.Trace(err)
		}

		return opFactory.NewRunHook(hook)
	}

	return nil, resolver.ErrNoOperation
}

func (r *createdRelationsResolver) nextHookForRelation(relationId int, relationSnapshot remotestate.RelationSnapshot) (hook.Info, error) {
	isImplicit, _ := r.stateTracker.IsImplicit(relationId)
	if r.stateTracker.RelationCreated(relationId) || isImplicit {
		return hook.Info{}, resolver.ErrNoOperation
	}

	return hook.Info{
		Kind:              hooks.RelationCreated,
		RelationId:        relationId,
		RemoteApplication: r.stateTracker.RemoteApplication(relationId),
	}, nil
}
