package parser

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var (
	// Matches names of types that are valid targets for the generator.
	rxTargetName = regexp.MustCompile(`^(?i:Select|Insert|Update|Delete|Filter)`)
)

type Target struct {
	Named *types.Named
	Pos   token.Pos
}

type File struct {
	FilePath  string
	Directory *Directory
	Targets   []*Target
}

type Directory struct {
	DirPath string
	FileSet *token.FileSet
	Package *ast.Package
	Info    *types.Info
	Files   []*File
}

// ParseDirectory parses and type-checks the directory at its given path.
func ParseDirectory(dirpath string) (*Directory, error) {
	directoryCache.RLock()
	dir := directoryCache.m[dirpath]
	directoryCache.RUnlock()
	if dir != nil {
		return dir, nil
	}

	dir = new(Directory)
	dir.DirPath = dirpath

	dir.FileSet = token.NewFileSet()
	pkgList, err := parser.ParseDir(dir.FileSet, dir.DirPath, ignoreTestFiles, parser.ParseComments)
	if err != nil {
		return nil, err
	} else if len(pkgList) != 1 {
		// This should not happen but it's here just to make sure everything
		// works as expected.
		//
		// Go allows only one package per directory, however it is possible to
		// have test files declare an additional xxx_test package, but since the
		// ignoreTestFiles was passed to ParseDir those test files and that
		// xxx_test package ought to be omitted by the parser.
		return nil, fmt.Errorf("unexpected number of parsed packages, want 1 got %d", len(pkgList))
	}

	// Turn the package's map of files into a slice of files for type checking.
	astFileList := []*ast.File{}
	for _, pkg := range pkgList {
		dir.Package = pkg
		for _, f := range pkg.Files {
			astFileList = append(astFileList, f)
		}
	}

	// Type checking of the package's files is done here becaue it is the type
	// checker that imports, and provides information on, all the referenced types
	// that we need for the subsequent analysis of the target types.
	conf := types.Config{Importer: importer.ForCompiler(dir.FileSet, "source", nil)}
	dir.Info = &types.Info{Defs: make(map[*ast.Ident]types.Object)}
	if _, err := conf.Check(dir.DirPath, dir.FileSet, astFileList, dir.Info); err != nil {
		return nil, err
	}

	directoryCache.Lock()
	directoryCache.m[dir.DirPath] = dir
	directoryCache.Unlock()
	return dir, nil
}

// ParseFileDirectories
func ParseFileDirectories(filepaths ...string) ([]*Directory, error) {
	dirs := []*Directory{}

	for _, fp := range filepaths {
		dirpath := filepath.Dir(fp)
		dir, err := ParseDirectory(dirpath)
		if err != nil {
			return nil, err
		}

		_ = FileWithTargetTypes(dir, fp)
		dirs = append(dirs, dir)
	}

	return dirs, nil
}

// FileWithTargetTypes aggregates *Target instances of all of the target types declared in the file.
func FileWithTargetTypes(dir *Directory, filepath string) *File {
	fileCache.RLock()
	f := fileCache.m[filepath]
	fileCache.RUnlock()
	if f != nil {
		return f
	}

	f = &File{FilePath: filepath, Directory: dir}
	dir.Files = append(dir.Files, f)

	for _, decl := range dir.Package.Files[f.FilePath].Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.TYPE || hasIgnoreDirective(genDecl.Doc) {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok || !rxTargetName.MatchString(typeSpec.Name.Name) || hasIgnoreDirective(typeSpec.Doc) {
				continue
			}

			if obj, ok := dir.Info.Defs[typeSpec.Name]; ok {
				if typeName, ok := obj.(*types.TypeName); ok {
					if named, ok := typeName.Type().(*types.Named); ok {
						f.Targets = append(f.Targets, &Target{Named: named, Pos: typeName.Pos()})
					}
				}

			}
		}
	}

	fileCache.Lock()
	fileCache.m[f.FilePath] = f
	fileCache.Unlock()
	return f
}

// ignoreTestFiles is intended to be passed in as the filter argument to the
// parser.ParseDir function so that it can ignore files ending in _test.go.
func ignoreTestFiles(fi os.FileInfo) bool {
	return !strings.HasSuffix(fi.Name(), "_test.go")
}

// hasIgnoreDirective reports whether or not the given documentation contains the "gosql:ignore" directive.
func hasIgnoreDirective(doc *ast.CommentGroup) bool {
	if doc != nil {
		for _, com := range doc.List {
			if strings.Contains(com.Text, "gosql:ignore") {
				return true
			}
		}
	}
	return false
}

var directoryCache = struct {
	sync.RWMutex
	m map[string]*Directory
}{m: make(map[string]*Directory)}

var fileCache = struct {
	sync.RWMutex
	m map[string]*File
}{m: make(map[string]*File)}
