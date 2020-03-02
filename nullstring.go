package datatypes

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgtype"
)

type NullString pgtype.Text

func (p *NullString) Scan(src interface{}) error {
	t := pgtype.Text(*p)
	err := t.Scan(src)
	*p = NullString(t)
	return err
}

func (p NullString) Value() (driver.Value, error) {
	return pgtype.Text(p).Value()
}

//MarshalJSON convert field value to JSON
func (src NullString) MarshalJSON() ([]byte, error) {
	//fmt.Printf("Marshaling NullString: %d\n", src.Status)

	switch src.Status {
	case pgtype.Present:
		return json.Marshal(src.String)
	case pgtype.Null:
		return json.Marshal("")
	case pgtype.Undefined:
		return json.Marshal("")
	}

	//fmt.Println("NullString is nil")
	return nil, errBadStatus
}

//UnmarshalJSON parse JSON valus and set into field
func (dst *NullString) UnmarshalJSON(b []byte) error {
	var s *string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	if s == nil {
		*dst = NullString{Status: pgtype.Null}
	} else {
		*dst = NullString{String: *s, Status: pgtype.Present}
	}

	return nil
}
