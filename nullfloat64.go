package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

type NullFloat64 struct {
	Float float64
	Valid bool
}

// Scan implements the database/sql Scanner interface.
func (dst *NullFloat64) Scan(src interface{}) error {
	if src == nil {
		*dst = NullFloat64{Valid: false}
		return nil
	}

	switch src := src.(type) {
	case float64:
		*dst = NullFloat64{Float: src, Valid: true}
		return nil
	case string:
		return dst.FromString(src)
	case []byte:
		//srcCopy := make([]byte, len(src))
		//copy(srcCopy, src)
		return dst.FromString(string(src))
	}

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src NullFloat64) Value() (driver.Value, error) {
	switch src.Valid {
	case true:
		return float64(src.Float), nil
	default:
		return nil, nil
	}
}

//MarshalJSON convert field value to JSON
func (src NullFloat64) MarshalJSON() ([]byte, error) {
	//fmt.Printf("Marshaling NullFloat64: %d\n", src.Valid)

	switch src.Valid {
	case true:
		return []byte(strconv.FormatFloat(float64(src.Float), 'f', -1, 64)), nil
	default:
		return []byte("0"), nil
	}
}

//UnmarshalJSON parse JSON valus and set into field
func (dst *NullFloat64) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		// donot raise error as nil is valid type
		return nil
	}

	var v *float64
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v == nil {
		*dst = NullFloat64{}
	} else {
		*dst = NullFloat64{Float: *v, Valid: true}
	}

	return nil
}

func (dst *NullFloat64) FromString(src string) error {
	if src == "" {
		*dst = NullFloat64{}
		return nil
	}

	n, err := strconv.ParseFloat(src, 64)
	if err != nil {
		*dst = NullFloat64{}
		return err
	}

	*dst = NullFloat64{Float: n, Valid: true}
	return nil
}
