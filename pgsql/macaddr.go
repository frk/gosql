package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"net"
)

// MACAddrFromHardwareAddr returns a driver.Valuer that produces a PostgreSQL macaddr from the given Go net.HardwareAddr.
func MACAddrFromHardwareAddr(val net.HardwareAddr) driver.Valuer {
	return macAddrFromHardwareAddr{val: val}
}

// MACAddrToHardwareAddr returns an sql.Scanner that converts a PostgreSQL macaddr into a Go net.HardwareAddr and sets it to val.
func MACAddrToHardwareAddr(val *net.HardwareAddr) sql.Scanner {
	return macAddrToHardwareAddr{val: val}
}

type macAddrFromHardwareAddr struct {
	val net.HardwareAddr
}

func (v macAddrFromHardwareAddr) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	}
	return []byte(v.val.String()), nil
}

type macAddrToHardwareAddr struct {
	val *net.HardwareAddr
}

func (v macAddrToHardwareAddr) Scan(src interface{}) error {
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

	*v.val = mac
	return nil
}
