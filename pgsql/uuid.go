package pgsql

import (
	"database/sql"
	"database/sql/driver"
)

// UUIDFromByteArray16 returns a driver.Valuer that produces a PostgreSQL uuid from the given Go [16]byte.
func UUIDFromByteArray16(val [16]byte) driver.Valuer {
	return uuidFromByteArray16{val: val}
}

// UUIDToByteArray16 returns an sql.Scanner that converts a PostgreSQL uuid into a Go [16]byte and sets it to val.
func UUIDToByteArray16(val *[16]byte) sql.Scanner {
	return uuidToByteArray16{val: val}
}

type uuidFromByteArray16 struct {
	val [16]byte
}

func (v uuidFromByteArray16) Value() (driver.Value, error) {
	return pgFormatUUID(v.val), nil
}

type uuidToByteArray16 struct {
	val *[16]byte
}

func (v uuidToByteArray16) Scan(src interface{}) error {
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

	*v.val = uuid
	return nil
}
