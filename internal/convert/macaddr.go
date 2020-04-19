package convert

import (
	"database/sql/driver"
	"net"
)

type MACAddrFromHardwareAddr struct {
	Val net.HardwareAddr
}

func (v MACAddrFromHardwareAddr) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	}
	return []byte(v.Val.String()), nil
}

type MACAddrToHardwareAddr struct {
	Val *net.HardwareAddr
}

func (v MACAddrToHardwareAddr) Scan(src interface{}) error {
	var str string
	switch s := src.(type) {
	case []byte:
		str = string(s)
	case string:
		str = s
	case nil:
		return nil
	}

	mac, err := net.ParseMAC(str)
	if err != nil {
		return err
	}

	*v.Val = mac
	return nil
}
