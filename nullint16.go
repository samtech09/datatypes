package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"strconv"

	"github.com/jackc/pgtype"
)

type NullInt16 pgtype.Int2

func (p *NullInt16) Scan(src interface{}) error {
	t := pgtype.Int2(*p)
	err := t.Scan(src)
	*p = NullInt16(t)
	return err
}

func (p NullInt16) Value() (driver.Value, error) {
	return pgtype.Int2(p).Value()
}

//MarshalJSON convert field value to JSON
func (src NullInt16) MarshalJSON() ([]byte, error) {
	//fmt.Printf("Marshaling NullInt16: %d\n", src.Status)

	switch src.Status {
	case pgtype.Present:
		return []byte(strconv.FormatInt(int64(src.Int), 10)), nil
	case pgtype.Null:
		return []byte("0"), nil
	case pgtype.Undefined:
		return []byte("0"), nil
	}

	//fmt.Println("NullInt16 is nil")
	return nil, errBadStatus
}

//UnmarshalJSON parse JSON valus and set into field
func (dst *NullInt16) UnmarshalJSON(b []byte) error {
	var v *int16
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	if v == nil {
		*dst = NullInt16{Status: pgtype.Null}
	} else {
		*dst = NullInt16{Int: *v, Status: pgtype.Present}
	}

	return nil
}
