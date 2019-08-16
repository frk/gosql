package errors

import (
	"fmt"
)

type errnum uint

func (e errnum) Error() string {
	return fmt.Sprintf("Error number #%d", uint(e))
}

const (
	// TODO(mkopriva): this is a preliminary solution to accommodate the
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
	BadUnaryCmpopError                         errnum = 8
	BadCmpopComboError                         errnum = 9
	ExtraScalarropError                        errnum = 10
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
	BadScalarFieldTypeError                    errnum = 25
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
)
