package golang

import (
	"github.com/frk/gosql/internal/writer"
)

type Decl interface {
	Node
	declNode()
}

type GENDECL_TOKEN string

const (
	GENDECL_CONST GENDECL_TOKEN = "const"
	GENDECL_TYPE  GENDECL_TOKEN = "type"
	GENDECL_VAR   GENDECL_TOKEN = "var"
)

type GenDecl struct {
	Token GENDECL_TOKEN
	Specs []Spec
	// Doc
}

func (d GenDecl) Walk(w *writer.Writer) {
	w.Write(string(d.Token))
	w.Write(" ")

	withParens := len(d.Specs) > 1
	if withParens {
		w.Write("(\n")
	}
	d.Specs[0].Walk(w)
	for _, spec := range d.Specs[1:] {
		w.Write("\n")
		spec.Walk(w)
	}
	if withParens {
		w.Write("\n)")
	}
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
	Recv ParamList // empty for functions, 1 elem for methods
	Name Ident
	Type FuncType
	Body BlockStmt
	// Doc
}

func (d FuncDecl) Walk(w *writer.Writer) {
	w.Write("func ")
	if len(d.Recv) > 0 {
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

func (d FuncDecl) WithReceiver(typ, ident string, star bool) FuncDecl {
	param := Param{Names: []Ident{{ident}}, Type: Ident{typ}}
	if star {
		param.Type = StarExpr{X: param.Type}
	}
	d.Recv = ParamList{param}
	return d
}
