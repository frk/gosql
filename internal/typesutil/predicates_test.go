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

			v := (named.Underlying().(*types.Struct)).Field(0)
			got := IsError(v)
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

			v := (named.Underlying().(*types.Struct)).Field(0)
			got := IsEmptyInterface(v)
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

			v := (named.Underlying().(*types.Struct)).Field(0)
			got := IsTime(v)
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

			v := (named.Underlying().(*types.Struct)).Field(0)
			got := IsSqlDriverValue(v)
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
