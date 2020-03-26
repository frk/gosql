package convert

import (
	"database/sql/driver"
	"net"
)

type CIDRFromIPNet struct {
	V net.IPNet
}

func (v CIDRFromIPNet) Value() (driver.Value, error) {
	// TODO
	return nil, nil
}

type CIDRToIPNet struct {
	V *net.IPNet
}

func (v CIDRToIPNet) Scan(src interface{}) error {
	// TODO
	return nil
}
