package analysis

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/frk/tagutil"
)

// analysis error
type anError struct {
	Code          errorCode
	PkgPath       string
	TargetName    string
	BlockName     string
	RelField      string
	RelType       RelType
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

func (e *anError) Error() string {
	sb := new(strings.Builder)
	if err := error_templates.ExecuteTemplate(sb, e.Code.name(), e); err != nil {
		panic(err)
	}
	return sb.String()
}

func (e *anError) FileAndLine() string {
	return e.FileName + ":" + strconv.Itoa(e.FileLine)
}

func (e *anError) IsDirective() bool {
	return e.FieldName == "_"
}

func (e *anError) IsSelectQueryKind() bool {
	return strings.HasPrefix(tolower(e.TargetName), "select")
}

func (e *anError) TagExprIsUnary() bool {
	_, _, _, rhs := parsePredicateExpr(e.TagExpr)
	return len(rhs) == 0 && e.IsDirective()
}

func (e *anError) TagExprHasRHS() bool {
	_, _, _, rhs := parsePredicateExpr(e.TagExpr)
	return len(rhs) > 0
}

func (e *anError) TagExprPredicate() string {
	_, op, _, _ := parsePredicateExpr(e.TagExpr)
	return op
}

func (e *anError) IsSequence() bool {
	return e.FieldTypeKind == "slice" || e.FieldTypeKind == "array"
}

func (e *anError) TagValueRel() string {
	return tagutil.New(e.TagString).First("rel")
}

func (e *anError) TagValueSql() string {
	return tagutil.New(e.TagString).Get("sql")
}

func (e *anError) TagValueSqlFirst() string {
	return tagutil.New(e.TagString).First("sql")
}

func (e *anError) TagValueSqlSecond() string {
	return tagutil.New(e.TagString).Second("sql")
}

func (e *anError) RelDefinition() string {
	return fmt.Sprintf("%s %s", e.RelField, e.RelTypeShort())
}

// Reformats the field type to make it more readable.
func (e *anError) RelTypeShort() string {
	// anon struct?
	if e.RelType.Base.Name == "" && e.RelType.Base.Kind == TypeKindStruct {
		return "struct{ ... }"
	}

	if e.RelType.Base.PkgPath == e.PkgPath {
		return e.RelType.Base.Name
	}
	return e.RelType.Base.PkgName + "." + e.RelType.Base.Name

	if i := strings.LastIndexByte(e.FieldType, '/'); i > -1 {
		return e.FieldType[i+1:]
	}
	return e.FieldType
}

func (e *anError) FieldDefinition() string {
	if len(e.TagString) > 0 {
		return fmt.Sprintf("%s %s `%s`", e.FieldName, e.FieldTypeShort(), e.TagString)
	}
	return fmt.Sprintf("%s %s", e.FieldName, e.FieldTypeShort())
}

// Reformats the field type to make it more readable.
func (e *anError) FieldTypeShort() string {
	// non-empty anon struct?
	if strings.HasPrefix(e.FieldType, "struct{") && e.FieldType != "struct{}" && e.FieldType != "struct{ }" {
		return "struct{ ... }"
	}

	if e.FieldTypePkgPath() == e.PkgPath {
		return e.FieldTypeName()
	}

	if i := strings.LastIndexByte(e.FieldType, '/'); i > -1 {
		return e.FieldType[i+1:]
	}
	return e.FieldType
}

func (e *anError) FieldTypeName() string {
	if i := strings.LastIndexByte(e.FieldType, '.'); i > -1 {
		return e.FieldType[i+1:]
	}
	return e.FieldType
}

func (e *anError) FieldTypePkgPath() string {
	if i := strings.LastIndexByte(e.FieldType, '.'); i > -1 {
		return e.FieldType[:i]
	}
	return ""
}

func (e *anError) FieldKind() (out string) {
	if e.FieldName == "_" {
		return "directive"
	}
	return "field"
}

func (e *anError) TargetXxx() (out string) {
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

func (e *anError) TargetKind() (out string) {
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
	errIllegalIteratorField
	errConflictingRelTag
	errConflictingRelName
	errConflictingRelAlias
	errConflictingWhere
	errConflictingOnConfictTarget
	errConflictingOnConfictAction
	errConflictingResultTarget
	errConflictingRelationDirective
	errConflictingFieldOrDirective
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
	errBadJoinConditionLHS
	errIllegalSliceUpdateModifier
	errIllegalListPredicate
	errIllegalUnaryPredicate
	errIllegalFieldQuantifier
	errIllegalPredicateQuantifier
	errUnknownColumnQualifier
	errColumnFieldUnknown
)

var error_template_string = `
{{ define "` + errBadFieldTypeInt.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad field type."}}
    Cannot use {{R .FieldTypeShort}} as the type of the {{Wb .FieldName}} field in {{Wb .TargetName}}.
    {{Wb "FIX:"}} change the field type to one of: {{Ci "int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64"}}.
{{ end }}

{{ define "` + errBadFieldTypeStruct.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad field type."}}
    Cannot use {{R .FieldTypeShort}} (kind {{R .FieldTypeKind}}) as the type of the {{Wb .FieldName}} field in {{Wb .TargetName}}.
    {{Wb "FIX:"}} change the field type to be of kind {{Ci "struct"}}.
{{ end }}

{{ define "` + errBadIterTypeInterface.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad iterator interface type."}}
    The type {{R .FieldTypeShort}} of the {{.TargetName}}.{{Wb .FieldName}} field {{Wu "IS NOT"}} a valid "{{W "iterator interface"}}" type.
    {{Wb "HINT:"}} a valid "{{W "iterator interface"}}" type {{Wu "MUST"}} satisfy all of the following requirements:
        - The interface MUST have exactly {{Wu "one"}} method.
        - The interface method MUST be {{Wu "accessible"}} from the generated code, i.e. the method must either be exported or, if unexported, the interface type itself must be anonymous or declared in the same package as the target query type.
        - The interface method's signature MUST have exactly {{Wu "one parameter"}} type and exactly {{Wu "one result"}} type.
        - The interface method's parameter MUST be of a {{Wu "named struct"}} type, or a pointer to a {{Wu "named struct"}} type.
        - The interface method's result type MUST be the standard {{Ci "error"}} type.
{{ end }}

{{ define "` + errBadIterTypeFunc.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad iterator func type."}}
    The type {{R .FieldTypeShort}} of the {{.TargetName}}.{{Wb .FieldName}} field {{Wu "IS NOT"}} a valid "{{W "iterator func"}}" type.
    {{Wb "HINT:"}} a valid "{{W "iterator func"}}" type {{Wu "MUST"}} satisfy all of the following requirements:
        - The function's signature MUST have exactly {{Wu "one parameter"}} type and exactly {{Wu "one result"}} type.
        - The function's parameter MUST be of a {{Wu "named struct"}} type, or a pointer to a {{Wu "named struct"}} type.
        - The function's result type MUST be the standard {{Ci "error"}} type.
{{ end }}

{{ define "` + errBadRelType.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad \"rel\" type."}}
    The type {{R .FieldTypeShort}} of the {{.TargetName}}.{{Wb .FieldName}} field {{Wu "IS NOT"}} a valid "{{W "iterator func"}}" type.
    {{Wb "HINT:"}} a valid "{{W "rel"}}" type {{Wu "MUST"}} be one of the following:
        - A {{Wu "named or unnamed struct"}} type.
        - A {{Wu "pointer to a named struct"}} type.
        - A {{Wu "slice of named structs"}} type.
        - A {{Wu "slice of pointers to named structs"}} type.
        - An {{Wu "array of named structs"}} type.
        - An {{Wu "array of pointers to named structs"}} type.
        - A valid {{Wu "iterator interface"}} type.
        - A valid {{Wu "iterator func"}} type.
{{ end }}

{{ define "` + errIllegalQueryField.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Illegal " .FieldKind "."}}
    The {{Wb .TargetXxx}} {{.TargetKind}} types {{Wu "DO NOT"}} support the {{R .FieldDefinition}} {{.FieldKind}}.
    {{Wb "FIX:"}} Remove the {{R .FieldDefinition}} {{.FieldKind}} from the {{Wb .TargetName}} {{.TargetKind}} type.
{{ end }}

{{ define "` + errIllegalStructDirective.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Illegal directive field."}}
    The {{Wb .BlockName}} struct {{Wu "DOES NOT"}} support the {{R .FieldName}} {{R .FieldTypeShort}} directive.
    {{Wb "FIX:"}} Remove the {{R .FieldName}} {{R .FieldTypeShort}} directive from the {{Wb .BlockName}} struct of the {{W .TargetName}} {{.TargetKind}} type.
{{ end }}

{{ define "` + errIllegalIteratorField.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Illegal iterator type field."}}
    The {{Wb .TargetXxx}} struct types {{Wu "DO NOT"}} allow for the "rel" field's type to be of the "{{Wi "iterator"}}" kind.
    {{Wb "FIX:"}} Change the type of the "{{R .FieldName}} {{R .FieldTypeShort}}" field in {{Wb .TargetName}} to a "{{Wi "non-iterator"}}" kind.
{{ end }}

{{ define "` + errConflictingRelTag.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Conflicting \"rel\" tag."}}
    The field {{R .FieldDefinition}} in {{Wb .TargetName}} is in conflict with another field that also has a "{{Wb "rel"}}" tag.
    {{Wb "FIX:"}} Make sure that the {{Wb .TargetName}} {{.TargetKind}} type has {{Wu "only one"}} field with the "{{Wb "rel"}}" tag. 
{{ end }}

{{ define "` + errConflictingRelAlias.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Conflicting relation alias."}}
    The relation alias "{{R .TagError}}" is used more than once in {{Wb .TargetName}} causing a name collision.
    {{Wb "HINT:"}} A {{Wi "relation_alias"}} MUST be unique within a query.
{{ end }}

{{ define "` + errConflictingRelName.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Conflicting relation name."}}
    The relation name "{{R .TagError}}" is used more than once in {{Wb .TargetName}} causing a name collision.
    {{Wb "HINT:"}} A {{Wi "relation_name"}} MUST be unique within a query or it MUST be hidden behind a unique {{Wi "relation_alias"}}.
{{ end }}

{{ define "` + errConflictingWhere.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Conflicting \"where\" fields / directives."}}
    {{if .IsDirective -}}
    The {{R .FieldTypeShort}} directive in {{Wb .TargetName}} is in conflict with another field or directive.
    {{else -}}
    The field {{R .FieldName}} (type {{R .FieldTypeShort}}) in {{Wb .TargetName}} is in conflict with another field or directive.
    {{end -}}
    {{if .IsSelectQueryKind -}}
    {{Wb "HINT:"}} The {{Wb .TargetXxx}} query types can have {{Wu "only one"}} WHERE producing field which MUST be one of the following:
        - A field of type {{Ci "gosql.Filter"}}.
        - A struct field named {{Wu "Where"}} (case insensitive).
    {{else -}}
    {{Wb "HINT:"}} The {{Wb .TargetXxx}} query types can have {{Wu "only one"}} WHERE producing field or directive which MUST be one of the following:
        - A directive field of type {{Ci "gosql.All"}}.
        - A field of type {{Ci "gosql.Filter"}}.
        - A struct field named {{Wu "Where"}} (case insensitive).
    {{end -}}
{{ end }}

{{ define "` + errConflictingOnConfictTarget.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Conflicting \"on conflict\" target directives."}}
    The {{R .FieldTypeShort}} directive is in conflict with another "target" directive of the {{Wb .BlockName}} struct in the {{Wb .TargetName}} query type.
    {{Wb "HINT:"}} The {{Wb .BlockName}} struct can have {{Wu "only one"}} of the following "target" directives:
        - The {{Ci "gosql.Column"}} directive.
        - The {{Ci "gosql.Index"}} directive.
        - The {{Ci "gosql.Constraint"}} directive.
{{ end }}

{{ define "` + errConflictingOnConfictAction.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Conflicting \"on conflict\" action directives."}}
    The {{R .FieldTypeShort}} directive is in conflict with another "action" directive of the {{Wb .BlockName}} struct in the {{Wb .TargetName}} query type.
    {{Wb "HINT:"}} The {{Wb .BlockName}} struct can have {{Wu "only one"}} of the following "action" directives:
        - The {{Ci "gosql.Ignore"}} directive.
        - The {{Ci "gosql.Update"}} directive.
{{ end }}

{{ define "` + errConflictingResultTarget.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Conflicting \"result\" fields / directives."}}
    {{if .IsDirective -}}
    The {{R .FieldTypeShort}} directive in {{Wb .TargetName}} is in conflict with another "result" field or directive.
    {{else -}}
    The field {{R .FieldName}} (type {{R .FieldTypeShort}}) in {{Wb .TargetName}} is in conflict with another "result" field or directive.
    {{end -}}
    {{Wb "HINT:"}} The {{Wb "InsertXxx"}}, {{Wb "UpdateXxx"}}, and {{Wb "DeleteXxx"}} query struct types can have {{Wu "only one"}} of the following:
        - The {{Ci "gosql.Return"}} directive.
	- A field named {{Wu "RowsAffected"}} (case insensitive) which MUST be of one of the following types: {{Ci "int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64"}}.
        - A field named {{Wu "Result"}} (case insensitive) whose type MUST be one of the following:
            - A {{Wu "named or unnamed struct"}} type.
            - A {{Wu "pointer to a named struct"}} type.
            - A {{Wu "slice of named structs"}} type.
            - A {{Wu "slice of pointers to named structs"}} type.
            - An {{Wu "array of named structs"}} type.
            - An {{Wu "array of pointers to named structs"}} type.
            - A valid {{Wu "iterator interface"}} type.
            - A valid {{Wu "iterator func"}} type.
{{ end }}

{{ define "` + errConflictingRelationDirective.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Conflicting \"gosql.Relation\" directive."}}
    The {{R .FieldTypeShort}} directive is in conflict with another directive of the same type.
    {{Wb "FIX:"}} Make sure that the {{Wb .BlockName}} struct field in {{Wb .TargetName}} has {{Wu "only one"}} {{R .FieldTypeShort}} directive. 
{{ end }}

{{ define "` + errConflictingFieldOrDirective.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Conflicting fields / directives."}}
    {{if .IsDirective -}}
    The {{R .FieldTypeShort}} directive is in conflict with another field or directive of its "kind".
    {{Wb "FIX:"}} Make sure that the {{Wb .TargetName}} struct type has, at most, {{Wu "only one"}} field / directive of the same "kind" as the {{R .FieldTypeShort}} directive.
    {{else -}}
    The field {{R .FieldName}} (type {{R .FieldTypeShort}}) in {{Wb .TargetName}} is in conflict with another field or directive of its "kind".
    {{Wb "FIX:"}} Make sure that the {{Wb .TargetName}} struct type has, at most, {{Wu "only one"}} field / directive of the same "kind" as the {{R .FieldName}} (type {{R .FieldTypeShort}}) field.
    {{end -}}
{{ end }}

{{ define "` + errMissingRelField.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Missing \"rel\" field."}}
    The {{R .TargetName}} struct type has no field with the "rel" tag.
    {{Wb "FIX:"}} Make sure that the {{R .TargetName}} struct type has {{Wu "exactly one"}} field with the "rel" tag. 
{{ end }}

{{ define "` + errMissingTagValue.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Missing \"sql\" tag value."}}
    The {{R .FieldName}} {{R .FieldTypeShort}} field / directive in {{Wb .TargetName}} is missing a value in its "sql" tag, or it is missing the tag completely.
    {{Wb "FIX:"}} Make sure that the {{R .FieldName}} {{R .FieldTypeShort}} field / directive in {{Wb .TargetName}} has the "sql" tag with a value that's valid for that field / directive. 
{{ end }}

{{ define "` + errMissingTagColumnList.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Missing \"sql\" tag with column list."}}
    The {{R .FieldName}} {{R .FieldTypeShort}} field / directive in {{Wb .TargetName}} is missing a list of columns in its "sql" tag, or it is missing the tag completely.
    {{Wb "FIX:"}} Make sure that the {{R .FieldName}} {{R .FieldTypeShort}} field / directive in {{Wb .TargetName}} has the "sql" tag with a valid list of columns. 
{{ end }}

{{ define "` + errMissingOnConflictTarget.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Missing on_conflict target directive."}}
    The {{R .FieldName}} struct field in {{Wb .TargetName}} is missing a "target" directive.
    {{Wb "FIX:"}} Make sure that {{R .FieldName}} in {{Wb .TargetName}} has {{Wb "exactly one"}} of the following "target" directives: 
        - The {{Ci "gosql.Column"}} directive.
        - The {{Ci "gosql.Index"}} directive.
        - The {{Ci "gosql.Constraint"}} directive.
{{ end }}

{{ define "` + errBadIdentTagValue.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad identifier value in tag."}}
    The "sql" tag value {{R .TagValueSqlFirst}} in {{R .FieldDefinition}} from {{Wb .TargetName}} is an invalid {{Wi "identifier"}}.
    {{Wb "HINT:"}} A valid {{Wi "identifier"}} MUST begin with a letter (a-z [case insensitive]) or an underscore (_), ` +
	`subsequent characters in the {{Wi "identifier"}} can be letters, underscores, and digits (0-9). Put another way, ` +
	`a valid {{Wi "identifier"}} MUST match the following regular expression: {{G ` + "`" + rxIdent.String() + "`" + `}}
{{ end }}

{{ define "` + errBadColIdTagValue.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad column_identifier value in tag."}}
    The "sql" tag value {{R .TagError}} in {{R .FieldDefinition}} from {{Wb .TargetName}} is an invalid {{Wi "column_identifier"}}.
    {{Wb "HINT:"}} A valid {{Wi "column_identifier"}} MUST be of the format {{Wi "[rel_alias.]column_name"}} where ` +
	`both the {{Wi "rel_alias"}} and the {{Wi "column_name"}} are valid {{Wi "identifiers"}}. A valid ` +
	`{{Wi "identifier"}} MUST begin with a letter (a-z [case insensitive]) or an underscore (_) and any subsequent ` +
	`characters in the {{Wi "identifier"}} can be letters, underscores, and digits (0-9). Put another way, a valid ` +
	`{{Wi "column_identifier"}} MUST match the following regular expression: {{G ` + "`" + rxColIdent.String() + "`" + `}}
{{ end }}

{{ define "` + errBadRelIdTagValue.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad relation_identifier value in tag."}}
    The "sql" tag value {{R .TagError}} in {{R .FieldDefinition}} from {{Wb .TargetName}} is an invalid {{Wi "relation_identifier"}}.
    {{Wb "HINT:"}} A valid {{Wi "relation_identifier"}} MUST be of the format {{Wi "[schema_name.]relation_name[:alias_name]"}} ` +
	`where {{Wi "schema_name"}}, {{Wi "relation_name"}}, and {{Wi "alias_name"}} are all valid {{Wi "identifiers"}}. ` +
	`A valid {{Wi "identifier"}} MUST begin with a letter (a-z [case insensitive]) or an underscore (_) and any subsequent ` +
	`characters in the {{Wi "identifier"}} can be letters, underscores, and digits (0-9). Put another way, a valid ` +
	`{{Wi "relation_identifier"}} MUST match the following regular expression: {{G ` + "`" + rxRelIdent.String() + "`" + `}}
{{ end }}

{{ define "` + errBadBoolTagValue.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad boolean value in tag."}}
    The "bool" tag value {{R .TagError}} in {{R .FieldDefinition}} from {{Wb .TargetName}} is an invalid {{Wi "boolean"}} value.
    {{Wb "HINT:"}} A valid {{Wi "boolean"}} value MUST be either {{G "AND"}} or {{G "OR"}} (case insensitive).
{{ end }}

{{ define "` + errBadUIntegerTagValue.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad unsigned_integer value in tag."}}
    The "sql" tag value {{R .TagError}} in {{R .FieldDefinition}} from {{Wb .TargetName}} is an invalid {{Wi "unsigned_integer"}} value.
    {{Wb "FIX:"}} Change {{R .TagError}} to a valid unsigned integer value.
{{ end }}

{{ define "` + errBadNullsOrderTagValue.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad nulls_order value in tag."}}
    The "sql" tag value {{R .TagError}} in {{R .FieldDefinition}} from {{Wb .TargetName}} is an invalid {{Wi "nulls_order"}} value.
    {{Wb "HINT:"}} A valid {{Wi "nulls_order"}} value MUST be either {{G "NULLSFIRST"}} or {{G "NULLSLAST"}} (case insensitive).
{{ end }}

{{ define "` + errBadOverrideTagValue.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad overriding_kind value in tag."}}
    The "sql" tag value {{R .TagError}} in {{R .FieldDefinition}} from {{Wb .TargetName}} is an invalid {{Wi "overriding_kind"}} value.
    {{Wb "HINT:"}} A valid {{Wi "overriding_kind"}} value MUST be either {{G "USER"}} or {{G "SYSTEM"}} (case insensitive).
{{ end }}

{{ define "` + errBadDirectiveBooleanExpr.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad boolean_expression in directive."}}
    The expression "{{R .TagError}}" in {{R .FieldDefinition}} from {{W .TargetName}}{{if .BlockName}}.{{Wb .BlockName}}{{end}} ` +
	`is an invalid {{Wi "boolean_expression"}} for a directive.
    {{Wb "HINT:"}} A valid {{Wi "boolean_expression"}} in a directive MUST be a {{Wu "complete"}} binary or unary expression.
{{ end }}

{{ define "` + errBadBetweenPredicate.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad between_predicate type."}}
    The {{R .FieldName}} {{R .FieldTypeShort}} field from {{Wb .TargetName}} has an invalid {{Wi "between_predicate"}} type.
    {{Wb "FIX:"}} A valid {{Wi "between_predicate"}} type MUST be a {{Ci "struct"}} type that satisfies the following requirements:
        - The {{Ci "struct"}} type MUST have {{Wu "exactly two"}} fields, no more and no less.
        - One of the fields MUST be tagged with the letter "{{Wb "x"}}" and the other field MUST be tagged with the letter "{{Wb "y"}}" (case insensitive).
        - The fields MUST be either of the {{Ci "gosql.Column"}} directive, or they MUST be named and of a normal non-directive type.
        {{Wb "EXAMPLE:"}} type T struct {
		F int          {{raw "sql:\"x\""}}
		_ gosql.Column {{raw "sql:\"t.col_name,y\""}}
	}
{{ end }}

{{ define "` + errBadJoinConditionLHS.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Bad join_condition LHS."}}
    The LHS "{{R .TagError}}" of the "{{R .TagExpr}}" expression in {{R (raw .TagString)}} is not allowed. 
    {{Wb "HINT:"}} The {{Wb "LHS"}} of a {{Wi "join_condition"}} expression MUST refer to a column of the relation being joined.
{{ end }}

{{ define "` + errIllegalSliceUpdateModifier.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Illegal slice update modifier."}}
    The {{R .TargetXxx}} query types with a "rel" field of {{Wb "slice"}} type {{Wu "DO NOT"}} support {{R .FieldDefinition}} {{.FieldKind}}s.
    {{Wb "FIX:"}} remove the {{R .FieldDefinition}} {{.FieldKind}} from the {{Wb .TargetName}} query type.
{{ end }}

{{ define "` + errIllegalListPredicate.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Illegal list predicate."}}
    The use of a {{Wi "list_predicate"}} with a field of non-sequence type is illegal ({{R .TagError}} in "{{R .TagExpr}}" from the field {{R .FieldDefinition}} in {{W .TargetName}}).
    {{Wb "FIX:"}} Change the {{Wi "list_predicate"}} {{R .TagError}} to a non-sequence predicate, or change the {{R .FieldName}} field's type {{R .FieldType}} to a {{Wu "sequence"}} type, i.e. {{Ci "slice"}}, or {{Ci "array"}}.
{{ end }}

{{ define "` + errIllegalUnaryPredicate.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Illegal unary predicate."}}
    {{if .IsDirective -}}
    The use of a {{Wi "unary_predicate"}} in a {{Wi "binary_expression"}} is illegal ({{R .TagError}} in "{{R .TagExpr}}" from the directive {{R .FieldDefinition}} in {{W .TargetName}}).
    {{Wb "FIX:"}} Change the {{Wi "binary_expression"}} "{{R .TagExpr}}" to a valid {{Wi "unary_expression"}}, or replace the {{Wi "unary_predicate"}} {{R .TagError}} with a {{Wi "binary_predicate"}}.
    {{else -}}
    The use of a {{Wi "unary_predicate"}} is illegal in the tag of a plain field ({{R .TagError}} in "{{R .TagExpr}}" from the field {{R .FieldDefinition}} in {{W .TargetName}}).
    {{Wb "FIX:"}} In the expression "{{R .TagExpr}}" remove the {{Wi "unary_predicate"}} {{R .TagError}} or replace it with a {{Wi "binary_predicate"}}.
    {{end -}}
{{ end }}

{{ define "` + errIllegalFieldQuantifier.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Illegal field quantifier."}}
    The use of a {{Wi "predicate_quantifier"}} with a {{Wu "non-sequence"}} field is illegal ({{R .TagError}} in "{{R .TagExpr}}" from the {{.FieldKind}} {{R .FieldDefinition}} in {{W .TargetName}}).
    {{Wb "FIX:"}} Remove the {{Wi "predicate_quantifier"}} {{R .TagError}} from the expression "{{R .TagExpr}}", or change the {{R .FieldName}} field's type {{R .FieldType}} to a {{Wu "sequence"}} type, i.e. {{Ci "slice"}}, or {{Ci "array"}}.
{{ end }}

{{ define "` + errIllegalPredicateQuantifier.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Illegal predicate quantifier."}}
    {{if .TagExprIsUnary -}}
    The use of a {{Wi "predicate_quantifier"}} in a directive's {{Wi "unary_expression"}} is illegal ({{R .TagError}} in "{{R .TagExpr}}" from the {{.FieldKind}} {{R .FieldDefinition}} in {{W .TargetName}}).
    {{else -}}
    The use of a {{Wi "predicate_quantifier"}} with a {{Wu "non-quantifiable"}} predicate is illegal ({{R .TagError}} in "{{R .TagExpr}}" from the {{.FieldKind}} {{R .FieldDefinition}} in {{W .TargetName}}).
    {{end -}}
    {{if (or .TagExprHasRHS (not .IsDirective)) -}}
    {{Wb "FIX:"}} Remove the {{Wi "predicate_quantifier"}} {{R .TagError}} from the expression "{{R .TagExpr}}", or change the {{Wu "non-quantifiable"}} predicate {{R .TagExprPredicate}} to a {{Wu "quantifiable"}} one.
    {{else -}}
    {{Wb "FIX:"}} Remove the {{Wi "predicate_quantifier"}} {{R .TagError}} from the expression "{{R .TagExpr}}".
    {{end -}}
{{ end }}

{{ define "` + errUnknownColumnQualifier.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Unknown column qualifier."}}
    The column qualifier "{{R .TagError}}" in {{Wb .TargetName}} field {{R .FieldDefinition}} references an unknown, as of yet unspecified relation.
    {{Wb "HINT:"}} Make sure that "{{R .TagError}}" does not contain any typos, and that it is referenced {{Wu "after"}} being first specified in the same tag or another tag in the same query. 
{{ end }}

{{ define "` + errColumnFieldUnknown.name() + `" -}}
{{Wb .FileAndLine}}: {{Y "Column has no matching field."}}
    The column "{{R .TagError}}" referenced in "{{R .FieldDefinition}}" has no matching field in "{{R .RelDefinition}}" type.
    {{Wb "HINT:"}} Columns referenced by the {{W .FieldTypeShort}} directive MUST have a matching field in the "{{W "rel"}}" type. 
{{ end }}

` // `

var error_templates = template.Must(template.New("t").Funcs(template.FuncMap{
	// white color (terminal)
	"w":  func(v ...string) string { return getcolor("\033[0;37m", v) },
	"wb": func(v ...string) string { return getcolor("\033[1;37m", v) },
	"wi": func(v ...string) string { return getcolor("\033[3;37m", v) },
	"wu": func(v ...string) string { return getcolor("\033[4;37m", v) },
	// cyan color (terminal)
	"c":  func(v ...string) string { return getcolor("\033[0;36m", v) },
	"cb": func(v ...string) string { return getcolor("\033[1;36m", v) },
	"ci": func(v ...string) string { return getcolor("\033[3;36m", v) },
	"cu": func(v ...string) string { return getcolor("\033[4;36m", v) },

	/////////////////////////////////////////////////////////////////////////
	// High Intensity
	/////////////////////////////////////////////////////////////////////////

	// red color HI (terminal)
	"R":  func(v ...string) string { return getcolor("\033[0;91m", v) },
	"Rb": func(v ...string) string { return getcolor("\033[1;91m", v) },
	"Ri": func(v ...string) string { return getcolor("\033[3;91m", v) },
	"Ru": func(v ...string) string { return getcolor("\033[4;91m", v) },
	// green color HI (terminal)
	"G":  func(v ...string) string { return getcolor("\033[0;92m", v) },
	"Gb": func(v ...string) string { return getcolor("\033[1;92m", v) },
	"Gi": func(v ...string) string { return getcolor("\033[3;92m", v) },
	"Gu": func(v ...string) string { return getcolor("\033[4;92m", v) },
	// yellow color HI (terminal)
	"Y":  func(v ...string) string { return getcolor("\033[0;93m", v) },
	"Yb": func(v ...string) string { return getcolor("\033[1;93m", v) },
	"Yi": func(v ...string) string { return getcolor("\033[3;93m", v) },
	"Yu": func(v ...string) string { return getcolor("\033[4;93m", v) },
	// blue color HI (terminal)
	"B":  func(v ...string) string { return getcolor("\033[0;94m", v) },
	"Bb": func(v ...string) string { return getcolor("\033[1;94m", v) },
	"Bi": func(v ...string) string { return getcolor("\033[3;94m", v) },
	"Bu": func(v ...string) string { return getcolor("\033[4;94m", v) },
	// cyan color HI (terminal)
	"C":  func(v ...string) string { return getcolor("\033[0;96m", v) },
	"Cb": func(v ...string) string { return getcolor("\033[1;96m", v) },
	"Ci": func(v ...string) string { return getcolor("\033[3;96m", v) },
	"Cu": func(v ...string) string { return getcolor("\033[4;96m", v) },
	// white color HI (terminal)
	"W":  func(v ...string) string { return getcolor("\033[0;97m", v) },
	"Wb": func(v ...string) string { return getcolor("\033[1;97m", v) },
	"Wi": func(v ...string) string { return getcolor("\033[3;97m", v) },
	"Wu": func(v ...string) string { return getcolor("\033[4;97m", v) },

	// no color (terminal)
	"off": func() string { return "\033[0m" },

	"raw": func(s string) string { return "`" + s + "`" },
	"Up":  strings.ToUpper,
}).Parse(error_template_string))

func getcolor(c string, v []string) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v\033[0m", c, stringsStringer(v))
	}
	return c
}

type stringsStringer []string

func (s stringsStringer) String() string {
	return strings.Join([]string(s), "")
}
