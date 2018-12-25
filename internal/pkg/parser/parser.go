package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
)

type Parser struct {
}

type Field struct {
	Name     string
	TypeName string
}

type StructDecl struct {
	Name   string
	Fields []Field
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(r io.Reader, allowedStructs []string) (string, []StructDecl, error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "", r, 0)
	if err != nil {
		return "", nil, err
	}

	packageName := astFile.Name.Name
	structs, err := p.parseStructs(astFile, allowedStructs)
	return packageName, structs, err
}

func (p *Parser) parseStructs(root *ast.File, allowedStructs []string) ([]StructDecl, error) {
	v := newVisitor(allowedStructs)
	ast.Walk(v, root)
	return v.structs, nil
}
