package golang

import (
	"github.com/frk/gosql/internal/x/writer"
)

// ValueSpec produces a value specification.
type ValueSpec struct {
	Doc     CommentNode   // associated documentation
	Names   IdentListNode // one or more identifiers for the values
	Type    TypeNode      // the type of the values; can be nil if the type is to be infered
	Values  ExprListNode  // list of expressions that produce the values
	Comment LineComment   // trailing line comment
}

func (s ValueSpec) Walk(w *writer.Writer) {
	if s.Doc != nil {
		s.Doc.Walk(w)
	}

	s.Names.Walk(w)
	if s.Type != nil {
		w.Write(" ")
		s.Type.Walk(w)
	}

	if s.Values != nil {
		w.Write(" = ")

		vals := s.Values.exprListNode()
		vals[0].Walk(w)
		for _, v := range vals[1:] {
			w.Write(", ")
			v.Walk(w)
		}
	}
	s.Comment.Walk(w)
}

// ValueSpecList produces a list of one or more value specification in parentheses.
type ValueSpecList []ValueSpec

func (ls ValueSpecList) Walk(w *writer.Writer) {
	withParens := len(ls) > 1
	if withParens {
		w.Write("(\n")
	}

	ls[0].Walk(w)
	for _, n := range ls[1:] {
		w.Write("\n")
		n.Walk(w)
	}

	if withParens {
		w.Write("\n)")
	}
}

// TypeSpec produces a type specification.
type TypeSpec struct {
	Doc     CommentNode // associated documentation
	Name    Ident       // the type's identifier
	Alias   bool        // if set to true the TypeSpec will produce an alias declaration.
	Type    TypeNode    // Ident, QualifiedIdent, or any of the XxxTypes
	Comment LineComment // trailing comment
}

func (s TypeSpec) Walk(w *writer.Writer) {
	if s.Doc != nil {
		s.Doc.Walk(w)
	}
	s.Name.Walk(w)
	if s.Alias {
		w.Write(" = ")
	} else {
		w.Write(" ")
	}
	s.Type.Walk(w)
	s.Comment.Walk(w)
}

// TypeSpecList produces a list of one or more type specification in parentheses.
type TypeSpecList []TypeSpec

func (list TypeSpecList) Walk(w *writer.Writer) {
	withParens := len(list) > 1
	if withParens {
		w.Write("(\n")
	}

	if len(list) > 0 {
		list[0].Walk(w)
		for _, n := range list[1:] {
			w.Write("\n")
			n.Walk(w)
		}
	}

	if withParens {
		w.Write("\n)")
	}
}

// implements ValueSpecNode
func (ValueSpec) valueSpecNode()     {}
func (ValueSpecList) valueSpecNode() {}

// implements TypeSpecNode
func (TypeSpec) typeSpecNode()     {}
func (TypeSpecList) typeSpecNode() {}
