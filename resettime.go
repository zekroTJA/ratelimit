package ratelimit

import "time"

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

// IsNil returns the boolean value if
// the inner Time value is equal
// an empty time object (time.Time{}).
func (r *ResetTime) IsNil() bool {
	return r.isNil
}

// Unix overwrites time.Time#Unix()
// so it will return 0 if IsNil is
// true. Else, it will behave like
// default.
func (r *ResetTime) Unix() int64 {
	if r.isNil {
		return 0
	}

	return r.Time.Unix()
}

// UnixNano overwrites time.Time#UnixNano()
// so it will return 0 if IsNil is
// true. Else, it will behave like
// default.
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
