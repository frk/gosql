package typesutil

import (
	"go/types"
	"strings"
)

// IsString reports whether or not the given type is the builtin "string" type.
func IsString(typ types.Type) bool {
	basic, ok := typ.(*types.Basic)
	if !ok {
		return false
	}
	return basic.Kind() == types.String && basic.Name() == "string"
}

// IsStringMap reports whether or not the given type is the "map[string]string" type.
func IsStringMap(typ types.Type) bool {
	m, ok := typ.(*types.Map)
	if !ok {
		return false
	}
	if k, ok := m.Key().(*types.Basic); !ok || k.Kind() != types.String || k.Name() != "string" {
		return false
	}
	if e, ok := m.Elem().(*types.Basic); !ok || e.Kind() != types.String || e.Name() != "string" {
		return false
	}
	return true
}

// IsStringColumnMap reports whether or not the given type is the "map[string]filter.Column" type.
func IsStringColumnMap(typ types.Type) bool {
	m, ok := typ.(*types.Map)
	if !ok {
		return false
	}
	if k, ok := m.Key().(*types.Basic); !ok || k.Kind() != types.String || k.Name() != "string" {
		return false
	}
	if !IsFilterColumn(m.Elem()) {
		return false
	}
	return true
}

// IsNiladicFunc reports whether or not the given type is the "func()" type.
func IsNiladicFunc(typ types.Type) bool {
	sig, ok := typ.(*types.Signature)
	if !ok {
		return false
	}
	return sig.Params().Len() == 0 && sig.Results().Len() == 0
}

// IsError reports whether or not the given type is the "error" type.
func IsError(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	pkg := named.Obj().Pkg()
	name := named.Obj().Name()
	return pkg == nil && name == "error"
}

// IsEmptyInterface reports whether or not the given type is the "interface{}" type.
func IsEmptyInterface(typ types.Type) bool {
	iface, ok := typ.(*types.Interface)
	if !ok {
		return false
	}
	return iface.NumMethods() == 0
}

// IsEmptyInterfaceSlice reports whether or not the given type is the "[]interface{}" type.
func IsEmptyInterfaceSlice(typ types.Type) bool {
	if s, ok := typ.(*types.Slice); ok {
		return IsEmptyInterface(s.Elem())
	}
	return false
}

// IsContext reports whether or not the given type is the standard "context.Context" type.
func IsContext(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	typeName := named.Obj()
	if pkg := typeName.Pkg(); pkg != nil {
		path := pkg.Path()
		name := typeName.Name()
		return path == "context" && name == "Context"
	}
	return false
}

// IsTime reports whether or not the given type is the "time.Time" type.
// IsTime returns true also for types that embed "time.Time" directly, this
// is to provide support for custom timestamp types.
func IsTime(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	typeName := named.Obj()
	if pkg := typeName.Pkg(); pkg != nil {
		path := pkg.Path()
		name := typeName.Name()
		if path == "time" && name == "Time" {
			return true
		}
	}

	if S, ok := named.Underlying().(*types.Struct); ok {
		for i := 0; i < S.NumFields(); i++ {
			f := S.Field(i)
			named, ok := f.Type().(*types.Named)
			if ok && f.Anonymous() {
				typeName := named.Obj()
				if pkg := typeName.Pkg(); pkg != nil {
					path := pkg.Path()
					name := typeName.Name()
					if path == "time" && name == "Time" {
						return true
					}
				}
			}
		}
	}
	return false
}

// IsSqlResult reports whether or not the given type is the standard "sql.Result" type.
func IsSqlResult(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	typeName := named.Obj()
	if pkg := typeName.Pkg(); pkg != nil {
		path := pkg.Path()
		name := typeName.Name()
		return path == "database/sql" && name == "Result"
	}
	return false
}

// IsSqlRowsPtr reports whether or not the given type is the standard "*sql.Rows" type.
func IsSqlRowsPtr(typ types.Type) bool {
	ptr, ok := typ.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}

	typeName := named.Obj()
	if pkg := typeName.Pkg(); pkg != nil {
		path := pkg.Path()
		name := typeName.Name()
		return path == "database/sql" && name == "Rows"
	}
	return false
}

// IsSqlRowPtr reports whether or not the given type is the standard "*sql.Row" type.
func IsSqlRowPtr(typ types.Type) bool {
	ptr, ok := typ.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}

	typeName := named.Obj()
	if pkg := typeName.Pkg(); pkg != nil {
		path := pkg.Path()
		name := typeName.Name()
		return path == "database/sql" && name == "Row"
	}
	return false
}

