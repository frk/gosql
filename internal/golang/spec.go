package golang

import (
	"github.com/frk/gosql/internal/writer"
)

type SpecNode interface {
	Node
	specNode()
}

type SpecList []SpecNode

func (ls SpecList) Walk(w *writer.Writer) {
	withParens := len(ls) > 1
	if withParens {
		w.Write("(\n")
	}

	ls[0].Walk(w)
	for _, n := range ls[1:] {
		w.Write("\n")
		n.Walk(w)
	}

	if withParens {
		w.Write("\n)")
	}
}

////////////////////////////////////////////////////////////////////////////////

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
	Names       IdentNode
	Type        ExprNode
	Values      ExprNodeList
	LineComment LineComment
	Doc         Comment
}

func (s ValueSpec) Walk(w *writer.Writer) {
	if s.Doc != nil {
		s.Doc.Walk(w)
	}

	s.Names.Walk(w)
	if s.Type != nil {
		w.Write(" ")
		s.Type.Walk(w)
	}

	if s.Values != nil {
		w.Write(" = ")

		vals := s.Values.exprNodeList()
		vals[0].Walk(w)
		for _, v := range vals[1:] {
			w.Write(", ")
			v.Walk(w)
		}
	}
	s.LineComment.Walk(w)
}

type TypeSpec struct {
	Name  Ident
	Alias bool
	Type  ExprNode // Ident, ParenExpr, SelectorExpr, StarExpr, or any of the XxxTypes
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

func (SpecList) specNode()   {}
func (ImportSpec) specNode() {}
func (ValueSpec) specNode()  {}
func (TypeSpec) specNode()   {}
