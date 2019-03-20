package ratelimit

import (
	"testing"
	"time"
)

func TestNewLimiter(t *testing.T) {
	const limit = 5 * time.Second
	const burst = 4

	l := NewLimiter(limit, burst)
	if l == nil {
		t.Error("NewLimiter() should not return nil")
	}

	if l.limit != limit {
		t.Errorf("limit should be %v but was %v", l.limit, limit)
	}

	if l.burst != burst {
		t.Errorf("burst should be %d but was %d", l.burst, burst)
	}

	if l.tokens != burst {
		t.Errorf("tokens should be %d but was %d", l.tokens, burst)
	}
}
