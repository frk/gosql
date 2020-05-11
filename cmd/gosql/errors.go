package main

import (
	"fmt"
)

type analysisErrorCode uint

const (
	_ analysisErrorCode = iota
	errFieldType
	errIterType
	errDataType
	errIllegalField
	errFieldBlock
	errFieldConflict
	errRelTagConflict // multiple fields with `rel` tag
	errNoTargetField
	errNoTagValue

	// invalid values in tags
	errBadTagValue
	errBadBoolTagValue
	errBadColIdTagValue
	errBadRelIdTagValue

	errIllegalUpdateModifier
	errIllegalPredicateQuantifier
	errIllegalUnaryPredicate
	errBadUnaryPredicate
	errBadBetweenPredicate
)

type analysisError struct {
	errorCode   analysisErrorCode
	packagePath string
	structName  string
	blockName   string
	fieldType   string
	fieldName   string
	tagValue    string
	fileName    string
	fileLine    int
}

func (e analysisError) Error() string {
	return fmt.Sprintf("%s:%d: [ TODO ERROR MESG ] ", e.fileName, e.fileLine)
}
