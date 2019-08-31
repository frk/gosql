package typesutil

import (
	"go/types"
	"testing"

	"github.com/frk/gosql/internal/testutil"
)

var tdata = testutil.ParseTestdata("testdata")

func TestIsError(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "IsErrorTest1", want: false},
		{name: "IsErrorTest2", want: false},
		{name: "IsErrorTest3", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			named := testutil.FindNamedType(tt.name, tdata)
			if named == nil {
				t.Errorf("%q named type not found", tt.name)
				return
			}

			typ := (named.Underlying().(*types.Struct)).Field(0).Type()
			got := IsError(typ)
			if got != tt.want {
				t.Errorf("got=%t; want=%t", got, tt.want)
			}
		})
	}
}

func TestIsEmptyInterface(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "IsEmptyInterfaceTest1", want: false},
		{name: "IsEmptyInterfaceTest2", want: false},
		{name: "IsEmptyInterfaceTest3", want: false},
		{name: "IsEmptyInterfaceTest4", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			named := testutil.FindNamedType(tt.name, tdata)
			if named == nil {
				t.Errorf("%q named type not found", tt.name)
				return
			}

			typ := (named.Underlying().(*types.Struct)).Field(0).Type()
			got := IsEmptyInterface(typ)
			if got != tt.want {
				t.Errorf("got=%t; want=%t", got, tt.want)
			}
		})
	}
}

func TestIsTime(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "IsTimeTest1", want: false},
		{name: "IsTimeTest2", want: false},
		{name: "IsTimeTest3", want: true},
		{name: "IsTimeTest4", want: true},
		{name: "IsTimeTest5", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			named := testutil.FindNamedType(tt.name, tdata)
			if named == nil {
				t.Errorf("%q named type not found", tt.name)
				return
			}

			typ := (named.Underlying().(*types.Struct)).Field(0).Type()
			got := IsTime(typ)
			if got != tt.want {
				t.Errorf("got=%t; want=%t", got, tt.want)
			}
		})
	}
}

func TestIsSqlDriverValue(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "IsSqlDriverValueTest1", want: false},
		{name: "IsSqlDriverValueTest2", want: false},
		{name: "IsSqlDriverValueTest3", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			named := testutil.FindNamedType(tt.name, tdata)
			if named == nil {
				t.Errorf("%q named type not found", tt.name)
				return
			}

			typ := (named.Underlying().(*types.Struct)).Field(0).Type()
			got := IsSqlDriverValue(typ)
			if got != tt.want {
				t.Errorf("got=%t; want=%t", got, tt.want)
			}
		})
	}
}

func TestIsDirectiveValue(t *testing.T) {
	tests := []struct {
		name  string
		ident string
		want  bool
	}{
		{name: "IsDirectiveTest1", ident: "Column", want: false},
		{name: "IsDirectiveTest2", ident: "Column", want: false},
		{name: "IsDirectiveTest3", ident: "Column", want: false},
		{name: "IsDirectiveTest4", ident: "Relation", want: false},
		{name: "IsDirectiveTest5", ident: "Column", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			named := testutil.FindNamedType(tt.name, tdata)
			if named == nil {
				t.Errorf("%q named type not found", tt.name)
				return
			}

			typ := (named.Underlying().(*types.Struct)).Field(0).Type()
			got := IsDirective(tt.ident, typ)
			if got != tt.want {
				t.Errorf("got=%t; want=%t", got, tt.want)
			}
		})
	}
}

func TestImplementsScanner(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "ImplementsScannerTest1", want: false},
		{name: "ImplementsScannerTest2", want: false},
		{name: "ImplementsScannerTest3", want: false},
		{name: "ImplementsScannerTest4", want: false},
		{name: "ImplementsScannerTest5", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			named := testutil.FindNamedType(tt.name, tdata)
			if named == nil {
				t.Errorf("%q named type not found", tt.name)
				return
			}

			got := ImplementsScanner(named)
			if got != tt.want {
				t.Errorf("got=%t; want=%t", got, tt.want)
			}
		})
	}
}

