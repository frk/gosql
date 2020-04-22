package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"net"
)

// MACAddr8FromHardwareAddr returns a driver.Valuer that produces a PostgreSQL macaddr8 from the given Go net.HardwareAddr.
func MACAddr8FromHardwareAddr(val net.HardwareAddr) driver.Valuer {
	return macAddrFromHardwareAddr{val: val}
}

// MACAddr8ToHardwareAddr returns an sql.Scanner that converts a PostgreSQL macaddr8 into a Go net.HardwareAddr and sets it to val.
func MACAddr8ToHardwareAddr(val *net.HardwareAddr) sql.Scanner {
	return macAddrToHardwareAddr{val: val}
}
