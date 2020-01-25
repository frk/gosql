package golang

import (
	"github.com/frk/gosql/internal/writer"
)

////////////////////////////////////////////////////////////////////////////////
// Comments
////////////////////////////////////////////////////////////////////////////////

// CommentNode interface represents a comment.
type CommentNode interface {
	Stmt
	commentNode()
}

////////////////////////////////////////////////////////////////////////////////
// Identifiers
////////////////////////////////////////////////////////////////////////////////

// IdentNode interface represents a identifier.
type IdentNode interface {
	Node
	identNode()
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
	exprNodeList() []ExprNode
}

// ExprNodes interface represents one or more expression nodes.
type ExprNodeList interface {
	Node
	exprNodeList() []ExprNode
}

// ExprList implements the ExprNodeList interface.
type ExprList []ExprNode

func (list ExprList) Walk(w *writer.Writer) {
	list[0].Walk(w)
	for _, x := range list[1:] {
		w.Write(", ")
		x.Walk(w)
	}
}

func (ls ExprList) exprNodeList() []ExprNode { return ls }

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
