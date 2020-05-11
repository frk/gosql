package errors

import (
	"fmt"
)

type errnum uint

func (e errnum) Error() string {
	return fmt.Sprintf("Error number #%d", uint(e))
}

const (
	// TODO(mkopriva): this is a temporary solution to accommodate the
	// large variety of errors. Eventually these values should be grouped
	// into as few as possible error types that will replace the values
	// and that will also be able to generate informative error messages.
	NoDBRelationError                     errnum = 62
	NoDBColumnError                       errnum = 63
	NoDBIndexError                        errnum = 64
	NoDBIndexForColumnListError           errnum = 65
	NoDBConstraintError                   errnum = 66
	BadDBColumnTypeError                  errnum = 67
	BadDBIndexError                       errnum = 68
	UnsupportedColumnTypeError            errnum = 69
	IllegalPtrFieldForNotNullColumnError  errnum = 70
	IllegalNullemptyForNotNullColumnError errnum = 71
	IllegalFieldTypeForQuantifierError    errnum = 72

	BadFieldToColumnTypeError            errnum = 73
	BadColumnTypeForDBFuncError          errnum = 74
	BadColumnTypeForQuantifierError      errnum = 76
	BadExpressionTypeForQuantifierError  errnum = 77
	BadLiteralExpressionError            errnum = 78
	UnknownPostgresTypeError             errnum = 79
	BadColumnTypeForUnaryOpError         errnum = 80
	BadColumnNULLSettingForNULLOpError   errnum = 81
	BadColumnToLiteralComparisonError    errnum = 82
	BadColumnToColumnTypeComparisonError errnum = 83
	BadTargetTableForDefaultError        errnum = 84
	NoColumnDefaultSetError              errnum = 87
	ReturnDirectiveWithNoDataFieldError  errnum = 88
	NoFieldColumnError                   errnum = 89
	UnsupportedFieldDataNopError         errnum = 93
	NoColumnFieldError                   errnum = 94
)