// IsSqlDriverValue reports whether or not the given type is the "database/sql/driver.Value" type.
func IsSqlDriverValue(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	typeName := named.Obj()
	if pkg := typeName.Pkg(); pkg != nil {
		path := pkg.Path()
		name := typeName.Name()
		return path == "database/sql/driver" && name == "Value"
	}
	return false
}

// IsFilterColumn reports whether or not the given type is the "github.com/frk/gosql/filter.Column" type.
func IsFilterColumn(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	name := named.Obj().Name()
	path := named.Obj().Pkg().Path()

	if name != "Column" {
		return false
	}
	if !strings.HasSuffix(path, "github.com/frk/gosql/filter") {
		return false
	}

	return true
}

// IsDirective reports whether or not the given type is a "github.com/frk/gosql" directive type.
func IsDirective(ident string, typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}
	name := named.Obj().Name()
	if name != ident {
		return false
	}

	// Compare the suffix only to allow for vendor imports.
	path := named.Obj().Pkg().Path()
	if !strings.HasSuffix(path, "github.com/frk/gosql") {
		return false
	}

	st, ok := named.Underlying().(*types.Struct)
	return ok && st.NumFields() == 1 && st.Field(0).Name() == "_isdir"
}

// ImplementsAfterScanner reports whether or not the given named type
// implements the "gosql.AfterScanner" interface.
func ImplementsAfterScanner(named *types.Named) bool {
	var sig *types.Signature
	for i := 0; i < named.NumMethods(); i++ {
		if m := named.Method(i); m.Name() == "AfterScan" {
			sig = m.Type().(*types.Signature)
			break
		}
	}
	return sig != nil && sig.Params().Len() == 0 && sig.Results().Len() == 0
}

// ImplementsErrorHandler reports whether or not the given named type
// implements the "gosql.ErrorHandler" interface.
func ImplementsErrorHandler(named *types.Named) bool {
	var sig *types.Signature
	for i := 0; i < named.NumMethods(); i++ {
		if m := named.Method(i); m.Name() == "HandleError" {
			sig = m.Type().(*types.Signature)
			break
		}
	}
	if sig == nil || sig.Params().Len() != 1 || sig.Results().Len() != 1 {
		return false
	}

	return IsError(sig.Params().At(0).Type()) && IsError(sig.Results().At(0).Type())
}

// ImplementsErrorInfoHandler reports whether or not the given named type
// implements the "gosql.ErrorInfoHandler" interface.
func ImplementsErrorInfoHandler(named *types.Named) bool {
	var sig *types.Signature
	for i := 0; i < named.NumMethods(); i++ {
		if m := named.Method(i); m.Name() == "HandleErrorInfo" {
			sig = m.Type().(*types.Signature)
			break
		}
	}
	if sig == nil || sig.Params().Len() != 1 || sig.Results().Len() != 1 {
		return false
	}
	if !IsError(sig.Results().At(0).Type()) {
		return false
	}

	// check that the method's argument is of type *gosql.ErrorInfo
	argtyp := sig.Params().At(0).Type()
	argptr, ok := argtyp.(*types.Pointer)
	if !ok {
		return false
	}
	argnamed, ok := argptr.Elem().(*types.Named)
	if !ok {
		return false
	}
	name := argnamed.Obj().Name()
	path := argnamed.Obj().Pkg().Path()
	return strings.HasSuffix(path, "github.com/frk/gosql") && name == "ErrorInfo"
}

