package golang

import (
	"strconv"

	"github.com/frk/gosql/internal/writer"
)

// From spec:
// - unary_op = "+" | "-" | "!" | "^" | "*" | "&" | "<-" .
type UNARY_OP string

const (
	UNARY_ADD  UNARY_OP = "+"
	UNARY_SUB  UNARY_OP = "-"
	UNARY_NOT  UNARY_OP = "!"
	UNARY_XOR  UNARY_OP = "^"
	UNARY_MUL  UNARY_OP = "*"
	UNARY_AMP  UNARY_OP = "&"
	UNARY_RECV UNARY_OP = "<-"
)

// From spec:
// - binary_op = "||" | "&&" | rel_op | add_op | mul_op .
// - rel_op    = "==" | "!=" | "<" | "<=" | ">" | ">=" .
// - add_op    = "+" | "-" | "|" | "^" .
// - mul_op    = "*" | "/" | "%" | "<<" | ">>" | "&" | "&^" .
type BINARY_OP string

const (
	BINARY_LOR     BINARY_OP = "||"
	BINARY_LAND    BINARY_OP = "&&"
	BINARY_EQL     BINARY_OP = "=="
	BINARY_NEQ     BINARY_OP = "!="
	BINARY_LSS     BINARY_OP = "<"
	BINARY_LEQ     BINARY_OP = "<="
	BINARY_GTR     BINARY_OP = ">"
	BINARY_GEQ     BINARY_OP = ">="
	BINARY_ADD     BINARY_OP = "+"
	BINARY_SUB     BINARY_OP = "-"
	BINARY_OR      BINARY_OP = "|"
	BINARY_XOR     BINARY_OP = "^"
	BINARY_MUL     BINARY_OP = "*"
	BINARY_QUO     BINARY_OP = "/"
	BINARY_REM     BINARY_OP = "%"
	BINARY_SHL     BINARY_OP = "<<"
	BINARY_SHR     BINARY_OP = ">>"
	BINARY_AND     BINARY_OP = "&"
	BINARY_AND_NOT BINARY_OP = "&^"
)

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

// A Field represents a field declaration in a field list in a struct type.
// Field.Names is nil for embedded struct fields.
type Field struct {
	Names IdentNode
	Type  ExprNode
	Tag   RawString
	// Doc, Comment
}

