// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package storage

import (
	"testing"

	"github.com/juju/tc"
)

// volumeAttachmentUUIDSuite is a suite of tests for asserting the behaviour of
// [VolumeAttachmentUUID].
type volumeAttachmentUUIDSuite struct{}

// TestVolumeAttachmentUUIDSuite runs all of the tests contained within
// [volumeAttachmentUUIDSuite].
func TestVolumeAttachmentUUIDSuite(t *testing.T) {
	tc.Run(t, volumeAttachmentUUIDSuite{})
}

// TestNew tests that constructing a new [VolumeAttachmentUUID] succeeds with no
// errors and the end result is valid.
func (volumeAttachmentUUIDSuite) TestNew(c *tc.C) {
	u, err := NewVolumeAttachmentUUID()
	c.Check(err, tc.ErrorIsNil)
	c.Check(u.Validate(), tc.ErrorIsNil)
}

// TestStringer asserts the [fmt.Stringer] interface of [VolumeAttachmentUUID]
// by making sure the correct string representation of the uuid is returned to
// the caller.
func (volumeAttachmentUUIDSuite) TestStringer(c *tc.C) {
	const validUUID = "0de7ed80-bfcf-49b1-876a-31462e940ca1"
	c.Check(VolumeAttachmentUUID(validUUID).String(), tc.Equals, validUUID)
}

// TestValidate asserts that a valid uuid passes validation with no errors.
func (volumeAttachmentUUIDSuite) TestValidate(c *tc.C) {
	const validUUID = "0de7ed80-bfcf-49b1-876a-31462e940ca1"
	c.Check(VolumeAttachmentUUID(validUUID).Validate(), tc.ErrorIsNil)
}

// TestValidateFail asserts that a bad uuid fails validation.
func (volumeAttachmentUUIDSuite) TestValidateFail(c *tc.C) {
	c.Check(VolumeAttachmentUUID("invalid").Validate(), tc.NotNil)
}
