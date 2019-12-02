package golang

import (
	"github.com/frk/gosql/internal/writer"
)

type File struct {
	PkgName  string
	Imports  ImportDecl
	Decls    []Decl
	Preamble Comment
}

func (f File) Walk(w *writer.Writer) {
	if f.Preamble != nil {
		f.Preamble.Walk(w)
		w.Write("\n\n")
	}

	w.Write("package ")
	w.Write(f.PkgName)

	if f.Imports != nil {
		w.Write("\n\n")
		f.Imports.Walk(w)
	}

	for _, d := range f.Decls {
		w.Write("\n\n")
		d.Walk(w)
	}
}
