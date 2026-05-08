// Copyright 2022 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package base

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/juju/tc"

	"github.com/juju/juju/internal/testhelpers"
)

type ChannelSuite struct {
	testhelpers.IsolationSuite
}

func TestChannelSuite(t *testing.T) {
	tc.Run(t, &ChannelSuite{})
}

func (s *ChannelSuite) TestParse(c *tc.C) {
	ch, err := ParseChannel("22.04")
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(ch, tc.DeepEquals, Channel{Track: "22.04"})
	ch, err = ParseChannel("22.04/edge")
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(ch, tc.DeepEquals, Channel{Track: "22.04", Risk: "edge"})
	ch, err = ParseChannel("all")
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(ch, tc.DeepEquals, Channel{Track: "all"})
}

func (s *ChannelSuite) TestParseError(c *tc.C) {
	_, err := ParseChannel("22.04/edge/foo")
	c.Assert(err, tc.ErrorMatches, `channel is malformed and has too many components "22.04/edge/foo"`)
	_, err = ParseChannel("22.04/foo")
	c.Assert(err, tc.ErrorMatches, `risk in channel "22.04/foo" not valid`)
}

func (s *ChannelSuite) TestParseNormalise(c *tc.C) {
	ch, err := ParseChannelNormalize("22.04")
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(ch, tc.DeepEquals, Channel{Track: "22.04", Risk: "stable"})
	ch, err = ParseChannelNormalize("22.04/edge")
	c.Assert(err, tc.ErrorIsNil)
	c.Assert(ch, tc.DeepEquals, Channel{Track: "22.04", Risk: "edge"})
}

func (s *ChannelSuite) TestMakeDefaultChannel(c *tc.C) {
	ch := MakeDefaultChannel("22.04")
	c.Assert(ch, tc.DeepEquals, Channel{Track: "22.04", Risk: "stable"})
}

func (s *ChannelSuite) TestString(c *tc.C) {
	c.Assert(Channel{Track: "22.04"}.String(), tc.Equals, "22.04")
	c.Assert(Channel{Track: "22.04", Risk: "edge"}.String(), tc.Equals, "22.04/edge")
}

func (s *ChannelSuite) TestDisplayString(c *tc.C) {
	c.Assert(Channel{Track: "18.04"}.DisplayString(), tc.Equals, "18.04")
	c.Assert(Channel{Track: "20.04", Risk: "stable"}.DisplayString(), tc.Equals, "20.04")
	c.Assert(Channel{Track: "22.04", Risk: "edge"}.DisplayString(), tc.Equals, "22.04/edge")
}

func (s *ChannelSuite) TestHasHigherPriorityThan(c *tc.C) {
	supportedLTSTrack := DefaultSupportedLTSBase().Channel.Track

	// Split base track to get year and month.
	parts := strings.Split(supportedLTSTrack, ".")
	c.Assert(len(parts), tc.Equals, 2)
	ltsYear, err := strconv.Atoi(parts[0])
	c.Assert(err, tc.ErrorIsNil)

	// Create future LTS base track.
	futureTrack := fmt.Sprintf("%02d.%s", ltsYear+2, parts[1])

	tests := []struct {
		name     string
		current  Channel
		other    Channel
		expected bool
		reason   string
	}{
		{
			name:     "LTS base track has highest priority",
			current:  Channel{Track: supportedLTSTrack, Risk: Stable},
			other:    Channel{Track: futureTrack, Risk: Stable},
			expected: true,
			reason:   "LTS base track should be preferred",
		},
		{
			name:     "non-LTS base track lower priority than LTS",
			current:  Channel{Track: futureTrack, Risk: Stable},
			other:    Channel{Track: supportedLTSTrack, Risk: Stable},
			expected: false,
			reason:   "non-LTS should not be preferred over LTS",
		},
		{
			name:     "higher version base track preferred when neither is LTS",
			current:  Channel{Track: "22.04", Risk: Stable},
			other:    Channel{Track: "20.04", Risk: Stable},
			expected: true,
			reason:   "22.04 > 20.04",
		},
		{
			name:     "lower version base track not preferred when neither is LTS",
			current:  Channel{Track: "20.04", Risk: Stable},
			other:    Channel{Track: "22.04", Risk: Stable},
			expected: false,
			reason:   "20.04 < 22.04",
		},
		{
			name:     "stable risk preferred over candidate for same base track",
			current:  Channel{Track: "22.04", Risk: Stable},
			other:    Channel{Track: "22.04", Risk: Candidate},
			expected: true,
			reason:   "stable > candidate",
		},
		{
			name:     "stable risk preferred over beta for same base track",
			current:  Channel{Track: "22.04", Risk: Stable},
			other:    Channel{Track: "22.04", Risk: Beta},
			expected: true,
			reason:   "stable > beta",
		},
		{
			name:     "stable risk preferred over edge for same base track",
			current:  Channel{Track: "22.04", Risk: Stable},
			other:    Channel{Track: "22.04", Risk: Edge},
			expected: true,
			reason:   "stable > edge",
		},
		{
			name:     "candidate risk preferred over beta for same base track",
			current:  Channel{Track: "22.04", Risk: Candidate},
			other:    Channel{Track: "22.04", Risk: Beta},
			expected: true,
			reason:   "candidate > beta",
		},
		{
			name:     "candidate risk preferred over edge for same base track",
			current:  Channel{Track: "22.04", Risk: Candidate},
			other:    Channel{Track: "22.04", Risk: Edge},
			expected: true,
			reason:   "candidate > edge",
		},
		{
			name:     "beta risk preferred over edge for same base track",
			current:  Channel{Track: "22.04", Risk: Beta},
			other:    Channel{Track: "22.04", Risk: Edge},
			expected: true,
			reason:   "beta > edge",
		},
		{
			name:     "edge risk not preferred over stable for same base track",
			current:  Channel{Track: "22.04", Risk: Edge},
			other:    Channel{Track: "22.04", Risk: Stable},
			expected: false,
			reason:   "edge < stable",
		},
		{
			name:     "base track priority overrides risk, higher base track is preferred even with worse risk",
			current:  Channel{Track: "22.04", Risk: Edge},
			other:    Channel{Track: "20.04", Risk: Stable},
			expected: true,
			reason:   "base track version takes precedence over risk",
		},
		{
			name:     "LTS base track with edge risk preferred over non-LTS with stable risk",
			current:  Channel{Track: supportedLTSTrack, Risk: Edge},
			other:    Channel{Track: futureTrack, Risk: Stable},
			expected: true,
			reason:   "LTS base track priority overrides risk differences",
		},
		{
			name:     "LTS base track with stable risk preferred over LTS with edge risk",
			current:  Channel{Track: supportedLTSTrack, Risk: Stable},
			other:    Channel{Track: supportedLTSTrack, Risk: Edge},
			expected: true,
			reason:   "LTS base track priority overrides risk differences",
		},
		{
			name:     "same base track and risk returns false",
			current:  Channel{Track: "22.04", Risk: Stable},
			other:    Channel{Track: "22.04", Risk: Stable},
			expected: false,
			reason:   "identical channels have no priority difference",
		},
	}

	for _, test := range tests {
		result := test.current.HasHigherPriorityThan(test.other)
		c.Check(result, tc.Equals, test.expected, tc.Commentf("test '%s' failed with reason: %s", test.name, test.reason))
	}
}
