package golang

import (
	"github.com/frk/gosql/internal/writer"
)

// ConstDecl produces as constant declaration.
//
//	const number = 182
type ConstDecl struct {
	Doc  CommentNode
	Spec ValueSpecNode
}

func (d ConstDecl) Walk(w *writer.Writer) {
	if d.Doc != nil {
		d.Doc.Walk(w)
		w.Write("\n")
	}
	w.Write("const ")
	d.Spec.Walk(w)
}

// VarDecl produces as variable declaration.
//
//	var name = "Jack"
type VarDecl struct {
	Doc  CommentNode
	Spec ValueSpecNode
}

func (d VarDecl) Walk(w *writer.Writer) {
	if d.Doc != nil {
		d.Doc.Walk(w)
		w.Write("\n")
	}
	w.Write("var ")
	d.Spec.Walk(w)
}

// TypeDecl produces as type declaration.
//
//	type T struct {
//		// ...
//	}
type TypeDecl struct {
	Doc  CommentNode
	Spec TypeSpecNode
}

func (d TypeDecl) Walk(w *writer.Writer) {
	if d.Doc != nil {
		d.Doc.Walk(w)
		w.Write("\n")
	}
	w.Write("type ")
	d.Spec.Walk(w)
}

// FuncDecl produces a func declaration.
//
//	func F() {
//		// ...
//	}
type FuncDecl struct {
	Doc  CommentNode
	Name Ident
	Type FuncType // the function signature
	Body BlockStmt
}

func (d FuncDecl) Walk(w *writer.Writer) {
	if d.Doc != nil {
		d.Doc.Walk(w)
		w.Write("\n")
	}
	w.Write("func ")

	d.Name.Walk(w)
	d.Type.Walk(w)

	w.Write(" ")
	d.Body.Walk(w)
}

func (d *FuncDecl) AddStmt(ss ...Stmt) {
	d.Body.List = append(d.Body.List, ss...)
}

// MethodDecl produces a method declaration.
//
//	func (t *T) M() {
//		// ...
//	}
type MethodDecl struct {
	Doc  CommentNode
	Recv RecvParam
	Name Ident
	Type FuncType // the function signature
	Body BlockStmt
}

func (d MethodDecl) Walk(w *writer.Writer) {
	if d.Doc != nil {
		d.Doc.Walk(w)
		w.Write("\n")
	}
	w.Write("func ")

	w.Write("(")
	d.Recv.Walk(w)
	w.Write(") ")

	d.Name.Walk(w)
	d.Type.Walk(w)

	w.Write(" ")
	d.Body.Walk(w)
}

// implements TopLevelDeclNode
func (ConstDecl) topLevelDeclNode()  {}
func (VarDecl) topLevelDeclNode()    {}
func (TypeDecl) topLevelDeclNode()   {}
func (FuncDecl) topLevelDeclNode()   {}
func (MethodDecl) topLevelDeclNode() {}

// implements DeclNode
func (ConstDecl) declNode() {}
func (VarDecl) declNode()   {}
func (TypeDecl) declNode()  {}
