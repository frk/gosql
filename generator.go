// TODO(mkopriva): currently handles only single file, update to handle a list of files, or a file pattern.
//
// https://golang.org/cmd/go/#hdr-Package_lists_and_patterns
// The filepath argument:
// - no argument: process package in current directory
// - single argument:
//	- file path: process just that file
//	- dir path: process package in that dir
//	- pattern: process packages matching that pattern
// - multiple arguments:
//	- iterate over each one and apply same rules as for "single argument"
//
// TODO(mkopriva): by default the generator will process only those files that import the gosql package but
// do provide an option flag to disable that requirement. This can be useful in the case where the programmer
// feeds a file pattern to the generator but some of the files matching that pattern are expected to be ignored
// because the programmer knows that in those files there are no types declared that would match a valid command...
package gosql

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/frk/gosql/internal/typesutil"
)

var (
	// Matches names of types that are valid targets for the generator.
	retargetname = regexp.MustCompile(`^(Select|Insert|Update|Delete|Filter)`)
)

func generate(file, dburl string) error {
	// Parse the whole directory in which the file is located, this is
	// required so that the type checker doesn't fail just because the
	// target file references an identifier declared in another file of
	// the same package.
	fset := token.NewFileSet()
	dirpath := filepath.Dir(file)
	pkgs, err := parser.ParseDir(fset, dirpath, filternotestfiles, parser.ParseComments)
	if err != nil {
		return err
	} else if len(pkgs) != 1 {
		// This should not happen but it's here just to make sure everything
		// works as expected.
		//
		// Go allows only one package per directory, however it is possible to
		// have test files declare an additional xxx_test package, but since the
		// filternotestfiles was passed to ParseDir those test files and by virtue
		// that package ought to be omitted by the parser.
		return fmt.Errorf("unexpected number parsed packages, want 1 got %d", len(pkgs))
	}

	// Get the AST of the target file and also turn the package's map of
	// files into a slice of files for type checking.
	var astfile *ast.File
	var files []*ast.File
	for _, p := range pkgs {
		astfile = p.Files[file]
		for _, f := range p.Files {
			files = append(files, f)
		}
	}

	// Type checking of the package's files is done here becaue it is the type
	// checker that imports, and provides information on, all the referenced types
	// that we need for the upcoming analysis of the target types.
	conf := types.Config{Importer: typesutil.NewImporter()}
	info := types.Info{Defs: make(map[*ast.Ident]types.Object)}
	pkg, err := conf.Check(dirpath, fset, files, &info)
	if err != nil {
		return err
	}

	g := new(generator)
	g.dburl = dburl
	g.pkgname = pkg.Name()

	// Aggregate *types.Named instances of all of the target types declared in the target file.
	for _, d := range astfile.Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			typ, ok := spec.(*ast.TypeSpec)
			if !ok || !retargetname.MatchString(typ.Name.Name) {
				continue
			}
			if obj, ok := info.Defs[typ.Name]; ok {
				if tn, ok := obj.(*types.TypeName); ok {
					if named, ok := tn.Type().(*types.Named); ok {
						g.types = append(g.types, named)
					}
				}

			}
		}
	}
	return g.run()
}

type generator struct {
	dburl   string
	pkgname string
	types   []*types.Named
	cmds    []*command
}

func (g *generator) run() error {
	if err := g.analyze(); err != nil {
		return err
	}
	if err := g.dbcheck(); err != nil {
		return err
	}
	if err := g.generate(); err != nil {
		return err
	}
	return nil
}

func (g *generator) analyze() error {
	for _, typ := range g.types {
		cmd, err := analyze(typ)
		if err != nil {
			return err
		}
		g.cmds = append(g.cmds, cmd)
	}
	return nil
}

func (g *generator) dbcheck() error {
	return pgcheck(g.dburl, g.cmds)
}

func (g *generator) generate() error {
	// TODO
	return nil
}

// filternotestfiles is intended to be passed in as the filter argument to
// the parser.ParseDir function so that it can ignore files ending in _test.go.
func filternotestfiles(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go")
}
