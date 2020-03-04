package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// taken from go-pg/types/time.go
const (
	dateFormat         = "2006-01-02"
	timeFormat         = "15:04:05.999999999"
	timestampFormat    = "2006-01-02 15:04:05.999999999"
	timestamptzFormat  = "2006-01-02 15:04:05.999999999-07:00:00"
	timestamptzFormat2 = "2006-01-02 15:04:05.999999999-07:00"
	timestamptzFormat3 = "2006-01-02 15:04:05.999999999-07"
)

//UnixTime to scan time values from database to unix timestamp
type UnixTime struct {
	Utime int64
	Valid bool // Valid is true if Time is not NULL
}

//NewUnixTime create new UnixTime struct for given time.Time
func NewUnixTime(t time.Time) UnixTime {
	return UnixTime{
		Utime: t.Unix(),
		Valid: true,
	}
}

// Scan implements the Scanner interface.
func (nt *UnixTime) Scan(value interface{}) error {
	var err error
	switch x := value.(type) {
	case time.Time:
		//fmt.Println("Unixtime.Scan: found time.Time ") // debug
		nt.Utime = x.Unix()
	case []uint8:
		// fmt.Println("Unixtime.Scan: called for %s ", string(x)) // debug
		t, err := parseTimeString(string(x))
		if err == nil {
			nt.Utime = t.Unix()
		}
	case string:
		t, err := parseTimeString(x)
		if err == nil {
			nt.Utime = t.Unix()
		}
	case nil:
		//fmt.Println("Unixtime.Scan: tmptime is nil") // debug
		nt.Valid = false
		return nil
	default:
		err = fmt.Errorf("null: cannot scan type %T into UnixTime: %v", value, value)
	}
	nt.Valid = err == nil
	return err
}

// Value implements the driver Valuer interface.
func (nt UnixTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return time.Unix(nt.Utime, 0), nil
}

func (nt *UnixTime) String() string {
	if !nt.Valid {
		return ""
	}
	return time.Unix(nt.Utime, 0).Format(time.RFC3339)
}

//Format formats time to given layout
func (nt *UnixTime) Format(layout string) string {
	if !nt.Valid {
		return ""
	}
	return time.Unix(nt.Utime, 0).Format(layout)
}

//Time returns value of time
func (nt *UnixTime) Time() time.Time {
	if !nt.Valid {
		return time.Time{}
	}
	return time.Unix(nt.Utime, 0)
}

//MarshalJSON convert field value to JSON
func (nt UnixTime) MarshalJSON() ([]byte, error) {
	//RFC3339 = 2006-01-02T15:04:05Z07:00

	//fmt.Printf("Marshaling UnixTime: %v\n", nt.Valid)

	if nt.Valid {
		//fmt.Printf("UnixTime Value: %d\n", nt.Utime)

		//return json.Marshal(nt.Time.Format(time.RFC3339))
		return json.Marshal(time.Unix(nt.Utime, 0).UTC().Format(time.RFC3339))

		// Note:
		// Zone 'Z' is retuned only for UTC time
		// for local time it returns + e.g. 2006-01-02T15:04:05+07:00
		// To fix it, replace + with Z
		//return json.Marshal(strings.Replace(time.Unix(nt.Utime, 0).Format(time.RFC3339), "+", "Z", -1))
	}

	return json.Marshal(time.Time{}.Format(time.RFC3339))

}

//UnmarshalJSON parse JSON valus and set into field
//  It is setter, so *UnixTime is required
func (nt *UnixTime) UnmarshalJSON(b []byte) (err error) {
	//fmt.Println("Unixtime.Un-MarshalJSON")
	// unquote - if present
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}
	// tmpTime, err := time.Parse(time.RFC3339, string(b))
	// if err != nil {
	// 	nt.Valid = false
	// 	return err
	// }
	// nt.Valid = true
	// nt.Utime = tmpTime.Unix()

	s := string(b)
	if s == "null" {
		// donot raise error as nil is valid type
		return nil
	}

	t, err := parseTimeString(s)
	if err == nil {
		nt.Valid = true
		nt.Utime = t.Unix()
	}

	return err
}
