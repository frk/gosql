package golang

import (
	"github.com/frk/gosql/internal/writer"
)

// An ArrayType node represents an array or slice type.
type ArrayType struct {
	Len ExprNode
	Elt ExprNode
}

func (at ArrayType) Walk(w *writer.Writer) {
	w.Write("[")
	if at.Len != nil {
		at.Len.Walk(w)
	}
	w.Write("]")
	at.Elt.Walk(w)
}

// A StructType node represents a struct type.
type StructType struct {
	Fields FieldList
}

func (st StructType) Walk(w *writer.Writer) {
	w.Write("struct")
	if len(st.Fields) == 0 {
		w.Write("{}")
		return
	}

	w.Write(" {\n")
	st.Fields.Walk(w)
	w.Write("\n}")
}

// A FuncType node represents a function type.
type FuncType struct {
	Params  ParamList
	Results ParamList
	Func    bool
}

func (ft FuncType) Walk(w *writer.Writer) {
	if ft.Func {
		w.Write("func")
	}
	w.Write("(")
	ft.Params.Walk(w)
	w.Write(")")

	if len(ft.Results) < 1 {
		return
	}
	w.Write(" ")

	withParens := len(ft.Results) > 1 || ft.Results[0].Names != nil

	if withParens {
		w.Write("(")
	}
	ft.Results.Walk(w)
	if withParens {
		w.Write(")")
	}

}

// An InterfaceType node represents an interface type.
type InterfaceType struct {
	Methods MethodList
}

func (it InterfaceType) Walk(w *writer.Writer) {
	w.Write("interface")
	if len(it.Methods) == 0 {
		w.Write("{}")
		return
	}

	w.Write(" {\n")
	it.Methods.Walk(w)
	w.Write("\n}")
}

// A MapType node represents a map type.
type MapType struct {
	Key   ExprNode
	Value ExprNode
}

func (mt MapType) Walk(w *writer.Writer) {
	w.Write("map[")
	mt.Key.Walk(w)
	w.Write("]")
	mt.Value.Walk(w)
}

type CHAN_DIR int

const (
	CHAN_BOTH CHAN_DIR = iota
	CHAN_RECV
	CHAN_SEND
)

// A ChanType node represents a channel type.
type ChanType struct {
	Dir   CHAN_DIR
	Value ExprNode
}

func (ct ChanType) Walk(w *writer.Writer) {
	if ct.Dir == CHAN_RECV {
		w.Write("<-")
	}
	w.Write("chan")
	if ct.Dir == CHAN_SEND {
		w.Write("<-")
	}
	w.Write(" ")
	ct.Value.Walk(w)
}

// implements ExprNode
func (ArrayType) exprNode()     {}
func (StructType) exprNode()    {}
func (FuncType) exprNode()      {}
func (InterfaceType) exprNode() {}
func (MapType) exprNode()       {}
func (ChanType) exprNode()      {}

// implements ExprNodeList
func (t ArrayType) exprNodeList() []ExprNode     { return []ExprNode{t} }
func (t StructType) exprNodeList() []ExprNode    { return []ExprNode{t} }
func (t FuncType) exprNodeList() []ExprNode      { return []ExprNode{t} }
func (t InterfaceType) exprNodeList() []ExprNode { return []ExprNode{t} }
func (t MapType) exprNodeList() []ExprNode       { return []ExprNode{t} }
func (t ChanType) exprNodeList() []ExprNode      { return []ExprNode{t} }
