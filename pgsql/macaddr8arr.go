package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"net"
)

// MACAddr8ArrayFromHardwareAddrSlice returns a driver.Valuer that produces a PostgreSQL macaddr8[] from the given Go []net.HardwareAddr.
func MACAddr8ArrayFromHardwareAddrSlice(val []net.HardwareAddr) driver.Valuer {
	return macAddrArrayFromHardwareAddrSlice{val: val}
}

// MACAddr8ArrayToHardwareAddrSlice returns an sql.Scanner that converts a PostgreSQL macaddr8[] into a Go []net.HardwareAddr and sets it to val.
func MACAddr8ArrayToHardwareAddrSlice(val *[]net.HardwareAddr) sql.Scanner {
	return macAddrArrayToHardwareAddrSlice{val: val}
}
