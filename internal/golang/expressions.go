package golang

import (
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

type Expr interface {
	Node
	exprNode()
}

type ExprList []Expr

func (list ExprList) Walk(w *writer.Writer) {
	list[0].Walk(w)
	for _, x := range list[1:] {
		w.Write(", ")
		x.Walk(w)
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

// A Field represents a field declaration in a field list in a struct type.
// Field.Names is nil for embedded struct fields.
type Field struct {
	Names []Ident
	Type  Expr
	Tag   RawString
	// Doc, Comment
}

func (f Field) Walk(w *writer.Writer) {
	num := len(f.Names)
	for i := 0; i < num; i++ {
		f.Names[i].Walk(w)
		if i < (num - 1) { // not last
			w.Write(", ")
		}
	}
	if num > 0 {
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
	List     []Expr
	Ellipsis bool
	// If 0 the arguments will be listed all on one line, if > 0 then the
	// arguments will be listed one per line starting from the argument whose
	// ordinal number matches the OnePerLine value.
	OnePerLine int
}

func (a ArgsList) Walk(w *writer.Writer) {
	num := len(a.List)
	for i := 0; i < num; i++ {
		if a.OnePerLine > 0 && i >= a.OnePerLine-1 {
			w.Write("\n")
		}
		a.List[i].Walk(w)
		if i < (num - 1) { // not last
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

func (a *ArgsList) AddExprs(xx ...Expr) {
	a.List = append(a.List, xx...)
}

// A Param represents a parameter/result declaration in a signature.
type Param struct {
	Names []Ident
	Type  Expr
	// Doc, Comment
}

func (p Param) Walk(w *writer.Writer) {
	num := len(p.Names)
	if num > 0 {
		p.Names[0].Walk(w)
		for _, n := range p.Names[1:] {
			w.Write(", ")
			n.Walk(w)
		}
		w.Write(" ")
	}
	p.Type.Walk(w)
}

type RecvParam struct {
	Name Ident
	Type Expr
}

func (p RecvParam) Walk(w *writer.Writer) {
	if len(p.Name.Name) > 0 {
		p.Name.Walk(w)
		w.Write(" ")
	}
	p.Type.Walk(w)
}

// An Ident node represents an identifier.
type Ident struct {
	Name string // identifier name
}

func (id Ident) Walk(w *writer.Writer) {
	w.Write(id.Name)
}

type IdentString string

func (id IdentString) Walk(w *writer.Writer) {
	w.Write(string(id))
}

// A StarExpr node represents an expression of the form "*" Expression.
// Semantically it could be a unary "*" expression, or a pointer type.
type StarExpr struct {
	X Expr
}

func (x StarExpr) Walk(w *writer.Writer) {
	w.Write("*")
	x.X.Walk(w)
}

// A UnaryExpr node represents a unary expression. Unary "*" expressions
// are represented via StarExpr nodes.
type UnaryExpr struct {
	Op UNARY_OP
	X  Expr
}

func (x UnaryExpr) Walk(w *writer.Writer) {
	w.Write(string(x.Op))
	x.X.Walk(w)
}

// A SelectorExpr node represents an expression followed by a selector.
type SelectorExpr struct {
	X   Expr
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
	X  Expr
	Op BINARY_OP
	Y  Expr
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
	Fun  Expr
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
	X Expr
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
	X     Expr
	Index Expr
}

func (x IndexExpr) Walk(w *writer.Writer) {
	x.X.Walk(w)
	w.Write("[")
	x.Index.Walk(w)
	w.Write("]")
}

// A KeyValueExpr node represents (key : value) pairs in composite literals.
type KeyValueExpr struct {
	Key   Expr
	Value Expr
}

func (x KeyValueExpr) Walk(w *writer.Writer) {
	x.Key.Walk(w)
	w.Write(": ")
	x.Value.Walk(w)
}

// An SliceExpr node represents an expression followed by slice indices.
type SliceExpr struct {
	X    Expr
	Low  Expr
	High Expr
	Max  Expr
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
	X    Expr
	Type Expr
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

type Ellipsis struct {
	Elt Expr
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
	Type    Expr
	Elts    []Expr
	Comma   bool
	Compact bool
}

func (lit CompositeLit) Walk(w *writer.Writer) {
	lit.Type.Walk(w)
	w.Write("{")
	for _, x := range lit.Elts {
		if !lit.Compact {
			w.Write("\n")
		}

		x.Walk(w)

		if lit.Comma {
			w.Write(",")
		}
	}
	if len(lit.Elts) > 0 && !lit.Compact {
		w.Write("\n")
	}
	w.Write("}")
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
	Exprs []Expr
}

func (m MultiLineExpr) Walk(w *writer.Writer) {
	for i, e := range m.Exprs {
		if i > 0 {
			w.Write(" ")
			w.Write(string(m.Op))
			w.NewLine()
		}
		e.Walk(w)
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

func (ExprList) exprNode()        {}
func (Ident) exprNode()           {}
func (IdentString) exprNode()     {}
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
