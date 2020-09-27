package gosql

import (
	"strconv"
	"strings"
)

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
