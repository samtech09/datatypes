package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"net"

	"github.com/pkg/errors"
)

type IPAddr struct {
	IPNet *net.IPNet
	Valid bool
}

func (dst *IPAddr) Scan(src interface{}) error {
	if src == nil {
		*dst = IPAddr{Valid: false}
		return nil
	}

	switch src := src.(type) {
	case string:
		return dst.FromText(src)
	case []byte:
		// srcCopy := make([]byte, len(src))
		// copy(srcCopy, src)
		return dst.FromText(string(src))
	}

	return errors.Errorf("cannot scan %T", src)
}

func (p IPAddr) Value() (driver.Value, error) {
	if p.Valid {
		return p.IPNet.IP.String(), nil
	}
	return nil, nil
}

//MarshalJSON convert field value to JSON
func (src IPAddr) MarshalJSON() ([]byte, error) {
	if src.Valid {
		return json.Marshal(src.IPNet.IP.String())
	}
	return json.Marshal("")
}

//UnmarshalJSON parse JSON valus and set into field
func (dst *IPAddr) UnmarshalJSON(b []byte) (err error) {
	if b[0] == '"' && b[len(b)-1] == '"' {
		b = b[1 : len(b)-1]
	}

	ip := net.ParseIP(string(b))
	msk := net.IPv4Mask(255, 255, 255, 0)
	if ip != nil {
		ipnet := net.IPNet{IP: ip, Mask: msk}
		*dst = IPAddr{IPNet: &ipnet, Valid: true}
		return nil
	}

	*dst = IPAddr{IPNet: nil, Valid: false}
	return nil
}

func (dst *IPAddr) FromText(src string) error {
	if src == "" {
		*dst = IPAddr{Valid: false}
		return nil
	}

	var ipnet *net.IPNet
	var err error

	if ip := net.ParseIP(src); ip != nil {
		ipv4 := ip.To4()
		if ipv4 != nil {
			ip = ipv4
		}
		bitCount := len(ip) * 8
		mask := net.CIDRMask(bitCount, bitCount)
		ipnet = &net.IPNet{Mask: mask, IP: ip}
	} else {
		_, ipnet, err = net.ParseCIDR(src)
		if err != nil {
			*dst = IPAddr{Valid: false}
			return err
		}
	}

	*dst = IPAddr{IPNet: ipnet, Valid: true}
	return nil
}
