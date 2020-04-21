package pg2go

import (
	"database/sql/driver"
	"net"
)

type CIDRArrayFromIPNetSlice struct {
	Val []net.IPNet
}

func (v CIDRArrayFromIPNetSlice) Value() (driver.Value, error) {
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

type CIDRArrayToIPNetSlice struct {
	Val *[]net.IPNet
}

func (v CIDRArrayToIPNetSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
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

	*v.Val = ipnets
	return nil
}
