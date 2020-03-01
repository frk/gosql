package golang

import (
	"github.com/frk/gosql/internal/writer"
)

// File node produces a go file.
type File struct {
	Doc     CommentNode        // associated documentation
	PkgName string             // name of the package for the package clause
	Imports []ImportDeclNode   // import declarations
	Decls   []TopLevelDeclNode // top level declarations

	// Preamble produces a comment above the associated documentation and
	// the package clause but separates it from them using a new line.
	// Can be useful for adding comments to the top of the file that should
	// not be rendered by godoc as the package's documentation.
	Preamble CommentNode
}

func (f File) Walk(w *writer.Writer) {
	if f.Preamble != nil {
		f.Preamble.Walk(w)
		w.Write("\n\n")
	}
	if f.Doc != nil {
		f.Doc.Walk(w)
		w.Write("\n")
	}
	w.Write("package ")
	w.Write(f.PkgName)

	if len(f.Imports) > 0 {
		w.Write("\n\n")
		for _, imp := range f.Imports {
			imp.Walk(w)
		}
	}

	for _, d := range f.Decls {
		w.Write("\n\n")
		d.Walk(w)
	}
}
