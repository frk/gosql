package gosql

import (
	"strings"
	"text/template"
)

var errortmpl *template.Template

func init() {
	str := strings.Replace(errortmplstring, "\\\n", "", -1)
	errortmpl = template.Must(template.New("x").Parse(str))
}

type errortype struct {
	code errcode
	args args
}

func (e *errortype) Error() string {
	tmpl, ok := errcode2tmplname[e.code]
	if !ok {
		return "unknown error"
	}

	var out strings.Builder
	if err := errortmpl.ExecuteTemplate(&out, tmpl, e.args); err != nil {
		panic(err)
	}
	return out.String()
}

type args map[string]interface{}

type errcode uint

const (
	_ errcode = iota
	errNoRelation
	errBadRelfieldType
	errBadIteratorType

	// illegal fields
	errIllegalLimitField
	errIllegalOffsetField
	errIllegalCountField
	errIllegalExistsField
	errIllegalNotExistsField
	errIllegalWhereField
	errIllegalFilterField
	errIllegalJoinField
	errIllegalFromField
	errIllegalUsingField
	errIllegalOnConflictField
	errIllegalResultField
	errIllegalRowsAffectedField

	// illegal directives
	errIllegalRelationDirective
	errIllegalRelationJoinDirective
	errIllegalAllDirective
	errIllegalDefaultDirective
	errIllegalReturnDirective
	errIllegalForceDirective
	errIllegalOrderByDirective
	errIllegalOverrideDirective
	errIllegalTextsearchDirective
	errIllegalDirectiveInCommand
	errIllegalDirectiveInJoinBlock
	errIllegalDirectiveInOnConflictBlock

	// conflicting fields / directives
	errConflictAllWhereFilter
	errConflictReturningResultRowsAffected
	errConflictErrorHandlerField
	errConflictTargetInOnConflict
	errConflictActionInOnConflict
	errConflictLimitField
	errConflictOffsetField

	errBadRelId
	errBadColId

	errEmptyColList
	errEmptyOrderByList
	errNoBetweenXYArgs
	errNoOnConflictTarget
	errNoLimitDirectiveValue
	errNoOffsetDirectiveValue

	errBadCmpopCombo
	errBadScalarFieldType
	errNotUnaryCmpop
	errExtraScalarrop

	/////////////
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
	errBadRowsAffectedType

	errBadKind
	errBadType
)

var errcode2tmplname = map[errcode]string{
	errNoRelation:               "missing_relation_field",
	errBadRelfieldType:          "bad_relfield_type",
	errBadIteratorType:          "bad_iterator_type",
	errIllegalLimitField:        "illegal_cmdtype_field_or_directive",
	errIllegalOffsetField:       "illegal_cmdtype_field_or_directive",
	errIllegalCountField:        "illegal_cmdtype_field_or_directive",
	errIllegalExistsField:       "illegal_cmdtype_field_or_directive",
	errIllegalNotExistsField:    "illegal_cmdtype_field_or_directive",
	errIllegalRelationDirective: "illegal_cmdtype_field_or_directive",
	errIllegalAllDirective:      "illegal_cmdtype_field_or_directive",
	errConflictAllWhereFilter:   "conflicting_all_where_filter",
	errBadRelId:                 "bad_relid",
	errBadColId:                 "bad_colid",
}

const errortmplstring = `

{{ define "valid_identifier_text" -}}
A valid identifier MUST begin with a letter (a-z) or an underscore (_), and the \
subsequent characters in the identifier can be letters, underscores, and digits (0-9).
{{- end }}

{{ define "missing_relation_field" -}}
The command type {{ .cmdname }} is missing a "relation" field. \
To fix the issue, make sure that {{ .cmdname }} contains a field marked \
with the ` + "`rel`" + ` tag.
{{- end }}

{{ define "bad_relfield_type" -}}
The {{ .cmdname }}.{{ .relfield }} relation field's type is invalid. \
The field's type MUST be either a struct, a pointer to a \
struct, a slice of structs, a slice of pointers to structs \
or, alternatively, an "iterator" over structs. If the field's \
type is a struct then the struct type CAN be unnamed, however \
if it is any other of the allowed types then the base struct \
type MUST be named.
{{- end }}

{{ define "bad_iterator_type" -}}
The {{ .cmdname }}.{{ .relfield }} relation field's type is an invalid "iterator". \
If the relation field's type is a function or an interface it is \
automatically assumed to be an iterator type, however, to be a valid \
iterator the function MUST take exactly one argument of a named struct \
type and it MUST return exactly one value of type error, when it's an \
interface, it MUST have exactly one method whose signature MUST be the \
same as that of the function described above.
{{- end }}

{{ define "illegal_cmdtype_field_or_directive" -}}
Illegal {{ .fieldname }} field or directive in {{ .cmdname }} command. \
{{ .fieldname }} is only allowed in {{ .allowed }} command(s).
{{- end }}

{{ define "bad_relid" -}}
The "{{ .tagvalue }}" tag value of the "{{ .field }}" field in {{ .cmdname }} command \
is not a valid "relid". A valid relid MUST have the following format: "[qualifier.]name[:alias]" \
where each of the three elements is a valid identifier. {{ template "valid_identifier_text" }}
{{- end }}

{{ define "bad_colid" -}}
The "{{ .tagvalue }}" tag value of the "{{ .field }}" field in {{ .cmdname }} command \
is not a valid "colid". A valid colid MUST have the following format: "[qualifier.]name" \
where each of the two elements is a valid identifier. {{ template "valid_identifier_text" }}
{{- end }}

{{ define "conflicting_all_where_filter" -}}
The "{{ .tagvalue }}" tag value of the "{{ .field }}" field in {{ .cmdname }} command \
is not a valid "colid". A valid colid MUST have the following format: "[qualifier.]name" \
where each of the two elements is a valid identifier. {{ template "valid_identifier_text" }}
{{- end }}
` //`
