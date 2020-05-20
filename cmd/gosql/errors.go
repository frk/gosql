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
	errorCode  analysisErrorCode
	pkgPath    string
	targetName string
	blockName  string
	fieldType  string
	fieldName  string
	tagValue   string
	fileName   string
	fileLine   int
}

func (e analysisError) Error() string {
	return fmt.Sprintf("%s:%d: [ TODO ERROR MESG ] ", e.fileName, e.fileLine)
}

type typeErrorCode uint

const (
	_ typeErrorCode = iota
	errNoDatabaseRelation
	errNoRelationColumn
	errNoUniqueIndex
	errNoUniqueConstraint
	errNoColumnDefault
	errBadColumnQualifier
	errNoColumnField
	errBadColumnReadType
	errBadColumnWriteType
	errBadColumnReadIfaceType
	errBadColumnPredicate
)

type typeError struct {
	errorCode    typeErrorCode
	pkgPath      string
	targetName   string
	fieldType    string
	fieldName    string
	tagValue     string
	dbName       string
	relQualifier string
	relName      string
	colQualifier string
	colName      string
	colType      string
	fileName     string
	fileLine     int
}

func (e typeError) Error() string {
	return fmt.Sprintf("%s:%d: [ TODO ERROR MESG ] ", e.fileName, e.fileLine)
}

type predicateErrorCode uint

const (
	_ predicateErrorCode = iota
	errBadUnaryPredicateType
	errBadPredicateQuantifierType
	errBadComparisonOperation
)

type predicateError struct {
	errorCode  predicateErrorCode
	pkgPath    string
	targetName string
	fieldType  string
	fieldName  string
	tagValue   string
	dbName     string
	lhsRel     string
	lhsExpr    string
	lhsType    string
	rhsRel     string
	rhsExpr    string
	rhsType    string
	predicate  string
	quantifier string
	fileName   string
	fileLine   int
}

func (e predicateError) Error() string {
	return fmt.Sprintf("%s:%d: [ TODO ERROR MESG ] ", e.fileName, e.fileLine)
}
