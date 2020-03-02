package datatypes

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgtype"
)

type NullBool pgtype.Bool

func (p *NullBool) Scan(src interface{}) error {
	t := pgtype.Bool(*p)
	err := t.Scan(src)
	*p = NullBool(t)
	return err
}

func (p NullBool) Value() (driver.Value, error) {
	return pgtype.Bool(p).Value()
}

//MarshalJSON convert field value to JSON
func (src NullBool) MarshalJSON() ([]byte, error) {
	switch src.Status {
	case pgtype.Present:
		if src.Bool {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	case pgtype.Null:
		return []byte("false"), nil
	case pgtype.Undefined:
		return []byte("false"), nil
	}

	return nil, errBadStatus
}

//UnmarshalJSON parse JSON valus and set into field
func (dst *NullBool) UnmarshalJSON(b []byte) error {
	var v *bool
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v == nil {
		*dst = NullBool{Status: pgtype.Null}
	} else {
		*dst = NullBool{Bool: *v, Status: pgtype.Present}
	}

	return nil
}
