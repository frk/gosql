package golang

import (
	"github.com/frk/gosql/internal/writer"
)

type BRANCH_TOKEN string

const (
	BRANCH_BREAK BRANCH_TOKEN = "break"
	BRANCH_CONT  BRANCH_TOKEN = "continue"
	BRANCH_GOTO  BRANCH_TOKEN = "goto"
	BRANCH_FALL  BRANCH_TOKEN = "fallthrough"
)

type Stmt interface {
	Node
	stmtNode()
}

// A DeclStmt node represents a declaration in a statement list.
type DeclStmt struct {
	Decl Decl
}

func (s DeclStmt) Walk(w *writer.Writer) {
	s.Decl.Walk(w)
}

// A LabeledStmt node represents a labeled statement.
type LabeledStmt struct {
	Label Ident
	Stmt  Stmt
}

func (s LabeledStmt) Walk(w *writer.Writer) {
	s.Label.Walk(w)
	w.Write(":\n")
	s.Stmt.Walk(w)
}

// An ExprStmt node represents a (stand-alone) expression in a statement list.
type ExprStmt struct {
	X ExprNode
}

func (s ExprStmt) Walk(w *writer.Writer) {
	s.X.Walk(w)
}

type SendStmt struct {
	Chan  ExprNode
	Value ExprNode
}

func (s SendStmt) Walk(w *writer.Writer) {
	s.Chan.Walk(w)
	w.Write(" <- ")
	s.Value.Walk(w)
}

type INCDEC_TOKEN string

const (
	INCDEC_INC INCDEC_TOKEN = "++"
	INCDEC_DEC INCDEC_TOKEN = "--"
)

type IncDecStmt struct {
	X     ExprNode
	Token INCDEC_TOKEN
}

func (s IncDecStmt) Walk(w *writer.Writer) {
	s.X.Walk(w)
	w.Write(string(s.Token))
}

type ASSIGN_TOKEN string

const (
	ASSIGN         ASSIGN_TOKEN = "="
	ASSIGN_ADD     ASSIGN_TOKEN = "+="
	ASSIGN_SUB     ASSIGN_TOKEN = "-="
	ASSIGN_MUL     ASSIGN_TOKEN = "*="
	ASSIGN_QUO     ASSIGN_TOKEN = "%="
	ASSIGN_REM     ASSIGN_TOKEN = "%="
	ASSIGN_AND     ASSIGN_TOKEN = "&="
	ASSIGN_OR      ASSIGN_TOKEN = "|="
	ASSIGN_XOR     ASSIGN_TOKEN = "^="
	ASSIGN_SHL     ASSIGN_TOKEN = "<<="
	ASSIGN_SHR     ASSIGN_TOKEN = ">>="
	ASSIGN_AND_NOT ASSIGN_TOKEN = "&^="
	ASSIGN_DEFINE  ASSIGN_TOKEN = ":="
)

type AssignStmt struct {
	Lhs   ExprNodeList // must be set, i.e. cannot be nil
	Rhs   ExprNodeList // must be set, i.e. cannot be nil
	Token ASSIGN_TOKEN
}

func (s AssignStmt) Walk(w *writer.Writer) {
	s.Lhs.Walk(w)
	w.Write(" ")
	w.Write(string(s.Token))
	w.Write(" ")
	s.Rhs.Walk(w)
}

func (s *AssignStmt) SetLhs(xx ...ExprNode) {
	s.Lhs = ExprList(xx)
}

func (s *AssignStmt) SetRhs(xx ...ExprNode) {
	s.Rhs = ExprList(xx)
}

// A GoStmt node represents a go statement.
type GoStmt struct {
	Call CallExpr
}

func (s GoStmt) Walk(w *writer.Writer) {
	w.Write("go ")
	s.Call.Walk(w)
}

type DeferStmt struct {
	Call CallExpr
}

func (s DeferStmt) Walk(w *writer.Writer) {
	w.Write("defer ")
	s.Call.Walk(w)
}

type ReturnStmt struct {
	Result ExprNodeList // can be nil
}

func (s ReturnStmt) Walk(w *writer.Writer) {
	w.Write("return")
	if s.Result != nil {
		w.Write(" ")
		s.Result.Walk(w)
	}
}

// A BranchStmt node represents a break, continue, goto, or fallthrough statement.
type BranchStmt struct {
	Token BRANCH_TOKEN
	Label Ident
}

func (s BranchStmt) Walk(w *writer.Writer) {
	w.Write(string(s.Token))
	if len(s.Label.Name) > 0 {
		w.Write(" ")
		s.Label.Walk(w)
	}
}

type BlockStmt struct {
	List []Stmt
}

func (s BlockStmt) Walk(w *writer.Writer) {
	w.Write("{")
	for _, stmt := range s.List {
		w.Write("\n")
		stmt.Walk(w)
	}
	if len(s.List) > 0 {
		w.Write("\n")
	}
	w.Write("}")
}

func (s *BlockStmt) Add(ss ...Stmt) {
	s.List = append(s.List, ss...)
}

type IfStmt struct {
	Init Stmt
	Cond ExprNode
	Body BlockStmt
	Else Stmt
}

func (s IfStmt) Walk(w *writer.Writer) {
	w.Write("if ")
	if s.Init != nil {
		s.Init.Walk(w)
		w.Write("; ")
	}

	s.Cond.Walk(w)
	w.Write(" ")
	s.Body.Walk(w)

	if s.Else != nil {
		w.Write(" else ")
		s.Else.Walk(w)
	}
}

