package typesutil

import (
	"errors"
	"go/types"
	"strings"
)

var ErrBadType = errors.New("bad type")

type NamedStruct struct {
	Named  *types.Named
	Struct *types.Struct
}

// GetStruct is a helper function that returns a *NamedStruct value that represents
// the struct type of the the given *types.Var. If the struct type is unnamed then
// the Named field of the *NamedStruct value will remain uninitialized. If the var's
// type is not a struct then GetStruct will return an error.
func GetStruct(v *types.Var) (*NamedStruct, error) {
	ns := new(NamedStruct)
	typ := v.Type()

	var ok bool
	if ns.Named, ok = typ.(*types.Named); ok {
		typ = ns.Named.Underlying()
	}

	if ns.Struct, ok = typ.(*types.Struct); !ok {
		return nil, ErrBadType
	}
	return ns, nil
}

// GetDirectiveName returns the name of the gosql directive type of the given
// variable. If the type of the given variable is not a directive an empty
// string will be returned instead.
func GetDirectiveName(v *types.Var) string {
	named, ok := v.Type().(*types.Named)
	if !ok {
		return ""
	}

	path := named.Obj().Pkg().Path()
	if !strings.HasSuffix(path, "github.com/frk/gosql") {
		return ""
	}

	st, ok := named.Underlying().(*types.Struct)
	if !ok || st.NumFields() != 1 {
		return ""
	} else if st.Field(0).Name() != "_isdir" {
		return ""
	}
	return named.Obj().Name()
}
