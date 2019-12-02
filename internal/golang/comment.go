package golang

import (
	"github.com/frk/gosql/internal/writer"
)

// All comment nodes implement the Comment interface.
type Comment interface {
	Stmt
	commentNode()
}

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

func (LineComment) commentNode()    {}
func (GeneralComment) commentNode() {}

func (LineComment) stmtNode()    {}
func (GeneralComment) stmtNode() {}
