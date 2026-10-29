package clock

import "time"

// Clock provides a time source that can be replaced in tests.
type Clock interface {
	Now() time.Time
}

// RealClock is the production clock backed by the system time.
type RealClock struct{}

func (RealClock) Now() time.Time {
	return time.Now().UTC()
}
