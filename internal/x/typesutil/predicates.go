package typesutil

import (
	"go/types"
	"strings"
)

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