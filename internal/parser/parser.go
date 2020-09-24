package parser

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
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo

var (
	// Matches names of types that are valid targets for the generator.
	rxTargetName = regexp.MustCompile(`^(?i:Select|Insert|Update|Delete|Filter)`)
)

type Target struct {
	Named *types.Named
	Pos   token.Pos
}

type File struct {
	Path    string
	Package *Package
	Targets []*Target
}

type Package struct {
	Name  string
	Path  string
	Fset  *token.FileSet
	Info  *types.Info
	Files []*File
}

// Parse parses Go packages at the given dir / pattern. Only packages that
// contain files with type declarations that match the standard gosql targets
// will be included in the returned slice.
func Parse(dir string, recursive bool, filter func(filePath string) bool) (out []*Package, err error) {
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

	// aggregate targets from all files in all packages
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

					target := new(Target)
					target.Named = named
					target.Pos = typeName.Pos()
					f.Targets = append(f.Targets, target)
				}
			}

			// add file only if it declares targets
			if len(f.Targets) > 0 {
				p.Files = append(p.Files, f)
			}
		}

		// add package only if it has files with targets
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
