package golang

import (
	"github.com/frk/gosql/internal/writer"
)

type Decl interface {
	Node
	declNode()
}

type DECL_TOKEN string

const (
	DECL_CONST DECL_TOKEN = "const"
	DECL_TYPE  DECL_TOKEN = "type"
	DECL_VAR   DECL_TOKEN = "var"
)

type GenDecl struct {
	Token DECL_TOKEN
	Specs SpecNode
	// Doc
}

func (d GenDecl) Walk(w *writer.Writer) {
	w.Write(string(d.Token))
	w.Write(" ")
	d.Specs.Walk(w)
}

type ImportDecl []ImportSpec

func (d ImportDecl) Walk(w *writer.Writer) {
	if len(d) == 0 {
		return
	}

	w.Write("import (\n")
	for _, spec := range d {
		spec.Walk(w)
		w.Write("\n")
	}
	w.Write(")")
}

func (d *ImportDecl) Add(path string) {
	for _, spec := range *d {
		if string(spec.Path) == path {
			return
		}
	}

	*d = append(*d, ImportSpec{Path: String(path)})
}

func (d *ImportDecl) NewLine() {
	*d = append(*d, ImportSpec{NewLine: true})
}

type FuncDecl struct {
	Recv RecvParam // empty for functions
	Name Ident
	Type FuncType
	Body BlockStmt
	// Doc
}

func (d FuncDecl) Walk(w *writer.Writer) {
	w.Write("func ")
	if d.Recv != (RecvParam{}) { // if not empty
		w.Write("(")
		d.Recv.Walk(w)
		w.Write(") ")
	}

	d.Name.Walk(w)
	d.Type.Walk(w)

	w.Write(" ")
	d.Body.Walk(w)
}

func (d *FuncDecl) AddStmt(ss ...Stmt) {
	d.Body.List = append(d.Body.List, ss...)
}

func (GenDecl) declNode()    {}
func (FuncDecl) declNode()   {}
func (ImportDecl) declNode() {}
