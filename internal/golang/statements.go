package golang

import (
	"github.com/frk/gosql/internal/writer"
)

type BranchToken string

const (
	BranchBreak BranchToken = "break"
	BranchCont  BranchToken = "continue"
	BranchGoto  BranchToken = "goto"
	BranchFall  BranchToken = "fallthrough"
)

type IncDecToken string

const (
	IncDecIncrement IncDecToken = "++"
	IncDecDecrement IncDecToken = "--"
)

type AssignToken string

const (
	Assign       AssignToken = "="
	AssignAdd    AssignToken = "+="
	AssignSub    AssignToken = "-="
	AssignMul    AssignToken = "*="
	AssignQuo    AssignToken = "%="
	AssignRem    AssignToken = "%="
	AssignAnd    AssignToken = "&="
	AssignOr     AssignToken = "|="
	AssignXOr    AssignToken = "^="
	AssignShl    AssignToken = "<<="
	AssignShr    AssignToken = ">>="
	AssignAndNot AssignToken = "&^="
	AssignDefine AssignToken = ":="
)

// DeclStmt node produces a declaration statement.
type DeclStmt struct {
	Decl DeclNode // the declaration
}

func (s DeclStmt) Walk(w *writer.Writer) {
	s.Decl.Walk(w)
}

// LabeledStmt node produces a labeled statement.
type LabeledStmt struct {
	Label Ident    // the label
	Stmt  StmtNode // the statement
}

func (s LabeledStmt) Walk(w *writer.Writer) {
	s.Label.Walk(w)
	w.Write(":\n")
	s.Stmt.Walk(w)
}

// ExprStmt node produces a (stand-alone) expression statement.
type ExprStmt struct {
	X ExprNode // the expression
}

func (s ExprStmt) Walk(w *writer.Writer) {
	s.X.Walk(w)
}

// SendStmt produces a channel send statement.
type SendStmt struct {
	Chan  ExprNode // the channel to which to send the value
	Value ExprNode // the value to send to the channel
}

func (s SendStmt) Walk(w *writer.Writer) {
	s.Chan.Walk(w)
	w.Write(" <- ")
	s.Value.Walk(w)
}

// IncDecStmt produces an increment / decrement statement.
type IncDecStmt struct {
	X     ExprNode
	Token IncDecToken
}

func (s IncDecStmt) Walk(w *writer.Writer) {
	s.X.Walk(w)
	w.Write(string(s.Token))
}

