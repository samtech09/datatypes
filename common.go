package datatypes

import (
	"time"

	errors "golang.org/x/xerrors"
)

var errBadStatus = errors.New("invalid status")

//taken from go-pg/types/time.go, Credit to go-pg for this time-parsing function
func parseTimeString(s string) (time.Time, error) {
	switch l := len(s); {
	case l <= len(timeFormat):
		if s[2] == ':' {
			return time.ParseInLocation(timeFormat, s, time.UTC)
		}
		return time.ParseInLocation(dateFormat, s, time.UTC)
	default:
		if s[10] == 'T' {
			return time.Parse(time.RFC3339Nano, s)
		}
		if c := s[l-9]; c == '+' || c == '-' {
			return time.Parse(timestamptzFormat, s)
		}
		if c := s[l-6]; c == '+' || c == '-' {
			return time.Parse(timestamptzFormat2, s)
		}
		if c := s[l-3]; c == '+' || c == '-' {
			return time.Parse(timestamptzFormat3, s)
		}
		return time.ParseInLocation(timestampFormat, s, time.UTC)
	}
}

func toString(t interface{}) (string, bool) {
	value, ok := t.(string)
	if ok {
		return value, true
	}
	if val, ok := t.([]byte); ok {
		return string(val), true
	}
	// return default value
	return "", false
}
