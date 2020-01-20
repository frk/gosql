package golang

import (
	"github.com/frk/gosql/internal/writer"
)

type IdentNode interface {
	Node
	identNode()
}

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

// implement IdentNode
func (Ident) identNode()     {}
func (IdentList) identNode() {}

// implements ExprNode
func (Ident) exprNode()     {}
func (IdentList) exprNode() {}

// implements ExprNodeList
func (i Ident) exprNodeList() []ExprNode     { return []ExprNode{i} }
func (i IdentList) exprNodeList() []ExprNode { return []ExprNode{i} }