func TestImplementsValuer(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "ImplementsValuerTest1", want: false},
		{name: "ImplementsValuerTest2", want: false},
		{name: "ImplementsValuerTest3", want: false},
		{name: "ImplementsValuerTest4", want: false},
		{name: "ImplementsValuerTest5", want: false},
		{name: "ImplementsValuerTest6", want: false},
		{name: "ImplementsValuerTest7", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			named := testutil.FindNamedType(tt.name, tdata)
			if named == nil {
				t.Errorf("%q named type not found", tt.name)
				return
			}

			got := ImplementsValuer(named)
			if got != tt.want {
				t.Errorf("got=%t; want=%t", got, tt.want)
			}
		})
	}
}

func TestImplementsAfterScanner(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "ImplementsAfterScannerTest1", want: false},
		{name: "ImplementsAfterScannerTest2", want: false},
		{name: "ImplementsAfterScannerTest3", want: false},
		{name: "ImplementsAfterScannerTest4", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			named := testutil.FindNamedType(tt.name, tdata)
			if named == nil {
				t.Errorf("%q named type not found", tt.name)
				return
			}

			got := ImplementsAfterScanner(named)
			if got != tt.want {
				t.Errorf("got=%t; want=%t", got, tt.want)
			}
		})
	}
}

func TestImplementsErrorHandler(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "ImplementsErrorHandlerTest1", want: false},
		{name: "ImplementsErrorHandlerTest2", want: false},
		{name: "ImplementsErrorHandlerTest3", want: false},
		{name: "ImplementsErrorHandlerTest4", want: false},
		{name: "ImplementsErrorHandlerTest5", want: false},
		{name: "ImplementsErrorHandlerTest6", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			named := testutil.FindNamedType(tt.name, tdata)
			if named == nil {
				t.Errorf("%q named type not found", tt.name)
				return
			}

			got := ImplementsErrorHandler(named)
			if got != tt.want {
				t.Errorf("got=%t; want=%t", got, tt.want)
			}
		})
	}
}

func TestImplementsJSONMarshaler(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "ImplementsJSONMarshalerTest1", want: false},
		{name: "ImplementsJSONMarshalerTest2", want: false},
		{name: "ImplementsJSONMarshalerTest3", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			named := testutil.FindNamedType(tt.name, tdata)
			if named == nil {
				t.Errorf("%q named type not found", tt.name)
				return
			}

			got := ImplementsJSONMarshaler(named)
			if got != tt.want {
				t.Errorf("got=%t; want=%t", got, tt.want)
			}
		})
	}
}

func TestImplementsJSONUnmarshaler(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "ImplementsJSONUnmarshalerTest1", want: false},
		{name: "ImplementsJSONUnmarshalerTest2", want: false},
		{name: "ImplementsJSONUnmarshalerTest3", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			named := testutil.FindNamedType(tt.name, tdata)
			if named == nil {
				t.Errorf("%q named type not found", tt.name)
				return
			}

			got := ImplementsJSONUnmarshaler(named)
			if got != tt.want {
				t.Errorf("got=%t; want=%t", got, tt.want)
			}
		})
	}
}

func TestImplementsXMLMarshaler(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "ImplementsXMLMarshalerTest1", want: false},
		{name: "ImplementsXMLMarshalerTest2", want: false},
		{name: "ImplementsXMLMarshalerTest3", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			named := testutil.FindNamedType(tt.name, tdata)
			if named == nil {
				t.Errorf("%q named type not found", tt.name)
				return
			}

			got := ImplementsXMLMarshaler(named)
			if got != tt.want {
				t.Errorf("got=%t; want=%t", got, tt.want)
			}
		})
	}
}

func TestImplementsXMLUnmarshaler(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "ImplementsXMLUnmarshalerTest1", want: false},
		{name: "ImplementsXMLUnmarshalerTest2", want: false},
		{name: "ImplementsXMLUnmarshalerTest3", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			named := testutil.FindNamedType(tt.name, tdata)
			if named == nil {
				t.Errorf("%q named type not found", tt.name)
				return
			}

			got := ImplementsXMLUnmarshaler(named)
			if got != tt.want {
				t.Errorf("got=%t; want=%t", got, tt.want)
			}
		})
	}
}
