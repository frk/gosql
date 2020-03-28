package convert

import (
	"database/sql/driver"
	"net"
)

type CIDRFromIPNet struct {
	V net.IPNet
}

func (v CIDRFromIPNet) Value() (driver.Value, error) {
	if len(v.V.IP) == 0 && len(v.V.Mask) == 0 {
		return nil, nil
	}
	out := v.V.String()
	return []byte(out), nil
}

type CIDRToIPNet struct {
	V *net.IPNet
}

func (v CIDRToIPNet) Scan(src interface{}) error {
	var str string
	switch s := src.(type) {
	case []byte:
		str = string(s)
	case string:
		str = s
	case nil:
		return nil
	}

	_, ipnet, err := net.ParseCIDR(str)
	if err != nil {
		return err
	}
	*v.V = *ipnet
	return nil
}
