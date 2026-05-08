// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package errors

import "github.com/juju/juju/internal/errors"

const (
	// AlreadyExists describes an error that occurs when a cloud with the same
	// name already exists.
	AlreadyExists = errors.ConstError("cloud already exists")

	// NotFound describes an error that occurs when the cloud being operated on
	// does not exist.
	NotFound = errors.ConstError("cloud not found")

	// CloudStillInUse describes an error that occurs when attempting to delete
	// a cloud that is still referenced by one or more models.
	CloudStillInUse = errors.ConstError("cloud still in use")
)
