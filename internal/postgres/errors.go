package postgres

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/frk/gosql/internal/analysis"
	"github.com/frk/tagutil"
)

type dbErrorCode uint

func (e dbErrorCode) name() string { return fmt.Sprintf("error_template_%d", e) }

const (
	_ dbErrorCode = iota
	errDatabaseOpen
	errDatabaseInit // TODO

	// catalog errors
	errCatalogTypeGet       // TODO
	errCatalogTypeScan      // TODO
	errCatalogOperatorGet   // TODO
	errCatalogOperatorScan  // TODO
	errCatalogCastGet       // TODO
	errCatalogCastScan      // TODO
	errCatalogProcedureGet  // TODO
	errCatalogProcedureScan // TODO

	// relation errors
	errRelationUnknown
	errRelationScan              // TODO
	errRelationColumnGet         // TODO
	errRelationColumnScan        // TODO
	errRelationColumnUnknownType // TODO
	errRelationConstraintGet     // TODO
	errRelationConstraintScan    // TODO
	errRelationIndexGet          // TODO
	errRelationIndexScan         // TODO

	// column errors
	errColumnUnknown
	errColumnDefaultUnset
	errColumnTextSearchType
	errColumnFieldTypeWrite
	errColumnFieldTypeRead
	errColumnFieldComparison
	errColumnQualifierUnknown
	errColumnComparison
	// predicate errors
	errPredicateOperandQuantifier
	errPredicateOperandArray
	errPredicateOperandBool
	errPredicateOperandNull
	errPredicateLiteralExpr
	// between errors
	errBetweenColumnComparison
	errBetweenFieldComparison
	// procedure errors
	errProcedureUnknown
	// on conflict errors
	errOnConflictIndexUnknown
	errOnConflictIndexNotUnique
	errOnConflictIndexColumnsUnknown
	errOnConflictIndexColumnsNotUnique
	errOnConflictConstraintUnknown
	errOnConflictConstraintNotUnique
)

type dbError struct {
	Code    dbErrorCode
	DB      dbInfo
	Target  targetInfo
	Field   fieldInfo
	WBField fieldInfo
	Rel     relInfo
	Col     colInfo
	RHSCol  colInfo
	RHSLit  exprInfo
	Pred    analysis.Predicate
	Quant   analysis.Quantifier
	Func    analysis.FuncName
	Err     error `cmp:"+"`
}

func (e *dbError) Error() string {
	e.Field.tpkg = e.Target.Pkg
	sb := new(strings.Builder)
	if err := error_templates.ExecuteTemplate(sb, e.Code.name(), *e); err != nil {
		panic(err)
	}
	return sb.String()
}

func (d dbError) FuncName() string {
	return string(d.Func)
}

func (d dbError) PredString() string {
	if d.Pred.IsUnary() {
		return d.Col.IdRef() + " " + d.PredQuant()
	}

	if len(d.RHSCol.Id.Name) > 0 {
		return d.Col.IdRef() + " " + d.PredQuant() + " " + d.RHSCol.IdRef()
	} else if len(d.RHSLit.Expr) > 0 {
		return d.Col.IdRef() + " " + d.PredQuant() + " " + d.RHSLit.Expr
	}
	return ""
}

func (d dbError) PredQuant() string {
	if d.Quant > 0 {
		return d.Pred.String() + d.Quant.String()
	}
	return d.Pred.String()
}

func (d dbError) WBFieldDefinition() string {
	if len(d.WBField.Tag) > 0 {
		return fmt.Sprintf("%s struct { ... %s ... } `%s`", d.WBField.Name, d.Field.Definition(), d.WBField.Tag)
	}
	return fmt.Sprintf("%s struct { ... %s ... }", d.WBField.Name, d.Field.Definition())
}

////////////////////////////////////////////////////////////////////////////////

type dbInfo struct {
	DSN        string
	Name       string
	User       string
	SearchPath string
}

type fileInfo struct {
	Name string
	Line int
}

func (f fileInfo) NameAndLine() string {
	return f.Name + ":" + strconv.Itoa(f.Line)
}

type targetInfo struct {
	Pkg  string
	Name string
	File fileInfo
}

type fieldInfo struct {
	// The name of the field.
	Name string
	// The field's type name, possibly qualified with package path.
	Type string
	Tag  string
	File fileInfo

	// target package
	tpkg string
}

