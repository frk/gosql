package golang

import (
	"strconv"

	"github.com/frk/gosql/internal/writer"
)

// IntLit produces an int literal.
type IntLit int64

func (i IntLit) Walk(w *writer.Writer) {
	w.Write(strconv.FormatInt(int64(i), 10))
}

// StringLit produces a string literal.
type StringLit string

func (s StringLit) Walk(w *writer.Writer) {
	w.Write(`"`)
	w.Write(string(s))
	w.Write(`"`)
}

// RawStringLit produces a raw string literal.
type RawStringLit string

func (s RawStringLit) Walk(w *writer.Writer) {
	w.Write("`")
	w.Write(string(s))
	w.Write("`")
}

// FuncLit produces a function literal.
type FuncLit struct {
	Type FuncType  // the function type
	Body BlockStmt // function's body
}

func (lit FuncLit) Walk(w *writer.Writer) {
	lit.Type.Walk(w)
	w.Write(" ")
	lit.Body.Walk(w)
}

// SliceLit produces a composite slice or array literal.
type SliceLit struct {
	Type    TypeNode     // the slice / array type
	Elems   ExprListNode // the list of elements
	Compact bool         // if set to true the product will be a one-liner
}

func (lit SliceLit) Walk(w *writer.Writer) {
	if lit.Type != nil {
		lit.Type.Walk(w)
	}
	w.Write("{")

	if lit.Elems != nil {
		elems := lit.Elems.exprListNode()
		for i, x := range elems {
			if !lit.Compact {
				w.Write("\n")
			}

			x.Walk(w)
			if !lit.Compact || (i+1) < len(elems) {
				w.Write(", ")
			}
		}
		if len(elems) > 0 && !lit.Compact {
			w.Write("\n")
		}
	}
	w.Write("}")
}

func (lit *SliceLit) AddElems(xx ...ExprNode) {
	if lit.Elems != nil {
		lit.Elems = ExprList(append(lit.Elems.exprListNode(), xx...))
	} else {
		lit.Elems = ExprList(xx)
	}
}

func (lit *SliceLit) NumElems() int {
	if lit.Elems != nil {
		return len(lit.Elems.exprListNode())
	}
	return 0
}

// StructLit produces a composite struct literal.
type StructLit struct {
	Type    TypeNode       // the struct type
	Elems   []FieldElement // the list of field-value pairs
	Compact bool           // if set to true the product will be a one-liner
}

func (lit StructLit) Walk(w *writer.Writer) {
	if lit.Type != nil {
		lit.Type.Walk(w)
	}
	w.Write("{")

	for i, x := range lit.Elems {
		if !lit.Compact {
			w.Write("\n")
		}

		x.Walk(w)
		if !lit.Compact || (i+1) < len(lit.Elems) {
			w.Write(", ")
		}
	}
	if len(lit.Elems) > 0 && !lit.Compact {
		w.Write("\n")
	}
	w.Write("}")
}

// FieldElement produces a field-value pair in a struct literal expression.
type FieldElement struct {
	Field string   // field name
	Value ExprNode // the field's value
}

func (e FieldElement) Walk(w *writer.Writer) {
	if len(e.Field) > 0 {
		w.Write(e.Field)
		w.Write(": ")
	}
	e.Value.Walk(w)
}

// MapLit produces a composite map literal.
type MapLit struct {
	Type    TypeNode // the map's type
	Elems   []KeyElement
	Compact bool
}

func (lit MapLit) Walk(w *writer.Writer) {
	if lit.Type != nil {
		lit.Type.Walk(w)
	}
	w.Write("{")

	for i, x := range lit.Elems {
		if !lit.Compact {
			w.Write("\n")
		}

		x.Walk(w)
		if !lit.Compact || (i+1) < len(lit.Elems) {
			w.Write(", ")
		}
	}
	if len(lit.Elems) > 0 && !lit.Compact {
		w.Write("\n")
	}
	w.Write("}")
}

// KeyElement produces a key-value pair in map composite literals.
type KeyElement struct {
	Key   ExprNode // the key expression
	Value ExprNode // the value expression
}

func (x KeyElement) Walk(w *writer.Writer) {
	x.Key.Walk(w)
	w.Write(": ")
	x.Value.Walk(w)
}

func (IntLit) exprNode()       {}
func (StringLit) exprNode()    {}
func (RawStringLit) exprNode() {}
func (SliceLit) exprNode()     {}
func (StructLit) exprNode()    {}
func (MapLit) exprNode()       {}
func (KeyElement) exprNode()   {}
func (FuncLit) exprNode()      {}

func (x IntLit) exprListNode() []ExprNode       { return []ExprNode{x} }
func (x StringLit) exprListNode() []ExprNode    { return []ExprNode{x} }
func (x RawStringLit) exprListNode() []ExprNode { return []ExprNode{x} }
func (x SliceLit) exprListNode() []ExprNode     { return []ExprNode{x} }
func (x StructLit) exprListNode() []ExprNode    { return []ExprNode{x} }
func (x MapLit) exprListNode() []ExprNode       { return []ExprNode{x} }
func (x KeyElement) exprListNode() []ExprNode   { return []ExprNode{x} }
func (x FuncLit) exprListNode() []ExprNode      { return []ExprNode{x} }
