package ratelimit

import (
	"testing"
	"time"
)

func TestIsNil(t *testing.T) {
	const limit = 100 * time.Second
	const burst = 2

	l := NewLimiter(limit, burst)

	_, res := l.Reserve()
	// Remaining tokens: 1
	// ok == true
	if !res.Reset.IsNil() {
		t.Error("res.Rest.IsNil should be true but was false")
	}

	_, res = l.Reserve()
	// Remaining tokens: 0
	// ok == true
	if res.Reset.IsNil() {
		t.Error("res.Rest.IsNil should be false but was true")
	}

	_, res = l.Reserve()
	// Remaining tokens: 0
	// ok == fasle
	if res.Reset.IsNil() {
		t.Error("res.Rest.IsNil should be false but was true")
	}
}

func TestUnix(t *testing.T) {
	const limit = 100 * time.Second
	const burst = 2

	l := NewLimiter(limit, burst)

	_, res := l.Reserve()
	// Remaining tokens: 1
	// ok == true
	if ti := res.Reset.Unix(); ti != 0 {
		t.Errorf("res.Reset.Unix() should be 0 but was %d", ti)
	}

	_, res = l.Reserve()
	// Remaining tokens: 0
	// ok == true
	ti := res.Reset.Unix()
	tm := res.Reset.Time.Unix()
	if ti != tm {
		t.Errorf("res.Reset.Unix() should be %d but was %d", tm, ti)
	}

	_, res = l.Reserve()
	// Remaining tokens: 0
	// ok == fasle
	ti = res.Reset.Unix()
	tm = res.Reset.Time.Unix()
	if ti != tm {
		t.Errorf("res.Reset.Unix() should be %d but was %d", tm, ti)
	}
}

func TestUnixNano(t *testing.T) {
	const limit = 100 * time.Second
	const burst = 2

	l := NewLimiter(limit, burst)

	_, res := l.Reserve()
	// Remaining tokens: 1
	// ok == true
	if ti := res.Reset.UnixNano(); ti != 0 {
		t.Errorf("res.Reset.Unix() should be 0 but was %d", ti)
	}

	_, res = l.Reserve()
	// Remaining tokens: 0
	// ok == true
	ti := res.Reset.UnixNano()
	tm := res.Reset.Time.UnixNano()
	if ti != tm {
		t.Errorf("res.Reset.Unix() should be %d but was %d", tm, ti)
	}

	_, res = l.Reserve()
	// Remaining tokens: 0
	// ok == fasle
	ti = res.Reset.UnixNano()
	tm = res.Reset.Time.UnixNano()
	if ti != tm {
		t.Errorf("res.Reset.Unix() should be %d but was %d", tm, ti)
	}
}

func TestFormat(t *testing.T) {
	const limit = 100 * time.Second
	const burst = 2
	const layout = time.RFC3339Nano
	const defF = "empty"

	l := NewLimiter(limit, burst)

	_, res := l.Reserve()
	// Remaining tokens: 1
	// ok == true
	if ti := res.Reset.Format(layout, defF); ti != defF {
		t.Errorf("res.Reset.Unix() should be '%s' but was '%s'", defF, ti)
	}
	if ti := res.Reset.Format(layout); ti != "" {
		t.Errorf("res.Reset.Unix() should be '' but was '%s'", ti)
	}

	_, res = l.Reserve()
	// Remaining tokens: 0
	// ok == true
	ti := res.Reset.Format(layout)
	tm := res.Reset.Time.Format(layout)
	if ti != tm {
		t.Errorf("res.Reset.Unix() should be %s but was %s", tm, ti)
	}

	_, res = l.Reserve()
	// Remaining tokens: 0
	// ok == fasle
	ti = res.Reset.Format(layout)
	tm = res.Reset.Time.Format(layout)
	if ti != tm {
		t.Errorf("res.Reset.Unix() should be %s but was %s", tm, ti)
	}
}
