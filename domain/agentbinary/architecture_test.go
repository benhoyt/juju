// Copyright 2026 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package agentbinary

import (
	"testing"

	"github.com/juju/tc"
)

// architectureSuite defines a suite of tests for asserting the interface on
// offer for the [Architecture] type.
type architectureSuite struct{}

// TestArchitectureSuite runs all of the tests contained in [architectureSuite].
func TestArchitectureSuite(t *testing.T) {
	tc.Run(t, architectureSuite{})
}

// TestArchitecturesNotIn is a happy path test for [ArchitectureNotIn] to ensure
// that it correctly returns architectures in a that are not present in b.
func (architectureSuite) TestArchitecturesNotIn(c *tc.C) {
	a := []Architecture{S390X, PPC64EL}
	b := []Architecture{S390X, AMD64}

	val := ArchitectureNotIn(a, b)
	c.Check(val, tc.DeepEquals, []Architecture{PPC64EL})
}

// TestArchitecturesNotInEqual is a happy path test for [ArchitectureNotIn] to
// ensure that given two slices with the same values [ArchitectureNotIn] returns
// an empty result slice.
func (architectureSuite) TestArchitecturesNotInEqual(c *tc.C) {
	a := []Architecture{S390X, AMD64}
	b := []Architecture{AMD64, S390X}

	val := ArchitectureNotIn(a, b)
	c.Check(val, tc.HasLen, 0)
}

// TestArchitecturesNotInEmpty ensures that [ArchitectureNotIn] behaves
// correctly when either a or b are nil and that for a nilness is preserved.
func (architectureSuite) TestArchitecturesNotInEmpty(c *tc.C) {
	c.Run("a and b nil", func(t *testing.T) {
		val := ArchitectureNotIn(nil, nil)
		tc.Check(t, val, tc.IsNil)
	})

	c.Run("a nil", func(t *testing.T) {
		val := ArchitectureNotIn(nil, []Architecture{AMD64})
		tc.Check(t, val, tc.IsNil)
	})

	c.Run("b nil", func(t *testing.T) {
		val := ArchitectureNotIn([]Architecture{AMD64, ARM64}, nil)
		tc.Check(t, val, tc.DeepEquals, []Architecture{AMD64, ARM64})
	})
}

// TestFromString tests all of the well known valid architecture
// string values to make sure that they convert correctly to an [Architecture]
// type.
func (architectureSuite) TestFromString(c *tc.C) {
	tests := []struct {
		E Architecture
		V string
	}{
		{
			E: AMD64,
			V: "amd64",
		},
		{
			E: ARM64,
			V: "arm64",
		},
		{
			E: PPC64EL,
			V: "ppc64el",
		},
		{
			E: S390X,
			V: "s390x",
		},
		{
			E: RISCV64,
			V: "riscv64",
		},
	}

	for _, t := range tests {
		c.Run(t.V, func(c *testing.T) {
			a, converted := ArchitectureFromString(t.V)
			tc.Check(c, converted, tc.IsTrue)
			tc.Check(c, a, tc.Equals, t.E)
		})
	}
}

// TestFromStringUnknown tests that calling [ArchitectureFromString] with an
// unknown architecture string returns false and an invalid [Architecture]
// value.
func (architectureSuite) TestFromStringUnknown(c *tc.C) {
	val, converted := ArchitectureFromString("unknown")
	c.Check(converted, tc.IsFalse)
	c.Check(val.IsValid(), tc.IsFalse)
}

// TestIsValid checks all of the defined [Architecture] constants report that
// they are a valid value with [Architecture.IsValid].
func (architectureSuite) TestIsValid(c *tc.C) {
	tests := []Architecture{
		AMD64, ARM64, PPC64EL, S390X, RISCV64,
	}
	for _, t := range tests {
		c.Run(t.String(), func(c *testing.T) {
			tc.Check(c, t.IsValid(), tc.IsTrue)
		})
	}
}

// TestIsValidFail checks that an invalid [Architecture] value returns false for
// isValid.
func (architectureSuite) TestIsValidFail(c *tc.C) {
	c.Check(Architecture(-10).IsValid(), tc.IsFalse)
}

// TestToString tests all of the well known valid [Architecture]
// constants to make sure they correctly convert to the correct string value.
func (architectureSuite) TestToString(c *tc.C) {
	tests := []struct {
		E string
		V Architecture
	}{
		{
			E: "amd64",
			V: AMD64,
		},
		{
			E: "arm64",
			V: ARM64,
		},
		{
			E: "ppc64el",
			V: PPC64EL,
		},
		{
			E: "s390x",
			V: S390X,
		},
		{
			E: "riscv64",
			V: RISCV64,
		},
	}

	for _, t := range tests {
		c.Run(t.E, func(c *testing.T) {
			tc.Check(c, t.V.String(), tc.Equals, t.E)
		})
	}
}

// TestToStringInvalid checks that an invlaid [Architecture] values String
// method returns a zero value string when the value is invalid.
func (architectureSuite) TestToStringInvalid(c *tc.C) {
	c.Check(Architecture(-10).String(), tc.Equals, "")
}