func (f fieldInfo) SqlTagFirst() string {
	return tagutil.New(f.Tag).First("sql")
}

func (f fieldInfo) SqlTag() string {
	return tagutil.New(f.Tag).Get("sql")
}

func (f fieldInfo) Definition() string {
	if len(f.Tag) > 0 {
		return fmt.Sprintf("%s %s `%s`", f.Name, f.TypeShort(), f.Tag)
	}
	return fmt.Sprintf("%s %s", f.Name, f.TypeShort())
}

// Reformats the field type to make it more readable.
func (f fieldInfo) TypeShort() string {
	// non-empty anon struct?
	if strings.HasPrefix(f.Type, "struct{") && f.Type != "struct{}" && f.Type != "struct{ }" {
		return "struct{ ... }"
	}

	if f.TypePkgPath() == f.tpkg {
		return f.TypeName()
	}

	if i := strings.LastIndexByte(f.Type, '/'); i > -1 {
		return f.Type[i+1:]
	}
	return f.Type
}

func (f fieldInfo) TypeName() string {
	if i := strings.LastIndexByte(f.Type, '.'); i > -1 {
		return f.Type[i+1:]
	}
	return f.Type
}

func (f fieldInfo) TypePkgPath() string {
	if i := strings.LastIndexByte(f.Type, '.'); i > -1 {
		return f.Type[:i]
	}
	return ""
}

func (f fieldInfo) IsDefaultDirective() bool {
	return f.Name == "_" && f.Type == "github.com/frk/gosql.Default"
}

type colInfo struct {
	Id analysis.ColIdent
	*Column
}

func (c colInfo) Ref() string {
	return c.Relation.Name + "." + c.Name
}

func (c colInfo) IdRef() string {
	if len(c.Id.Qualifier) > 0 {
		return c.Id.Qualifier + "." + c.Id.Name
	}
	return c.Id.Name
}

type relInfo struct {
	Id analysis.RelIdent
	*Relation
}

func (r relInfo) Ref() string {
	if len(r.Id.Qualifier) > 0 {
		return r.Id.Qualifier + "." + r.Id.Name
	}
	return r.Id.Name
}

type exprInfo struct {
	Expr string
	Type *Type
}

type relIdent analysis.RelIdent

func (id relIdent) Ref() string {
	if len(id.Qualifier) > 0 {
		return id.Qualifier + "." + id.Name
	}
	return id.Name
}

type colIdent analysis.ColIdent

func (id colIdent) Ref() string {
	if len(id.Qualifier) > 0 {
		return id.Qualifier + "." + id.Name
	}
	return id.Name
}

