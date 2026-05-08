// Copyright 2022 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package secrets

import "time"

// RotatePolicy defines a policy for how often
// to rotate a secret.
type RotatePolicy string

const (
	RotateNever     = RotatePolicy("never")
	RotateHourly    = RotatePolicy("hourly")
	RotateDaily     = RotatePolicy("daily")
	RotateWeekly    = RotatePolicy("weekly")
	RotateMonthly   = RotatePolicy("monthly")
	RotateQuarterly = RotatePolicy("quarterly")
	RotateYearly    = RotatePolicy("yearly")
)

const (
	// RotateRetryDelay is how long to wait to re-run the rotate hook
	// if the secret was not updated.
	RotateRetryDelay = 5 * time.Minute

	// ExpireRetryDelay is how long to wait to re-run the expire hook
	// if the expired secret revision was not removed.
	ExpireRetryDelay = 5 * time.Minute
)

func (p RotatePolicy) String() string {
	if p == "" {
		return string(RotateNever)
	}
	return string(p)
}

// WillRotate returns true if the policy is not RotateNever.
func (p *RotatePolicy) WillRotate() bool {
	return p != nil && *p != "" && *p != RotateNever
}

// IsValid returns true if p is a valid rotate policy.
func (p RotatePolicy) IsValid() bool {
	switch p {
	case RotateNever, RotateHourly, RotateDaily, RotateWeekly,
		RotateMonthly, RotateQuarterly, RotateYearly:
		return true
	}
	return false
}

// NextRotateTime returns when the policy dictates a secret should be next
// rotated given the last rotation time.
func (p RotatePolicy) NextRotateTime(lastRotated time.Time) *time.Time {
	var result time.Time
	switch p {
	case RotateNever:
		return nil
	case RotateHourly:
		result = lastRotated.Add(time.Hour)
	case RotateDaily:
		result = lastRotated.AddDate(0, 0, 1)
	case RotateWeekly:
		result = lastRotated.AddDate(0, 0, 7)
	case RotateMonthly:
		result = lastRotated.AddDate(0, 1, 0)
	case RotateQuarterly:
		result = lastRotated.AddDate(0, 3, 0)
	case RotateYearly:
		result = lastRotated.AddDate(1, 0, 0)
	}
	return &result
}

// LessThan returns true if the policy is more frequent than other.
func (p *RotatePolicy) LessThan(other RotatePolicy) bool {
	if !p.WillRotate() && !other.WillRotate() {
		return false // those are basically equal
	} else if !p.WillRotate() {
		return false // current won't rotate make it less frequent
	} else if !other.WillRotate() {
		return true // if the other doesn't rotate, current will be more frequent
	}
	now := time.Now()
	return p.NextRotateTime(now).Before(*other.NextRotateTime(now))
}
