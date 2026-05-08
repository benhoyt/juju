// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import (
	"testing"

	"github.com/juju/tc"
)

// volumeUUIDSuite is a suite of tests for asserting the behaviour of
// [VolumeUUID].
type volumeUUIDSuite struct{}

// TestVolumeUUIDSuite runs all of the tests contained within [volumeUUIDSuite].
func TestVolumeUUIDSuite(t *testing.T) {
	tc.Run(t, volumeUUIDSuite{})
}

// TestNew tests that constructing a new volume uuid succeeds with no errors and
// the end result is valid.
func (volumeUUIDSuite) TestNew(c *tc.C) {
	u, err := NewVolumeUUID()
	c.Check(err, tc.ErrorIsNil)
	c.Check(u.Validate(), tc.ErrorIsNil)
}

// TestStringer asserts the [fmt.Stringer] interface of [VolumeUUID] by making
// sure the correct string representation of the uuid is returned to the caller.
func (volumeUUIDSuite) TestStringer(c *tc.C) {
	const validUUID = "0de7ed80-bfcf-49b1-876a-31462e940ca1"
	c.Check(VolumeUUID(validUUID).String(), tc.Equals, validUUID)
}

// TestValidate asserts that a valid uuid passes validation with no errors.
func (volumeUUIDSuite) TestValidate(c *tc.C) {
	const validUUID = "0de7ed80-bfcf-49b1-876a-31462e940ca1"
	c.Check(VolumeUUID(validUUID).Validate(), tc.ErrorIsNil)
}

// TestValidateFail asserts that a bad uuid fails validation.
func (volumeUUIDSuite) TestValidateFail(c *tc.C) {
	c.Check(VolumeUUID("invalid").Validate(), tc.NotNil)
}
