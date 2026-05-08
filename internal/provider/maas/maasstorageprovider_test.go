// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package maas

import (
	"testing"

	"github.com/juju/tc"

	coreerrors "github.com/juju/juju/core/errors"
	internalstorage "github.com/juju/juju/internal/storage"
)

// maasStorageProviderSuite provides a suite of tests for asserting the
// contracts and behaviour of the [maasStorageProvider].
type maasStorageProviderSuite struct{}

// TestMAASStorageProviderSuite runs the tests contained in
// [maasStorageProviderSuite].
func TestMAASStorageProviderSuite(t *testing.T) {
	tc.Run(t, &maasStorageProviderSuite{})
}

// TestTagsFromAttributesStringSlice tests that when
// [maasStorageProvider.TagsFromAttributes] processes a slice of string tags
// correctly.
func (maasStorageProviderSuite) TestTagsFromAttributesStringSlice(c *tc.C) {
	provider := &maasStorageProvider{}
	input := map[string]any{
		"tags": []string{"tag1", "tag2", "tag3"},
	}

	output, err := provider.TagsFromAttributes(input)
	c.Check(err, tc.ErrorIsNil)
	c.Check(output, tc.SameContents, []string{"tag1", "tag2", "tag3"})
}

// TestTagsFromAttributesAnyStringSlice tests that when
// [maasStorageProvider.TagsFromAttributes] processes a slice of any values that
// are strings the correct tag values are returned.
func (maasStorageProviderSuite) TestTagsFromAttributesAnyStringSlice(c *tc.C) {
	provider := &maasStorageProvider{}
	input := map[string]any{
		"tags": []any{"tag1", "tag2", "tag3"},
	}

	output, err := provider.TagsFromAttributes(input)
	c.Check(err, tc.ErrorIsNil)
	c.Check(output, tc.SameContents, []string{"tag1", "tag2", "tag3"})
}

// TestTagsFromAttributesStringSliceTrimsWhitespace tests that when
// [maasStorageProvider.TagsFromAttributes] processes a slice of strings tags
// that contain either leading or trailing whitespace, the whitespace is removed
// from the output result.
func (maasStorageProviderSuite) TestTagsFromAttributesStringSliceTrimsWhitespace(c *tc.C) {
	provider := &maasStorageProvider{}
	input := map[string]any{
		"tags": []string{"tag1", " spaceprefix", "tabsuffix	"},
	}

	output, err := provider.TagsFromAttributes(input)
	c.Check(err, tc.ErrorIsNil)
	c.Check(output, tc.SameContents, []string{"tag1", "spaceprefix", "tabsuffix"})
}

// TestTagsFromAttributesAnyStringSliceTrimsWhitespace tests that when
// [maasStorageProvider.TagsFromAttributes] processes a slice of any value that
// are string tags containing either leading or trailing whitespace, the
// whitespace is removed from the output result.
func (maasStorageProviderSuite) TestTagsFromAttributesAnyStringSliceTrimsWhitespace(c *tc.C) {
	provider := &maasStorageProvider{}
	input := map[string]any{
		"tags": []any{"tag1", " spaceprefix", "tabsuffix	"},
	}

	output, err := provider.TagsFromAttributes(input)
	c.Check(err, tc.ErrorIsNil)
	c.Check(output, tc.SameContents, []string{"tag1", "spaceprefix", "tabsuffix"})
}

// TestTagsFromAttributesAnySliceNotAllStrings tests that when
// [maasStorageProvider.TagsFromAttributes] processes a slice of any value that
// are not all string values the caller gets back an error satisfying
// [coreerrors.NotValid].
func (maasStorageProviderSuite) TestTagsFromAttributesAnySliceNotAllStrings(c *tc.C) {
	provider := &maasStorageProvider{}
	input := map[string]any{
		"tags": []any{"tag1", 123, "tag3"},
	}

	_, err := provider.TagsFromAttributes(input)
	c.Check(err, tc.ErrorIs, coreerrors.NotValid)
}

// TestTagsFromAttributesStringSliceWithWhitespace tests that when supplied with
// a string slice of tags and one of the tags contains whitespace that is not
// leading or trailing the caller gets back an error satisfying
// [coreerrors.NotSupported].
func (maasStorageProviderSuite) TestTagsFromAttributesStringSliceWithWhitespace(c *tc.C) {
	provider := &maasStorageProvider{}
	input := map[string]any{
		"tags": []string{"tag1", "has whitespace", "tag3"},
	}

	_, err := provider.TagsFromAttributes(input)
	c.Check(err, tc.ErrorIs, coreerrors.NotSupported)
}

