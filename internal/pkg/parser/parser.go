package parser

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
)

type Parser struct {
	unexported         bool
	invertStructsMatch bool
}

type Field struct {
	Name     string
	TypeName string
}

type StructDecl struct {
	Name   string
	Fields []Field
}

func NewParser(unexported, invertMatch bool) *Parser {
	return &Parser{
		unexported:         unexported,
		invertStructsMatch: invertMatch,
	}
}

func (p *Parser) Parse(r io.Reader, allowedStructs []string) (map[string][]StructDecl, error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, "", r, 0)
	if err != nil {
		return nil, err
	}

	packageName := astFile.Name.Name
	structs, err := p.parseStructs(astFile, allowedStructs)
	return map[string][]StructDecl{
		packageName: structs,
	}, err
}

func (p *Parser) ParseDir(path string, allowedStructs []string) (map[string][]StructDecl, error) {
	fset := token.NewFileSet()
	astPackages, err := parser.ParseDir(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	result := map[string][]StructDecl{}

	for packageName, astPackage := range astPackages {
		structs := []StructDecl{}

		for _, astFile := range astPackage.Files {
			s, err := p.parseStructs(astFile, allowedStructs)
			if err != nil {
				return nil, err
			}

			structs = append(structs, s...)
		}

		result[packageName] = structs
	}
	return result, nil
}

func (p *Parser) parseStructs(root *ast.File, allowedStructs []string) ([]StructDecl, error) {
	v, err := newVisitor(allowedStructs, p.unexported, p.invertStructsMatch)
	if err != nil {
		return nil, err
	}
	ast.Walk(v, root)
	return v.structs, nil
}
