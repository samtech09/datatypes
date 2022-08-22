package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

//NullTime
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Value implements the driver Valuer interface.
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	//return strings.Replace(nt.Time.Format(time.RFC3339), "+", "Z", -1), nil
	return nt.Time.Format(time.RFC3339), nil
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	// if value is nil, false
	if value == nil {
		*nt = NullTime{Time: time.Time{}}
		return nil
	}

	var err error
	switch x := value.(type) {
	case time.Time:
		nt.Time = x
	case string:
		t, err := parseTimeString(x)
		if err == nil {
			nt.Time = t
		}
	case []byte:
		t, err := parseTimeString(string(x))
		if err == nil {
			nt.Time = t
		}
	case nil:
		nt.Valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into NullTime: %v", value, value)
	}
	nt.Valid = err == nil
	return err
}

func (nt *NullTime) String() string {
	if !nt.Valid {
		return ""
	}
	return nt.Time.Format(time.RFC3339)
}

//MarshalJSON convert field value to JSON
func (nt NullTime) MarshalJSON() ([]byte, error) {
	//RFC3339 = 2006-01-02T15:04:05Z07:00
	if nt.Valid {
		return json.Marshal(nt.Time.Format(time.RFC3339))
	}
	return json.Marshal(time.Time{}.Format(time.RFC3339))
}

//UnmarshalJSON parse JSON valus and set into field
func (nt *NullTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	s := string(b)
	if s == "null" {
		// donot raise error as nil is valid type
		return nil
	}

	tmpTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		*nt = NullTime{Time: time.Time{}}
	} else {
		*nt = NullTime{Time: tmpTime, Valid: true}
	}
	return err
}

func (ns *NullTime) SetValue(val time.Time) {
	ns.Time = val
	ns.Valid = !val.IsZero()
}

func NewNullTime(val time.Time) NullTime {
	ns := NullTime{}
	ns.Time = val
	ns.Valid = !val.IsZero()
	return ns
}
