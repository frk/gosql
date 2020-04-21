package pg2go

import (
	"database/sql/driver"
)

type UUIDFromByteArray16 struct {
	Val [16]byte
}

func (v UUIDFromByteArray16) Value() (driver.Value, error) {
	return pgFormatUUID(v.Val), nil
}

type UUIDToByteArray16 struct {
	Val *[16]byte
}

func (v UUIDToByteArray16) Scan(src interface{}) error {
	data, err := srcbytes(src)
	if err != nil {
		return err
	} else if data == nil {
		return nil
	}

	uuid, err := pgParseUUID(data)
	if err != nil {
		return err
	}

	*v.Val = uuid
	return nil
}
