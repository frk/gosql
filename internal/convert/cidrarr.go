package convert

import (
	"database/sql/driver"
	"net"
)

type CIDRArrayFromIPNetSlice struct {
	V []net.IPNet
}

func (v CIDRArrayFromIPNetSlice) Value() (driver.Value, error) {
	if v.V == nil {
		return nil, nil
	} else if len(v.V) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, ipnet := range v.V {
		out = append(out, []byte(ipnet.String())...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type CIDRArrayToIPNetSlice struct {
	V *[]net.IPNet
}

func (v CIDRArrayToIPNetSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.V = nil
		return nil
	}

	elems := pgparsearray1(data)
	ipnets := make([]net.IPNet, len(elems))
	for i := 0; i < len(elems); i++ {
		_, ipnet, err := net.ParseCIDR(string(elems[i]))
		if err != nil {
			return err
		}
		ipnets[i] = *ipnet
	}

	*v.V = ipnets
	return nil
}
