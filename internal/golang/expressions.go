package golang

import (
	"github.com/frk/gosql/internal/writer"
)

type UnaryOp string

const (
	UnaryAdd  UnaryOp = "+"  // no-op
	UnarySub  UnaryOp = "-"  // negate a numeric expression
	UnaryNot  UnaryOp = "!"  // negate a boolean expression
	UnaryXOr  UnaryOp = "^"  // bitwise complement
	UnaryMul  UnaryOp = "*"  // pointer indirection
	UnaryAmp  UnaryOp = "&"  // address operation
	UnaryRecv UnaryOp = "<-" // channel receive
)

type BinaryOp string

const (
	BinaryLOr    BinaryOp = "||" // logical or
	BinaryLAnd   BinaryOp = "&&" // logical and
	BinaryEql    BinaryOp = "==" // is equal
	BinaryNeq    BinaryOp = "!=" // not equal
	BinaryLss    BinaryOp = "<"  // less than
	BinaryLeq    BinaryOp = "<=" // less than or equal
	BinaryGtr    BinaryOp = ">"  // greater than
	BinaryGeq    BinaryOp = ">=" // greater than or equal
	BinaryAdd    BinaryOp = "+"  // sum
	BinarySub    BinaryOp = "-"  // difference
	BinaryOr     BinaryOp = "|"  // bitwise or
	BinaryXOr    BinaryOp = "^"  // bitwise xor
	BinaryMul    BinaryOp = "*"  // product
	BinaryQuo    BinaryOp = "/"  // quotient
	BinaryRem    BinaryOp = "%"  // remainder
	BinaryShl    BinaryOp = "<<" // left shift
	BinaryShr    BinaryOp = ">>" // right shift
	BinaryAnd    BinaryOp = "&"  // bitwise and
	BinaryAndNot BinaryOp = "&^" // bit clear
)

// PointerIndirectionExpr produces a pointer indirection expression.
type PointerIndirectionExpr struct {
	X ExprNode // the operand
}

func (x PointerIndirectionExpr) Walk(w *writer.Writer) {
	w.Write("*")
	x.X.Walk(w)
}

// UnaryExpr produces a unary expression.
type UnaryExpr struct {
	Op UnaryOp  // the operator
	X  ExprNode // the operand
}

func (x UnaryExpr) Walk(w *writer.Writer) {
	w.Write(string(x.Op))
	x.X.Walk(w)
}

// SelectorExpr produces a selector expression, i.e. a primary expression followed
// by a dot followed by a selector. The selector must be a field or a method identifier.
type SelectorExpr struct {
	X   ExprNode // the primary expression
	Sel Ident    // the selector, i.e. field or method
}

func (x SelectorExpr) Walk(w *writer.Writer) {
	x.X.Walk(w)
	w.Write(".")
	x.Sel.Walk(w)
}

// BinaryExpr produces a binary expression.
type BinaryExpr struct {
	X  ExprNode // left operand
	Op BinaryOp // operator
	Y  ExprNode // right operand
}

func (x BinaryExpr) Walk(w *writer.Writer) {
	x.X.Walk(w)
	w.Write(" ")
	w.Write(string(x.Op))
	w.Write(" ")
	x.Y.Walk(w)
}

// CallExpr node produces a function/method call expression.
type CallExpr struct {
	Fun  ExprNode // function expression
	Args ArgsList // function arguments
}

func (x CallExpr) Walk(w *writer.Writer) {
	x.Fun.Walk(w)
	w.Write("(")
	x.Args.Walk(w)
	w.Write(")")
}

// ArgsList produces a list of arguments in a call expression.
type ArgsList struct {
	List     ExprListNode // the list of arguments
	Ellipsis bool         // if set, adds "..." to the last argument in the list

	// If 0 the arguments will be listed all on one line, if > 0 then the
	// arguments will be listed one per line starting from the argument whose
	// ordinal number matches the OnePerLine value.
	OnePerLine int
}

func (a ArgsList) Walk(w *writer.Writer) {
	if a.List != nil {
		list := a.List.exprListNode()
		length := len(list)

		for i := 0; i < length; i++ {
			if a.OnePerLine > 0 && i >= a.OnePerLine-1 {
				w.Write("\n")
			}
			list[i].Walk(w)
			if i < (length - 1) { // not last
				w.Write(", ")
			} else { // last
				if a.Ellipsis {
					w.Write("...")
				}
				if a.OnePerLine > 0 {
					w.Write(",\n")
				}
			}
		}
	}
}

// AddExprs appends the given ExprNodes to the argument list.
func (a *ArgsList) AddExprs(xx ...ExprNode) {
	if a.List != nil {
		a.List = ExprList(append(a.List.exprListNode(), xx...))
	} else {
		a.List = ExprList(xx)
	}
}

// Len returns the number of arguments in the list.
func (a *ArgsList) Len() int {
	if a.List != nil {
		return len(a.List.exprListNode())
	}
	return 0
}

