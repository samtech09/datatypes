package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"math"
	"strconv"

	"github.com/pkg/errors"
)

type NullInt32 struct {
	Int   int32
	Valid bool
}

// Scan implements the database/sql Scanner interface.
func (dst *NullInt32) Scan(src interface{}) error {
	if src == nil {
		*dst = NullInt32{Valid: false}
		return nil
	}

	switch src := src.(type) {
	case int64:
		if src < math.MinInt32 {
			return errors.Errorf("%d is greater than maximum value for Int4", src)
		}
		if src > math.MaxInt32 {
			return errors.Errorf("%d is greater than maximum value for Int4", src)
		}
		*dst = NullInt32{Int: int32(src), Valid: true}
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

func (src NullInt32) Value() (driver.Value, error) {
	switch src.Valid {
	case true:
		return int64(src.Int), nil
	default:
		return nil, nil
	}
}

//MarshalJSON convert field value to JSON
func (src NullInt32) MarshalJSON() ([]byte, error) {
	//fmt.Printf("Marshaling NullInt32: %d\n", src.Status)

	switch src.Valid {
	case true:
		return []byte(strconv.FormatInt(int64(src.Int), 10)), nil
	default:
		return []byte("0"), nil
	}
}

//UnmarshalJSON parse JSON valus and set into field
func (dst *NullInt32) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		// donot raise error as nil is valid type
		return nil
	}

	var v *int32
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v == nil {
		*dst = NullInt32{Valid: false}
	} else {
		*dst = NullInt32{Int: *v, Valid: true}
	}

	return nil
}

func (dst *NullInt32) FromString(src string) error {
	if src == "" {
		*dst = NullInt32{Valid: false}
		return nil
	}

	n, err := strconv.ParseInt(src, 10, 32)
	if err != nil {
		*dst = NullInt32{Valid: false}
		return err
	}

	*dst = NullInt32{Int: int32(n), Valid: true}
	return nil
}

func (ns *NullInt32) SetValue(val int32) {
	ns.Int = val
	ns.Valid = (val != 0)
}
