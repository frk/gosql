package sqlang

import (
	"bytes"
	"io"

	"github.com/frk/gosql/internal/x/writer"
)

type Node interface {
	Walk(w *writer.Writer)
}

type Expr interface {
	Node
	exprNode()
}

func Write(n Node, w io.Writer) error {
	out := writer.NewWriter(w)
	out.Indent()
	n.Walk(out)
	return out.Err()
}

func ToString(n Node) (string, error) {
	b := new(bytes.Buffer)
	if err := Write(n, b); err != nil {
		return "", nil
	}
	return b.String(), nil
}