var error_template_string = `
--------------------------------------------------------------------------------
Helper templates
--------------------------------------------------------------------------------
{{ define "external_error_issue" }}
    - original error message: "{{Wb .Error.Error}}"
    - please consider opening an issue if the error persists.
{{ end }}

{{ define "connection_info" }}
    - database: {{W .DB.Name}} (dsn "{{W .DB.DSN}}")
    - relation: {{W .Rel.Name}} (schema "{{W .Rel.Schema}}")
    - target type: {{W .Target.Name}} (source "{{W .Target.File.NameAndLine}}")
{{ end }}

--------------------------------------------------------------------------------
Database error templates
--------------------------------------------------------------------------------

{{ define "` + errDatabaseOpen.name() + `" -}}
ERROR: {{Y "Failed to initialize database connection."}}
    An error occurred during the initialization of the database connection for the provided DSN "{{R .DB.DSN}}".
    - original error message: "{{Wb .Error.Error}}"
{{ end }}

{{ define "` + errDatabaseInit.name() + `" -}}
ERROR: {{Y "Failed to initialize database connection."}}
    An error occurred during the initialization of the database connection for the provided DSN "{{R .DB.DSN}}".
    {{- template "external_error_issue" . -}}
{{ end }}

--------------------------------------------------------------------------------
Catalog error templates
--------------------------------------------------------------------------------

{{ define "` + errCatalogTypeGet.name() + `" -}}
{{Y "Failed to retrieve pg_type information."}}
    An error occurred during the retrieval of the {{Wi "pg_type"}} information from the "{{R .DB.Name}}" database's catalog.
    {{- template "external_error_issue" . }}
{{ end }}

{{ define "` + errCatalogTypeScan.name() + `" -}}
{{Y "Failed to scan pg_type information."}}
    An error occurred during the scanning of the {{Wi "pg_type"}} information from the "{{R .DB.Name}}" database's catalog.
    {{- template "external_error_issue" . }}
{{ end }}

{{ define "` + errCatalogOperatorGet.name() + `" -}}
{{Y "Failed to retrieve pg_operator information."}}
    An error occurred during the retrieval of the {{Wi "pg_operator"}} information from the "{{R .DB.Name}}" database's catalog.
    {{- template "external_error_issue" . }}
{{ end }}

{{ define "` + errCatalogOperatorScan.name() + `" -}}
{{Y "Failed to scan pg_operator information."}}
    An error occurred during the scanning of the {{Wi "pg_operator"}} information from the "{{R .DB.Name}}" database's catalog.
    {{- template "external_error_issue" . }}
{{ end }}

{{ define "` + errCatalogCastGet.name() + `" -}}
{{Y "Failed to retrieve pg_cast information."}}
    An error occurred during the retrieval of the {{Wi "pg_cast"}} information from the "{{R .DB.Name}}" database's catalog.
    {{- template "external_error_issue" . }}
{{ end }}

{{ define "` + errCatalogCastScan.name() + `" -}}
{{Y "Failed to scan pg_cast information."}}
    An error occurred during the scanning of the {{Wi "pg_cast"}} information from the "{{R .DB.Name}}" database's catalog.
    {{- template "external_error_issue" . }}
{{ end }}

{{ define "` + errCatalogProcedureGet.name() + `" -}}
{{Y "Failed to retrieve pg_proc information."}}
    An error occurred during the retrieval of the {{Wi "pg_proc"}} information from the "{{R .DB.Name}}" database's catalog.
    {{- template "external_error_issue" . }}
{{ end }}

{{ define "` + errCatalogProcedureScan.name() + `" -}}
{{Y "Failed to scan pg_proc information."}}
    An error occurred during the scanning of the {{Wi "pg_proc"}} information from the "{{R .DB.Name}}" database's catalog.
    {{- template "external_error_issue" . }}
{{ end }}

--------------------------------------------------------------------------------
Relation error templates
--------------------------------------------------------------------------------

{{ define "` + errRelationUnknown.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Unknown relation."}}
    The relation "{{R .Rel.Ref}}" referenced in "{{R .Field.Definition}}" does not exist.
    - database: {{W .DB.Name}} (dsn "{{W .DB.DSN}}")
    {{if not .Rel.Id.Qualifier -}}
    - search path: {{Wb .DB.SearchPath}}
    {{end -}}
    - target type: {{W .Target.Name}} (source "{{W .Target.File.NameAndLine}}")
{{ end }}

{{ define "` + errRelationScan.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Failed to scan relation."}}
    An error occurred during the scanning of the "{{R .Rel.Ref}}" relation's information.
    {{- template "external_error_issue" . }}
{{ end }}

{{ define "` + errRelationColumnGet.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Failed to retrieve relation columns."}}
    An error occurred during the retrieval of the "{{R .Rel.Ref}}" relation's column information.
    {{- template "external_error_issue" . }}
{{ end }}

{{ define "` + errRelationColumnScan.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Failed to scan relation columns."}}
    An error occurred during the scanning of the "{{R .Rel.Ref}}" relation's column information.
    {{- template "external_error_issue" . }}
{{ end }}

{{ define "` + errRelationColumnUnknownType.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Unknown column type."}}
    {{W "github.com/frk/gosql"}} is unable to retrieve the type information for column "{{R .Col.Name}}" in relation "{{R .Rel.Ref}}".
    The column type's OID: {{R .Col.TypeOID}}
    Please consider opening an issue if the error persists.
{{ end }}

{{ define "` + errRelationConstraintGet.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Failed to retrieve relation constraints."}}
    An error occurred during the retrieval of the "{{R .Rel.Ref}}" relation's constraint information.
    {{- template "external_error_issue" . }}
{{ end }}

{{ define "` + errRelationConstraintScan.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Failed to scan relation constraints."}}
    An error occurred during the scanning of the "{{R .Rel.Ref}}" relation's constraint information.
    {{- template "external_error_issue" . }}
{{ end }}

{{ define "` + errRelationIndexGet.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Failed to retrieve relation indexes."}}
    An error occurred during the retrieval of the "{{R .Rel.Ref}}" relation's index information.
    {{- template "external_error_issue" . }}
{{ end }}

