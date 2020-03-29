package convert

import (
	"database/sql/driver"
	"time"
)

type DateArrayFromTimeSlice struct {
	V []time.Time
}

func (v DateArrayFromTimeSlice) Value() (driver.Value, error) {
	if v.V == nil {
		return nil, nil
	} else if len(v.V) == 0 {
		return []byte{'{', '}'}, nil
	}

	out := []byte{'{'}

	for _, t := range v.V {
		out = append(out, []byte(t.Format(dateLayout))...)
		out = append(out, ',')
	}

	out[len(out)-1] = '}' // replace last ";" with "}"
	return out, nil
}

type DateArrayToTimeSlice struct {
	V *[]time.Time
}

func (v DateArrayToTimeSlice) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		*v.V = nil
		return nil
	}

	elems := pgparsearray1(data)
	dates := make([]time.Time, len(elems))
	for i := 0; i < len(elems); i++ {
		t, err := pgparsedate(elems[i])
		if err != nil {
			return err
		}
		dates[i] = t
	}

	*v.V = dates
	return nil
}
