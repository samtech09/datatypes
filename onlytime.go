package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	errors "golang.org/x/xerrors"
)

//OnlyTime represents only time portion of datetime
type OnlyTime time.Time

// Value - Implementation of valuer for database/sql
func (ot OnlyTime) Value() (driver.Value, error) {
	t := time.Time(ot).Format("15:04:05")
	return t, nil
}

//Scan - Implementation of scanner for database/sql
func (ot *OnlyTime) Scan(value interface{}) error {
	// if value is nil, false
	if value == nil {
		*ot = OnlyTime(time.Time{})
		return nil
	}

	//fmt.Println("value: ", value)

	if v, ok := value.(time.Time); ok {
		*ot = OnlyTime(v)
		return nil
	}

	switch src := value.(type) {
	case string:
		return ot.parseTime([]byte(src))
	case []byte:
		srcCopy := make([]byte, len(src))
		copy(srcCopy, src)
		return ot.parseTime(srcCopy)
	case time.Time:
		*ot = OnlyTime(src)
		return nil
	}

	return errors.New("failed to scan OnlyTime")
}

func (ot *OnlyTime) String() string {
	return fmt.Sprintf("%s", time.Time(*ot).Format("15:04:05"))
}

//MarshalJSON convert field value to JSON
func (ot OnlyTime) MarshalJSON() ([]byte, error) {
	//RFC3339 = 2006-01-02T15:04:05Z07:00
	//return json.Marshal(time.Time(ot).Format(time.RFC3339))
	return json.Marshal(time.Time(ot).Format("15:04:05"))
}

//UnmarshalJSON parse JSON valus and set into field
func (ot *OnlyTime) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	s := string(b)
	if s == "null" {
		// donot raise error as nil is valid type
		return nil
	}

	ot.parseTime(b)
	// if err != nil {
	// 	tmpTime, _ := time.Parse(time.RFC3339, s)
	// 	*ot = OnlyTime(tmpTime)
	// }
	return
}

func (ot *OnlyTime) parseTime(src []byte) error {
	var t time.Time
	var err error
	s := string(src)
	if len(src) > 19 {
		t, err = time.Parse(time.RFC3339, s)
	} else if len(src) > 9 {
		t, err = time.Parse("15:04:05.999999999", s)
	} else {
		t, err = time.Parse("15:04:05", s)
	}
	if err != nil {
		return err
	}
	*ot = OnlyTime(t)
	return nil
}
