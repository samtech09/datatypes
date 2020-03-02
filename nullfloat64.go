package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/jackc/pgtype"
)

type NullFloat64 pgtype.Float8

func (p *NullFloat64) Scan(src interface{}) error {
	t := pgtype.Float8(*p)
	err := t.Scan(src)
	*p = NullFloat64(t)
	return err
}

func (p NullFloat64) Value() (driver.Value, error) {
	return pgtype.Float8(p).Value()
}

//MarshalJSON convert field value to JSON
func (src NullFloat64) MarshalJSON() ([]byte, error) {
	//fmt.Printf("Marshaling NullFloat64: %d\n", src.Status)

	switch src.Status {
	case pgtype.Present:
		//fmt.Println(src.Float)
		//fmt.Println(strconv.FormatFloat(float64(src.Float), 'f', 6, 64))
		return []byte(strconv.FormatFloat(float64(src.Float), 'f', -1, 64)), nil
		//return json.Marshal(strconv.FormatFloat(float64(src.Float), 'f', -1, 64))
	case pgtype.Null:
		return []byte("0"), nil
	case pgtype.Undefined:
		return []byte("0"), nil
	}

	//fmt.Println("NullFloat64 is nil")
	return nil, errBadStatus
}

//UnmarshalJSON parse JSON valus and set into field
func (dst *NullFloat64) UnmarshalJSON(b []byte) error {
	var v *float64
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v == nil {
		*dst = NullFloat64{Status: pgtype.Null}
	} else {
		*dst = NullFloat64{Float: *v, Status: pgtype.Present}
	}

	return nil
}