// TestTagsFromAttributesAnySliceWithWhitespace tests that when supplied with
// a any slice of string tags and one of the tags contains whitespace that is
// not leading or trailing the caller gets back an error satisfying
// [coreerrors.NotSupported].
func (maasStorageProviderSuite) TestTagsFromAttributesAnySliceWithWhitespace(c *tc.C) {
	provider := &maasStorageProvider{}
	input := map[string]any{
		"tags": []any{"tag1", "has whitespace", "tag3"},
	}

	_, err := provider.TagsFromAttributes(input)
	c.Check(err, tc.ErrorIs, coreerrors.NotSupported)
}

// TestTagsFromAttributesStringSingle tests that when
// [maasStorageProvider.TagsFromAttributes] is given a single tag as a string
// value it is correctly returned in the slice of tags.
func (maasStorageProviderSuite) TestTagsFromAttributesStringSingle(c *tc.C) {
	provider := &maasStorageProvider{}
	input := map[string]any{
		"tags": "single🏷️",
	}

	output, err := provider.TagsFromAttributes(input)
	c.Check(err, tc.ErrorIsNil)
	c.Check(output, tc.SameContents, []string{"single🏷️"})
}

// TestTagsFromAttributesString tests that when
// [maasStorageProvider.TagsFromAttributes] is supplied a set of comma separated
// tags as a string values the slice of tags is correctly returned.
func (maasStorageProviderSuite) TestTagsFromAttributesString(c *tc.C) {
	provider := &maasStorageProvider{}
	input := map[string]any{
		"tags": "tag1,tag2,tag3",
	}

	output, err := provider.TagsFromAttributes(input)
	c.Check(err, tc.ErrorIsNil)
	c.Check(output, tc.SameContents, []string{"tag1", "tag2", "tag3"})
}

// TestTagsFromAttributesStringTrimsWhitespace tests that when
// [maasStorageProvider.TagsFromAttributes] is supplied a set of comma separated
// tags that contain either leading or trailing whitespace, the whitespace is
// removed from the output result.
func (maasStorageProviderSuite) TestTagsFromAttributesStringTrimsWhitespace(c *tc.C) {
	provider := &maasStorageProvider{}
	input := map[string]any{
		"tags": "tag1, spaceprefix,tabsuffix	",
	}

	output, err := provider.TagsFromAttributes(input)
	c.Check(err, tc.ErrorIsNil)
	c.Check(output, tc.SameContents, []string{"tag1", "spaceprefix", "tabsuffix"})
}

// TestTagsFromAttributesInt tests that when
// [maasStorageProvider.TagsFromAttributes] is supplied a non-string value (int),
// the caller gets back an error satisfying [coreerrors.NotValid].
func (maasStorageProviderSuite) TestTagsFromAttributesInt(c *tc.C) {
	provider := &maasStorageProvider{}
	input := map[string]any{
		"tags": 10,
	}

	_, err := provider.TagsFromAttributes(input)
	c.Check(err, tc.ErrorIs, coreerrors.NotValid)
}

// TestTagsFromAttributesStringWithWhitespace tests that when supplied with
// a string value containing tags that contains white space which is not leading
// or trailing the caller gets back an error satisfying
// [coreerrors.NotSupported].
func (maasStorageProviderSuite) TestTagsFromAttributesStringWithWhitespace(c *tc.C) {
	provider := &maasStorageProvider{}
	input := map[string]any{
		"tags": []any{"tag1,has whitespace,tag3"},
	}

	_, err := provider.TagsFromAttributes(input)
	c.Check(err, tc.ErrorIs, coreerrors.NotSupported)
}