func (f Field) Walk(w *writer.Writer) {
	if f.Names != nil {
		f.Names.Walk(w)
		w.Write(" ")
	}
	f.Type.Walk(w)

	if len(f.Tag) > 0 {
		w.Write(" ")
		f.Tag.Walk(w)
	}
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

// A Method represents a method declaration in a method list in an interface type.
type Method struct {
	Name Ident
	Type FuncType
	// Doc, Comment
}

func (m Method) Walk(w *writer.Writer) {
	m.Name.Walk(w)
	m.Type.Walk(w)
}

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

type ArgsList struct {
	List     ExprNodeList
	Ellipsis bool
	// If 0 the arguments will be listed all on one line, if > 0 then the
	// arguments will be listed one per line starting from the argument whose
	// ordinal number matches the OnePerLine value.
	OnePerLine int
}

func (a ArgsList) Walk(w *writer.Writer) {
	if a.List != nil {
		list := a.List.exprNodeList()
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

func (a *ArgsList) AddExprs(xx ...ExprNode) {
	if a.List != nil {
		a.List = ExprList(append(a.List.exprNodeList(), xx...))
	} else {
		a.List = ExprList(xx)
	}
}

func (a *ArgsList) Len() int {
	if a.List != nil {
		return len(a.List.exprNodeList())
	}
	return 0
}

// A Param represents a parameter/result declaration in a signature.
type Param struct {
	Names IdentNode
	Type  ExprNode
	// Doc, Comment
}

func (p Param) Walk(w *writer.Writer) {
	if p.Names != nil {
		p.Names.Walk(w)
		w.Write(" ")
	}
	p.Type.Walk(w)
}

type RecvParam struct {
	Name Ident
	Type ExprNode
}

func (p RecvParam) Walk(w *writer.Writer) {
	if len(p.Name.Name) > 0 {
		p.Name.Walk(w)
		w.Write(" ")
	}
	p.Type.Walk(w)
}

// A StarExpr node represents an expression of the form "*" Expression.
// Semantically it could be a unary "*" expression, or a pointer type.
type StarExpr struct {
	X ExprNode
}

func (x StarExpr) Walk(w *writer.Writer) {
	w.Write("*")
	x.X.Walk(w)
}

// A UnaryExpr node represents a unary expression. Unary "*" expressions
// are represented via StarExpr nodes.
type UnaryExpr struct {
	Op UNARY_OP
	X  ExprNode
}

func (x UnaryExpr) Walk(w *writer.Writer) {
	w.Write(string(x.Op))
	x.X.Walk(w)
}

// A SelectorExpr node represents an expression followed by a selector.
type SelectorExpr struct {
	X   ExprNode
	Sel Ident
}

func (x SelectorExpr) Walk(w *writer.Writer) {
	x.X.Walk(w)
	w.Write(".")
	x.Sel.Walk(w)
}

// A BinaryExpr node represents a unary expression. Unary "*" expressions
// are represented via StarExpr nodes.
type BinaryExpr struct {
	X  ExprNode
	Op BINARY_OP
	Y  ExprNode
}

func (x BinaryExpr) Walk(w *writer.Writer) {
	x.X.Walk(w)
	w.Write(" ")
	w.Write(string(x.Op))
	w.Write(" ")
	x.Y.Walk(w)
}

// A CallExpr node represents an expression followed by an argument list.
type CallExpr struct {
	Fun  ExprNode
	Args ArgsList
}

func (x CallExpr) Walk(w *writer.Writer) {
	x.Fun.Walk(w)
	w.Write("(")
	x.Args.Walk(w)
	w.Write(")")
}

// A ParenExpr node represents a parenthesized expression.
type ParenExpr struct {
	X ExprNode
}

func (x ParenExpr) Walk(w *writer.Writer) {
	w.Write("(")
	if x.X != nil {
		x.X.Walk(w)
	}
	w.Write(")")
}

// An IndexExpr node represents an expression followed by an index.
type IndexExpr struct {
	X     ExprNode
	Index ExprNode
}

func (x IndexExpr) Walk(w *writer.Writer) {
	x.X.Walk(w)
	w.Write("[")
	x.Index.Walk(w)
	w.Write("]")
}

// A KeyValueExpr node represents (key : value) pairs in composite literals.
type KeyValueExpr struct {
	Key   ExprNode
	Value ExprNode
}

func (x KeyValueExpr) Walk(w *writer.Writer) {
	x.Key.Walk(w)
	w.Write(": ")
	x.Value.Walk(w)
}

// An SliceExpr node represents an expression followed by slice indices.
type SliceExpr struct {
	X    ExprNode
	Low  ExprNode
	High ExprNode
	Max  ExprNode
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

// A TypeAssertExpr node represents an expression followed by a type assertion.
type TypeAssertExpr struct {
	X    ExprNode
	Type ExprNode
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

type Int int

func (i Int) Walk(w *writer.Writer) {
	w.Write(strconv.Itoa(int(i)))
}

type String string

func (s String) Walk(w *writer.Writer) {
	w.Write(`"`)
	w.Write(string(s))
	w.Write(`"`)
}

type RawString string

func (s RawString) Walk(w *writer.Writer) {
	w.Write("`")
	w.Write(string(s))
	w.Write("`")
}

type RawStringImplant struct {
	X ExprNode
}

func (s RawStringImplant) Walk(w *writer.Writer) {
	w.Write("` + ")
	s.X.Walk(w)
	w.Write(" + `")
}

type Ellipsis struct {
	Elt ExprNode
}

func (e Ellipsis) Walk(w *writer.Writer) {
	w.Write("...")
	if e.Elt != nil {
		e.Elt.Walk(w)
	}
}

type BasicLit struct {
	Value string
}

func (lit BasicLit) Walk(w *writer.Writer) {
	w.Write(lit.Value)
}

type CompositeLit struct {
	Type    ExprNode
	Elts    ExprNodeList
	Comma   bool
	Compact bool
}

func (lit CompositeLit) Walk(w *writer.Writer) {
	lit.Type.Walk(w)
	w.Write("{")

	if lit.Elts != nil {
		elts := lit.Elts.exprNodeList()
		for _, x := range elts {
			if !lit.Compact {
				w.Write("\n")
			}

			x.Walk(w)

			if lit.Comma {
				w.Write(",")
			}
		}
		if len(elts) > 0 && !lit.Compact {
			w.Write("\n")
		}
	}
	w.Write("}")
}

func (lit *CompositeLit) AddElts(xx ...ExprNode) {
	if lit.Elts != nil {
		lit.Elts = ExprList(append(lit.Elts.exprNodeList(), xx...))
	} else {
		lit.Elts = ExprList(xx)
	}
}

func (lit *CompositeLit) NumElts() int {
	if lit.Elts != nil {
		return len(lit.Elts.exprNodeList())
	}
	return 0
}

type FuncLit struct {
	Type FuncType
	Body BlockStmt
}

func (lit FuncLit) Walk(w *writer.Writer) {
	lit.Type.Walk(w)
	lit.Body.Walk(w)
}

type MultiLineExpr struct {
	Op    BINARY_OP
	Exprs ExprNodeList
}

func (m MultiLineExpr) Walk(w *writer.Writer) {
	if m.Exprs != nil {
		for i, x := range m.Exprs.exprNodeList() {
			if i > 0 {
				w.Write(" ")
				w.Write(string(m.Op))
				w.NewLine()
			}
			x.Walk(w)
		}
	}
}

type StringNode struct {
	Node Node
}

func (s StringNode) Walk(w *writer.Writer) {
	w.Write(`"`)
	s.Node.Walk(w)
	w.Write(`"`)
}

type RawStringNode struct {
	Node Node
}

func (s RawStringNode) Walk(w *writer.Writer) {
	w.Write("`")
	s.Node.Walk(w)
	w.Write("`")
}

type AffixStringNode struct {
	Prefix string
	Node   Node
	Suffix string
}

func (s AffixStringNode) Walk(w *writer.Writer) {
	w.Write(s.Prefix)
	s.Node.Walk(w)
	w.Write(s.Suffix)
}

func (StarExpr) exprNode()        {}
func (UnaryExpr) exprNode()       {}
func (SelectorExpr) exprNode()    {}
func (BinaryExpr) exprNode()      {}
func (CallExpr) exprNode()        {}
func (ParenExpr) exprNode()       {}
func (IndexExpr) exprNode()       {}
func (KeyValueExpr) exprNode()    {}
func (SliceExpr) exprNode()       {}
func (TypeAssertExpr) exprNode()  {}
func (Int) exprNode()             {}
func (String) exprNode()          {}
func (RawString) exprNode()       {}
func (Ellipsis) exprNode()        {}
func (BasicLit) exprNode()        {}
func (CompositeLit) exprNode()    {}
func (FuncLit) exprNode()         {}
func (MultiLineExpr) exprNode()   {}
func (StringNode) exprNode()      {}
func (RawStringNode) exprNode()   {}
func (AffixStringNode) exprNode() {}

func (x StarExpr) exprNodeList() []ExprNode        { return []ExprNode{x} }
func (x UnaryExpr) exprNodeList() []ExprNode       { return []ExprNode{x} }
func (x SelectorExpr) exprNodeList() []ExprNode    { return []ExprNode{x} }
func (x BinaryExpr) exprNodeList() []ExprNode      { return []ExprNode{x} }
func (x CallExpr) exprNodeList() []ExprNode        { return []ExprNode{x} }
func (x ParenExpr) exprNodeList() []ExprNode       { return []ExprNode{x} }
func (x IndexExpr) exprNodeList() []ExprNode       { return []ExprNode{x} }
func (x KeyValueExpr) exprNodeList() []ExprNode    { return []ExprNode{x} }
func (x SliceExpr) exprNodeList() []ExprNode       { return []ExprNode{x} }
func (x TypeAssertExpr) exprNodeList() []ExprNode  { return []ExprNode{x} }
func (x Int) exprNodeList() []ExprNode             { return []ExprNode{x} }
func (x String) exprNodeList() []ExprNode          { return []ExprNode{x} }
func (x RawString) exprNodeList() []ExprNode       { return []ExprNode{x} }
func (x Ellipsis) exprNodeList() []ExprNode        { return []ExprNode{x} }
func (x BasicLit) exprNodeList() []ExprNode        { return []ExprNode{x} }
func (x CompositeLit) exprNodeList() []ExprNode    { return []ExprNode{x} }
func (x FuncLit) exprNodeList() []ExprNode         { return []ExprNode{x} }
func (x MultiLineExpr) exprNodeList() []ExprNode   { return []ExprNode{x} }
func (x StringNode) exprNodeList() []ExprNode      { return []ExprNode{x} }
func (x RawStringNode) exprNodeList() []ExprNode   { return []ExprNode{x} }
func (x AffixStringNode) exprNodeList() []ExprNode { return []ExprNode{x} }
