package testutil

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
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

	conf := types.Config{Importer: testingImporter{def: importer.Default(), src: importer.For("source", nil)}}
	info := types.Info{Defs: make(map[*ast.Ident]types.Object)}
	if _, err = conf.Check("path/to/test", fset, files, &info); err != nil {
		log.Fatal(err)
	}
	return Testdata{
		Files: files,
		Defs:  info.Defs,
	}
}

type testingImporter struct {
	def types.Importer // default
	src types.Importer // fallback
}

func (i testingImporter) Import(path string) (*types.Package, error) {
	// Try using the default, if that fails then fallback to source based import.
	pkg, err := i.def.Import(path)
	if err != nil {
		return i.src.Import(path)
	}
	return pkg, nil
}