// TestValidateConfigInvalid tests [maasStorageProvider.ValidateConfig] with a
// set of known invalid storage provider configurations to make sure validation
// fails with the correct error return.
func (maasStorageProviderSuite) TestValidateConfigInvalid(c *tc.C) {
	tests := []struct {
		Attributes    map[string]any
		ExpectedError error
		Reason        string
	}{
		{
			Attributes: map[string]any{
				"tags": 10,
			},
			ExpectedError: coreerrors.NotValid,
			Reason:        "tags must be a string value",
		},
		{
			Attributes: map[string]any{
				"tags": []int{10, -1},
			},
			ExpectedError: coreerrors.NotValid,
			Reason:        "tags must be a string slice",
		},
		{
			Attributes: map[string]any{
				"tags": []any{10, -1},
			},
			ExpectedError: coreerrors.NotValid,
			Reason:        "tags must be an any slice of strings",
		},
		{
			Attributes: map[string]any{
				"tags": "has space",
			},
			ExpectedError: coreerrors.NotSupported,
			Reason:        "tag string must not contain white space",
		},
		{
			Attributes: map[string]any{
				"tags": []string{"has whitespace"},
			},
			ExpectedError: coreerrors.NotSupported,
			Reason:        "tag string slice must not contain white space",
		},
	}

	for _, test := range tests {
		c.Run(test.Reason, func(t *testing.T) {
			provider := maasStorageProvider{}
			cfg, err := internalstorage.NewConfig(
				"maasstorage", "maas", test.Attributes,
			)
			c.Assert(err, tc.ErrorIsNil)
			err = provider.ValidateConfig(cfg)
			c.Check(err, tc.ErrorIs, test.ExpectedError)
		})
	}
}

// TestValidateConfigValid tests [maasStorageProvider.ValidateConfig] with a
// set of known valid storage provider configurations to make sure validation
// succeeds.
func (maasStorageProviderSuite) TestValidateConfigValid(c *tc.C) {
	tests := []struct {
		Attributes map[string]any
		Name       string
	}{
		{
			Attributes: map[string]any{},
			Name:       "no tags",
		},
		{
			Attributes: map[string]any{
				"tags": "tag1",
			},
			Name: "single string tag",
		},
		{
			Attributes: map[string]any{
				"tags": []string{"tag1"},
			},
			Name: "single string slice tag",
		},
		{
			Attributes: map[string]any{
				"tags": []any{"tag1"},
			},
			Name: "single any slice string tag",
		},
		{
			Attributes: map[string]any{
				"tags": "tag1, leadingspace,trailingtab	",
			},
			Name: "many string tags",
		},
		{
			Attributes: map[string]any{
				"tags": []string{"tag1", " leadingspace", "trailingtab	"},
			},
			Name: "many string slice tags",
		},
		{
			Attributes: map[string]any{
				"tags": []any{"tag1", " leadingspace", "trailingtab	"},
			},
			Name: "many any slice tags",
		},
		{
			Attributes: map[string]any{
				"unknown": 10,
			},
			Name: "unknown attributes are ignored",
		},
	}

	for _, test := range tests {
		c.Run(test.Name, func(t *testing.T) {
			provider := maasStorageProvider{}
			cfg, err := internalstorage.NewConfig(
				"maasstorage", "maas", test.Attributes,
			)
			c.Assert(err, tc.ErrorIsNil)
			err = provider.ValidateConfig(cfg)
			c.Check(err, tc.ErrorIsNil)
		})
	}
}

// TestSupportsAllStorageKindsFalse asserts that the [maasStorageProvider] does
// not support any of the available storage kinds. This test is important as it
// ensures that a storage pool made from this provider cannot be used to create
// storage for a charm.
//
// [maasStorageProvider] only supports creating storage for root disks of new
// machines being provisioned in MAAS.
func (maasStorageProviderSuite) TestSupportsAllStorageKindsFalse(c *tc.C) {
	c.Run("block", func(c *testing.T) {
		p := maasStorageProvider{}
		tc.Assert(c, p.Supports(internalstorage.StorageKindBlock), tc.IsFalse)
	})
	c.Run("filesystem", func(c *testing.T) {
		p := maasStorageProvider{}
		tc.Assert(c, p.Supports(internalstorage.StorageKindFilesystem), tc.IsFalse)
	})
}

// TestScope asserts that the [maasStorageProvider] scope is
// always [internalstorage.ScopeEnviron].
func (maasStorageProviderSuite) TestScope(c *tc.C) {
	p := maasStorageProvider{}
	c.Check(p.Scope(), tc.Equals, internalstorage.ScopeEnviron)
}

// TestFilesystemSourceNotSupported asserts that the [maasStorageProvider] does
// not support supplying a filesystem source.
func (maasStorageProviderSuite) TestFilesystemSourceNotSupported(c *tc.C) {
	p := maasStorageProvider{}
	_, err := p.FilesystemSource(nil)
	c.Check(err, tc.ErrorIs, coreerrors.NotSupported)
}

// TestVolumeSourceNotSupported asserts that the [maasStorageProvider] does
// not support supplying a volume source.
func (maasStorageProviderSuite) TestVolumeSourceNotSupported(c *tc.C) {
	p := maasStorageProvider{}
	_, err := p.VolumeSource(nil)
	c.Check(err, tc.ErrorIs, coreerrors.NotSupported)
}
