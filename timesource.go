package ratelimit

import (
	"time"
)

// TimeSource is a function returning the current time
type TimeSource func() time.Time

type testTimeSource struct {
	currentTime time.Time
}

func (t *testTimeSource) Now() time.Time {
	return t.currentTime
}

func (t *testTimeSource) Advance(d time.Duration) {
	t.currentTime = t.currentTime.Add(d)
}
