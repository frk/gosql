package golang

import (
	"github.com/frk/gosql/internal/writer"
)

// ConstDecl produces as constant declaration.
//
//	const number = 182
type ConstDecl struct {
	Doc  CommentNode   // associated documentation
	Spec ValueSpecNode // the value spec(s)
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
	Doc  CommentNode   // associated documentation
	Spec ValueSpecNode // the value spec(s)
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
	Doc  CommentNode  // associated documentation
	Spec TypeSpecNode // the type spec(s)
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
	Doc  CommentNode // associated documentation
	Name Ident       // the function's name
	Type Signature   // the function signature
	Body BlockStmt   // the body of the function
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

func (d *FuncDecl) AddStmt(ss ...StmtNode) {
	d.Body.List = append(d.Body.List, ss...)
}

// MethodDecl produces a method declaration.
//
//	func (t *T) M() {
//		// ...
//	}
type MethodDecl struct {
	Doc  CommentNode // associated documentation
	Recv RecvParam   // the receiver parameter
	Name Ident       // method name
	Type Signature   // the function signature
	Body BlockStmt   // the function body
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

// RecvParam produces the receiver parameter section of a method declaration.
type RecvParam struct {
	Name Ident        // name of the receiver
	Type RecvTypeNode // type of the receiver
}

func (p RecvParam) Walk(w *writer.Writer) {
	if len(p.Name.Name) > 0 {
		p.Name.Walk(w)
		w.Write(" ")
	}
	p.Type.Walk(w)
}

// PointerRecvType produces a pointer type in the receiver parameter
// section of a method declaration.
type PointerRecvType struct {
	Type string // the base type
}

func (p PointerRecvType) Walk(w *writer.Writer) {
	w.Write("*")
	w.Write(p.Type)
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

// implements RecvTypeNode
func (PointerRecvType) recvTypeNode()              {}
func (PointerRecvType) typeNode()                  {}
func (t PointerRecvType) typeListNode() []TypeNode { return []TypeNode{t} }