// ImplementsGosqlConn reports whether or not the given named type implements the "gosql.Conn" interface.
func ImplementsGosqlConn(named *types.Named) bool {
	mm := Methoder(named)
	if iface, ok := named.Underlying().(*types.Interface); ok {
		mm = iface
	}

	var hasExec, hasQuery, hasQueryRow bool
	var hasExecContext, hasQueryContext, hasQueryRowContext bool

	for i := 0; i < mm.NumMethods(); i++ {
		m := mm.Method(i)
		switch m.Name() {
		case "Exec": // Exec(query string, args ...interface{}) (sql.Result, error)
			sig := m.Type().(*types.Signature)
			p, r := sig.Params(), sig.Results()
			if p.Len() != 2 || !sig.Variadic() || r.Len() != 2 {
				return false
			}
			if !IsString(p.At(0).Type()) || !IsEmptyInterfaceSlice(p.At(1).Type()) {
				return false
			}
			if !IsSqlResult(r.At(0).Type()) || !IsError(r.At(1).Type()) {
				return false
			}
			hasExec = true
		case "Query": // Query(query string, args ...interface{}) (*sql.Rows, error)
			sig := m.Type().(*types.Signature)
			p, r := sig.Params(), sig.Results()
			if p.Len() != 2 || !sig.Variadic() || r.Len() != 2 {
				return false
			}
			if !IsString(p.At(0).Type()) || !IsEmptyInterfaceSlice(p.At(1).Type()) {
				return false
			}
			if !IsSqlRowsPtr(r.At(0).Type()) || !IsError(r.At(1).Type()) {
				return false
			}
			hasQuery = true
		case "QueryRow": // QueryRow(query string, args ...interface{}) *sql.Row
			sig := m.Type().(*types.Signature)
			p, r := sig.Params(), sig.Results()
			if p.Len() != 2 || !sig.Variadic() || r.Len() != 1 {
				return false
			}
			if !IsString(p.At(0).Type()) || !IsEmptyInterfaceSlice(p.At(1).Type()) {
				return false
			}
			if !IsSqlRowPtr(r.At(0).Type()) {
				return false
			}
			hasQueryRow = true
		case "ExecContext": // ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
			sig := m.Type().(*types.Signature)
			p, r := sig.Params(), sig.Results()
			if p.Len() != 3 || !sig.Variadic() || r.Len() != 2 {
				return false
			}
			if !IsContext(p.At(0).Type()) || !IsString(p.At(1).Type()) || !IsEmptyInterfaceSlice(p.At(2).Type()) {
				return false
			}
			if !IsSqlResult(r.At(0).Type()) || !IsError(r.At(1).Type()) {
				return false
			}
			hasExecContext = true
		case "QueryContext": // QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
			sig := m.Type().(*types.Signature)
			p, r := sig.Params(), sig.Results()
			if p.Len() != 3 || !sig.Variadic() || r.Len() != 2 {
				return false
			}
			if !IsContext(p.At(0).Type()) || !IsString(p.At(1).Type()) || !IsEmptyInterfaceSlice(p.At(2).Type()) {
				return false
			}
			if !IsSqlRowsPtr(r.At(0).Type()) || !IsError(r.At(1).Type()) {
				return false
			}
			hasQueryContext = true
		case "QueryRowContext": // QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
			sig := m.Type().(*types.Signature)
			p, r := sig.Params(), sig.Results()
			if p.Len() != 3 || !sig.Variadic() || r.Len() != 1 {
				return false
			}
			if !IsContext(p.At(0).Type()) || !IsString(p.At(1).Type()) || !IsEmptyInterfaceSlice(p.At(2).Type()) {
				return false
			}
			if !IsSqlRowPtr(r.At(0).Type()) {
				return false
			}
			hasQueryRowContext = true
		}
	}

	return hasExec && hasQuery && hasQueryRow &&
		hasExecContext && hasQueryContext && hasQueryRowContext
}

// ImplementsGosqlFilterConstructor reports whether or not the given named
// type implements the "gosql.FilterConstructor" interface.
func ImplementsGosqlFilterConstructor(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	mm := Methoder(named)
	if iface, ok := named.Underlying().(*types.Interface); ok {
		mm = iface
	}

	var hasInit, hasInitV2, hasCol, hasAnd, hasOr bool

	for i := 0; i < mm.NumMethods(); i++ {
		m := mm.Method(i)
		switch m.Name() {
		case "Init": // Init(colmap map[string]string, tscol string)
			sig := m.Type().(*types.Signature)
			p, r := sig.Params(), sig.Results()
			if p.Len() != 2 || r.Len() != 0 {
				return false
			}
			if !IsStringMap(p.At(0).Type()) || !IsString(p.At(1).Type()) {
				return false
			}
			hasInit = true
		case "InitV2": // InitV2(colmap map[string]filter.Column, tscol string)
			sig := m.Type().(*types.Signature)
			p, r := sig.Params(), sig.Results()
			if p.Len() != 2 || r.Len() != 0 {
				return false
			}
			if !IsStringColumnMap(p.At(0).Type()) || !IsString(p.At(1).Type()) {
				return false
			}
			hasInitV2 = true
		case "Col": // Col(column string, op string, value interface{})
			sig := m.Type().(*types.Signature)
			p, r := sig.Params(), sig.Results()
			if p.Len() != 3 || r.Len() != 0 {
				return false
			}
			if !IsString(p.At(0).Type()) || !IsString(p.At(1).Type()) || !IsEmptyInterface(p.At(2).Type()) {
				return false
			}
			hasCol = true
		case "And": // And(nest func())
			sig := m.Type().(*types.Signature)
			p, r := sig.Params(), sig.Results()
			if p.Len() != 1 || r.Len() != 0 {
				return false
			}
			if !IsNiladicFunc(p.At(0).Type()) {
				return false
			}
			hasAnd = true
			break
		case "Or": // Or(nest func())
			sig := m.Type().(*types.Signature)
			p, r := sig.Params(), sig.Results()
			if p.Len() != 1 || r.Len() != 0 {
				return false
			}
			if !IsNiladicFunc(p.At(0).Type()) {
				return false
			}
			hasOr = true
			break
		}
	}

	return hasInit && hasInitV2 && hasCol && hasAnd && hasOr
}

