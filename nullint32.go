package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/jackc/pgtype"
)

type NullInt32 pgtype.Int4

func (p *NullInt32) Scan(src interface{}) error {
	t := pgtype.Int4(*p)
	err := t.Scan(src)
	*p = NullInt32(t)
	return err
}

func (p NullInt32) Value() (driver.Value, error) {
	return pgtype.Int4(p).Value()
}

//MarshalJSON convert field value to JSON
func (src NullInt32) MarshalJSON() ([]byte, error) {
	//fmt.Printf("Marshaling NullInt32: %d\n", src.Status)

	switch src.Status {
	case pgtype.Present:
		return []byte(strconv.FormatInt(int64(src.Int), 10)), nil
	case pgtype.Null:
		return []byte("0"), nil
	case pgtype.Undefined:
		return []byte("0"), nil
	}

	//fmt.Println("NullInt32 is nil")
	return nil, errBadStatus
}

//UnmarshalJSON parse JSON valus and set into field
func (dst *NullInt32) UnmarshalJSON(b []byte) error {
	var v *int32
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v == nil {
		*dst = NullInt32{Status: pgtype.Null}
	} else {
		*dst = NullInt32{Int: *v, Status: pgtype.Present}
	}

	return nil
}
