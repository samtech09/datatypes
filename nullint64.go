package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/pkg/errors"
)

type NullInt64 struct {
	Int   int64
	Valid bool
}

// Scan implements the database/sql Scanner interface.
func (dst *NullInt64) Scan(src interface{}) error {
	if src == nil {
		*dst = NullInt64{Valid: false}
		return nil
	}

	switch src := src.(type) {
	case int64:
		*dst = NullInt64{Int: src, Valid: true}
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

func (src NullInt64) Value() (driver.Value, error) {
	switch src.Valid {
	case true:
		return src.Int, nil
	default:
		return nil, nil
	}
}

//MarshalJSON convert field value to JSON
func (src NullInt64) MarshalJSON() ([]byte, error) {
	//fmt.Printf("Marshaling NullInt64: %d\n", src.Status)

	switch src.Valid {
	case true:
		return []byte(strconv.FormatInt(src.Int, 10)), nil
	default:
		return []byte("0"), nil
	}
}

//UnmarshalJSON parse JSON valus and set into field
func (dst *NullInt64) UnmarshalJSON(b []byte) error {
	if string(b) == "null" {
		// donot raise error as nil is valid type
		return nil
	}

	var v *int64
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v == nil {
		*dst = NullInt64{Valid: false}
	} else {
		*dst = NullInt64{Int: *v, Valid: true}
	}

	return nil
}

func (dst *NullInt64) FromString(src string) error {
	if src == "" {
		*dst = NullInt64{Valid: false}
		return nil
	}

	n, err := strconv.ParseInt(src, 10, 64)
	if err != nil {
		*dst = NullInt64{Valid: false}
		return err
	}

	*dst = NullInt64{Int: n, Valid: true}
	return nil
}

func (ns *NullInt64) SetValue(val int64) {
	ns.Int = val
	ns.Valid = (val != 0)
}
