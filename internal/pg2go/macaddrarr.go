package pg2go

import (
	"database/sql/driver"
	"net"
)

type MACAddrArrayFromHardwareAddrSlice struct {
	Val []net.HardwareAddr
}

func (v MACAddrArrayFromHardwareAddrSlice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, mac := range v.Val {
		out = append(out, []byte(mac.String())...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type MACAddrArrayToHardwareAddrSlice struct {
	Val *[]net.HardwareAddr
}

func (v MACAddrArrayToHardwareAddrSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseCommaArray(data)
	macs := make([]net.HardwareAddr, len(elems))
	for i := 0; i < len(elems); i++ {
		if macs[i], err = net.ParseMAC(string(elems[i])); err != nil {
			return err
		}
	}

	*v.Val = macs
	return nil
}