// ImplementsScanner reports whether or not the given named type implements the
// "database/sql.Scanner" interface. If the named type's underlying type is an
// interface type, ImplementsScanner will report whether or not that interface
// type declares, or embeds, the "Scan(interface{}) error" method.
func ImplementsScanner(named *types.Named) bool {
	if iface, ok := named.Underlying().(*types.Interface); ok {
		return IsScanner(iface)
	}
	return IsScanner(named)
}

// IsScanner reports whether or not the given Methoder type declares,
// or embeds, the "Scan(interface{}) error" method.
func IsScanner(mm Methoder) bool {
	var sig *types.Signature
	for i := 0; i < mm.NumMethods(); i++ {
		if m := mm.Method(i); m.Name() == "Scan" {
			sig = m.Type().(*types.Signature)
			break
		}
	}
	if sig == nil || sig.Params().Len() != 1 || sig.Results().Len() != 1 {
		return false
	}

	if !IsEmptyInterface(sig.Params().At(0).Type()) {
		return false
	}
	return IsError(sig.Results().At(0).Type())
}

// ImplementsValuer reports whether or not the given named type implements the
// "database/sql/driver.Valuer" interface. If the named type's underlying type is
// an interface type, ImplementsValuer will report whether or not that interface
// type declares, or embeds, the "Value() (driver.Value, error)" method.
func ImplementsValuer(named *types.Named) bool {
	if iface, ok := named.Underlying().(*types.Interface); ok {
		return IsValuer(iface)
	}
	return IsValuer(named)
}

// IsValuer reports whether or not the given Methoder type declares,
// or embeds, the "Value() (driver.Value, error)" method.
func IsValuer(mm Methoder) bool {
	var sig *types.Signature
	for i := 0; i < mm.NumMethods(); i++ {
		if m := mm.Method(i); m.Name() == "Value" {
			sig = m.Type().(*types.Signature)
			break
		}
	}
	if sig == nil || sig.Params().Len() > 0 || sig.Results().Len() != 2 {
		return false
	}

	if !IsSqlDriverValue(sig.Results().At(0).Type()) {
		return false
	}
	return IsError(sig.Results().At(1).Type())
}

// ImplementsJSONMarshaler reports whether or not the given named type implements
// the "encoding/json.Marshaler" interface. If the named type's underlying type
// is an interface type, ImplementsJSONMarshaler will report whether or not that
// interface type declares, or embeds, the "MarshalJSON() ([]byte, error)" method.
func ImplementsJSONMarshaler(named *types.Named) bool {
	if iface, ok := named.Underlying().(*types.Interface); ok {
		return IsJSONMarshaler(iface)
	}
	return IsJSONMarshaler(named)
}

// IsJSONMarshaler reports whether or not the given Methoder type declares,
// or embeds, the "MarshalJSON() ([]byte, error)" method.
func IsJSONMarshaler(mm Methoder) bool {
	var sig *types.Signature
	for i := 0; i < mm.NumMethods(); i++ {
		if m := mm.Method(i); m.Name() == "MarshalJSON" {
			sig = m.Type().(*types.Signature)
			break
		}
	}
	if sig == nil || sig.Params().Len() != 0 || sig.Results().Len() != 2 {
		return false
	}

	if !isbyteslice(sig.Results().At(0).Type()) {
		return false
	}
	return IsError(sig.Results().At(1).Type())
}

// ImplementsJSONUnmarshaler reports whether or not the given named type implements
// the "encoding/json.Unmarshaler" interface. If the named type's underlying type
// is an interface type, ImplementsJSONUnmarshaler will report whether or not that
// interface type declares, or embeds, the "UnmarshalJSON([]byte) (error)" method.
func ImplementsJSONUnmarshaler(named *types.Named) bool {
	if iface, ok := named.Underlying().(*types.Interface); ok {
		return IsJSONUnmarshaler(iface)
	}
	return IsJSONUnmarshaler(named)
}

