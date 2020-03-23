package gosql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/frk/gosql/internal/convert"
)

type Conn interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type AfterScanner interface {
	AfterScan()
}

type ErrorHandler interface {
	HandleError(err error) error
}

type ErrorInfo struct {
	Error     error
	Query     string
	SpecName  string
	SpecKind  string
	SpecValue interface{}
}

type ErrorInfoHandler interface {
	HandleErrorInfo(info *ErrorInfo) error
}

func InValueList(num, pos int) string {
	var b strings.Builder

	// write the first parameter
	if num > 0 {
		b.WriteString(OrdinalParameters[pos])
	}

	// write the rest with a comma
	for i := 1; i < num; i++ {
		b.WriteByte(',')
		b.WriteString(OrdinalParameters[pos+i])
	}

	return b.String()
}

var OrdinalParameters = func() (a [65535]string) {
	for i := 0; i < len(a); i++ {
		a[i] = "$" + strconv.Itoa(i+1)
	}
	return a
}()

func IntSliceToIntArray(s []int) driver.Valuer {
	return convert.IntSlice2IntArray{S: s}
}

func StringSliceToTextArray(s []string) driver.Valuer {
	return nil // TODO
}

type scannervaluer interface {
	driver.Valuer
	sql.Scanner
}

func JSON(v interface{}) scannervaluer {
	return jsontype{v: v}
}

type jsontype struct {
	v interface{}
}

func (j jsontype) Value() (driver.Value, error) {
	b, err := json.Marshal(j.v)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (j jsontype) Scan(src interface{}) error {
	if b, ok := src.([]byte); ok {
		return json.Unmarshal(b, j.v)
	}
	return nil
}
