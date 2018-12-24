package parser

import (
	"go/ast"
)

type Parser struct {
}

type field struct {
	name     string
	typeName string
}

type structDecl struct {
	name   string
	fields []field
}

func NewParser() (*Parser, error) {
	return &Parser{}, nil
}

func (p *Parser) parseStructs(root *ast.File, allowedStructs []string) ([]structDecl, error) {
	v := newVisitor(allowedStructs)
	ast.Walk(v, root)
	return v.structs, nil
}
