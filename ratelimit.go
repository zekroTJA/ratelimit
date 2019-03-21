// Package ratelimit provides a verry simple,
// token bukkit based request rate limiter.
package ratelimit

import (
	"sync"
	"time"
)

// ResetTime inherits from time.Time and
// contains an extra-vaule 'isNil', which
// identifies if Time == time.Time{}.
// Also, the functions Unix, UnixNano
// and Format are overwritten to take
// notice of this extra-value.
type ResetTime struct {
	time.Time

	isNil bool
}

// A Limiter controls how frequently accesses
// should be allowed to happen. It implements
// the principle of the token bukkit, which
// defines a bukkit with an exact size of tokens.
// Also, a rate is defined after exactly 1 token
// will be added to the bukkits volume, if the
// bukkit is not "full" (nVolume == nSize).
// The amount of tokens in the bukkit are
// defined as ability to perform an action,
// which then reduces the volume of the bukkit
// by n tickets.
type Limiter struct {
	mu sync.Mutex

	limit time.Duration
	burst int

	tokens int
	last   time.Time
}

// Reservation contains the pre-defined burst rate
// of the Limiter, the amount of remaining tickets
// and the time until a new token will be added to
// the bukkit if Remaining == 0. Else, reset will
// be time.Time{} (0001-01-01 00:00:00 +0000 UTC).
//
// This struct contains JSON tags, so it can be
// easily parsed to JSON.
type Reservation struct {
	Burst     int       `json:"burst"`
	Remaining int       `json:"remaining"`
	Reset     ResetTime `json:"reset"`
}

// NewLimiter returns a new instance of Limiter
// with a burst rate of b and a limit time
// of l until a new token will be generated.
func NewLimiter(l time.Duration, b int) *Limiter {
	return &Limiter{
		limit:  l,
		burst:  b,
		tokens: b,
	}
}

// ReserveN checks if an amount of n tickets are
// currently available. If this is the case, true
// will be returned with a Reservation object as
// status information of the Limiter and n tokens
// will be consumed.
// If there are not enough tokens available to
// satisfy the reservation, false will be returned
// with a Reservation object containing the Limiters
// status containing the time until next token
// generation.
func (l *Limiter) ReserveN(n int) (bool, *Reservation) {
	if n <= 0 {
		return true, nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	res := &Reservation{
		Burst: l.burst,
		Reset: ResetTime{
			isNil: true,
		},
	}

	tokensSinceLast := int(time.Since(l.last) / l.limit)
	l.tokens += tokensSinceLast
	if l.tokens > l.burst {
		l.tokens = l.burst
	}

	if l.tokens >= n {
		l.tokens -= n
		l.last = time.Now()
		res.Remaining = l.tokens

		if l.tokens == 0 {
			res.Reset.Time = l.last.Add(l.limit)
			res.Reset.isNil = false
		}

		return true, res
	}

	res.Remaining = l.tokens
	res.Reset.Time = l.last.Add(l.limit)
	res.Reset.isNil = false

	return false, res
}

// Reserve is shorthand for ReserveN(1).
func (l *Limiter) Reserve() (bool, *Reservation) {
	return l.ReserveN(1)
}

// AllowN is shorthand for reserveN(n) but only
// returning a boolean which exposes the
// succeed of the reservation.
func (l *Limiter) AllowN(n int) bool {
	ok, _ := l.ReserveN(n)
	return ok
}

// Allow is shorthand for AllowN(1).
func (l *Limiter) Allow() bool {
	return l.AllowN(1)
}

// Limit returns the defined limit duration
// after which a new token will be generated.
//
// This function does not consume tokens.
func (l *Limiter) Limit() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.limit
}

// Burst returns the defined burst value.
//
// This function does not consume tokens.
func (l *Limiter) Burst() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.burst
}

// Tokens returns the current available
// tokens. This value is unequal to the
// actual value of tokens, because this
// value is only refreshed after token
// consumption. So the returned value
// is the actial value of tokens plus
// the calculated amount of tokens which
// are virtually generated after last
// consumption.
//
// This function does not consume tokens.
func (l *Limiter) Tokens() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	t := l.tokens + int(time.Since(l.last)/l.limit)
	if t > l.burst {
		return l.burst
	}

	return t
}

// IsNil returns the boolean value if
// the inner Time value is equal
// an empty time object (time.Time{}).
func (r *ResetTime) IsNil() bool {
	return r.isNil
}

// Unix overwrites time.Time#Unix()
// so it will return 0 if IsNil is
// true. Else, it will behave like
// defualt.
func (r *ResetTime) Unix() int64 {
	if r.isNil {
		return 0
	}

	return r.Time.Unix()
}

// UnixNano overwrites time.Time#UnixNano()
// so it will return 0 if IsNil is
// true. Else, it will behave like
// defualt.
func (r *ResetTime) UnixNano() int64 {
	if r.isNil {
		return 0
	}

	return r.Time.UnixNano()
}

// Format overwrites time.Time#Format()
// so it will return the value of def
// (if not defined it will return "") if
// IsNil is true. Else, it will behave
// like default.
func (r *ResetTime) Format(layout string, def ...string) string {
	if r.isNil {
		if len(def) > 0 {
			return def[0]
		}
		return ""
	}

	return r.Time.Format(layout)
}
