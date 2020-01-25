package golang

import (
	"github.com/frk/gosql/internal/writer"
)

// An Ident node represents an identifier.
type Ident struct {
	Name string // identifier name
}

func (id Ident) Walk(w *writer.Writer) {
	w.Write(id.Name)
}

// IdentList holds one or more identifiers.
type IdentList []Ident

func (list IdentList) Walk(w *writer.Writer) {
	list[0].Walk(w)
	for _, id := range list[1:] {
		w.Write(", ")
		id.Walk(w)
	}
}

// A QualifiedIdent node produces a qualified identifier.
type QualifiedIdent struct {
	Package string // package name
	Name    string // identifier name
}

func (id QualifiedIdent) Walk(w *writer.Writer) {
	w.Write(id.Package)
	w.Write(".")
	w.Write(id.Name)
}

// implement IdentNode
func (Ident) identNode()          {}
func (IdentList) identNode()      {}
func (QualifiedIdent) identNode() {}

// implements ExprNode
func (Ident) exprNode()          {}
func (IdentList) exprNode()      {}
func (QualifiedIdent) exprNode() {}

// implements ExprNodeList
func (i Ident) exprNodeList() []ExprNode          { return []ExprNode{i} }
func (i IdentList) exprNodeList() []ExprNode      { return []ExprNode{i} }
func (i QualifiedIdent) exprNodeList() []ExprNode { return []ExprNode{i} }

// implements TypeNode
func (Ident) typeNode()          {}
func (QualifiedIdent) typeNode() {}
