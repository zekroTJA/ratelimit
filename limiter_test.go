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
		t.Fatal("NewLimiter() should not return nil")
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

	ts := &testTimeSource{}
	l := NewLimiterWithTimeSource(ts.Now, limit, burst)

	ok, res := l.ReserveN(0)
	if !ok || (res != Reservation{}) {
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

	ts.Advance(3 * limit)

	if l.Tokens() != burst {
		t.Errorf("recovered amount of tokens should be %d but was %d",
			burst, l.Tokens())
	}

	// -----------------------------------------

	ts = &testTimeSource{}
	l = NewLimiterWithTimeSource(ts.Now, limit, burst)

	for i := 0; i < 3; i++ {
		ok, _ = l.ReserveN(1)
		if !ok {
			t.Fatalf("Reservation was not successful")
		}
	}

	ok, _ = l.ReserveN(1)
	if ok {
		t.Fatalf("Reservation was successful even though it should not")
	}

	ts.Advance(limit)
	ok, _ = l.ReserveN(1)
	if !ok {
		t.Fatalf("Reservation was not successful")
	}

	ok, _ = l.ReserveN(1)
	if ok {
		t.Fatalf("Reservation was successful even though it should not")
	}

	ts.Advance(2 * limit)
	ok, _ = l.ReserveN(3)
	if ok {
		t.Fatalf("Reservation was successful even though it should not")
	}
	ok, _ = l.ReserveN(2)
	if !ok {
		t.Fatalf("Reservation was not successful")
	}
	ok, _ = l.ReserveN(1)
	if ok {
		t.Fatalf("Reservation was successful even though it should not")
	}
}

func TestReserve(t *testing.T) {
	const limit = 100 * time.Millisecond
	const burst = 3

	ts := &testTimeSource{}
	l := NewLimiterWithTimeSource(ts.Now, limit, burst)

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

	ts.Advance(3 * limit)

	if l.Tokens() != burst {
		t.Errorf("recovered amount of tokens should be %d but was %d",
			burst, l.Tokens())
	}

	// -----------------------------------------

	ts = &testTimeSource{}
	l = NewLimiterWithTimeSource(ts.Now, limit, burst)

	for i := 0; i < 3; i++ {
		ok, _ = l.Reserve()
		if !ok {
			t.Fatalf("Reservation was not successful")
		}
	}

	ok, _ = l.Reserve()
	if ok {
		t.Fatalf("Reservation was successful even though it should not")
	}

	ts.Advance(limit)
	ok, _ = l.Reserve()
	if !ok {
		t.Fatalf("Reservation was not successful")
	}

	ok, _ = l.Reserve()
	if ok {
		t.Fatalf("Reservation was successful even though it should not")
	}
}

func TestAllowN(t *testing.T) {
	const limit = 100 * time.Millisecond
	const burst = 3

	ts := &testTimeSource{}
	l := NewLimiterWithTimeSource(ts.Now, limit, burst)

	ok := l.AllowN(1)
	if !ok {
		t.Error("returned false but should return true")
	}

	ts.Advance(3 * limit)

	if l.Tokens() != burst {
		t.Errorf("recovered amount of tokens should be %d but was %d",
			burst, l.Tokens())
	}

	// -----------------------------------------

	ts = &testTimeSource{}
	l = NewLimiterWithTimeSource(ts.Now, limit, burst)

	for i := 0; i < 3; i++ {
		ok = l.AllowN(1)
		if !ok {
			t.Fatalf("Reservation was not successful")
		}
	}
	// 0 tokens remaining

	ok = l.AllowN(1)
	if ok {
		t.Fatalf("Reservation was successful even though it should not")
	}

	ts.Advance(limit)
	// 1 tokens remaining
	ok = l.AllowN(1)
	if !ok {
		t.Fatalf("Reservation was not successful")
	}
	// 0 tokens remaining

	ok = l.AllowN(1)
	if ok {
		t.Fatalf("Reservation was successful even though it should not")
	}

	ts.Advance(2 * limit)
	// 2 tokens remaining
	ok = l.AllowN(3)
	if ok {
		t.Fatalf("Reservation was successful even though it should not")
	}
	// 2 tokens remaining
	ok = l.AllowN(2)
	if !ok {
		t.Fatalf("Reservation was not successful")
	}
	// 0 tokens remaining
	ok = l.AllowN(1)
	if ok {
		t.Fatalf("Reservation was successful even though it should not")
	}
}

func TestAllow(t *testing.T) {
	const limit = 100 * time.Millisecond
	const burst = 3

	ts := &testTimeSource{}
	l := NewLimiterWithTimeSource(ts.Now, limit, burst)

	ok := l.Allow()
	if !ok {
		t.Error("returned false but should return true")
	}

	ts.Advance(3 * limit)

	if l.Tokens() != burst {
		t.Errorf("recovered amount of tokens should be %d but was %d",
			burst, l.Tokens())
	}

	// -----------------------------------------

	ts = &testTimeSource{}
	l = NewLimiterWithTimeSource(ts.Now, limit, burst)

	for i := 0; i < 3; i++ {
		ok = l.Allow()
		if !ok {
			t.Fatalf("Reservation was not successful")
		}
	}

	ok = l.Allow()
	if ok {
		t.Fatalf("Reservation was successful even though it should not")
	}

	ts.Advance(limit)
	ok = l.Allow()
	if !ok {
		t.Fatalf("Reservation was not successful")
	}

	ok = l.Allow()
	if ok {
		t.Fatalf("Reservation was successful even though it should not")
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

func TestReserve_zero(t *testing.T) {
	l := NewLimiter(0, 10)

	ok, res := l.ReserveN(1)
	if ok || (res != Reservation{}) {
		t.Error("ReserveN should return false when l is 0")
	}

	l = NewLimiter(100*time.Millisecond, 0)

	ok, res = l.ReserveN(1)
	if ok || (res != Reservation{}) {
		t.Error("ReserveN should return false when b is 0")
	}

	l = NewLimiter(-1, 10)

	ok, res = l.ReserveN(1)
	if ok || (res != Reservation{}) {
		t.Error("ReserveN should return false when b is -1")
	}

	l = NewLimiter(100*time.Millisecond, -1)

	ok, res = l.ReserveN(1)
	if ok || (res != Reservation{}) {
		t.Error("ReserveN should return false when b is -1")
	}
}

// Test for issue #1
// https://github.com/zekroTJA/ratelimit/issues/1
func TestReserve_reset(t *testing.T) {
	const limit = 100 * time.Millisecond
	const burst = 1

	ts := &testTimeSource{}
	l := NewLimiterWithTimeSource(ts.Now, limit, burst)

	ok, _ := l.Reserve()
	if !ok {
		t.Fatal("First Reserve should return true")
	}

	ts.Advance(limit / 2)
	ok, _ = l.Reserve()
	if ok {
		t.Fatal("Second Reserve should return false")
	}

	ts.Advance(limit / 2)
	ok, _ = l.Reserve()
	if !ok {
		t.Fatal("Third Reserve should return true")
	}
}
