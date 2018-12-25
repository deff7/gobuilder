package parser

import (
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
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := collectTypeName(tc.expr)
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func newASTStringPtr() ast.Expr {
	return &ast.StarExpr{
		X: &ast.Ident{
			Name: "string",
		},
	}
}
