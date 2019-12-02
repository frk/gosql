package golang

import (
	"github.com/frk/gosql/internal/writer"
)

type Spec interface {
	Node
	specNode()
}

type ImportSpec struct {
	Name Ident
	Path String
	// Doc, Comment
	NewLine bool
}

func (s ImportSpec) Walk(w *writer.Writer) {
	if s.NewLine {
		w.Write("\n")
		return
	}

	if len(s.Name.Name) > 0 {
		s.Name.Walk(w)
		w.Write(" ")
	}
	s.Path.Walk(w)
}

type ValueSpec struct {
	Names       []Ident
	Type        Expr
	Values      []Expr
	LineComment LineComment
	Doc         Comment
}

func (s ValueSpec) Walk(w *writer.Writer) {
	if s.Doc != nil {
		s.Doc.Walk(w)
	}

	s.Names[0].Walk(w)
	for _, name := range s.Names[1:] {
		w.Write(", ")
		name.Walk(w)
	}
	if s.Type != nil {
		w.Write(" ")
		s.Type.Walk(w)
	}
	if len(s.Values) > 0 {
		w.Write(" = ")
		s.Values[0].Walk(w)
		for _, value := range s.Values[1:] {
			w.Write(", ")
			value.Walk(w)
		}
	}
	s.LineComment.Walk(w)
}

type TypeSpec struct {
	Name  Ident
	Alias bool
	Type  Expr // Ident, ParenExpr, SelectorExpr, StarExpr, or any of the XxxTypes
	// Doc, Comment
}

func (s TypeSpec) Walk(w *writer.Writer) {
	s.Name.Walk(w)
	if s.Alias {
		w.Write(" = ")
	} else {
		w.Write(" ")
	}
	s.Type.Walk(w)
}

func (ImportSpec) specNode() {}
func (ValueSpec) specNode()  {}
func (TypeSpec) specNode()   {}
