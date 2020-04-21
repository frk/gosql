package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"net"
)

// InetFromIPNet returns a driver.Valuer that produces a PostgreSQL inet from the given Go net.IPNet.
func InetFromIPNet(val net.IPNet) driver.Valuer {
	return inetFromIPNet{val: val}
}

// InetToIPNet returns an sql.Scanner that converts a PostgreSQL inet into a Go net.IPNet and sets it to val.
func InetToIPNet(val *net.IPNet) sql.Scanner {
	return inetToIPNet{val: val}
}

// InetFromIP returns a driver.Valuer that produces a PostgreSQL inet from the given Go net.IP.
func InetFromIP(val net.IP) driver.Valuer {
	return inetFromIP{val: val}
}

// InetToIP returns an sql.Scanner that converts a PostgreSQL inet into a Go net.IP and sets it to val.
func InetToIP(val *net.IP) sql.Scanner {
	return inetToIP{val: val}
}

type inetFromIPNet struct {
	val net.IPNet
}

func (v inetFromIPNet) Value() (driver.Value, error) {
	if v.val.IP == nil && v.val.Mask == nil {
		return nil, nil
	}
	return []byte(v.val.String()), nil
}

type inetToIPNet struct {
	val *net.IPNet
}

func (v inetToIPNet) Scan(src interface{}) error {
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
	*v.val = *ipnet
	return nil
}

type inetFromIP struct {
	val net.IP
}

func (v inetFromIP) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	}
	return []byte(v.val.String()), nil
}

type inetToIP struct {
	val *net.IP
}

func (v inetToIP) Scan(src interface{}) error {
	var str string
	switch s := src.(type) {
	case []byte:
		str = string(s)
	case string:
		str = s
	case nil:
		return nil
	}

	*v.val = net.ParseIP(str)
	return nil
}
