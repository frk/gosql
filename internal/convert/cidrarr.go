package convert

import (
	"database/sql/driver"
	"net"
)

type CIDRArrayFromIPNetSlice struct {
	V net.IPNet
}

func (v CIDRArrayFromIPNetSlice) Value() (driver.Value, error) {
	// TODO
	return nil, nil
}

type CIDRArrayToIPNetSlice struct {
	V *net.IPNet
}

func (v CIDRArrayToIPNetSlice) Scan(src interface{}) error {
	// TODO
	return nil
}
