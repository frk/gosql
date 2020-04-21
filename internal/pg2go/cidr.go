package pg2go

import (
	"database/sql/driver"
	"net"
)

type CIDRFromIPNet struct {
	Val net.IPNet
}

func (v CIDRFromIPNet) Value() (driver.Value, error) {
	if len(v.Val.IP) == 0 && len(v.Val.Mask) == 0 {
		return nil, nil
	}
	out := v.Val.String()
	return []byte(out), nil
}

type CIDRToIPNet struct {
	Val *net.IPNet
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
	*v.Val = *ipnet
	return nil
}