// A CaseClause represents a case of an expression or type switch statement.
type CaseClause struct {
	List ExprNodeList // list of expressions or types; nil means default case
	Body []Stmt
}

func (s CaseClause) Walk(w *writer.Writer) {
	if s.List != nil {
		w.Write("case ")
		s.List.Walk(w)
		w.Write(":")
	} else {
		w.Write("default:")
	}

	for _, stmt := range s.Body {
		w.Write("\n")
		stmt.Walk(w)
	}
}

type SwitchStmt struct {
	Init Stmt
	Tag  ExprNode
	Body []CaseClause
}

func (s SwitchStmt) Walk(w *writer.Writer) {
	w.Write("switch ")
	if s.Init != nil {
		s.Init.Walk(w)
		w.Write("; ")
	}
	if s.Tag != nil {
		s.Tag.Walk(w)
		w.Write(" ")
	}
	w.Write("{")

	if len(s.Body) > 0 {
		for _, cc := range s.Body {
			w.Write("\n")
			cc.Walk(w)
		}
		w.Write("\n")
	}
	w.Write("}")
}

type TypeSwitchGuard struct {
	Name Ident
	X    ExprNode
}

func (g TypeSwitchGuard) Walk(w *writer.Writer) {
	if len(g.Name.Name) > 0 {
		g.Name.Walk(w)
		w.Write(" := ")
	}
	g.X.Walk(w)
	w.Write(".(type)")
}

type TypeSwitchStmt struct {
	Init  Stmt
	Guard TypeSwitchGuard
	Body  []CaseClause
}

func (s TypeSwitchStmt) Walk(w *writer.Writer) {
	w.Write("switch ")
	if s.Init != nil {
		s.Init.Walk(w)
		w.Write("; ")
	}
	s.Guard.Walk(w)
	w.Write(" {")

	if len(s.Body) > 0 {
		for _, cc := range s.Body {
			w.Write("\n")
			cc.Walk(w)
		}
		w.Write("\n")
	}
	w.Write("}")
}

type CommClause struct {
	Comm Stmt // send or receive statement; nil means default case
	Body []Stmt
}

func (s CommClause) Walk(w *writer.Writer) {
	if s.Comm != nil {
		w.Write("case ")
		s.Comm.Walk(w)
		w.Write(":")
	} else {
		w.Write("default:")
	}

	for _, stmt := range s.Body {
		w.Write("\n")
		stmt.Walk(w)
	}
}

type SelectStmt struct {
	Body []CommClause
}

func (s SelectStmt) Walk(w *writer.Writer) {
	w.Write("select ")
	if len(s.Body) == 0 {
		w.Write("{}")
		return
	}
	w.Write("{")
	for _, comm := range s.Body {
		w.Write("\n")
		comm.Walk(w)
	}
	w.Write("\n}")
}

type ForStmt struct {
	Init Stmt
	Cond ExprNode
	Post Stmt
	Body BlockStmt
}

func (s ForStmt) Walk(w *writer.Writer) {
	w.Write("for ")
	if s.Init == nil && s.Post == nil {
		if s.Cond != nil {
			s.Cond.Walk(w)
			w.Write(" ")
		}
		s.Body.Walk(w)
		return
	}

	s.Init.Walk(w)
	w.Write("; ")
	s.Cond.Walk(w)
	w.Write("; ")
	s.Post.Walk(w)
	w.Write(" ")
	s.Body.Walk(w)
}

type RangeStmt struct {
	Key    ExprNode
	Value  ExprNode
	Define bool
	X      ExprNode
	Body   BlockStmt
}

func (s RangeStmt) Walk(w *writer.Writer) {
	w.Write("for ")
	if s.Key != nil {
		s.Key.Walk(w)

		if s.Value != nil {
			w.Write(", ")
			s.Value.Walk(w)
		}

		if s.Define {
			w.Write(" := ")
		} else {
			w.Write(" = ")
		}
	}

	w.Write("range ")
	s.X.Walk(w)
	w.Write(" ")
	s.Body.Walk(w)
}

func (s *RangeStmt) SetVariables(key, val string) {
	s.Key = Ident{key}
	s.Value = Ident{val}
	s.Define = true
}

func (s *RangeStmt) SetExpression(x ExprNode) {
	s.X = x
}

func (s *RangeStmt) AddStmt(ss ...Stmt) {
	s.Body.List = append(s.Body.List, ss...)
}

func (DeclStmt) stmtNode()       {} // done
func (LabeledStmt) stmtNode()    {}
func (ExprStmt) stmtNode()       {}
func (SendStmt) stmtNode()       {}
func (IncDecStmt) stmtNode()     {} // done
func (AssignStmt) stmtNode()     {} // done
func (GoStmt) stmtNode()         {} // done
func (DeferStmt) stmtNode()      {} // done
func (ReturnStmt) stmtNode()     {} // done
func (BranchStmt) stmtNode()     {} // done
func (BlockStmt) stmtNode()      {} // done
func (IfStmt) stmtNode()         {} // done
func (CaseClause) stmtNode()     {} // done
func (SwitchStmt) stmtNode()     {} // done
func (TypeSwitchStmt) stmtNode() {} // done
func (CommClause) stmtNode()     {} // done
func (SelectStmt) stmtNode()     {} // done
func (ForStmt) stmtNode()        {} // done
func (RangeStmt) stmtNode()      {} // done
