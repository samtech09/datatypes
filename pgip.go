package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"net"

	"github.com/jackc/pgtype"
)

type IPAddr pgtype.Inet

func (p *IPAddr) Scan(src interface{}) error {
	inet := pgtype.Inet(*p)
	err := inet.Scan(src)
	*p = IPAddr(inet)
	return err
}

func (p IPAddr) Value() (driver.Value, error) {
	return pgtype.Inet(p).Value()
}

//MarshalJSON convert field value to JSON
func (p IPAddr) MarshalJSON() ([]byte, error) {
	if p.Status == pgtype.Present {
		return json.Marshal(p.IPNet.IP.String())
	}
	return json.Marshal("")
}

//UnmarshalJSON parse JSON valus and set into field
func (p *IPAddr) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	ip := net.ParseIP(string(b))
	msk := net.IPv4Mask(255, 255, 255, 0)
	if ip != nil {
		ipnet := net.IPNet{IP: ip, Mask: msk}
		*p = IPAddr{IPNet: &ipnet, Status: pgtype.Present}
		return nil
	}

	*p = IPAddr{IPNet: nil, Status: pgtype.Null}
	return nil
}