{{ define "` + errRelationIndexScan.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Failed to scan relation indexes."}}
    An error occurred during the scanning of the "{{R .Rel.Ref}}" relation's index information.
    {{- template "external_error_issue" . }}
{{ end }}

--------------------------------------------------------------------------------
Column error templates
--------------------------------------------------------------------------------

{{ define "` + errColumnUnknown.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Column not found."}}
    The column "{{R .Col.IdRef}}" referenced in "{{R .Field.Definition}}" does not exist in the relation "{{R .Rel.Name}}".
    {{- template "connection_info" . -}}
{{ end }}

{{ define "` + errColumnFieldTypeWrite.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Column not compatible with field."}}
    The "{{R .Col.Name}}" column's type "{{R .Col.Type.NameFmt}}" is not compatible for {{wu "writing"}}` +
	` with the "{{R .Field.Name}}" field's type "{{R .Field.Type}}".
{{ end }}

{{ define "` + errColumnFieldTypeRead.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Column not compatible with field."}}
    The "{{R .Col.Name}}" column's type "{{R .Col.Type.NameFmt}}" is not compatible for {{wu "reading"}}` +
	` with the "{{R .Field.Name}}" field's type "{{R .Field.Type}}".
{{ end }}

{{ define "` + errColumnTextSearchType.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Bad column type for text search."}}
    The column "{{R .Col.IdRef}}" referenced in "{{R .Field.Definition}}" is of type "{{R .Col.Type.NameFmt}}".
    - a column referenced by a {{W .Field.TypeShort}} directive MUST be of type {{Ci "tsvector"}} to support {{wu "full text search"}}.
{{ end }}

{{ define "` + errColumnDefaultUnset.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Column DEFAULT not set."}}
    The column "{{R .Col.IdRef}}" referenced in "{{R .Field.Definition}}" has no {{Wb "DEFAULT"}} constraint.
    {{if .Field.IsDefaultDirective -}}
    - a column referenced by a {{W .Field.TypeShort}} directive MUST have a {{Wb "DEFAULT"}} constraint.
    {{else -}}
    - a column referenced together with the "{{W "default"}}" option in the tag MUST have a {{Wb "DEFAULT"}} constraint.
    {{end -}}
{{ end }}

{{ define "` + errColumnQualifierUnknown.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Column qualifier relation not found."}}
    The qualifier "{{R .Col.Id.Qualifier}}" matches no relation.
    - please consider opening an issue if the error persists.
{{ end }}

--------------------------------------------------------------------------------
Pred error templates
--------------------------------------------------------------------------------

{{ define "` + errPredicateOperandQuantifier.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Bad operand type in quantified predicate."}}
    {{if .RHSCol.Id.Name -}}
    The predicate "{{R .PredString}}" used in "{{R .Field.Definition}}" is not valid.
    - the column "{{R .RHSCol.IdRef}}" on the right side is of type {{Wi .RHSCol.Type.NameFmt}}.
    - the quantifier {{R .Quant.String}} requires an {{Wi "array"}} on the right side.
    {{else -}}
    The predicate "{{R .PredString}}" used in "{{R .Field.Definition}}" is not valid.
    - the expression "{{R .RHSLit.Expr}}" on the right side is of type {{Wi .RHSLit.Type.NameFmt}}.
    - the quantifier {{R .Quant.String}} requires an {{Wi "array"}} on the right side.
    {{end -}}
{{ end }}

{{ define "` + errPredicateOperandArray.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Bad operand type in array predicate."}}
    {{if .RHSCol.Id.Name -}}
    The predicate "{{R .PredString}}" used in "{{R .Field.Definition}}" is not valid.
    - the column "{{R .RHSCol.IdRef}}" on the right side is of type {{Wi .RHSCol.Type.NameFmt}}.
    - the operator {{R .Pred.String}} requires an {{Wi "array"}} on the right side.
    {{else -}}
    The predicate "{{R .PredString}}" used in "{{R .Field.Definition}}" is not valid.
    - the expression "{{R .RHSLit.Expr}}" on the right side is of type {{Wi .RHSLit.Type.NameFmt}}.
    - the operator {{R .Pred.String}} requires an {{Wi "array"}} on the right side.
    {{end -}}
{{ end }}

