package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"net"
)

// CIDRArrayFromIPNetSlice returns a driver.Valuer that produces a PostgreSQL cidr[] from the given Go []net.IPNet.
func CIDRArrayFromIPNetSlice(val []net.IPNet) driver.Valuer {
	return cidrArrayFromIPNetSlice{val: val}
}

// CIDRArrayToIPNetSlice returns an sql.Scanner that converts a PostgreSQL cidr[] into a Go []net.IPNet and sets it to val.
func CIDRArrayToIPNetSlice(val *[]net.IPNet) sql.Scanner {
	return cidrArrayToIPNetSlice{val: val}
}

type cidrArrayFromIPNetSlice struct {
	val []net.IPNet
}

func (v cidrArrayFromIPNetSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for i := 0; i < len(v.val); i++ {
		out = append(out, []byte(v.val[i].String())...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type cidrArrayToIPNetSlice struct {
	val *[]net.IPNet
}

func (v cidrArrayToIPNetSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(data)
	ipnets := make([]net.IPNet, len(elems))
	for i := 0; i < len(elems); i++ {
		_, ipnet, err := net.ParseCIDR(string(elems[i]))
		if err != nil {
			return err
		}
		ipnets[i] = *ipnet
	}

	*v.val = ipnets
	return nil
}
