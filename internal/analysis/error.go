package analysis

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/frk/tagutil"
)

type Error struct {
	Code          errorCode
	PkgPath       string
	TargetName    string
	BlockName     string
	FieldType     string
	FieldTypeKind string
	FieldName     string
	// The full tag string.
	TagString string
	// The specific tag expression that caused the error.
	TagExpr string
	// The specific tag value that caused the error.
	TagError string
	// The name of the file in which the error occurs.
	FileName string
	// The specific line of the file at which the error occurs.
	FileLine int
}

func (e *Error) Error() string {
	sb := new(strings.Builder)
	if err := error_templates.ExecuteTemplate(sb, e.Code.name(), e); err != nil {
		panic(err)
	}
	return sb.String()
}

func (e *Error) IsDirective() bool {
	return e.FieldName == "_"
}

func (e *Error) IsSelectQueryKind() bool {
	return strings.HasPrefix(tolower(e.TargetName), "select")
}

func (e *Error) TagExprIsUnary() bool {
	_, _, _, rhs := parsePredicateExpr(e.TagExpr)
	return len(rhs) == 0 && e.IsDirective()
}

func (e *Error) TagExprHasRHS() bool {
	_, _, _, rhs := parsePredicateExpr(e.TagExpr)
	return len(rhs) > 0
}

func (e *Error) TagExprPredicate() string {
	_, op, _, _ := parsePredicateExpr(e.TagExpr)
	return op
}

func (e *Error) IsSequence() bool {
	return e.FieldTypeKind == "slice" || e.FieldTypeKind == "array"
}

func (e *Error) TagValueRel() string {
	return tagutil.New(e.TagString).First("rel")
}

func (e *Error) TagValueSql() string {
	return tagutil.New(e.TagString).Get("sql")
}

func (e *Error) TagValueSqlFirst() string {
	return tagutil.New(e.TagString).First("sql")
}

func (e *Error) TagValueSqlSecond() string {
	return tagutil.New(e.TagString).Second("sql")
}

func (e *Error) FieldDefinition() string {
	if len(e.TagString) > 0 {
		return fmt.Sprintf("%s %s `%s`", e.FieldName, e.FieldTypeR(), e.TagString)
	}
	return fmt.Sprintf("%s %s", e.FieldName, e.FieldTypeR())
}

// Reformats the field type to make it more readable.
func (e *Error) FieldTypeR() string {
	// non-empty anon struct?
	if strings.HasPrefix(e.FieldType, "struct{") && e.FieldType != "struct{}" && e.FieldType != "struct{ }" {
		return "struct{ ... }"
	}

	// At this point, AFAICT, only anonymous structs are a pain to read in the
	// error messages. If, later on, other types reveal themselves as unsuitable
	// for error messages, this method can be extended to handle them.
	return e.FieldType
}

func (e *Error) FieldKind() (out string) {
	if e.FieldName == "_" {
		return "directive"
	}
	return "field"
}

func (e *Error) TargetXxx() (out string) {
	key := tolower(e.TargetName)
	if len(key) > 5 {
		key = key[:6]
	}
	switch key {
	case "insert":
		out = "InsertXxx"
	case "update":
		out = "UpdateXxx"
	case "select":
		out = "SelectXxx"
	case "delete":
		out = "DeleteXxx"
	case "filter":
		out = "FilterXxx"
	default:
		out = "<unknown>"
	}
	return out
}

func (e *Error) TargetKind() (out string) {
	key := tolower(e.TargetName)
	if len(key) > 5 {
		key = key[:6]
	}
	switch key {
	case "insert", "update", "select", "delete":
		out = "query"
	case "filter":
		out = "filter"
	default:
		out = "<unknown>"
	}
	return out
}

type errorCode uint8

func (e errorCode) name() string { return fmt.Sprintf("error_template_%d", e) }

