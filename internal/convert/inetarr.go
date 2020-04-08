package convert

import (
	"database/sql/driver"
	"net"
)

type InetArrayFromIPNetSlice struct {
	Val []net.IPNet
}

func (v InetArrayFromIPNetSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, ipnet := range v.Val {
		out = append(out, []byte(ipnet.String())...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type InetArrayToIPNetSlice struct {
	Val *[]net.IPNet
}

func (v InetArrayToIPNetSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgparsearray1(data)
	ipnets := make([]net.IPNet, len(elems))
	for i := 0; i < len(elems); i++ {
		_, ipnet, err := net.ParseCIDR(string(elems[i]))
		if err != nil {
			ipnet = &net.IPNet{IP: net.ParseIP(string(elems[i]))}
		}
		ipnets[i] = *ipnet
	}

	*v.Val = ipnets
	return nil
}

type InetArrayFromIPSlice struct {
	Val []net.IP
}

func (v InetArrayFromIPSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, ip := range v.Val {
		out = append(out, []byte(ip.String())...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type InetArrayToIPSlice struct {
	Val *[]net.IP
}

func (v InetArrayToIPSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgparsearray1(data)
	ips := make([]net.IP, len(elems))
	for i := 0; i < len(elems); i++ {
		ips[i] = net.ParseIP(string(elems[i]))
	}

	*v.Val = ips
	return nil
}
