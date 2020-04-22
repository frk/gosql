package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"net"
)

// MACAddrArrayFromHardwareAddrSlice returns a driver.Valuer that produces a PostgreSQL macaddr[] from the given Go []net.HardwareAddr.
func MACAddrArrayFromHardwareAddrSlice(val []net.HardwareAddr) driver.Valuer {
	return macAddrArrayFromHardwareAddrSlice{val: val}
}

// MACAddrArrayToHardwareAddrSlice returns an sql.Scanner that converts a PostgreSQL macaddr[] into a Go []net.HardwareAddr and sets it to val.
func MACAddrArrayToHardwareAddrSlice(val *[]net.HardwareAddr) sql.Scanner {
	return macAddrArrayToHardwareAddrSlice{val: val}
}

type macAddrArrayFromHardwareAddrSlice struct {
	val []net.HardwareAddr
}

func (v macAddrArrayFromHardwareAddrSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, mac := range v.val {
		out = append(out, []byte(mac.String())...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type macAddrArrayToHardwareAddrSlice struct {
	val *[]net.HardwareAddr
}

func (v macAddrArrayToHardwareAddrSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(data)
	macs := make([]net.HardwareAddr, len(elems))
	for i := 0; i < len(elems); i++ {
		if macs[i], err = net.ParseMAC(string(elems[i])); err != nil {
			return err
		}
	}

	*v.val = macs
	return nil
}