// AssignStmt node produces an assignment statement.
type AssignStmt struct {
	Lhs   ExprListNode // left hand side operand
	Rhs   ExprListNode // right hand side operand
	Token AssignToken  // assignment token
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

// GoStmt node produces a go statement.
type GoStmt struct {
	Call CallExpr // the call
}

func (s GoStmt) Walk(w *writer.Writer) {
	w.Write("go ")
	s.Call.Walk(w)
}

// DeferStmt produces a defer statement.
type DeferStmt struct {
	Call CallExpr // the call
}

func (s DeferStmt) Walk(w *writer.Writer) {
	w.Write("defer ")
	s.Call.Walk(w)
}

// ReturnStmt node produces a return statement.
type ReturnStmt struct {
	Result ExprListNode // the expression to return; or nil
}

func (s ReturnStmt) Walk(w *writer.Writer) {
	w.Write("return")
	if s.Result != nil {
		w.Write(" ")
		s.Result.Walk(w)
	}
}

// BranchStmt node produces a break, continue, goto, or fallthrough statement.
type BranchStmt struct {
	Token BranchToken // the type of the branch
	Label Ident       // optional label
}

func (s BranchStmt) Walk(w *writer.Writer) {
	w.Write(string(s.Token))
	if len(s.Label.Name) > 0 {
		w.Write(" ")
		s.Label.Walk(w)
	}
}

// BlockStmt node produces a block statement.
type BlockStmt struct {
	List []StmtNode // statements inside the block
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

func (s *BlockStmt) Add(ss ...StmtNode) {
	s.List = append(s.List, ss...)
}

// IfStmt produces an if statement.
type IfStmt struct {
	Init StmtNode  // an optional simple statement
	Cond ExprNode  // the condition of the if statement
	Body BlockStmt // the statments to be executed when the condition is met
	Else ElseNode  // the optional else, or else-if, node
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

// SwitchStmt produces a switch statment.
type SwitchStmt struct {
	Init  StmtNode     // an optional simple statement
	X     ExprNode     // the switch expression (optional)
	Cases []CaseClause // a list of the switch cases
}

func (s SwitchStmt) Walk(w *writer.Writer) {
	w.Write("switch ")
	if s.Init != nil {
		s.Init.Walk(w)
		w.Write("; ")
	}
	if s.X != nil {
		s.X.Walk(w)
		w.Write(" ")
	}
	w.Write("{")

	if len(s.Cases) > 0 {
		for _, cc := range s.Cases {
			w.Write("\n")
			cc.Walk(w)
		}
		w.Write("\n")
	}
	w.Write("}")
}

// CaseClause produces a case clause in an expression switch statement.
type CaseClause struct {
	List ExprListNode // list of expressions; nil means default case
	Body []StmtNode   // list of the case's statements to be executed
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

// TypeSwitchStmt produces a type swtich.
type TypeSwitchStmt struct {
	Init  StmtNode         // an optional simple statement
	Guard TypeSwitchGuard  // the type switch guard
	Cases []TypeCaseClause // list of the switch's cases
}

func (s TypeSwitchStmt) Walk(w *writer.Writer) {
	w.Write("switch ")
	if s.Init != nil {
		s.Init.Walk(w)
		w.Write("; ")
	}
	s.Guard.Walk(w)
	w.Write(" {")

	if len(s.Cases) > 0 {
		for _, cc := range s.Cases {
			w.Write("\n")
			cc.Walk(w)
		}
		w.Write("\n")
	}
	w.Write("}")
}

// TypeSwitchGuard produces the special swtich expression of a type swtich.
type TypeSwitchGuard struct {
	Name Ident    // optional identifier for short-variable-declaration
	X    ExprNode // the primary expression
}

func (g TypeSwitchGuard) Walk(w *writer.Writer) {
	if len(g.Name.Name) > 0 {
		g.Name.Walk(w)
		w.Write(" := ")
	}
	g.X.Walk(w)
	w.Write(".(type)")
}

// TypeCaseClause produces a case clause in a type switch statement.
type TypeCaseClause struct {
	List TypeListNode // list of types; nil means default case
	Body []StmtNode   // list of the case's statements to be executed
}

func (s TypeCaseClause) Walk(w *writer.Writer) {
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

// SelectStmt produces a select statement.
type SelectStmt struct {
	Cases []CommClause // list of the select's cases
}

func (s SelectStmt) Walk(w *writer.Writer) {
	w.Write("select ")
	if len(s.Cases) == 0 {
		w.Write("{}")
		return
	}
	w.Write("{")
	for _, comm := range s.Cases {
		w.Write("\n")
		comm.Walk(w)
	}
	w.Write("\n}")
}

// CommClause produces a case clause in a select statement.
type CommClause struct {
	Comm StmtNode   // the communication operation; nil means default case
	Body []StmtNode // list of the case's statements to be executed
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

// ForStmt produces a for statement.
type ForStmt struct {
	Clause ForClauseNode // the clause of the for statement [condition | for-clause | range-clause]
	Body   BlockStmt     // the code block to be executed by the for statement
}

func (s ForStmt) Walk(w *writer.Writer) {
	w.Write("for ")
	if s.Clause != nil {
		s.Clause.Walk(w)
		w.Write(" ")
	}
	s.Body.Walk(w)
}

// ForCondition produces the expression to be evaluated before each iteration in a for statement.
type ForCondition struct {
	X ExprNode // the boolean expression to be evaluated
}

func (f ForCondition) Walk(w *writer.Writer) {
	if f.X != nil {
		f.X.Walk(w)
	}
}

// ForClause produces the init-cond-post clause in a for statement.
type ForClause struct {
	Init StmtNode // an optional simple initialization statement, e.g. assignment
	Cond ExprNode // the boolean expression to be evaluated
	Post StmtNode // an optional simple post statement, e.g. increment/decrement
}

func (f ForClause) Walk(w *writer.Writer) {
	if f.Init != nil {
		f.Init.Walk(w)
	}
	w.Write("; ")
	if f.Cond != nil {
		f.Cond.Walk(w)
	}
	w.Write("; ")
	if f.Post != nil {
		f.Post.Walk(w)
	}
}

// ForRangeClause produces a range clause in a for statement.
type ForRangeClause struct {
	Key    ExprNode // the optional iteration variable
	Value  ExprNode // the 2nd optional iteration variable
	X      ExprNode // the range expression
	Define bool     // indicates whether the iteration variables should be defined or assigned
}

func (f ForRangeClause) Walk(w *writer.Writer) {
	if f.Key != nil {
		f.Key.Walk(w)
	} else if f.Value != nil {
		w.Write("_")
	}

	if f.Value != nil {
		w.Write(", ")
		f.Value.Walk(w)
	}

	if f.Key != nil || f.Value != nil {
		if f.Define {
			w.Write(" := ")
		} else {
			w.Write(" = ")
		}
	}

	w.Write("range ")
	f.X.Walk(w)
}

func (DeclStmt) stmtNode()       {}
func (LabeledStmt) stmtNode()    {}
func (ExprStmt) stmtNode()       {}
func (SendStmt) stmtNode()       {}
func (IncDecStmt) stmtNode()     {}
func (AssignStmt) stmtNode()     {}
func (GoStmt) stmtNode()         {}
func (DeferStmt) stmtNode()      {}
func (ReturnStmt) stmtNode()     {}
func (BranchStmt) stmtNode()     {}
func (BlockStmt) stmtNode()      {}
func (IfStmt) stmtNode()         {}
func (CaseClause) stmtNode()     {}
func (SwitchStmt) stmtNode()     {}
func (TypeSwitchStmt) stmtNode() {}
func (CommClause) stmtNode()     {}
func (SelectStmt) stmtNode()     {}
func (ForStmt) stmtNode()        {}

func (IfStmt) elseNode()    {}
func (BlockStmt) elseNode() {}

func (ForCondition) forClauseNode()   {}
func (ForClause) forClauseNode()      {}
func (ForRangeClause) forClauseNode() {}
