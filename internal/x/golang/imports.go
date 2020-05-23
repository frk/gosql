package golang

import (
	"github.com/frk/gosql/internal/x/writer"
)

// ImportDecl produces an import declaration for one or more packages in parentheses.
//
//	import (
//		"foo/bar/baz"
//		"foo/bar/baz"
//	)
type ImportDecl struct {
	Doc   CommentNode // associated documentation; or nil
	Specs []ImportSpec
}

func (d ImportDecl) Walk(w *writer.Writer) {
	if len(d.Specs) > 0 {
		if d.Doc != nil {
			d.Doc.Walk(w)
			w.Write("\n")
		}

		w.Write("import (\n")
		for _, s := range d.Specs {
			s.Walk(w)
			w.Write("\n")
		}
		w.Write(")")
	}
}

// SingleImportDecl produces an import declaration for a single package.
//
//	import "foo/bar/baz"
type SingleImportDecl struct {
	Doc  CommentNode // associated documentation; or nil
	Spec ImportSpec
}

func (d SingleImportDecl) Walk(w *writer.Writer) {
	if d.Doc != nil {
		d.Doc.Walk(w)
		w.Write("\n")
	}
	if d.Spec.Doc != nil {
		d.Spec.Doc.Walk(w)
		w.Write("\n")
		d.Spec.Doc = nil
	}

	w.Write("import ")
	d.Spec.Walk(w)
}

// ImportSpec produces the import path in an import declaration.
type ImportSpec struct {
	Doc     CommentNode // associated documentation; or nil
	Name    Ident       // local package name (including "." and "_"); or empty
	Path    StringLit   // import path
	Comment LineComment // trailing comment
}

func (s ImportSpec) Walk(w *writer.Writer) {
	if s.Doc != nil {
		s.Doc.Walk(w)
		w.Write("\n")
	}
	if len(s.Name.Name) > 0 {
		s.Name.Walk(w)
		w.Write(" ")
	}
	s.Path.Walk(w)
	s.Comment.Walk(w)
}

// implements ImportDeclNode
func (ImportDecl) importDeclNode()       {}
func (SingleImportDecl) importDeclNode() {}
