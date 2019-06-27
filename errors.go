package gosql

import (
	"fmt"
)

type args []interface{} // convenience type

type analysisError struct {
	code errorCode
	args args
}

func (e *analysisError) Error() string {
	return fmt.Sprintf(errorCodeMessageFormats[e.code], e.args...)
}

type errorCode uint

const (
	badTypeError errorCode = iota + 1
	badTypeKindError
	badCmdTypeError
	badCmdNameError
	noRecordError
	manyRecordError
	badRecordTypeError
	badIteratorTypeError
)

var errorCodeMessageFormats = map[errorCode]string{
	badTypeError:     "",
	badTypeKindError: "",
	badCmdTypeError: "foo" +
		"bar",
	badCmdNameError:      "",
	noRecordError:        "",
	manyRecordError:      "",
	badRecordTypeError:   "",
	badIteratorTypeError: "",
}
