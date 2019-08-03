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

// ImplementsScanner reports whether or not the given named type implements
// the "database/sql.Scanner" interface.
func ImplementsScanner(named *types.Named) bool {
	var sig *types.Signature
	for i := 0; i < named.NumMethods(); i++ {
		if m := named.Method(i); m.Name() == "Scan" {
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

// ImplementsValuer reports whether or not the given named type implements
// the "database/sql/driver.Valuer" interface.
func ImplementsValuer(named *types.Named) bool {
	var sig *types.Signature
	for i := 0; i < named.NumMethods(); i++ {
		if m := named.Method(i); m.Name() == "Value" {
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
