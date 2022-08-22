package datatypes

import (
	"database/sql/driver"
	"encoding/json"
)

// NullString to scan nil string values from database
type NullString struct {
	String string
	Valid  bool // Valid is true if string is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullString) Scan(value interface{}) error {
	if value == nil {
		ns.String, ns.Valid = "", false
		return nil
	}
	ns.String, ns.Valid = toString(value)
	return nil
}

// Value implements the driver Valuer interface.
func (ns NullString) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return ns.String, nil
}

// MarshalJSON implements json.Marshaler.
// It will encode null if this String is null.
func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		ns.String = ""
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON implements json.Unmarshaler.
// It supports string and null input. Blank string input does not produce a null String.
// It also supports unmarshalling a sql.NullString.
func (ns *NullString) UnmarshalJSON(data []byte) error {
	var err error
	var v *string
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	if v == nil {
		*ns = NullString{"", false}
	} else {
		*ns = NullString{*v, true}
	}

	return nil
}

func (ns *NullString) SetValue(val string) {
	ns.String = val
	ns.Valid = (val != "")
}

func NewNullString(val string) NullString {
	ns := NullString{}
	ns.String = val
	ns.Valid = (val != "")
	return ns
}
