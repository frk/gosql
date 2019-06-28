package typesutil

import (
	"go/types"
)

// IsError reports whether or not the given variable is of type "error".
func IsError(v *types.Var) bool {
	named, ok := v.Type().(*types.Named)
	if !ok {
		return false
	}
	pkg := named.Obj().Pkg()
	name := named.Obj().Name()
	return pkg == nil && name == "error"
}

// IsEmptyInterface reports whether or not the given variable is of type "interface{}".
func IsEmptyInterface(v *types.Var) bool {
	iface, ok := v.Type().(*types.Interface)
	if !ok {
		return false
	}
	return iface.NumMethods() == 0
}

// IsTime reports whether or not the given variable is of type "time.Time".
// IsTime returns true also for types that embed "time.Time" directly, this
// is to provide support for custom timestamp types.
func IsTime(v *types.Var) bool {
	named, ok := v.Type().(*types.Named)
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

// IsSqlDriverValue reports whether or not the given variable is
// of type "database/sql/driver.Value".
func IsSqlDriverValue(v *types.Var) bool {
	named, ok := v.Type().(*types.Named)
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

	if !IsEmptyInterface(sig.Params().At(0)) {
		return false
	}
	return IsError(sig.Results().At(0))
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

	if !IsSqlDriverValue(sig.Results().At(0)) {
		return false
	}
	return IsError(sig.Results().At(1))
}
