package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/jackc/pgtype"
)

type NullInt64 pgtype.Int8

func (p *NullInt64) Scan(src interface{}) error {
	t := pgtype.Int8(*p)
	err := t.Scan(src)
	*p = NullInt64(t)
	return err
}

func (p NullInt64) Value() (driver.Value, error) {
	return pgtype.Int8(p).Value()
}

//MarshalJSON convert field value to JSON
func (src NullInt64) MarshalJSON() ([]byte, error) {
	//fmt.Printf("Marshaling NullInt64: %d\n", src.Status)

	switch src.Status {
	case pgtype.Present:
		//return []byte(strconv.FormatInt(src.Int, 10)), nil
		return json.Marshal(strconv.FormatInt(src.Int, 10))
	case pgtype.Null:
		return []byte("0"), nil
	case pgtype.Undefined:
		return []byte("0"), nil
	}

	//fmt.Println("NullInt64 is nil")
	return nil, errBadStatus
}

//UnmarshalJSON parse JSON valus and set into field
func (dst *NullInt64) UnmarshalJSON(b []byte) error {
	var v *int64
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v == nil {
		*dst = NullInt64{Status: pgtype.Null}
	} else {
		*dst = NullInt64{Int: *v, Status: pgtype.Present}
	}

	return nil
}
