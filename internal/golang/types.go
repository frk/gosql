package golang

import (
	"github.com/frk/gosql/internal/writer"
)

// ArrayType produces an array type literal.
type ArrayType struct {
	Len  ExprNode
	Elem TypeNode
}

func (at ArrayType) Walk(w *writer.Writer) {
	w.Write("[")
	at.Len.Walk(w)
	w.Write("]")
	at.Elem.Walk(w)
}

// SliceType produces a slice type literal.
type SliceType struct {
	Elem TypeNode
}

func (st SliceType) Walk(w *writer.Writer) {
	w.Write("[]")
	st.Elem.Walk(w)
}

// StructType produces a struct type literal.
type StructType struct {
	Fields FieldNode
}

func (st StructType) Walk(w *writer.Writer) {
	w.Write("struct")
	if st.Fields == nil {
		w.Write("{}")
		return
	}

	w.Write(" {\n")
	st.Fields.Walk(w)
	w.Write("\n}")
}

// A Field produces a field declaration in a field list of a struct type.
// Field.Names is nil for embedded struct fields.
type Field struct {
	Doc     CommentNode   // associated documentation
	Names   IdentListNode // the name of the field, if nil the field will be embedded
	Type    TypeNode      // field's type.
	Tag     RawStringLit  // the field's tag
	Comment CommentNode   // trailing comment
}

func (f Field) Walk(w *writer.Writer) {
	if f.Doc != nil {
		f.Doc.Walk(w)
		w.Write("\n")
	}
	if f.Names != nil {
		f.Names.Walk(w)
		w.Write(" ")
	}
	if f.Type != nil {
		f.Type.Walk(w)
	}

	if len(f.Tag) > 0 {
		w.Write(" ")
		f.Tag.Walk(w)
	}
	if f.Comment != nil {
		w.Write(" ")
		f.Comment.Walk(w)
	}
}

type FieldList []Field

func (list FieldList) Walk(w *writer.Writer) {
	if len(list) < 1 {
		return
	}

	list[0].Walk(w)
	for _, f := range list[1:] {
		w.Write("\n")
		f.Walk(w)
	}
}

// FuncType produces a function signature.
type FuncType Signature

func (ft FuncType) Walk(w *writer.Writer) {
	w.Write("func")
	Signature(ft).Walk(w)
}

// Signature produces a function signature.
type Signature struct {
	Params  ParamList
	Results ParamList
}

func (sig Signature) Walk(w *writer.Writer) {
	w.Write("(")
	sig.Params.Walk(w)
	w.Write(")")

	if len(sig.Results) < 1 {
		return
	}
	w.Write(" ")

	withParens := len(sig.Results) > 1 || sig.Results[0].Names != nil

	if withParens {
		w.Write("(")
	}
	sig.Results.Walk(w)
	if withParens {
		w.Write(")")
	}

}

// ParamList produces a list of parameters/results in a function signature.
type ParamList []Param

func (list ParamList) Walk(w *writer.Writer) {
	if len(list) < 1 {
		return
	}

	list[0].Walk(w)
	for _, m := range list[1:] {
		w.Write(", ")
		m.Walk(w)
	}
}

// Param produces a parameter/result declaration in a signature.
type Param struct {
	Names    IdentListNode
	Type     TypeNode
	Variadic bool
}

func (p Param) Walk(w *writer.Writer) {
	if p.Names != nil {
		p.Names.Walk(w)
		w.Write(" ")
	}
	if p.Variadic {
		w.Write("...")
	}
	p.Type.Walk(w)
}

// InterfaceType produces an interface type literal.
type InterfaceType struct {
	Methods MethodNode
}

func (it InterfaceType) Walk(w *writer.Writer) {
	w.Write("interface")
	if it.Methods == nil {
		w.Write("{}")
		return
	}

	w.Write(" {\n")
	it.Methods.Walk(w)
	w.Write("\n}")
}

// A Method produces a method declaration in a method list of an interface type.
type Method struct {
	Doc     CommentNode // associated documentation
	Name    Ident       // the method's name
	Type    Signature   // the method's signature
	Comment CommentNode // trailing comment
}

func (m Method) Walk(w *writer.Writer) {
	m.Name.Walk(w)
	m.Type.Walk(w)
}

type MethodList []Method

func (list MethodList) Walk(w *writer.Writer) {
	if len(list) < 1 {
		return
	}

	list[0].Walk(w)
	for _, m := range list[1:] {
		w.Write("\n")
		m.Walk(w)
	}
}

// MapType produces a map type literal.
type MapType struct {
	Key   TypeNode // the map's key type
	Value TypeNode // the map's value type
}

func (mt MapType) Walk(w *writer.Writer) {
	w.Write("map[")
	mt.Key.Walk(w)
	w.Write("]")
	mt.Value.Walk(w)
}

type ChanDir int

const (
	ChanBoth ChanDir = iota
	ChanRecv
	ChanSend
)

// ChanType produces a channel type literal.
type ChanType struct {
	Dir  ChanDir  // channel direction
	Elem TypeNode // the element type of the channel
}

func (ct ChanType) Walk(w *writer.Writer) {
	if ct.Dir == ChanRecv {
		w.Write("<-")
	}
	w.Write("chan")
	if ct.Dir == ChanSend {
		w.Write("<-")
	}
	w.Write(" ")
	ct.Elem.Walk(w)
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

// implements TypeListNode
func (t ArrayType) typeListNode() []TypeNode     { return []TypeNode{t} }
func (t SliceType) typeListNode() []TypeNode     { return []TypeNode{t} }
func (t StructType) typeListNode() []TypeNode    { return []TypeNode{t} }
func (t FuncType) typeListNode() []TypeNode      { return []TypeNode{t} }
func (t InterfaceType) typeListNode() []TypeNode { return []TypeNode{t} }
func (t MapType) typeListNode() []TypeNode       { return []TypeNode{t} }
func (t ChanType) typeListNode() []TypeNode      { return []TypeNode{t} }
func (t PointerType) typeListNode() []TypeNode   { return []TypeNode{t} }

// implements ExprNode
func (ArrayType) exprNode()     {}
func (SliceType) exprNode()     {}
func (StructType) exprNode()    {}
func (FuncType) exprNode()      {}
func (InterfaceType) exprNode() {}
func (MapType) exprNode()       {}
func (ChanType) exprNode()      {}
func (PointerType) exprNode()   {}

// implements ExprListNode
func (t ArrayType) exprListNode() []ExprNode     { return []ExprNode{t} }
func (t SliceType) exprListNode() []ExprNode     { return []ExprNode{t} }
func (t StructType) exprListNode() []ExprNode    { return []ExprNode{t} }
func (t FuncType) exprListNode() []ExprNode      { return []ExprNode{t} }
func (t InterfaceType) exprListNode() []ExprNode { return []ExprNode{t} }
func (t MapType) exprListNode() []ExprNode       { return []ExprNode{t} }
func (t ChanType) exprListNode() []ExprNode      { return []ExprNode{t} }
func (t PointerType) exprListNode() []ExprNode   { return []ExprNode{t} }

// implements FieldNode
func (Field) fieldNode()     {}
func (FieldList) fieldNode() {}

// implements MethodNode
func (Method) methodNode()     {}
func (MethodList) methodNode() {}
