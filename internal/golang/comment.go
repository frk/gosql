package golang

import (
	"github.com/frk/gosql/internal/writer"
)

// A LineComment node represents a list of //-style comments.
type LineComment []string

func (list LineComment) Walk(w *writer.Writer) {
	if len(list) == 0 {
		return
	}

	w.Write("//")
	w.Write(list[0])
	for _, c := range list[1:] {
		w.Write("\n//")
		w.Write(c)
	}
}

// A GeneralComment node represents a /*-style comment where each item
// in the slice represents an individual line of the comment.
type GeneralComment []string

func (gc GeneralComment) Walk(w *writer.Writer) {
	if len(gc) == 0 {
		return
	}

	w.Write("/*")
	if len(gc) == 1 {
		w.Write(gc[0])
		w.Write("*/")
		return
	}

	w.Indent()
	for _, c := range gc {
		w.NewLine()
		w.Write(c)
	}
	w.Unindent()
	w.NewLine()
	w.Write("*/")
}

type CommentNodeList []CommentNode

func (list CommentNodeList) Walk(w *writer.Writer) {
	for _, n := range list {
		n.Walk(w)
		w.Write("\n")
	}
}

// implements StmtNode
func (LineComment) stmtNode()     {}
func (GeneralComment) stmtNode()  {}
func (CommentNodeList) stmtNode() {}

// implements CommentNode
func (LineComment) commentNode()     {}
func (GeneralComment) commentNode()  {}
func (CommentNodeList) commentNode() {}
