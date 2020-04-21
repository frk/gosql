package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"net"
)

// CIDRFromIPNet returns a driver.Valuer that produces a PostgreSQL cidr from the given Go net.IPNet.
func CIDRFromIPNet(val net.IPNet) driver.Valuer {
	return cidrFromIPNet{val: val}
}

// CIDRToIPNet returns an sql.Scanner that converts a PostgreSQL cidr into a Go net.IPNet and sets it to val.
func CIDRToIPNet(val *net.IPNet) sql.Scanner {
	return cidrToIPNet{val: val}
}

type cidrFromIPNet struct {
	val net.IPNet
}

func (v cidrFromIPNet) Value() (driver.Value, error) {
	if len(v.val.IP) == 0 && len(v.val.Mask) == 0 {
		return nil, nil
	}
	out := v.val.String()
	return []byte(out), nil
}

type cidrToIPNet struct {
	val *net.IPNet
}

func (v cidrToIPNet) Scan(src interface{}) error {
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
	*v.val = *ipnet
	return nil
}
