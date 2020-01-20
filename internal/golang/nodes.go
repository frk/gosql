package golang

import (
	"github.com/frk/gosql/internal/writer"
)

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
