package golang

import (
	"github.com/frk/gosql/internal/writer"
)

// ArrayType produces an array type literal.
type ArrayType struct {
	Len ExprNode
	Elt TypeNode
}

func (at ArrayType) Walk(w *writer.Writer) {
	w.Write("[")
	at.Len.Walk(w)
	w.Write("]")
	at.Elt.Walk(w)
}

// SliceType produces a slice type literal.
type SliceType struct {
	Elt TypeNode
}

func (st SliceType) Walk(w *writer.Writer) {
	w.Write("[]")
	st.Elt.Walk(w)
}

// StructType produces a struct type literal.
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

// FuncType produces a function type literal.
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

// InterfaceType produces an interface type literal.
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

// MapType produces a map type literal.
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

// ChanType produces a channel type literal.
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

// PointerType produces a pointer type literal.
type PointerType struct {
	Elem TypeNode
}

func (pt PointerType) Walk(w *writer.Writer) {
	w.Write("*")
	pt.Elem.Walk(w)
}

// implements TypeNode
func (ArrayType) typeNode()     {}
func (SliceType) typeNode()     {}
func (StructType) typeNode()    {}
func (FuncType) typeNode()      {}
func (InterfaceType) typeNode() {}
func (MapType) typeNode()       {}
func (ChanType) typeNode()      {}
func (PointerType) typeNode()   {}

// implements ExprNode
func (ArrayType) exprNode()     {}
func (SliceType) exprNode()     {}
func (StructType) exprNode()    {}
func (FuncType) exprNode()      {}
func (InterfaceType) exprNode() {}
func (MapType) exprNode()       {}
func (ChanType) exprNode()      {}
func (PointerType) exprNode()   {}

// implements ExprNodeList
func (t ArrayType) exprNodeList() []ExprNode     { return []ExprNode{t} }
func (t SliceType) exprNodeList() []ExprNode     { return []ExprNode{t} }
func (t StructType) exprNodeList() []ExprNode    { return []ExprNode{t} }
func (t FuncType) exprNodeList() []ExprNode      { return []ExprNode{t} }
func (t InterfaceType) exprNodeList() []ExprNode { return []ExprNode{t} }
func (t MapType) exprNodeList() []ExprNode       { return []ExprNode{t} }
func (t ChanType) exprNodeList() []ExprNode      { return []ExprNode{t} }
func (t PointerType) exprNodeList() []ExprNode   { return []ExprNode{t} }