// IsJSONUnmarshaler reports whether or not the given Methoder type declares,
// or embeds, the "UnmarshalJSON([]byte) (error)" method.
func IsJSONUnmarshaler(mm Methoder) bool {
	var sig *types.Signature
	for i := 0; i < mm.NumMethods(); i++ {
		if m := mm.Method(i); m.Name() == "UnmarshalJSON" {
			sig = m.Type().(*types.Signature)
			break
		}
	}
	if sig == nil || sig.Params().Len() != 1 || sig.Results().Len() != 1 {
		return false
	}

	if !isbyteslice(sig.Params().At(0).Type()) {
		return false
	}
	return IsError(sig.Results().At(0).Type())
}

// ImplementsXMLMarshaler reports whether or not the given named type implements the
// "encoding/xml.Marshaler" interface. If the named type's underlying type is an interface
// type, ImplementsXMLMarshaler will report whether or not that interface type declares,
// or embeds, the "MarshalXML(*xml.Encoder, xml.StartElement) (error)" method.
func ImplementsXMLMarshaler(named *types.Named) bool {
	if iface, ok := named.Underlying().(*types.Interface); ok {
		return IsXMLMarshaler(iface)
	}
	return IsXMLMarshaler(named)
}

// IsXMLMarshaler reports whether or not the given Methoder type declares,
// or embeds, the "MarshalXML(*xml.Encoder, xml.StartElement) (error)" method.
func IsXMLMarshaler(mm Methoder) bool {
	var sig *types.Signature
	for i := 0; i < mm.NumMethods(); i++ {
		if m := mm.Method(i); m.Name() == "MarshalXML" {
			sig = m.Type().(*types.Signature)
			break
		}
	}
	if sig == nil || sig.Params().Len() != 2 || sig.Results().Len() != 1 {
		return false
	}

	param := sig.Params()
	if !isxmlencoder(param.At(0).Type()) || !isxmlstartelem(param.At(1).Type()) {
		return false
	}
	return IsError(sig.Results().At(0).Type())
}

// ImplementsXMLUnmarshaler reports whether or not the given named type implements the
// "encoding/xml.Unmarshaler" interface. If the named type's underlying type is an interface
// type, ImplementsXMLUnmarshaler will report whether or not that interface type declares,
// or embeds, the "UnmarshalXML(*xml.Decoder, xml.StartElement) (error)" method.
func ImplementsXMLUnmarshaler(named *types.Named) bool {
	if iface, ok := named.Underlying().(*types.Interface); ok {
		return IsXMLUnmarshaler(iface)
	}
	return IsXMLUnmarshaler(named)
}

// IsXMLUnmarshaler reports whether or not the given Methoder type declares,
// or embeds, the "UnmarshalXML(*xml.Decoder, xml.StartElement) (error)" method.
func IsXMLUnmarshaler(mm Methoder) bool {
	var sig *types.Signature
	for i := 0; i < mm.NumMethods(); i++ {
		if m := mm.Method(i); m.Name() == "UnmarshalXML" {
			sig = m.Type().(*types.Signature)
			break
		}
	}
	if sig == nil || sig.Params().Len() != 2 || sig.Results().Len() != 1 {
		return false
	}

	param := sig.Params()
	if !isxmldecoder(param.At(0).Type()) || !isxmlstartelem(param.At(1).Type()) {
		return false
	}
	return IsError(sig.Results().At(0).Type())
}

// Methoder represents a type with methods. It is implicitly implemented
// by *types.Interface and *types.Named.
type Methoder interface {
	NumMethods() int
	Method(i int) *types.Func
}

// isbyteslice reports whether or not the given type is a []byte type.
func isbyteslice(t types.Type) bool {
	if s, ok := t.(*types.Slice); ok {
		if e, ok := s.Elem().(*types.Basic); ok && e.Kind() == types.Byte {
			return true
		}
	}
	return false
}

// isxmlencoder reports whether or not the given type is an *encoding/xml.Encoder type.
func isxmlencoder(typ types.Type) bool {
	ptr, ok := typ.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}

	name := named.Obj().Name()
	path := named.Obj().Pkg().Path()
	return name == "Encoder" && path == "encoding/xml"
}

// isxmldecoder reports whether or not the given type is a *encoding/xml.Decoder type.
func isxmldecoder(typ types.Type) bool {
	ptr, ok := typ.(*types.Pointer)
	if !ok {
		return false
	}
	named, ok := ptr.Elem().(*types.Named)
	if !ok {
		return false
	}

	name := named.Obj().Name()
	path := named.Obj().Pkg().Path()
	return name == "Decoder" && path == "encoding/xml"
}

// isxmlstartelem reports whether or not the given type is a encoding/xml.StartElement type.
func isxmlstartelem(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	name := named.Obj().Name()
	path := named.Obj().Pkg().Path()
	return name == "StartElement" && path == "encoding/xml"
}