{{ define "` + errPredicateOperandBool.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Bad operand type in boolean predicate."}}
    The predicate "{{R .PredString}}" used in "{{R .Field.Definition}}" is not valid.
    - the column "{{R .Col.IdRef}}" is of type {{Wi .Col.Type.NameFmt}}.
    - the operator {{R .Pred.String}} requires a {{Wi "boolean"}} operand.
{{ end }}

{{ define "` + errPredicateOperandNull.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Bad column operand in NULL predicate."}}
    The predicate "{{R .PredString}}" used in "{{R .Field.Definition}}" is not valid.
    - the column "{{R .Col.IdRef}}" has the {{Wi "NOT NULL"}} constraint.
    - the operator {{R .Pred.String}} requires a column without the {{Wi "NOT NULL"}} constraint.
{{ end }}

{{ define "` + errColumnComparison.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Bad column comparison."}}
    {{if .RHSCol.Id.Name -}}
    The operator "{{R .Col.Type.NameFmt .Pred.String .RHSCol.Type.NameFmt}}" referenced in "{{R .Field.Definition}}" does not exist.
    - column "{{R .Col.IdRef}}" is of type "{{R .Col.Type.NameFmt}}".
    - column "{{R .RHSCol.IdRef}}" is of type "{{R .RHSCol.Type.NameFmt}}".
    {{else if .RHSLit.Expr -}}
    The operator "{{R .Col.Type.NameFmt .Pred.String .RHSLit.Type.NameFmt}}" referenced in "{{R .Field.Definition}}" does not exist.
    - column "{{R .Col.IdRef}}" is of type "{{R .Col.Type.NameFmt}}".
    - expression "{{R .RHSLit.Expr}}" is of type "{{R .RHSLit.Type.NameFmt}}".
    {{end -}}
{{ end }}

{{ define "` + errColumnFieldComparison.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Bad column field comparison."}}
    The operator "{{R .Field.TypeShort .Pred.String .Col.Type.NameFmt}}" referenced in "{{R .Field.Definition}}" does not exist.
    - field "{{R .Field.Name}}" is of type "{{R .Field.TypeShort}}".
    - column "{{R .Col.IdRef}}" is of type "{{R .Col.Type.NameFmt}}".
{{ end }}

{{ define "` + errPredicateLiteralExpr.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Bad literal expression."}}
    The literal expression "{{R .RHSLit.Expr}}" used in "{{R .Field.Definition}}" is not valid.
{{ end }}

--------------------------------------------------------------------------------
Between error templates
--------------------------------------------------------------------------------

{{ define "` + errBetweenColumnComparison.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Bad column between comparison."}}
    The operator "{{R .Col.Type.NameFmt .Pred.String .RHSCol.Type.NameFmt}}" represented by "{{R .WBFieldDefinition}}" does not exist.
    - column "{{R .Col.IdRef}}" is of type "{{R .Col.Type.NameFmt}}".
    - column "{{R .RHSCol.IdRef}}" is of type "{{R .RHSCol.Type.NameFmt}}".
{{ end }}

{{ define "` + errBetweenFieldComparison.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Bad column field between comparison."}}
    The operator "{{R .Field.TypeShort .Pred.String .Col.Type.NameFmt}}" represented by "{{R .WBFieldDefinition}}" does not exist.
    - column "{{R .Col.IdRef}}" is of type "{{R .Col.Type.NameFmt}}".
    - field "{{R .Field.Name}}" is of type "{{R .Field.TypeShort}}".
{{ end }}

--------------------------------------------------------------------------------
Procedure error templates
--------------------------------------------------------------------------------

{{ define "` + errProcedureUnknown.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Procedure not found."}}
    No function with the name "{{R .FuncName}}" and argument type "{{R .Col.Type.NameFmt}}" as referenced in "{{R .Field.Definition}}"
    exists in the database "{{W .DB.Name}}" (search_path: {{Wb .DB.SearchPath}}).
{{ end }}

--------------------------------------------------------------------------------
On Conflict error templates
--------------------------------------------------------------------------------

{{ define "` + errOnConflictIndexUnknown.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Index not found."}}
    The index "{{R .Field.SqlTagFirst}}" referenced by "{{R .Field.Definition}}" does not exist.
    - directive "{{W .Field.TypeShort}}" requires an index present on the target relation "{{W .Rel.Id.Name}}".
{{ end }}

