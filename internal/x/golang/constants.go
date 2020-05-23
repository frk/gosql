package golang

import (
	"github.com/frk/gosql/internal/x/writer"
)

type stringnode string

const (
	Ellipsis stringnode = "..."
	True     stringnode = "true"
	False    stringnode = "false"
)

func (s stringnode) Walk(w *writer.Writer) {
	w.Write(string(s))
}

func (stringnode) exprNode() {}

func (x stringnode) exprListNode() []ExprNode { return []ExprNode{x} }
