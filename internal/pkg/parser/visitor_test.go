package parser

import (
	"regexp"
	"testing"

	"go/ast"
)

func TestCollectTypeName(t *testing.T) {
	for _, tc := range []struct {
		name string
		expr ast.Expr
		want string
	}{
		{
			name: "simple type",
			expr: &ast.Ident{
				Name: "string",
			},
			want: "string",
		},
		{
			name: "pointer type",
			expr: newASTStringPtr(),
			want: "*string",
		},
		{
			name: "slice type",
			expr: &ast.ArrayType{
				Len: nil,
				Elt: newASTStringPtr(),
			},
			want: "[]*string",
		},
		{
			name: "array type",
			expr: &ast.ArrayType{
				Len: &ast.BasicLit{Value: "3"},
				Elt: newASTStringPtr(),
			},
			want: "[3]*string",
		},
		{
			name: "slice of types from imported package",
			expr: &ast.ArrayType{
				Len: nil,
				Elt: &ast.StarExpr{
					X: &ast.SelectorExpr{
						Sel: &ast.Ident{Name: "Bar"},
						X:   &ast.Ident{Name: "imported"},
					},
				},
			},
			want: "[]*imported.Bar",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := collectTypeName(tc.expr)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestCheckStructName(t *testing.T) {
	t.Run("with provided regexp", func(t *testing.T) {
		v := &visitor{}
		v.allowedStructs = []*regexp.Regexp{regexp.MustCompile(".oo")}

		t.Run("with name that satisfy pattern returns true", func(t *testing.T) {
			got := v.checkStructName("Foo")

			if got != true {
				t.Error("expect true")
			}
		})

		t.Run("with name that not satisfy pattern returns false", func(t *testing.T) {
			got := v.checkStructName("Bar")

			if got != false {
				t.Error("expect false")
			}
		})
	})
}

func newASTStringPtr() ast.Expr {
	return &ast.StarExpr{
		X: &ast.Ident{
			Name: "string",
		},
	}
}
