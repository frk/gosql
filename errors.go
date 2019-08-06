package gosql

import (
	"fmt"
)

type analysisError struct {
	code errorCode
	args []interface{}
}

func newerr(code errorCode, args ...interface{}) error {
	return &analysisError{code: code, args: args}
}

func (e *analysisError) Error() string {
	return fmt.Sprintf(errorCodeMessageFormats[e.code], e.args...)
}

type errorCode uint

const (
	errNoRelation errorCode = iota + 1
	errBadRelationType
	errBadIteratorType
	errBadObjId
	errBadColId
	errBadBoolTag
	errBadBetweenType
	errBadDistinctPredicate
	errBadLimitType
	errBadLimitValue
	errBadOffsetType
	errBadOffsetValue
	errBadNullsOrderOption
	errBadOverrideKind
	errBadIndexIdentifier
	errBadConstraintIdentifier

	errBadKind
	errBadType
	// NOTE(mkopriva): the two error codes below can probably be removed
	// since only *types.Named values that ALREADY satisfy the type and
	// name requirements will be passed in to the analysis. Struct types
	// with an invalid name or non-struct types with a valid name should
	// simply be skipped and not bothered with...
	errBadCmdType
	errBadCmdName
)

var errorCodeMessageFormats = map[errorCode]string{
	errNoRelation: "The command type %[1]s is missing a \"relation\" field. " +
		"To fix the issue, make sure that %[1]s contains a field marked " +
		"with the `rel` tag.",

	errBadRelationType: "The %s.%s relation field's type is invalid. " +
		"The field's type MUST be either a struct, a pointer to a " +
		"struct, a slice of structs, a slice of pointers to structs " +
		"or, alternatively, an \"iterator\" over structs. If the field's " +
		"type is a struct then the struct type can be named or unnamed, " +
		"if it is any other of the allowed types then the base struct " +
		"type MUST be named.",

	errBadIteratorType: "The %s.%s relation field's type is an invalid \"iterator\". " +
		"If the relation field's type is a function or an interface it is " +
		"automatically assumed to be an iterator type, however, to be a valid " +
		"iterator the function MUST take exactly one argument of a named struct " +
		"type and it MUST return exactly one value of type error, when it's an " +
		"interface, it MUST have exactly one method whose signature MUST be the " +
		"same as that of the function described above.",

	errBadObjId:       "bad object id",
	errBadColId:       "bad column id",
	errBadBoolTag:     "bad boolean tag",
	errBadBetweenType: "bad between type",
	errBadLimitType:   "bad limit type",

	errBadType: "bad type",
	errBadKind: "bad kind",

	errBadCmdType: "The %s command type must be a struct, instead got %s. " +
		"If this type should be ignored by gosql you can add the \"ignore\" comment marker " +
		"( //gosql:ignore ) right above the type declaration.",

	errBadCmdName: "%q is an invalid command type name. Command type names must " +
		"begin with one of the following verbs: \"insert\", \"update\", " +
		"\"select\", \"delete\", or \"filter\".",
}
