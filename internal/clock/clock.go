package clock

import "time"

// Clock provides the wall-clock dependency required only at protocol ingress
// and mining boundaries. Historical consensus validation must not use it.
type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func (RealClock) Now() time.Time { return time.Now() }

type FixedClock struct{ Current time.Time }

func (c FixedClock) Now() time.Time { return c.Current }

// PeerMedian supplies an authenticated/legitimate peer median when the
// networking layer has one. Build 1 deliberately does not invent peer time.
type PeerMedian interface {
	MedianTime() (time.Time, bool)
}

type UnavailablePeerMedian struct{}

func (UnavailablePeerMedian) MedianTime() (time.Time, bool) { return time.Time{}, false }