const (
	_ errorCode = iota
	errBadFieldTypeInt
	errBadFieldTypeStruct
	errBadIterTypeInterface
	errBadIterTypeFunc
	errBadRelType
	errIllegalQueryField
	errIllegalStructDirective
	errConflictingWhere
	errConflictingOnConfictTarget
	errConflictingOnConfictAction
	errConflictingResultTarget
	errConflictingRelationDirective
	errConflictingFieldOrDirective
	errConflictingRelTag
	errMissingRelField
	errMissingTagValue
	errMissingTagColumnList
	errMissingOnConflictTarget
	errBadIdentTagValue
	errBadColIdTagValue
	errBadRelIdTagValue
	errBadBoolTagValue
	errBadUIntegerTagValue
	errBadNullsOrderTagValue
	errBadOverrideTagValue
	errBadDirectiveBooleanExpr
	errBadBetweenPredicate
	errIllegalSliceUpdateModifier
	errIllegalListPredicate
	errIllegalUnaryPredicate
	errIllegalFieldQuantifier
	errIllegalPredicateQuantifier
)

var error_template_string = `
{{ define "` + errBadFieldTypeInt.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Bad field type.{{off}}
    Cannot use {{R .FieldTypeR}} as the type of the {{Wb .FieldName}} field in {{Wb .TargetName}}.
    {{Wb}}FIX:{{off}} change the field type to one of: {{Ci}}int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64{{off}}.
{{ end }}

