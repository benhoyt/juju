// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package migration

//go:generate go run go.uber.org/mock/mockgen -typed -package migration -destination package_mock_test.go github.com/juju/juju/domain/model/service/migration State
