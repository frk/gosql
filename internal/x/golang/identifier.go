package golang

import (
	"github.com/frk/gosql/internal/x/writer"
)

// Ident produces an identifier.
type Ident struct {
	Name string // identifier name
}

func (id Ident) Walk(w *writer.Writer) {
	w.Write(id.Name)
}

// IdentList holds one or more identifiers.
type IdentList []Ident

func (list IdentList) Walk(w *writer.Writer) {
	if len(list) > 0 {
		list[0].Walk(w)
		for _, id := range list[1:] {
			w.Write(", ")
			id.Walk(w)
		}
	}
}

// QualifiedIdent node produces a qualified identifier.
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
func (QualifiedIdent) identNode() {}

// implement IdentListNode
func (Ident) identListNode()          {}
func (IdentList) identListNode()      {}
func (QualifiedIdent) identListNode() {}

// implements ExprNode
func (Ident) exprNode()          {}
func (IdentList) exprNode()      {}
func (QualifiedIdent) exprNode() {}

// implements ExprListNode
func (i Ident) exprListNode() []ExprNode          { return []ExprNode{i} }
func (i IdentList) exprListNode() []ExprNode      { return []ExprNode{i} }
func (i QualifiedIdent) exprListNode() []ExprNode { return []ExprNode{i} }

// implements TypeNode
func (Ident) typeNode()          {}
func (QualifiedIdent) typeNode() {}

// implements TypeListNode
func (t Ident) typeListNode() []TypeNode          { return []TypeNode{t} }
func (t QualifiedIdent) typeListNode() []TypeNode { return []TypeNode{t} }

// implements RecvTypeNode
func (Ident) recvTypeNode() {}

// implements MethodNode
func (Ident) methodNode()          {}
func (QualifiedIdent) methodNode() {}
