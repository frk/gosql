package golang

import (
	"bytes"
	"testing"

	"github.com/frk/compare"
)

func TestFile(t *testing.T) {
	tests := []struct {
		file File
		want string
	}{{
		file: File{PkgName: "app"},
		want: "package app",
	}, {
		file: File{
			PkgName: "app",
			Imports: []ImportDeclNode{
				ImportDecl{Specs: []ImportSpec{{Path: "fmt"}}},
			},
		},
		want: "package app\n\nimport (\n\"fmt\"\n)",
	}, {
		file: File{
			PkgName: "app",
			Imports: []ImportDeclNode{
				ImportDecl{Specs: []ImportSpec{{Path: "fmt"}}},
			},
			Decls: []TopLevelDeclNode{FuncDecl{
				Name: Ident{"main"},
			}},
		},
		want: "package app\n\nimport (\n\"fmt\"\n)\n\nfunc main() {}",
	}, {
		file: File{
			PkgName: "app",
			Imports: []ImportDeclNode{
				ImportDecl{Specs: []ImportSpec{{Path: "fmt"}}},
			},
			Decls: []TopLevelDeclNode{FuncDecl{
				Name: Ident{"main"},
				Body: BlockStmt{List: []StmtNode{
					ExprStmt{X: CallExpr{
						Fun:  SelectorExpr{X: Ident{"fmt"}, Sel: Ident{"Println"}},
						Args: ArgsList{List: StringLit("Hello, 世界")},
					}},
				}},
			}},
		},
		want: "package app\n\n" +
			"import (\n\"fmt\"\n)\n\n" +
			"func main() {\n" +
			"fmt.Println(\"Hello, 世界\")\n" +
			"}",
	}}

	for _, tt := range tests {
		w := new(bytes.Buffer)

		if err := Write(tt.file, w); err != nil {
			t.Error(err)
		}

		got := w.String()
		if err := compare.Compare(got, tt.want); err != nil {
			t.Error(err)
		}
	}
}
