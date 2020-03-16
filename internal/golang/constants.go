package golang

import (
	"github.com/frk/gosql/internal/writer"
)

type stringnode string

const (
	Ellipsis stringnode = "..."
)

func (s stringnode) Walk(w *writer.Writer) {
	w.Write(string(s))
}

func (stringnode) exprNode() {}

func (x stringnode) exprListNode() []ExprNode { return []ExprNode{x} }