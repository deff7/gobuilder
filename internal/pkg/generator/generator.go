package generator

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"text/template"
)

var (
	ErrEmptyTypeName  = errors.New("no type name is provided")
	ErrEmptyFieldName = errors.New("empty field name")
	ErrEmptyFieldType = errors.New("empty field name")
)

type field struct {
	name     string
	typeName string
}

type Generator struct {
	typeName    string
	packageName string
	fields      []field
}

func NewGenerator(typeName string) (*Generator, error) {
	if typeName == "" {
		return nil, ErrEmptyTypeName
	}

	return &Generator{}, nil
}

func (g *Generator) AddField(fieldName, fieldType string) {
	g.fields = append(g.fields, field{
		name:     fieldName,
		typeName: fieldType,
	})
}

func (g *Generator) Generate() (string, error) {
	var (
		part  string
		err   error
		parts = []string{}
	)

	part, err = g.generateDeclaration()
	if err != nil {
		return "", fmt.Errorf("declaration: %s", err)
	}
	parts = append(parts, part)

	for _, field := range g.fields {
		part, err = g.generateSetMethod(field)
		if err != nil {
			return "", fmt.Errorf("setter method %q: %s", field.name, err)
		}
		parts = append(parts, part)
	}

	part, err = g.generateBuildPointer()
	if err != nil {
		return "", fmt.Errorf("build pointer: %s", err)
	}
	parts = append(parts, part)

	part, err = g.generateBuildValue()
	if err != nil {
		return "", fmt.Errorf("build value: %s", err)
	}
	parts = append(parts, part)

	return strings.Join(parts, "\n"), nil
}

func (g *Generator) generateSetMethod(field field) (string, error) {
	if field.name == "" {
		return "", ErrEmptyFieldName
	}

	if field.typeName == "" {
		return "", ErrEmptyFieldType
	}

	var buf = new(bytes.Buffer)
	err := setMethodTmpl.Execute(buf, struct {
		StructType, FieldName, FieldType string
	}{
		StructType: g.typeName,
		FieldName:  field.name,
		FieldType:  field.typeName,
	})

	return buf.String(), err
}

func (g *Generator) generateByTemplate(tmpl *template.Template) (string, error) {
	var (
		packageName   = g.packageName
		newFuncPrefix = "New"
	)

	if packageName != "" {
		packageName = packageName + "."
		newFuncPrefix = ""
	}

	var buf = new(bytes.Buffer)
	err := tmpl.Execute(buf, struct {
		StructType    string
		PackageName   string
		NewFuncPrefix string
	}{
		StructType:    g.typeName,
		PackageName:   packageName,
		NewFuncPrefix: newFuncPrefix,
	})

	return buf.String(), err

}

func (g *Generator) generateBuildValue() (string, error) {
	return g.generateByTemplate(buildValueTmpl)
}

func (g *Generator) generateBuildPointer() (string, error) {
	return g.generateByTemplate(buildPointerTmpl)
}

func (g *Generator) generateDeclaration() (string, error) {
	return g.generateByTemplate(declarationTmpl)
}
