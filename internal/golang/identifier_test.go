package golang

import (
	"bytes"
	"testing"

	"github.com/frk/compare"
)

func TestIdent(t *testing.T) {
	tests := []struct {
		id   Ident
		want string
	}{{
		id:   Ident{},
		want: "",
	}, {
		id:   Ident{"Name"},
		want: "Name",
	}, {
		id:   Ident{"value"},
		want: "value",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.id, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestIdentList(t *testing.T) {
	tests := []struct {
		list IdentList
		want string
	}{{
		list: IdentList{},
		want: "",
	}, {
		list: IdentList{{"foo"}, {"bar"}},
		want: "foo, bar",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.list, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}

func TestQualifiedIdent(t *testing.T) {
	tests := []struct {
		id   QualifiedIdent
		want string
	}{{
		id:   QualifiedIdent{"abc", "Name"},
		want: "abc.Name",
	}, {
		id:   QualifiedIdent{"packagename", "SomeType"},
		want: "packagename.SomeType",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.id, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}
