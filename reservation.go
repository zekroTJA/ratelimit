package ratelimit

// Reservation contains the pre-defined burst rate
// of the Limiter, the amount of remaining tickets
// and the time until a new token will be added to
// the bucket if Remaining == 0. Else, reset will
// be time.Time{} (0001-01-01 00:00:00 +0000 UTC).
//
// This struct contains JSON tags, so it can be
// easily parsed to JSON.
type Reservation struct {
	Burst     int       `json:"burst"`
	Remaining int       `json:"remaining"`
	Reset     ResetTime `json:"reset"`
}
