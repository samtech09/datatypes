package datatypes

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

type NullBool bool

func (dst *NullBool) Scan(src interface{}) error {
	if src == nil {
		*dst = false
		return nil
	}

	switch src := src.(type) {
	case bool:
		*dst = NullBool(src)
		return nil
	case string:
		return dst.FromString(src)
	case []byte:
		return dst.FromString(string(src))
	}

	return errors.Errorf("cannot scan %T", src)
}

func (src NullBool) Value() (driver.Value, error) {
	if src {
		return true, nil
	}
	return false, nil
}

//MarshalJSON convert field value to JSON
func (src NullBool) MarshalJSON() ([]byte, error) {
	if src {
		return []byte("true"), nil
	}
	return []byte("false"), nil
}

//UnmarshalJSON parse JSON valus and set into field
func (dst *NullBool) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		// donot raise error as nil is valid type
		return nil
	}

	var v *bool
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v == nil {
		*dst = NullBool(false)
	} else {
		*dst = NullBool(*v)
	}

	return nil
}

func (dst *NullBool) FromString(src string) error {
	if src == "" {
		*dst = false
		return nil
	}

	if len(src) != 1 {
		*dst = false
		return errors.Errorf("invalid length for bool: %v", len(src))
	}

	*dst = NullBool(src[0] == 't')
	return nil
}
