package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/jackc/pgtype"
)

type NullFloat32 pgtype.Float4

func (p *NullFloat32) Scan(src interface{}) error {
	t := pgtype.Float4(*p)
	err := t.Scan(src)
	*p = NullFloat32(t)
	return err
}

func (p NullFloat32) Value() (driver.Value, error) {
	return pgtype.Float4(p).Value()
}

//MarshalJSON convert field value to JSON
func (src NullFloat32) MarshalJSON() ([]byte, error) {
	//fmt.Printf("Marshaling NullFloat32: %d\n", src.Status)

	switch src.Status {
	case pgtype.Present:
		return []byte(strconv.FormatFloat(float64(src.Float), 'f', -1, 64)), nil
		//return json.Marshal(strconv.FormatFloat(float64(src.Float), 'f', -1, 64))
	case pgtype.Null:
		return []byte("0"), nil
	case pgtype.Undefined:
		return []byte("0"), nil
	}

	//fmt.Println("NullFloat32 is nil")
	return nil, errBadStatus
}

//UnmarshalJSON parse JSON valus and set into field
func (dst *NullFloat32) UnmarshalJSON(b []byte) error {
	var v *float32
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v == nil {
		*dst = NullFloat32{Status: pgtype.Null}
	} else {
		*dst = NullFloat32{Float: *v, Status: pgtype.Present}
	}

	return nil
}
