package pg2go

import (
	"database/sql/driver"
	"time"
)

type TsRangeArrayFromTimeArray2Slice struct {
	Val [][2]time.Time
}

func (v TsRangeArrayFromTimeArray2Slice) Value() (driver.Value, error) {
	if v.Val == nil {
		return nil, nil
	} else if len(v.Val) == 0 {
		return []byte{'{', '}'}, nil
	}

	size := (len(v.Val) * 51) + // len(`"[\"yyyy-mm-dd hh:mm:ss\",\"yyyy-mm-dd hh:mm:ss\")"`) == 51
		(len(v.Val) - 1) + // number of commas between array elements
		2 // surrounding curly braces

	out := make([]byte, 1, size)
	out[0] = '{'

	for _, a := range v.Val {
		out = append(out, '"', '[', '\\', '"')
		out = append(out, a[0].Format(timestampLayout)...)
		out = append(out, '\\', '"', ',', '\\', '"')
		out = append(out, a[1].Format(timestampLayout)...)
		out = append(out, '\\', '"', ')', '"', ',')
	}

	out[len(out)-1] = '}' // replace last "," with "}"
	return out, nil
}

type TsRangeArrayToTimeArray2Slice struct {
	Val *[][2]time.Time
}

func (v TsRangeArrayToTimeArray2Slice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.Val = nil
		return nil
	}

	elems := pgParseQuotedStringArray(data)
	slice := make([][2]time.Time, len(elems))
	for i := 0; i < len(elems); i++ {
		a := pgParseRange(elems[i])

		// drop surrounding escaped double quotes
		a[0] = a[0][2 : len(a[0])-2]
		a[1] = a[1][2 : len(a[1])-2]

		var t0, t1 time.Time
		t0, err = time.ParseInLocation(timestampLayout, string(a[0]), noZone)
		if err != nil {
			return err
		}
		t1, err = time.ParseInLocation(timestampLayout, string(a[1]), noZone)
		if err != nil {
			return err
		}

		slice[i][0] = t0
		slice[i][1] = t1
	}

	*v.Val = slice
	return nil
}