// CallNewExpr produces a call expression to the "new" builtin function.
type CallNewExpr struct {
	Type TypeNode // the argument
}

func (x CallNewExpr) Walk(w *writer.Writer) {
	w.Write("new(")
	x.Type.Walk(w)
	w.Write(")")
}

// CallMakeExpr produces a call expression to the "make" builtin function.
type CallMakeExpr struct {
	Type TypeNode // the argument
	Size ExprNode // a slice's len, or a map's size, or a chan's buffer cap
	Cap  ExprNode // a slice's cap
}

func (x CallMakeExpr) Walk(w *writer.Writer) {
	w.Write("make(")
	x.Type.Walk(w)
	if x.Size != nil {
		w.Write(", ")
		x.Size.Walk(w)
		if x.Cap != nil {
			w.Write(", ")
			x.Cap.Walk(w)
		}
	}
	w.Write(")")
}

// CallLenExpr produces a call expression to the "len" builtin function.
type CallLenExpr struct {
	X ExprNode // the argument
}

func (x CallLenExpr) Walk(w *writer.Writer) {
	w.Write("len(")
	x.X.Walk(w)
	w.Write(")")
}

// ParenExpr produces a parenthesized expression. ParenExpr can be used
// to enforce evaluation order.
type ParenExpr struct {
	X ExprNode // the expression
}

func (x ParenExpr) Walk(w *writer.Writer) {
	w.Write("(")
	if x.X != nil {
		x.X.Walk(w)
	}
	w.Write(")")
}

// IndexExpr produces an index expression.
type IndexExpr struct {
	X     ExprNode // the primary expression
	Index ExprNode // the index / map key
}

func (x IndexExpr) Walk(w *writer.Writer) {
	x.X.Walk(w)
	w.Write("[")
	x.Index.Walk(w)
	w.Write("]")
}

// SliceExpr produces a slicing expression.
type SliceExpr struct {
	X    ExprNode // the operand
	Low  ExprNode // low bound
	High ExprNode // high bound
	Max  ExprNode // capacity bound
}

func (x SliceExpr) Walk(w *writer.Writer) {
	x.X.Walk(w)
	w.Write("[")
	if x.Low != nil {
		x.Low.Walk(w)
	}
	w.Write(":")
	if x.High != nil {
		x.High.Walk(w)
	}
	if x.Max != nil {
		w.Write(":")
		x.Max.Walk(w)
	}
	w.Write("]")
}

// TypeAssertExpr produces a type assertion expression.
type TypeAssertExpr struct {
	X    ExprNode // the primary expression
	Type TypeNode // the type to assert to
}

func (x TypeAssertExpr) Walk(w *writer.Writer) {
	x.X.Walk(w)
	w.Write(".(")
	if x.Type != nil {
		x.Type.Walk(w)
	} else {
		w.Write("type") // type switch
	}
	w.Write(")")
}

type MultiLineExpr struct {
	Op    BinaryOp
	Exprs []ExprNode
}

func (xx MultiLineExpr) Walk(w *writer.Writer) {
	for i, x := range xx.Exprs {
		if i > 0 {
			w.Write(" ")
			w.Write(string(xx.Op))
			w.NewLine()
		}
		x.Walk(w)
	}
}

func (PointerIndirectionExpr) exprNode() {}
func (UnaryExpr) exprNode()              {}
func (SelectorExpr) exprNode()           {}
func (BinaryExpr) exprNode()             {}
func (CallExpr) exprNode()               {}
func (CallNewExpr) exprNode()            {}
func (CallMakeExpr) exprNode()           {}
func (CallLenExpr) exprNode()            {}
func (ParenExpr) exprNode()              {}
func (IndexExpr) exprNode()              {}
func (SliceExpr) exprNode()              {}
func (TypeAssertExpr) exprNode()         {}
func (MultiLineExpr) exprNode()          {}

func (x PointerIndirectionExpr) exprListNode() []ExprNode { return []ExprNode{x} }
func (x UnaryExpr) exprListNode() []ExprNode              { return []ExprNode{x} }
func (x SelectorExpr) exprListNode() []ExprNode           { return []ExprNode{x} }
func (x BinaryExpr) exprListNode() []ExprNode             { return []ExprNode{x} }
func (x CallExpr) exprListNode() []ExprNode               { return []ExprNode{x} }
func (x CallNewExpr) exprListNode() []ExprNode            { return []ExprNode{x} }
func (x CallMakeExpr) exprListNode() []ExprNode           { return []ExprNode{x} }
func (x CallLenExpr) exprListNode() []ExprNode            { return []ExprNode{x} }
func (x ParenExpr) exprListNode() []ExprNode              { return []ExprNode{x} }
func (x IndexExpr) exprListNode() []ExprNode              { return []ExprNode{x} }
func (x SliceExpr) exprListNode() []ExprNode              { return []ExprNode{x} }
func (x TypeAssertExpr) exprListNode() []ExprNode         { return []ExprNode{x} }
func (x MultiLineExpr) exprListNode() []ExprNode          { return []ExprNode{x} }
