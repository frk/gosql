package convert

import (
	"database/sql/driver"
	"net"
)

type InetFromIPNet struct {
	Val net.IPNet
}

func (v InetFromIPNet) Value() (driver.Value, error) {
	if v.Val.IP == nil && v.Val.Mask == nil {
		return nil, nil
	}
	return []byte(v.Val.String()), nil
}

type InetToIPNet struct {
	Val *net.IPNet
}

func (v InetToIPNet) Scan(src interface{}) error {
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
		ipnet = &net.IPNet{IP: net.ParseIP(str)}
	}
	*v.Val = *ipnet
	return nil
}

type InetFromIP struct {
	Val net.IP
}

func (v InetFromIP) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	}
	return []byte(v.Val.String()), nil
}

type InetToIP struct {
	Val *net.IP
}

func (v InetToIP) Scan(src interface{}) error {
	var str string
	switch s := src.(type) {
	case []byte:
		str = string(s)
	case string:
		str = s
	case nil:
		return nil
	}

	*v.Val = net.ParseIP(str)
	return nil
}
