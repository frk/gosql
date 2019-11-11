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
	"bytes"
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

type directory struct {
	path    string
	name    string
	fset    *token.FileSet
	files   []*ast.File
	astpkg  *ast.Package
	typpkg  *types.Package
	typinfo *types.Info
}

// parsedir parses and type-checks the directory at the given path.
func parsedir(path string) (*directory, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, filternotestfiles, parser.ParseComments)
	if err != nil {
		return nil, err
	} else if len(pkgs) != 1 {
		// This should not happen but it's here just to make sure everything
		// works as expected.
		//
		// Go allows only one package per directory, however it is possible to
		// have test files declare an additional xxx_test package, but since the
		// filternotestfiles was passed to ParseDir those test files and by virtue
		// that package ought to be omitted by the parser.
		return nil, fmt.Errorf("unexpected number parsed packages, want 1 got %d", len(pkgs))
	}

	dir := new(directory)
	dir.path = path
	dir.name = filepath.Base(path)
	dir.fset = fset

	// Turn the package's map of files into a slice of files for type checking.
	for _, pkg := range pkgs {
		dir.astpkg = pkg
		for _, f := range pkg.Files {
			dir.files = append(dir.files, f)
		}
	}

	// Type checking of the package's files is done here becaue it is the type
	// checker that imports, and provides information on, all the referenced types
	// that we need for the subsequent analysis of the target types.
	conf := types.Config{Importer: typesutil.NewImporter()}
	info := types.Info{Defs: make(map[*ast.Ident]types.Object)}
	pkg, err := conf.Check(path, fset, dir.files, &info)
	if err != nil {
		return nil, err
	}

	dir.typpkg = pkg
	dir.typinfo = &info
	return dir, nil
}

type generator struct {
	file  string // the target file
	dir   *directory
	types []*types.Named
	cmds  []*command
	pg    *postgres
	buf   bytes.Buffer
}

func generate(file, dburl string) error {
	// Parse the whole directory in which the file is located, this is to
	// ensure that the generator has all the type info it may need.
	dirpath := filepath.Dir(file)
	dir, err := parsedir(dirpath)
	if err != nil {
		return err
	}

	pg := &postgres{url: dburl}
	if err := pg.init(); err != nil {
		return err
	}
	defer pg.close()

	g := &generator{file: file, dir: dir, pg: pg}
	return g.run()
}

func (g *generator) init() {
	// aggregates *types.Named instances of all of the target types
	// declared in the target file.
	for _, d := range g.dir.astpkg.Files[g.file].Decls {
		gd, ok := d.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			typ, ok := spec.(*ast.TypeSpec)
			if !ok || !retargetname.MatchString(typ.Name.Name) {
				continue
			}
			if obj, ok := g.dir.typinfo.Defs[typ.Name]; ok {
				if tn, ok := obj.(*types.TypeName); ok {
					if named, ok := tn.Type().(*types.Named); ok {
						g.types = append(g.types, named)
					}
				}

			}
		}
	}
}

func (g *generator) run() error {
	g.init()
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
	return pgcheck(g.pg, g.cmds)
}

func (g *generator) generate() error {
	return nil
}

// filternotestfiles is intended to be passed in as the filter argument to
// the parser.ParseDir function so that it can ignore files ending in _test.go.
func filternotestfiles(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go")
}
