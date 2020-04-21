package pg2go

import (
	"database/sql"
	"database/sql/driver"
	"net"
)

// InetArrayFromIPNetSlice returns a driver.Valuer that produces a PostgreSQL inet[] from the given Go []net.IPNet.
func InetArrayFromIPNetSlice(val []net.IPNet) driver.Valuer {
	return inetArrayFromIPNetSlice{val: val}
}

// InetArrayToIPNetSlice returns an sql.Scanner that converts a PostgreSQL inet[] into a Go []net.IPNet and sets it to val.
func InetArrayToIPNetSlice(val *[]net.IPNet) sql.Scanner {
	return inetArrayToIPNetSlice{val: val}
}

// InetArrayFromIPSlice returns a driver.Valuer that produces a PostgreSQL inet[] from the given Go []net.IP.
func InetArrayFromIPSlice(val []net.IP) driver.Valuer {
	return inetArrayFromIPSlice{val: val}
}

// InetArrayToIPSlice returns an sql.Scanner that converts a PostgreSQL inet[] into a Go []net.IP and sets it to val.
func InetArrayToIPSlice(val *[]net.IP) sql.Scanner {
	return inetArrayToIPSlice{val: val}
}

type inetArrayFromIPNetSlice struct {
	val []net.IPNet
}

func (v inetArrayFromIPNetSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, ipnet := range v.val {
		out = append(out, []byte(ipnet.String())...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type inetArrayToIPNetSlice struct {
	val *[]net.IPNet
}

func (v inetArrayToIPNetSlice) Scan(src interface{}) error {
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
			ipnet = &net.IPNet{IP: net.ParseIP(string(elems[i]))}
		}
		ipnets[i] = *ipnet
	}

	*v.val = ipnets
	return nil
}

type inetArrayFromIPSlice struct {
	val []net.IP
}

func (v inetArrayFromIPSlice) Value() (driver.Value, error) {
	if v.val == nil {
		return nil, nil
	} else if len(v.val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, ip := range v.val {
		out = append(out, []byte(ip.String())...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type inetArrayToIPSlice struct {
	val *[]net.IP
}

func (v inetArrayToIPSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.val = nil
		return nil
	}

	elems := pgParseCommaArray(data)
	ips := make([]net.IP, len(elems))
	for i := 0; i < len(elems); i++ {
		ips[i] = net.ParseIP(string(elems[i]))
	}

	*v.val = ips
	return nil
}
