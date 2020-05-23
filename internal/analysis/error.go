package analysis

import (
	"fmt"
)

type ErrorCode uint

const (
	_ ErrorCode = iota

	ErrFieldType
	ErrIterType
	ErrDataType
	ErrIllegalField
	ErrFieldBlock
	ErrFieldConflict
	ErrRelTagConflict // multiple fields with `rel` tag
	ErrNoRelField
	ErrNoTagValue
	ErrNoConflictTarget
	ErrBadTagValue // invalid values in tags
	ErrBadBoolTagValue
	ErrBadColIdTagValue
	ErrBadRelIdTagValue
	ErrIllegalUpdateModifier
	ErrIllegalPredicateQuantifier
	ErrIllegalUnaryPredicate
	ErrBadUnaryPredicate
	ErrBadBetweenPredicate
)

type Error struct {
	Code       ErrorCode
	PkgPath    string
	TargetName string
	BlockName  string
	FieldType  string
	FieldName  string
	TagValue   string
	FileName   string
	FileLine   int
}

func (e Error) Error() string {
	return fmt.Sprintf("%s:%d: [ TODO ERROR MESG ] ", e.FileName, e.FileLine)
}
