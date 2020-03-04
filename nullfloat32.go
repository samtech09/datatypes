package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

type NullFloat32 struct {
	Float float32
	Valid bool
}

// Scan implements the database/sql Scanner interface.
func (dst *NullFloat32) Scan(src interface{}) error {
	if src == nil {
		*dst = NullFloat32{Valid: false}
		return nil
	}

	switch src := src.(type) {
	case float64:
		*dst = NullFloat32{Float: float32(src), Valid: true}
		return nil
	case string:
		return dst.FromString(src)
	case []byte:
		// srcCopy := make([]byte, len(src))
		// copy(srcCopy, src)
		return dst.FromString(string(src))
	}

	return errors.Errorf("cannot scan %T", src)
}

// Value implements the database/sql/driver Valuer interface.
func (src NullFloat32) Value() (driver.Value, error) {
	switch src.Valid {
	case true:
		return float64(src.Float), nil
	default:
		return nil, nil
	}
}

//MarshalJSON convert field value to JSON
func (src NullFloat32) MarshalJSON() ([]byte, error) {
	//fmt.Printf("Marshaling NullFloat32: %d\n", src.Valid)

	switch src.Valid {
	case true:
		return []byte(strconv.FormatFloat(float64(src.Float), 'f', -1, 64)), nil
	default:
		return []byte("0"), nil
	}
}

//UnmarshalJSON parse JSON valus and set into field
func (dst *NullFloat32) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		// donot raise error as nil is valid type
		return nil
	}

	var v *float32
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v == nil {
		*dst = NullFloat32{}
	} else {
		*dst = NullFloat32{Float: *v, Valid: true}
	}

	return nil
}

func (dst *NullFloat32) FromString(src string) error {
	if src == "" {
		*dst = NullFloat32{}
		return nil
	}

	n, err := strconv.ParseFloat(src, 32)
	if err != nil {
		*dst = NullFloat32{}
		return err
	}

	*dst = NullFloat32{Float: float32(n), Valid: true}
	return nil
}
