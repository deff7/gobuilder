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
		for _, tc := range []struct {
			name        string
			structName  string
			want        bool
			invertMatch bool
		}{
			{
				name:       "with name that satisfy pattern returns true",
				structName: "Foo",
				want:       true,
			},
			{
				name:       "with name that not satisfy pattern returns false",
				structName: "Bar",
				want:       false,
			},
			{
				name:        "with invert match with name that satisfy pattern returns false",
				structName:  "Foo",
				want:        false,
				invertMatch: true,
			},
			{
				name:        "with invert match with name that not satisfy pattern returns true",
				structName:  "Bar",
				want:        true,
				invertMatch: true,
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				v := &visitor{}
				v.allowedStructs = []*regexp.Regexp{regexp.MustCompile(".oo")}
				v.invertStructsMatch = tc.invertMatch
				got := v.checkStructName(tc.structName)

				if got != tc.want {
					t.Errorf("expect %T", tc.want)
				}
			})
		}

	})

}

func newASTStringPtr() ast.Expr {
	return &ast.StarExpr{
		X: &ast.Ident{
			Name: "string",
		},
	}
}
