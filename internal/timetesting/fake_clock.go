// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timetesting

import (
	"time"
)

type FakeClock struct {
	now time.Time
}

func NewFakeClock(now time.Time) *FakeClock {
	return &FakeClock{
		now: now,
	}
}

func (clock *FakeClock) Now() time.Time {
	return clock.now
}

func (clock *FakeClock) Since(t time.Time) time.Duration {
	return clock.Now().Sub(t)
}

func (clock *FakeClock) Increment(duration time.Duration) {
	now := clock.now.Add(duration)
	clock.now = now
}

func (clock *FakeClock) IncrementDate(years int, months int, days int) {
	now := clock.now.AddDate(years, months, days)
	clock.now = now
}
