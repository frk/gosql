package search

import (
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/tools/go/packages"
)

const loadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedExportFile |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedTypesSizes |
	packages.NeedModule |
	packages.NeedEmbedFiles |
	packages.NeedEmbedPatterns |
	packages.NeedTarget

var (
	// Matches names of types that are valid targets for the generator.
	rxTargetName = regexp.MustCompile(`^(?i:Select|Insert|Update|Delete|Filter)`)
)

// Match holds information on a matched query struct type.
type Match struct {
	// The go/types.Named representation of the matched type.
	Named *types.Named
	// The source position of the matched type.
	Pos token.Pos
}

// File represents a Go file that contains one or more matching query struct types.
type File struct {
	Path    string
	Package *Package
	Matches []*Match
}

// Package represents a Go package that contains one or more matching query struct types.
type Package struct {
	Name  string
	Path  string
	Fset  *token.FileSet
	Info  *types.Info
	Files []*File
}

// Search
func Search(dir string, recursive bool, filter func(filePath string) bool) (out []*Package, err error) {
	// resolve absolute dir path
	if dir, err = filepath.Abs(dir); err != nil {
		return nil, err
	}

	// if no filter was provided, pass all files
	if filter == nil {
		filter = func(string) bool { return true }
	}

	// initialize the pattern to use with packages.Load
	pattern := "."
	if recursive {
		pattern = "./..."
	}

	loadConfig := new(packages.Config)
	loadConfig.Mode = loadMode
	loadConfig.Dir = dir
	loadConfig.Fset = token.NewFileSet()
	pkgs, err := packages.Load(loadConfig, pattern)
	if err != nil {
		return nil, err
	}

	// aggregate matches from all files in all packages
	for _, pkg := range pkgs {
		p := new(Package)
		p.Name = pkg.Name
		p.Path = pkg.PkgPath
		p.Fset = pkg.Fset
		p.Info = pkg.TypesInfo

		for i, syn := range pkg.Syntax {
			// ignore file?
			if filePath := pkg.CompiledGoFiles[i]; !filter(filePath) {
				continue
			}

			f := new(File)
			f.Path = pkg.CompiledGoFiles[i]
			f.Package = p
			for _, dec := range syn.Decls {
				gd, ok := dec.(*ast.GenDecl)
				if !ok || gd.Tok != token.TYPE || hasIgnoreDirective(gd.Doc) {
					continue
				}

				for _, spec := range gd.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok || !rxTargetName.MatchString(typeSpec.Name.Name) || hasIgnoreDirective(typeSpec.Doc) {
						continue
					}

					obj, ok := p.Info.Defs[typeSpec.Name]
					if !ok {
						continue
					}
					typeName, ok := obj.(*types.TypeName)
					if !ok {
						continue
					}
					named, ok := typeName.Type().(*types.Named)
					if !ok {
						continue
					}

					match := new(Match)
					match.Named = named
					match.Pos = typeName.Pos()
					f.Matches = append(f.Matches, match)
				}
			}

			// add file only if a match was found
			if len(f.Matches) > 0 {
				p.Files = append(p.Files, f)
			}
		}

		// add package only if it has files with matches
		if len(p.Files) > 0 {
			out = append(out, p)
		}
	}

	return out, nil
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
