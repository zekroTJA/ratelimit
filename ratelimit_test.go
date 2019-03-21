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
		t.Errorf("limit should be %v but was %v", limit, l.limit)
	}

	if l.burst != burst {
		t.Errorf("burst should be %d but was %d", burst, l.burst)
	}

	if l.tokens != burst {
		t.Errorf("tokens should be %d but was %d", burst, l.tokens)
	}
}

func TestReserveN(t *testing.T) {
	const limit = 100 * time.Millisecond
	const burst = 3

	l := NewLimiter(limit, burst)

	ok, res := l.ReserveN(0)
	if !ok || res != nil {
		t.Errorf(
			"ReserveN(0) should return (true, nil) but returned (%t, %+v)",
			ok, res)
	}

	ok, res = l.ReserveN(2)
	if !ok {
		t.Error("returned false but should return true")
	}
	if res.Burst != burst {
		t.Errorf("res.Burst should be %d but was %d",
			burst, res.Burst)
	}
	if res.Remaining != 1 {
		t.Errorf("res.Remaining should be %d but was %d",
			1, res.Remaining)
	}

	time.Sleep(310 * time.Millisecond)

	if l.Tokens() != burst {
		t.Errorf("recovered amount of tokens should be %d but was %d",
			burst, l.Tokens())
	}

	// -----------------------------------------

	l = NewLimiter(limit, burst)

	for i := 0; i < 14; i++ {
		ok, _ := l.ReserveN(1)

		switch i {
		case 0, 1, 2:
			if !ok {
				t.Errorf("ROUND %d | ok was false but should be true", i)
			}
		case 3, 4, 5, 6, 7, 8, 9, 10, 11:
			if ok {
				t.Errorf("ROUND %d | ok was true but should be false", i)
			}
		case 12:
			if !ok {
				t.Errorf("ROUND %d | ok was false but should be true", i)
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func TestReserve(t *testing.T) {
	const limit = 100 * time.Millisecond
	const burst = 3

	l := NewLimiter(limit, burst)

	ok, res := l.Reserve()
	if !ok {
		t.Error("returned false but should return true")
	}
	if res.Burst != burst {
		t.Errorf("res.Burst should be %d but was %d",
			burst, res.Burst)
	}
	if res.Remaining != 2 {
		t.Errorf("res.Remaining should be %d but was %d",
			2, res.Remaining)
	}

	time.Sleep(310 * time.Millisecond)

	if l.Tokens() != burst {
		t.Errorf("recovered amount of tokens should be %d but was %d",
			burst, l.Tokens())
	}

	// -----------------------------------------

	l = NewLimiter(limit, burst)

	for i := 0; i < 14; i++ {
		ok, _ := l.ReserveN(1)

		switch i {
		case 0, 1, 2:
			if !ok {
				t.Errorf("ROUND %d | ok was false but should be true", i)
			}
		case 3, 4, 5, 6, 7, 8, 9, 10, 11:
			if ok {
				t.Errorf("ROUND %d | ok was true but should be false", i)
			}
		case 12:
			if !ok {
				t.Errorf("ROUND %d | ok was false but should be true", i)
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func TestAllowN(t *testing.T) {
	const limit = 100 * time.Millisecond
	const burst = 3

	l := NewLimiter(limit, burst)

	ok := l.AllowN(1)
	if !ok {
		t.Error("returned false but should return true")
	}

	time.Sleep(310 * time.Millisecond)

	if l.Tokens() != burst {
		t.Errorf("recovered amount of tokens should be %d but was %d",
			burst, l.Tokens())
	}

	// -----------------------------------------

	l = NewLimiter(limit, burst)

	for i := 0; i < 14; i++ {
		ok := l.AllowN(1)

		switch i {
		case 0, 1, 2:
			if !ok {
				t.Errorf("ROUND %d | ok was false but should be true", i)
			}
		case 3, 4, 5, 6, 7, 8, 9, 10, 11:
			if ok {
				t.Errorf("ROUND %d | ok was true but should be false", i)
			}
		case 12:
			if !ok {
				t.Errorf("ROUND %d | ok was false but should be true", i)
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func TestAllow(t *testing.T) {
	const limit = 100 * time.Millisecond
	const burst = 3

	l := NewLimiter(limit, burst)

	ok := l.Allow()
	if !ok {
		t.Error("returned false but should return true")
	}

	time.Sleep(310 * time.Millisecond)

	if l.Tokens() != burst {
		t.Errorf("recovered amount of tokens should be %d but was %d",
			burst, l.Tokens())
	}

	// -----------------------------------------

	l = NewLimiter(limit, burst)

	for i := 0; i < 14; i++ {
		ok := l.Allow()

		switch i {
		case 0, 1, 2:
			if !ok {
				t.Errorf("ROUND %d | ok was false but should be true", i)
			}
		case 3, 4, 5, 6, 7, 8, 9, 10, 11:
			if ok {
				t.Errorf("ROUND %d | ok was true but should be false", i)
			}
		case 12:
			if !ok {
				t.Errorf("ROUND %d | ok was false but should be true", i)
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func TestLimit(t *testing.T) {
	const limit = 100 * time.Millisecond
	const burst = 3

	l := NewLimiter(limit, burst)

	if m := l.Limit(); m != limit {
		t.Errorf("l.Burst() should be %s but was %s", limit, m)
	}
}

func TestBurst(t *testing.T) {
	const limit = 100 * time.Millisecond
	const burst = 3

	l := NewLimiter(limit, burst)

	if b := l.Burst(); b != burst {
		t.Errorf("l.Burst() should be %d but was %d", burst, b)
	}
}

func TestTokens(t *testing.T) {
	const limit = 100 * time.Millisecond
	const burst = 3

	l := NewLimiter(limit, burst)

	if tg := l.Tokens(); tg != burst {
		t.Errorf("l.Tokens() should be %d but was %d", burst, tg)
	}

	l.Reserve()

	if tg := l.Tokens(); tg != burst-1 {
		t.Errorf("l.Tokens() should be %d but was %d", burst-1, tg)
	}
}

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
