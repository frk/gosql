package golang

import (
	"github.com/frk/gosql/internal/x/writer"
)

// LineComment produces a single //-style comment.
type LineComment struct {
	Text string
}

func (lc LineComment) Walk(w *writer.Writer) {
	if len(lc.Text) == 0 {
		return
	}

	w.Write("//")
	w.Write(lc.Text)
}

// LineCommentList produces a list of //-style comments.
type LineCommentList []string

func (list LineCommentList) Walk(w *writer.Writer) {
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

// BlockComment produces a /*-style comment where each item
// in the slice represents an individual line of the comment.
type BlockComment []string

func (gc BlockComment) Walk(w *writer.Writer) {
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

// implements StmtNode
func (LineComment) stmtNode()     {}
func (LineCommentList) stmtNode() {}
func (BlockComment) stmtNode()    {}

// implements CommentNode
func (LineComment) commentNode()     {}
func (LineCommentList) commentNode() {}
func (BlockComment) commentNode()    {}
