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
	NoRelfieldError                            errnum = 1
	NoOnConflictTargetError                    errnum = 2
	NoLimitDirectiveValueError                 errnum = 3
	NoOffsetDirectiveValueError                errnum = 4
	NoBetweenXYArgsError                       errnum = 5
	EmptyColListError                          errnum = 6
	EmptyOrderByListError                      errnum = 7
	BadUnaryPredicateError                     errnum = 8
	BadPredicateComboError                     errnum = 9
	ExtraQuantifierError                       errnum = 10
	BadRelIdError                              errnum = 11
	BadColIdError                              errnum = 12
	BadBoolTagValueError                       errnum = 13
	BadLimitValueError                         errnum = 14
	BadOffsetValueError                        errnum = 15
	BadIndexIdentifierValueError               errnum = 16
	BadConstraintIdentifierValueError          errnum = 17
	BadOverrideKindValueError                  errnum = 18
	BadNullsOrderOptionValueError              errnum = 19
	BadRelfieldTypeError                       errnum = 20
	BadIteratorTypeError                       errnum = 21
	BadBetweenTypeError                        errnum = 22
	BadLimitTypeError                          errnum = 23
	BadOffsetTypeError                         errnum = 24
	BadQuantifierFieldTypeError                errnum = 25
	BadRowsAffectedTypeError                   errnum = 26
	IllegalCommandDirectiveError               errnum = 27
	IllegalJoinBlockDirectiveError             errnum = 28
	IllegalJoinBlockRelationDirectiveError     errnum = 29
	IllegalOnConflictBlockDirectiveError       errnum = 30
	IllegalFilterFieldError                    errnum = 31
	IllegalCountFieldError                     errnum = 32
	IllegalExistsFieldError                    errnum = 33
	IllegalNotExistsFieldError                 errnum = 34
	IllegalResultFieldError                    errnum = 35
	IllegalRowsAffectedFieldError              errnum = 36
	IllegalLimitFieldOrDirectiveError          errnum = 37
	IllegalOffsetFieldOrDirectiveError         errnum = 38
	IllegalAllDirectiveError                   errnum = 39
	IllegalDefaultDirectiveError               errnum = 40
	IllegalForceDirectiveError                 errnum = 41
	IllegalOrderByDirectiveError               errnum = 42
	IllegalOverrideDirectiveError              errnum = 43
	IllegalRelationDirectiveError              errnum = 44
	IllegalReturnDirectiveError                errnum = 45
	IllegalTextSearchDirectiveError            errnum = 46
	IllegalFromBlockError                      errnum = 47
	IllegalJoinBlockError                      errnum = 48
	IllegalOnConflictBlockError                errnum = 49
	IllegalUsingBlockError                     errnum = 50
	IllegalWhereBlockError                     errnum = 51
	ConflictWhereProducerError                 errnum = 52
	ConflictResultProducerError                errnum = 53
	ConflictErrorHandlerFieldError             errnum = 54
	ConflictJoinBlockRelationDirectiveError    errnum = 55
	ConflictOnConflictBlockTargetProducerError errnum = 56
	ConflictOnConflictBlockActionProducerError errnum = 57
	ConflictLimitProducerError                 errnum = 58
	ConflictOffsetProducerError                errnum = 59
	BadWhereBlockTypeError                     errnum = 60
	BadJoinBlockTypeError                      errnum = 61
	BadOnConflictBlockTypeError                errnum = 61
	NoDBRelationError                          errnum = 62
	NoDBColumnError                            errnum = 63
	NoDBIndexError                             errnum = 64
	NoDBIndexForColumnListError                errnum = 65
	NoDBConstraintError                        errnum = 66
	BadDBColumnTypeError                       errnum = 67
	BadDBIndexError                            errnum = 68
	UnsupportedColumnTypeError                 errnum = 69
	IllegalPtrFieldForNotNullColumnError       errnum = 70
	IllegalNullemptyForNotNullColumnError      errnum = 71
	IllegalFieldTypeForQuantifierError         errnum = 72
	BadFieldToColumnTypeError                  errnum = 73
	BadColumnTypeForDBFuncError                errnum = 74
	IllegalUnaryPredicateError                 errnum = 75
	BadColumnTypeForQuantifierError            errnum = 76
	BadExpressionTypeForQuantifierError        errnum = 77
	BadLiteralExpressionError                  errnum = 78
	UnknownPostgresTypeError                   errnum = 79
	BadColumnTypeForUnaryOpError               errnum = 80
	BadColumnNULLSettingForNULLOpError         errnum = 81
	BadColumnToLiteralComparisonError          errnum = 82
	BadColumnToColumnTypeComparisonError       errnum = 83
	BadTargetTableForDefaultError              errnum = 84
	BadUseJSONTargetColumnError                errnum = 85
	BadUseXMLTargetColumnError                 errnum = 86
	NoColumnDefaultSetError                    errnum = 87
	ReturnDirectiveWithNoRelfieldError         errnum = 88
	NoFieldColumnError                         errnum = 89
	MultipleRelfieldsError                     errnum = 90
	IllegalIteratorRecordError                 errnum = 91
)
