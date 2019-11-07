package testutil

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"log"

	"github.com/frk/gosql/internal/typesutil"
)

type Testdata struct {
	Files []*ast.File
	Defs  map[*ast.Ident]types.Object
}

// ParseTestdata parses go files in the specified directory. Regardless of the
// package path, the package name of those files should be "testdata".
func ParseTestdata(dir string) Testdata {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		log.Fatal(err)
	}

	var files []*ast.File
	for _, f := range pkgs["testdata"].Files {
		files = append(files, f)
	}

	conf := types.Config{Importer: typesutil.NewImporter()}
	info := types.Info{Defs: make(map[*ast.Ident]types.Object)}
	if _, err = conf.Check("path/to/test", fset, files, &info); err != nil {
		log.Fatal(err)
	}
	return Testdata{
		Files: files,
		Defs:  info.Defs,
	}
}

func FindNamedType(name string, tdata Testdata) *types.Named {
	for _, f := range tdata.Files {
		for _, decl := range f.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.TYPE {
				continue
			}

			for _, spec := range gen.Specs {
				typ, ok := spec.(*ast.TypeSpec)
				if !ok || typ.Name.Name != name {
					continue
				}
				if obj, ok := tdata.Defs[typ.Name]; ok {
					if tn, ok := obj.(*types.TypeName); ok {
						if named, ok := tn.Type().(*types.Named); ok {
							return named
						}
					}
				}
			}
		}
	}

	return nil
}
