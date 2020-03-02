package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"math"
	"strconv"

	"github.com/pkg/errors"
)

type NullInt16 struct {
	Int   int16
	Valid bool
}

// Scan implements the database/sql Scanner interface.
func (dst *NullInt16) Scan(src interface{}) error {
	if src == nil {
		*dst = NullInt16{Valid: false}
		return nil
	}

	switch src := src.(type) {
	case int64:
		if src < math.MinInt16 {
			return errors.Errorf("%d is smaller than minimum value for Int2", src)
		}
		if src > math.MaxInt16 {
			return errors.Errorf("%d is greater than maximum value for Int2", src)
		}
		*dst = NullInt16{Int: int16(src), Valid: true}
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

func (src NullInt16) Value() (driver.Value, error) {
	switch src.Valid {
	case true:
		return int64(src.Int), nil
	default:
		return nil, nil
	}
}

//MarshalJSON convert field value to JSON
func (src NullInt16) MarshalJSON() ([]byte, error) {
	//fmt.Printf("Marshaling NullInt16: %d\n", src.Status)

	switch src.Valid {
	case true:
		return []byte(strconv.FormatInt(int64(src.Int), 10)), nil
	default:
		return []byte("0"), nil
	}
}

//UnmarshalJSON parse JSON valus and set into field
func (dst *NullInt16) UnmarshalJSON(b []byte) error {
	var v *int16
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v == nil {
		*dst = NullInt16{Valid: false}
	} else {
		*dst = NullInt16{Int: *v, Valid: true}
	}

	return nil
}

func (dst *NullInt16) FromString(src string) error {
	if src == "" {
		*dst = NullInt16{Valid: false}
		return nil
	}

	n, err := strconv.ParseInt(src, 10, 16)
	if err != nil {
		*dst = NullInt16{Valid: false}
		return err
	}

	*dst = NullInt16{Int: int16(n), Valid: true}
	return nil
}
