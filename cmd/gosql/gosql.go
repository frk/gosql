package main

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
// because the programmer knows that in those files there are no types declared that would match a valid TypeSpec...

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
	"sync"

	"github.com/frk/gosql/internal/typesutil"
)

var (
	// Matches names of types that are valid targets for the generator.
	retargetname = regexp.MustCompile(`^(?i:Select|Insert|Update|Delete|Filter)`)
)

type file struct {
	path  string
	dir   *directory
	types []*types.Named
}

type directory struct {
	path  string
	pkg   *ast.Package
	info  *types.Info
	files []*file
}

type command struct {
	pg   *postgres
	dirs []*directory

	// options
	// TODO(mkopriva): this is just getting copied verbatim to the the generator
	// which is ugly, I probably need a config struct that's passed to from the
	// command to the analyzer, type checker, and generator....
	keytag  string
	keysep  string
	keybase bool
}

func (cmd *command) exec(dburl string, files ...string) error {
	for _, file := range files {
		// Parse the whole directory in which the file is located, this is to
		// ensure that the generator has all the type info it may need.
		dirpath := filepath.Dir(file)
		dir, err := cmd.parsedir(dirpath)
		if err != nil {
			return err
		}
		_ = cmd.aggtypes(dir, file)
	}

	cmd.pg = &postgres{url: dburl}
	if err := cmd.pg.init(); err != nil {
		return err
	}
	defer cmd.pg.close()

	//
	for _, dir := range cmd.dirs {
		for _, f := range dir.files {
			b, err := cmd.run(f)
			if err != nil {
				return err
			}

			// TODO(mkopriva): write the buffered data into a file
			_ = b
		}
	}
	return nil
}

func (cmd *command) run(f *file) (*bytes.Buffer, error) {
	var infos []*targetInfo

	// analyze named types
	for _, typ := range f.types {
		a := &analyzer{named: typ}
		if err := a.run(); err != nil {
			return nil, err
		}
		infos = append(infos, a.targetInfo())
	}

	// type-check specs against the db
	for _, ti := range infos {
		c := pgchecker{pg: cmd.pg, ti: ti}
		if err := c.run(); err != nil {
			return nil, err
		}
	}

	// generate code
	g := &generator{infos: infos}
	g.keytag = cmd.keytag
	g.keysep = cmd.keysep
	g.keybase = cmd.keybase
	if err := g.run(f.dir.pkg.Name); err != nil {
		return nil, err
	}
	return &g.buf, nil
}

// parsedir parses and type-checks the directory at its given path.
func (cmd *command) parsedir(path string) (*directory, error) {
	directoryCache.RLock()
	dir := directoryCache.m[path]
	directoryCache.RUnlock()
	if dir != nil {
		return dir, nil
	}

	dir = new(directory)
	dir.path = path
	cmd.dirs = append(cmd.dirs, dir)

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir.path, cmd.notestfiles, parser.ParseComments)
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
		return nil, fmt.Errorf("unexpected number of parsed packages, want 1 got %d", len(pkgs))
	}

	// Turn the package's map of files into a slice of files for type checking.
	var astfiles []*ast.File
	for _, pkg := range pkgs {
		dir.pkg = pkg
		for _, f := range pkg.Files {
			astfiles = append(astfiles, f)
		}
	}

	// Type checking of the package's files is done here becaue it is the type
	// checker that imports, and provides information on, all the referenced types
	// that we need for the subsequent analysis of the target types.
	conf := types.Config{Importer: typesutil.NewImporter()}
	dir.info = &types.Info{Defs: make(map[*ast.Ident]types.Object)}
	if _, err := conf.Check(dir.path, fset, astfiles, dir.info); err != nil {
		return nil, err
	}

	directoryCache.Lock()
	directoryCache.m[dir.path] = dir
	directoryCache.Unlock()
	return dir, nil
}

// aggtypes aggregates *types.Named instances of all of the target types declared in the file.
func (cmd *command) aggtypes(dir *directory, path string) *file {
	fileCache.RLock()
	f := fileCache.m[path]
	fileCache.RUnlock()
	if f != nil {
		return f
	}

	f = new(file)
	f.path = path
	f.dir = dir
	dir.files = append(dir.files, f)

	for _, decl := range dir.pkg.Files[f.path].Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.TYPE {
			continue
		}
		for _, spec := range gd.Specs {
			typ, ok := spec.(*ast.TypeSpec)
			if !ok || !retargetname.MatchString(typ.Name.Name) {
				continue
			}
			if obj, ok := dir.info.Defs[typ.Name]; ok {
				if tn, ok := obj.(*types.TypeName); ok {
					if named, ok := tn.Type().(*types.Named); ok {
						f.types = append(f.types, named)
					}
				}

			}
		}
	}

	fileCache.Lock()
	fileCache.m[f.path] = f
	fileCache.Unlock()
	return f
}

// notestfiles is intended to be passed in as the filter argument to the
// parser.ParseDir function so that it can ignore files ending in _test.go.
func (cmd *command) notestfiles(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go")
}

var directoryCache = struct {
	sync.RWMutex
	m map[string]*directory
}{m: make(map[string]*directory)}

var fileCache = struct {
	sync.RWMutex
	m map[string]*file
}{m: make(map[string]*file)}
