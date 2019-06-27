package typesutil

import (
	"go/types"
)

// Reports whether or not the given variable is of type error.
func IsError(v *types.Var) bool {
	named, ok := v.Type().(*types.Named)
	if !ok {
		return false
	}
	pkg := named.Obj().Pkg()
	name := named.Obj().Name()
	return pkg == nil && name == "error"
}
