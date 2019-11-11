package golang

import (
	"github.com/frk/gosql/internal/writer"
)

type CommentList []Comment

func (list CommentList) Walk(w *writer.Writer) {
	if len(list) < 1 {
		return
	}

	w.Write(" ")
	list[0].Walk(w)
	for _, c := range list[1:] {
		w.Write("\n")
		c.Walk(w)
	}
}

type Comment struct {
	Text string
}

func (c Comment) Walk(w *writer.Writer) {
	w.Write("//")
	w.Write(c.Text)
}

func (CommentList) stmtNode() {}
func (Comment) stmtNode()     {}