{{ define "` + errBadFieldTypeStruct.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Bad field type.{{off}}
    Cannot use {{R .FieldTypeR}} (kind {{R .FieldTypeKind}}) as the type of the {{Wb .FieldName}} field in {{Wb .TargetName}}.
    {{Wb}}FIX:{{off}} change the field type to be of kind {{Ci}}struct{{off}}.
{{ end }}

{{ define "` + errBadIterTypeInterface.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Bad iterator interface type.{{off}}
    The type {{R .FieldTypeR}} of the {{.TargetName}}.{{Wb .FieldName}} field {{Wu}}IS NOT{{off}} a valid {{W}}"iterator interface"{{off}} type.
    {{Wb}}HINT:{{off}} a valid {{W}}"iterator interface"{{off}} type {{Wu}}MUST{{off}} satisfy all of the following requirements:
        - The interface MUST have exactly {{Wu}}one{{off}} method.
        - The interface method MUST be {{Wu}}accessible{{off}} from the generated code, i.e. the method must either be exported or, if unexported, the interface type itself must be anonymous or declared in the same package as the target query type.
        - The interface method's signature MUST have exactly {{Wu}}one parameter{{off}} type and exactly {{Wu}}one result{{off}} type.
        - The interface method's parameter MUST be of a {{Wu}}named struct{{off}} type, or a pointer to a {{Wu}}named struct{{off}} type.
        - The interface method's result type MUST be the standard {{Ci}}error{{off}} type.
{{ end }}

{{ define "` + errBadIterTypeFunc.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Bad iterator func type.{{off}}
    The type {{R .FieldTypeR}} of the {{.TargetName}}.{{Wb .FieldName}} field {{Wu}}IS NOT{{off}} a valid {{W}}"iterator func"{{off}} type.
    {{Wb}}HINT:{{off}} a valid {{W}}"iterator func"{{off}} type {{Wu}}MUST{{off}} satisfy all of the following requirements:
        - The function's signature MUST have exactly {{Wu}}one parameter{{off}} type and exactly {{Wu}}one result{{off}} type.
        - The function's parameter MUST be of a {{Wu}}named struct{{off}} type, or a pointer to a {{Wu}}named struct{{off}} type.
        - The function's result type MUST be the standard {{Ci}}error{{off}} type.
{{ end }}

{{ define "` + errBadRelType.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Bad "rel" type.{{off}}
    The type {{R .FieldTypeR}} of the {{.TargetName}}.{{Wb .FieldName}} field {{Wu}}IS NOT{{off}} a valid {{W}}"iterator func"{{off}} type.
    {{Wb}}HINT:{{off}} a valid {{W}}"rel"{{off}} type {{Wu}}MUST{{off}} be one of the following:
        - A {{Wu}}named or unnamed struct{{off}} type.
        - A {{Wu}}pointer to a named struct{{off}} type.
        - A {{Wu}}slice of named structs{{off}} type.
        - A {{Wu}}slice of pointers to named structs{{off}} type.
        - An {{Wu}}array of named structs{{off}} type.
        - An {{Wu}}array of pointers to named structs{{off}} type.
        - A valid {{Wu}}iterator interface{{off}} type.
        - A valid {{Wu}}iterator func{{off}} type.
{{ end }}

{{ define "` + errIllegalQueryField.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Illegal {{.FieldKind}}.{{off}}
    The {{Wb .TargetXxx}} {{.TargetKind}} types {{Wu}}DO NOT{{off}} support the {{R .FieldDefinition}} {{.FieldKind}}.
    {{Wb}}FIX:{{off}} remove the {{R .FieldDefinition}} {{.FieldKind}} from the {{Wb .TargetName}} {{.TargetKind}} type.
{{ end }}

{{ define "` + errIllegalStructDirective.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Illegal directive field.{{off}}
    The {{Wb .BlockName}} struct {{Wu}}DOES NOT{{off}} support the {{R .FieldName}} {{R .FieldTypeR}} directive.
    {{Wb}}FIX:{{off}} remove the {{R .FieldName}} {{R .FieldTypeR}} directive from the {{Wb .BlockName}} struct of the {{W .TargetName}} {{.TargetKind}} type.
{{ end }}

{{ define "` + errConflictingWhere.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Conflicting "where" fields / directives.{{off}}
    {{if .IsDirective -}}
    The {{R .FieldTypeR}} directive in {{Wb .TargetName}} is in conflict with another field or directive.
    {{else -}}
    The field {{R .FieldName}} (type {{R .FieldTypeR}}) in {{Wb .TargetName}} is in conflict with another field or directive.
    {{end -}}
    {{if .IsSelectQueryKind -}}
    {{Wb}}HINT:{{off}} The {{Wb .TargetXxx}} query types can have {{Wu}}only one{{off}} WHERE producing field which MUST be one of the following:
        - A field of type {{Ci}}gosql.Filter{{off}}.
        - A struct field named {{Wu}}Where{{off}} (case insensitive).
    {{else -}}
    {{Wb}}HINT:{{off}} The {{Wb .TargetXxx}} query types can have {{Wu}}only one{{off}} WHERE producing field or directive which MUST be one of the following:
        - A directive field of type {{Ci}}gosql.All{{off}}.
        - A field of type {{Ci}}gosql.Filter{{off}}.
        - A struct field named {{Wu}}Where{{off}} (case insensitive).
    {{end -}}
{{ end }}

{{ define "` + errConflictingOnConfictTarget.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Conflicting "on conflict" target directives.{{off}}
    The {{R .FieldTypeR}} directive is in conflict with another "target" directive of the {{Wb .BlockName}} struct in the {{Wb .TargetName}} query type.
    {{Wb}}HINT:{{off}} The {{Wb .BlockName}} struct can have {{Wu}}only one{{off}} of the following "target" directives:
        - The {{Ci}}gosql.Column{{off}} directive.
        - The {{Ci}}gosql.Index{{off}} directive.
        - The {{Ci}}gosql.Constraint{{off}} directive.
{{ end }}

{{ define "` + errConflictingOnConfictAction.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Conflicting "on conflict" action directives.{{off}}
    The {{R .FieldTypeR}} directive is in conflict with another "action" directive of the {{Wb .BlockName}} struct in the {{Wb .TargetName}} query type.
    {{Wb}}HINT:{{off}} The {{Wb .BlockName}} struct can have {{Wu}}only one{{off}} of the following "action" directives:
        - The {{Ci}}gosql.Ignore{{off}} directive.
        - The {{Ci}}gosql.Update{{off}} directive.
{{ end }}

{{ define "` + errConflictingResultTarget.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Conflicting "result" fields / directives.{{off}}
    {{if .IsDirective -}}
    The {{R .FieldTypeR}} directive in {{Wb .TargetName}} is in conflict with another "result" field or directive.
    {{else -}}
    The field {{R .FieldName}} (type {{R .FieldTypeR}}) in {{Wb .TargetName}} is in conflict with another "result" field or directive.
    {{end -}}
    {{Wb}}HINT:{{off}} The {{Wb}}InsertXxx{{off}}, {{Wb}}UpdateXxx{{off}}, and {{Wb}}DeleteXxx{{off}} query struct types can have {{Wu}}only one{{off}} of the following:
        - The {{Ci}}gosql.Return{{off}} directive.
	- A field named {{Wu}}RowsAffected{{off}} (case insensitive) which MUST be of one of the following types: {{Ci}}int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64{{off}}.
        - A field named {{Wu}}Result{{off}} (case insensitive) whose type MUST be one of the following:
            - A {{Wu}}named or unnamed struct{{off}} type.
            - A {{Wu}}pointer to a named struct{{off}} type.
            - A {{Wu}}slice of named structs{{off}} type.
            - A {{Wu}}slice of pointers to named structs{{off}} type.
            - An {{Wu}}array of named structs{{off}} type.
            - An {{Wu}}array of pointers to named structs{{off}} type.
            - A valid {{Wu}}iterator interface{{off}} type.
            - A valid {{Wu}}iterator func{{off}} type.
{{ end }}

{{ define "` + errConflictingRelationDirective.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Conflicting "gosql.Relation" directive.{{off}}
    The {{R .FieldTypeR}} directive is in conflict with another directive of the same type.
    {{Wb}}FIX:{{off}} Make sure that the {{Wb .BlockName}} struct field in {{Wb .TargetName}} has {{Wu}}only one{{off}} {{R .FieldTypeR}} directive. 
{{ end }}

{{ define "` + errConflictingFieldOrDirective.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Conflicting fields / directives.{{off}}
    {{if .IsDirective -}}
    The {{R .FieldTypeR}} directive is in conflict with another field or directive of its "kind".
    {{Wb}}FIX:{{off}} Make sure that the {{Wb .TargetName}} struct type has, at most, {{Wu}}only one{{off}} field / directive of the same "kind" as the {{R .FieldTypeR}} directive.
    {{else -}}
    The field {{R .FieldName}} (type {{R .FieldTypeR}}) in {{Wb .TargetName}} is in conflict with another field or directive of its "kind".
    {{Wb}}FIX:{{off}} Make sure that the {{Wb .TargetName}} struct type has, at most, {{Wu}}only one{{off}} field / directive of the same "kind" as the {{R .FieldName}} (type {{R .FieldTypeR}}) field.
    {{end -}}
{{ end }}

{{ define "` + errConflictingRelTag.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Conflicting "rel" tag.{{off}}
    The field {{R .FieldDefinition}} in {{Wb .TargetName}} is in conflict with another field that also has a "rel" tag.
    {{Wb}}FIX:{{off}} Make sure that the {{Wb .TargetName}} query type has {{Wu}}only one{{off}} field with the "rel" tag. 
{{ end }}

{{ define "` + errMissingRelField.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Missing "rel" field.{{off}}
    The {{R .TargetName}} struct type has no field with the "rel" tag.
    {{Wb}}FIX:{{off}} Make sure that the {{R .TargetName}} struct type has {{Wu}}exactly one{{off}} field with the "rel" tag. 
{{ end }}

{{ define "` + errMissingTagValue.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Missing "sql" tag value.{{off}}
    The {{R .FieldName}} {{R .FieldTypeR}} field / directive in {{Wb .TargetName}} is missing a value in its "sql" tag, or it is missing the tag completely.
    {{Wb}}FIX:{{off}} Make sure that the {{R .FieldName}} {{R .FieldTypeR}} field / directive in {{Wb .TargetName}} has the "sql" tag with a value that's valid for that field / directive. 
{{ end }}

{{ define "` + errMissingTagColumnList.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Missing "sql" tag with column list.{{off}}
    The {{R .FieldName}} {{R .FieldTypeR}} field / directive in {{Wb .TargetName}} is missing a list of columns in its "sql" tag, or it is missing the tag completely.
    {{Wb}}FIX:{{off}} Make sure that the {{R .FieldName}} {{R .FieldTypeR}} field / directive in {{Wb .TargetName}} has the "sql" tag with a valid list of columns. 
{{ end }}

{{ define "` + errMissingOnConflictTarget.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Missing on_conflict target directive.{{off}}
    The {{R .FieldName}} struct field in {{Wb .TargetName}} is missing a "target" directive.
    {{Wb}}FIX:{{off}} Make sure that {{R .FieldName}} in {{Wb .TargetName}} has {{Wb}}exactly one{{off}} of the following "target" directives: 
        - The {{Ci}}gosql.Column{{off}} directive.
        - The {{Ci}}gosql.Index{{off}} directive.
        - The {{Ci}}gosql.Constraint{{off}} directive.
{{ end }}

{{ define "` + errBadIdentTagValue.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Bad identifier value in tag.{{off}}
    The "sql" tag value {{R .TagValueSqlFirst}} in {{R .FieldDefinition}} from {{Wb .TargetName}} is an invalid {{Wi}}identifier{{off}}.
    {{Wb}}HINT:{{off}} A valid {{Wi}}identifier{{off}} MUST begin with a letter (a-z [case insensitive]) or an underscore (_), ` +
	`subsequent characters in the {{Wi}}identifier{{off}} can be letters, underscores, and digits (0-9). Put another way, ` +
	`a valid {{Wi}}identifier{{off}} MUST match the following regular expression: {{G}}` + rxIdent.String() + `{{off}}
{{ end }}

{{ define "` + errBadColIdTagValue.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Bad column_identifier value in tag.{{off}}
    The "sql" tag value {{R .TagError}} in {{R .FieldDefinition}} from {{Wb .TargetName}} is an invalid {{Wi}}column_identifier{{off}}.
    {{Wb}}HINT:{{off}} A valid {{Wi}}column_identifier{{off}} MUST be of the format {{Wi}}[rel_alias.]column_name{{off}} where ` +
	`both the {{Wi}}rel_alias{{off}} and the {{Wi}}column_name{{off}} are valid {{Wi}}identifiers{{off}}. A valid ` +
	`{{Wi}}identifier{{off}} MUST begin with a letter (a-z [case insensitive]) or an underscore (_) and any subsequent ` +
	`characters in the {{Wi}}identifier{{off}} can be letters, underscores, and digits (0-9). Put another way, a valid ` +
	`{{Wi}}column_identifier{{off}} MUST match the following regular expression: {{G}}` + rxColIdent.String() + `{{off}}
{{ end }}

{{ define "` + errBadRelIdTagValue.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Bad relation_identifier value in tag.{{off}}
    The "sql" tag value {{R .TagError}} in {{R .FieldDefinition}} from {{Wb .TargetName}} is an invalid {{Wi}}relation_identifier{{off}}.
    {{Wb}}HINT:{{off}} A valid {{Wi}}relation_identifier{{off}} MUST be of the format {{Wi}}[schema_name.]relation_name[:alias_name]{{off}} ` +
	`where {{Wi}}schema_name{{off}}, {{Wi}}relation_name{{off}}, and {{Wi}}alias_name{{off}} are all valid {{Wi}}identifiers{{off}}. ` +
	`A valid {{Wi}}identifier{{off}} MUST begin with a letter (a-z [case insensitive]) or an underscore (_) and any subsequent ` +
	`characters in the {{Wi}}identifier{{off}} can be letters, underscores, and digits (0-9). Put another way, a valid ` +
	`{{Wi}}relation_identifier{{off}} MUST match the following regular expression: {{G}}` + rxRelIdent.String() + `{{off}}
{{ end }}

{{ define "` + errBadBoolTagValue.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Bad boolean value in tag.{{off}}
    The "bool" tag value {{R .TagError}} in {{R .FieldDefinition}} from {{Wb .TargetName}} is an invalid {{Wi}}boolean{{off}} value.
    {{Wb}}HINT:{{off}} A valid {{Wi}}boolean{{off}} value MUST be either {{G}}AND{{off}} or {{G}}OR{{off}} (case insensitive).
{{ end }}

{{ define "` + errBadUIntegerTagValue.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Bad unsigned_integer value in tag.{{off}}
    The "sql" tag value {{R .TagError}} in {{R .FieldDefinition}} from {{Wb .TargetName}} is an invalid {{Wi}}unsigned_integer{{off}} value.
    {{Wb}}FIX:{{off}} Change {{R .TagError}} to a valid unsigned integer value.
{{ end }}

{{ define "` + errBadNullsOrderTagValue.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Bad nulls_order value in tag.{{off}}
    The "sql" tag value {{R .TagError}} in {{R .FieldDefinition}} from {{Wb .TargetName}} is an invalid {{Wi}}nulls_order{{off}} value.
    {{Wb}}HINT:{{off}} A valid {{Wi}}nulls_order{{off}} value MUST be either {{G}}NULLSFIRST{{off}} or {{G}}NULLSLAST{{off}} (case insensitive).
{{ end }}

{{ define "` + errBadOverrideTagValue.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Bad overriding_kind value in tag.{{off}}
    The "sql" tag value {{R .TagError}} in {{R .FieldDefinition}} from {{Wb .TargetName}} is an invalid {{Wi}}overriding_kind{{off}} value.
    {{Wb}}HINT:{{off}} A valid {{Wi}}overriding_kind{{off}} value MUST be either {{G}}USER{{off}} or {{G}}SYSTEM{{off}} (case insensitive).
{{ end }}

{{ define "` + errBadDirectiveBooleanExpr.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Bad boolean_expression in directive.{{off}}
    The expression "{{R .TagError}}" in {{R .FieldDefinition}} from {{W .TargetName}}{{if .BlockName}}.{{Wb .BlockName}}{{end}} ` +
	`is an invalid {{Wi}}boolean_expression{{off}} for a directive.
    {{Wb}}HINT:{{off}} A valid {{Wi}}boolean_expression{{off}} in a directive MUST be a {{Wu}}complete{{off}} binary or unary expression.
{{ end }}

{{ define "` + errBadBetweenPredicate.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Bad between_predicate type.{{off}}
    The {{R .FieldName}} {{R .FieldTypeR}} field from {{Wb .TargetName}} has an invalid {{Wi}}between_predicate{{off}} type.
    {{Wb}}FIX:{{off}} A valid {{Wi}}between_predicate{{off}} type MUST be a {{Ci}}struct{{off}} type that satisfies the following requirements:
        - The {{Ci}}struct{{off}} type MUST have {{Wu}}exactly two{{off}} fields, no more and no less.
        - One of the fields MUST be tagged with the letter "{{Wb}}x{{off}}" and the other field MUST be tagged with the letter "{{Wb}}y{{off}}" (case insensitive).
        - The fields MUST be either of the {{Ci}}gosql.Column{{off}} directive, or they MUST be named and of a normal non-directive type.
        {{Wb}}EXAMPLE:{{off}} type T struct {
		F int          {{raw "sql:\"x\""}}
		_ gosql.Column {{raw "sql:\"t.col_name,y\""}}
	}
{{ end }}

{{ define "` + errIllegalSliceUpdateModifier.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Illegal slice update modifier.{{off}}
    The {{R .TargetXxx}} query types with a "rel" field of {{Wb}}slice{{off}} type {{Wu}}DO NOT{{off}} support {{R .FieldDefinition}} {{.FieldKind}}s.
    {{Wb}}FIX:{{off}} remove the {{R .FieldDefinition}} {{.FieldKind}} from the {{Wb .TargetName}} query type.
{{ end }}

{{ define "` + errIllegalListPredicate.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Illegal list predicate.{{off}}
    The use of a {{Wi}}list_predicate{{off}} with a field of non-sequence type is illegal ({{R .TagError}} in "{{R .TagExpr}}" from the field {{R .FieldDefinition}} in {{W .TargetName}}).
    {{Wb}}FIX:{{off}} Change the {{Wi}}list_predicate{{off}} {{R .TagError}} to a non-sequence predicate, or change the {{R .FieldName}} field's type {{R .FieldType}} to a {{Wu}}sequence{{off}} type, i.e. {{Ci}}slice{{off}}, or {{Ci}}array{{off}}.
{{ end }}

{{ define "` + errIllegalUnaryPredicate.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Illegal unary predicate.{{off}}
    {{if .IsDirective -}}
    The use of a {{Wi}}unary_predicate{{off}} in a {{Wi}}binary_expression{{off}} is illegal ({{R .TagError}} in "{{R .TagExpr}}" from the directive {{R .FieldDefinition}} in {{W .TargetName}}).
    {{Wb}}FIX:{{off}} Change the {{Wi}}binary_expression{{off}} "{{R .TagExpr}}" to a valid {{Wi}}unary_expression{{off}}, or replace the {{Wi}}unary_predicate{{off}} {{R .TagError}} with a {{Wi}}binary_predicate{{off}}.
    {{else -}}
    The use of a {{Wi}}unary_predicate{{off}} is illegal in the tag of a plain field ({{R .TagError}} in "{{R .TagExpr}}" from the field {{R .FieldDefinition}} in {{W .TargetName}}).
    {{Wb}}FIX:{{off}} In the expression "{{R .TagExpr}}" remove the {{Wi}}unary_predicate{{off}} {{R .TagError}} or replace it with a {{Wi}}binary_predicate{{off}}.
    {{end -}}
{{ end }}

{{ define "` + errIllegalFieldQuantifier.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Illegal field quantifier.{{off}}
    The use of a {{Wi}}predicate_quantifier{{off}} with a {{Wu}}non-sequence{{off}} field is illegal ({{R .TagError}} in "{{R .TagExpr}}" from the {{.FieldKind}} {{R .FieldDefinition}} in {{W .TargetName}}).
    {{Wb}}FIX:{{off}} Remove the {{Wi}}predicate_quantifier{{off}} {{R .TagError}} from the expression "{{R .TagExpr}}", or change the {{R .FieldName}} field's type {{R .FieldType}} to a {{Wu}}sequence{{off}} type, i.e. {{Ci}}slice{{off}}, or {{Ci}}array{{off}}.
{{ end }}

{{ define "` + errIllegalPredicateQuantifier.name() + `" -}}
{{Wb .FileName}}:{{Wb .FileLine}}: {{Y}}Illegal predicate quantifier.{{off}}
    {{if .TagExprIsUnary -}}
    The use of a {{Wi}}predicate_quantifier{{off}} in a directive's {{Wi}}unary_expression{{off}} is illegal ({{R .TagError}} in "{{R .TagExpr}}" from the {{.FieldKind}} {{R .FieldDefinition}} in {{W .TargetName}}).
    {{else -}}
    The use of a {{Wi}}predicate_quantifier{{off}} with a {{Wu}}non-quantifiable{{off}} predicate is illegal ({{R .TagError}} in "{{R .TagExpr}}" from the {{.FieldKind}} {{R .FieldDefinition}} in {{W .TargetName}}).
    {{end -}}
    {{if (or .TagExprHasRHS (not .IsDirective)) -}}
    {{Wb}}FIX:{{off}} Remove the {{Wi}}predicate_quantifier{{off}} {{R .TagError}} from the expression "{{R .TagExpr}}", or change the {{Wu}}non-quantifiable{{off}} predicate {{R .TagExprPredicate}} to a {{Wu}}quantifiable{{off}} one.
    {{else -}}
    {{Wb}}FIX:{{off}} Remove the {{Wi}}predicate_quantifier{{off}} {{R .TagError}} from the expression "{{R .TagExpr}}".
    {{end -}}
{{ end }}
` // `

var error_templates = template.Must(template.New("t").Funcs(template.FuncMap{
	// white color (terminal)
	"w":  func(v ...interface{}) string { return getcolor("\033[0;37m", v) },
	"wb": func(v ...interface{}) string { return getcolor("\033[1;37m", v) },
	"wi": func(v ...interface{}) string { return getcolor("\033[3;37m", v) },
	"wu": func(v ...interface{}) string { return getcolor("\033[4;37m", v) },
	// cyan color (terminal)
	"c":  func(v ...interface{}) string { return getcolor("\033[0;36m", v) },
	"cb": func(v ...interface{}) string { return getcolor("\033[1;36m", v) },
	"ci": func(v ...interface{}) string { return getcolor("\033[3;36m", v) },
	"cu": func(v ...interface{}) string { return getcolor("\033[4;36m", v) },

	/////////////////////////////////////////////////////////////////////////
	// High Intensity
	/////////////////////////////////////////////////////////////////////////

	// red color HI (terminal)
	"R":  func(v ...interface{}) string { return getcolor("\033[0;91m", v) },
	"Rb": func(v ...interface{}) string { return getcolor("\033[1;91m", v) },
	"Ri": func(v ...interface{}) string { return getcolor("\033[3;91m", v) },
	"Ru": func(v ...interface{}) string { return getcolor("\033[4;91m", v) },
	// green color HI (terminal)
	"G":  func(v ...interface{}) string { return getcolor("\033[0;92m", v) },
	"Gb": func(v ...interface{}) string { return getcolor("\033[1;92m", v) },
	"Gi": func(v ...interface{}) string { return getcolor("\033[3;92m", v) },
	"Gu": func(v ...interface{}) string { return getcolor("\033[4;92m", v) },
	// yellow color HI (terminal)
	"Y":  func(v ...interface{}) string { return getcolor("\033[0;93m", v) },
	"Yb": func(v ...interface{}) string { return getcolor("\033[1;93m", v) },
	"Yi": func(v ...interface{}) string { return getcolor("\033[3;93m", v) },
	"Yu": func(v ...interface{}) string { return getcolor("\033[4;93m", v) },
	// blue color HI (terminal)
	"B":  func(v ...interface{}) string { return getcolor("\033[0;94m", v) },
	"Bb": func(v ...interface{}) string { return getcolor("\033[1;94m", v) },
	"Bi": func(v ...interface{}) string { return getcolor("\033[3;94m", v) },
	"Bu": func(v ...interface{}) string { return getcolor("\033[4;94m", v) },
	// cyan color HI (terminal)
	"C":  func(v ...interface{}) string { return getcolor("\033[0;96m", v) },
	"Cb": func(v ...interface{}) string { return getcolor("\033[1;96m", v) },
	"Ci": func(v ...interface{}) string { return getcolor("\033[3;96m", v) },
	"Cu": func(v ...interface{}) string { return getcolor("\033[4;96m", v) },
	// white color HI (terminal)
	"W":  func(v ...interface{}) string { return getcolor("\033[0;97m", v) },
	"Wb": func(v ...interface{}) string { return getcolor("\033[1;97m", v) },
	"Wi": func(v ...interface{}) string { return getcolor("\033[3;97m", v) },
	"Wu": func(v ...interface{}) string { return getcolor("\033[4;97m", v) },

	// no color (terminal)
	"off": func() string { return "\033[0m" },

	"raw": func(s string) string { return "`" + s + "`" },
	"Up":  strings.ToUpper,
}).Parse(error_template_string))

func getcolor(c string, v []interface{}) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v\033[0m", c, v[0])
	}
	return c
}
