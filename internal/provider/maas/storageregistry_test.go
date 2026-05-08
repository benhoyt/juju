// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package maas

import (
	"testing"

	"github.com/juju/tc"

	internalstorage "github.com/juju/juju/internal/storage"
)

// storageRegistrySuite is a testing suite for asserting the behaviour of the
// [internalstorage.ProviderRegistry] implementation on the maas environ.
type storageRegistrySuite struct {
	maasSuite
}

func TestStorageRegistrySuite(t *testing.T) {
	tc.Run(t, &storageRegistrySuite{})
}

// TestRecommendedPoolForKind is a regression test to show that existing
// behavior of storage pool selection is maintained between 3.x and 4.x
// releases. This bug was introduced with the move in 4.0 where storage
// registries supply default storage pools to the model and also recommend
// storage pools to use for different storage kinds
//
// With this change we ended up recommending the maas storage provider for
// filesystem charm storage when the provider in not capable of provisioning
// any storage on behalf of a charm. In 3.x the recommended provider for
// filesystems on a MAAS cloud would have always been rootfs.
//
// This test ensures that the correct recommended pool and provider is always
// being recommended. While the MAAS storage provider is unable to provision
// storage for charms it SHOULD never be returned as recommended.
func (s *storageRegistrySuite) TestRecommendedPoolForKind(c *tc.C) {
	env := s.maasSuite.makeEnviron(c, newFakeController())

	c.Run("filesystem", func(c *testing.T) {
		bPool := env.RecommendedPoolForKind(internalstorage.StorageKindFilesystem)
		tc.Check(c, bPool.Name(), tc.Equals, "rootfs")
		tc.Check(c, bPool.Provider().String(), tc.Equals, "rootfs")
	})

	c.Run("block", func(c *testing.T) {
		bPool := env.RecommendedPoolForKind(internalstorage.StorageKindBlock)
		tc.Check(c, bPool.Name(), tc.Equals, "loop")
		tc.Check(c, bPool.Provider().String(), tc.Equals, "loop")
	})
}
