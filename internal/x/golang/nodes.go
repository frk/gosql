package golang

import (
	"github.com/frk/gosql/internal/x/writer"
)

////////////////////////////////////////////////////////////////////////////////
// Comments
////////////////////////////////////////////////////////////////////////////////

// CommentNode interface represents a comment.
type CommentNode interface {
	StmtNode
	commentNode()
}

////////////////////////////////////////////////////////////////////////////////
// Identifiers
////////////////////////////////////////////////////////////////////////////////

// IdentNode interface represents a single identifier.
type IdentNode interface {
	Node
	identNode()
}

// IdentListNode interface represents a list of 0 or more identifiers.
type IdentListNode interface {
	Node
	identListNode()
}

////////////////////////////////////////////////////////////////////////////////
// Declarations
////////////////////////////////////////////////////////////////////////////////

// ImportDeclNode interface represents an import declaration node.
type ImportDeclNode interface {
	Node
	importDeclNode()
}

// TopLevelDeclNode interface represents a top level declaration node.
// - type, const, var, func, method
type TopLevelDeclNode interface {
	Node
	topLevelDeclNode()
}

// DeclNode interface represents a declaration node.
// - type, const, var
type DeclNode interface {
	Node
	declNode()
}

////////////////////////////////////////////////////////////////////////////////
// Specs
////////////////////////////////////////////////////////////////////////////////

// ValueSpecNode interface represents a value specification node.
type ValueSpecNode interface {
	Node
	valueSpecNode()
}

// TypeSpecNode interface represents a value specification node.
type TypeSpecNode interface {
	Node
	typeSpecNode()
}

////////////////////////////////////////////////////////////////////////////////
// Expressions
////////////////////////////////////////////////////////////////////////////////

// ExprNode interface represents a single expression node.
type ExprNode interface {
	Node
	exprNode()
	exprListNode() []ExprNode
}

// ExprListNode interface represents one or more expression nodes.
type ExprListNode interface {
	Node
	exprListNode() []ExprNode
}

// ExprList implements the ExprListNode interface.
type ExprList []ExprNode

func (list ExprList) Walk(w *writer.Writer) {
	list[0].Walk(w)
	for _, x := range list[1:] {
		w.Write(", ")
		x.Walk(w)
	}
}

func (list ExprList) exprListNode() []ExprNode { return list }

////////////////////////////////////////////////////////////////////////////////
// Types
////////////////////////////////////////////////////////////////////////////////

// TypeNode interface represents a type literal node, or type name node.
// - TypeName = Ident | QualifiedIdent
// - TypeLit = ArrayType | StructType | PointerType | FuncType | InterfaceType |
//	    SliceType | MapType | ChanType .
type TypeNode interface {
	Node
	typeNode()
}

// TypeListNode interface represents one or more type nodes.
type TypeListNode interface {
	Node
	typeListNode() []TypeNode
}

// TypeList implements the TypeListNode interface.
type TypeList []TypeNode

func (list TypeList) Walk(w *writer.Writer) {
	list[0].Walk(w)
	for _, x := range list[1:] {
		w.Write(", ")
		x.Walk(w)
	}
}

func (list TypeList) typeListNode() []TypeNode { return list }

// FieldNode interface represents a struct type's field node.
type FieldNode interface {
	Node
	fieldNode()
}

// MethodNode interface represents a interface type's method node.
type MethodNode interface {
	Node
	methodNode()
}

// RecvTypeNode interface represents a receiver's type node.
type RecvTypeNode interface {
	Node
	recvTypeNode()
}

////////////////////////////////////////////////////////////////////////////////
// Statements
////////////////////////////////////////////////////////////////////////////////

type StmtNode interface {
	Node
	stmtNode()
}

type ElseNode interface {
	Node
	elseNode()
}

type ForClauseNode interface {
	Node
	forClauseNode()
}

////////////////////////////////////////////////////////////////////////////////
// helpers
////////////////////////////////////////////////////////////////////////////////

// StringNode produces a string from the underlying Node.
type StringNode struct {
	Prefix  string
	N       Node
	Suffix  string
	Comment *LineComment
}

func (n StringNode) Walk(w *writer.Writer) {
	w.Write(`"`)
	w.Write(n.Prefix)
	n.N.Walk(w)
	w.Write(n.Suffix)
	w.Write(`"`)
	if n.Comment != nil {
		w.Write(" ")
		n.Comment.Walk(w)
	}
}

// RawStringNode produces a raw string from the underlying Node.
type RawStringNode struct {
	Prefix  string
	N       Node
	Suffix  string
	Comment *LineComment
}

func (n RawStringNode) Walk(w *writer.Writer) {
	w.Write("`")
	w.Write(n.Prefix)
	n.N.Walk(w)
	w.Write(n.Suffix)
	w.Write("`")
	if n.Comment != nil {
		w.Write(" ")
		n.Comment.Walk(w)
	}
}

func (StringNode) exprNode()    {}
func (RawStringNode) exprNode() {}

func (x StringNode) exprListNode() []ExprNode    { return []ExprNode{x} }
func (x RawStringNode) exprListNode() []ExprNode { return []ExprNode{x} }