{{ define "` + errOnConflictIndexNotUnique.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Index not unique."}}
    The index "{{R .Field.SqlTagFirst}}" referenced by "{{R .Field.Definition}}" is not unique / primary key.
    - directive "{{W .Field.TypeShort}}" requires an index of type: "{{W "unique"}}", or "{{W "primary key"}}".
{{ end }}

{{ define "` + errOnConflictIndexColumnsUnknown.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Column index not found."}}
    The columns "{{R .Field.SqlTag}}" referenced by "{{R .Field.Definition}}" do not match an existing index.
    - directive "{{W .Field.TypeShort}}" requires the column(s) that match those of an index present on the target relation "{{W .Rel.Id.Name}}".
{{ end }}

{{ define "` + errOnConflictIndexColumnsNotUnique.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Column index not unique."}}
    The columns "{{R .Field.SqlTag}}" referenced by "{{R .Field.Definition}}" match no unique / primary key index.
    - directive "{{W .Field.TypeShort}}" requires the column(s) that match those of an index of type: "{{W "unique"}}", or "{{W "primary key"}}".
{{ end }}

{{ define "` + errOnConflictConstraintUnknown.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Constraint not found."}}
    The constraint "{{R .Field.SqlTagFirst}}" referenced by "{{R .Field.Definition}}" does not exist.
    - directive "{{W .Field.TypeShort}}" requires a constraint present on the target relation "{{W .Rel.Id.Name}}".
{{ end }}

{{ define "` + errOnConflictConstraintNotUnique.name() + `" -}}
{{Wb .Field.File.NameAndLine}}: {{Y "Constraint not unique."}}
    The constraint "{{R .Field.SqlTagFirst}}" referenced by "{{R .Field.Definition}}" is not unique / primary key.
    - directive "{{W .Field.TypeShort}}" requires a constraint of type: "{{W "unique"}}", or "{{W "primary key"}}".
{{ end }}

` // `

var error_templates = template.Must(template.New("t").Funcs(template.FuncMap{
	// red color (terminal)
	"r":  func(v ...string) string { return getcolor("\033[0;31m", v) },
	"rb": func(v ...string) string { return getcolor("\033[1;31m", v) },
	"ri": func(v ...string) string { return getcolor("\033[3;31m", v) },
	"ru": func(v ...string) string { return getcolor("\033[4;31m", v) },
	// green color (terminal)
	"g":  func(v ...string) string { return getcolor("\033[0;32m", v) },
	"gb": func(v ...string) string { return getcolor("\033[1;32m", v) },
	"gi": func(v ...string) string { return getcolor("\033[3;32m", v) },
	"gu": func(v ...string) string { return getcolor("\033[4;32m", v) },
	// yellow color (terminal)
	"y":  func(v ...string) string { return getcolor("\033[0;33m", v) },
	"yb": func(v ...string) string { return getcolor("\033[1;33m", v) },
	"yi": func(v ...string) string { return getcolor("\033[3;33m", v) },
	"yu": func(v ...string) string { return getcolor("\033[4;33m", v) },
	// blue color (terminal)
	"b":  func(v ...string) string { return getcolor("\033[0;34m", v) },
	"bb": func(v ...string) string { return getcolor("\033[1;34m", v) },
	"bi": func(v ...string) string { return getcolor("\033[3;34m", v) },
	"bu": func(v ...string) string { return getcolor("\033[4;34m", v) },
	// cyan color (terminal)
	"c":  func(v ...string) string { return getcolor("\033[0;36m", v) },
	"cb": func(v ...string) string { return getcolor("\033[1;36m", v) },
	"ci": func(v ...string) string { return getcolor("\033[3;36m", v) },
	"cu": func(v ...string) string { return getcolor("\033[4;36m", v) },
	// white color (terminal)
	"w":  func(v ...string) string { return getcolor("\033[0;37m", v) },
	"wb": func(v ...string) string { return getcolor("\033[1;37m", v) },
	"wi": func(v ...string) string { return getcolor("\033[3;37m", v) },
	"wu": func(v ...string) string { return getcolor("\033[4;37m", v) },

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

	"Up":  strings.ToUpper,
	"raw": func(s string) string { return "`" + s + "`" },
}).Parse(error_template_string))

func getcolor(c string, v []string) string {
	if len(v) > 0 {
		return fmt.Sprintf("%s%v\033[0m", c, stringsStringer(v))
	}
	return c
}

type stringsStringer []string

func (s stringsStringer) String() string {
	return strings.Join([]string(s), " ")
}
