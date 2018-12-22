package generator

import (
	"bytes"
	"errors"
)

var (
	ErrEmptyFields    = errors.New("no fields are provided")
	ErrEmptyTypeName  = errors.New("no type name is provided")
	ErrEmptyFieldName = errors.New("empty field name")
	ErrEmptyFieldType = errors.New("empty field name")
)

type field struct {
	name     string
	typeName string
}

type Generator struct {
	typeName string
	fields   []field
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
	return "", ErrEmptyFields
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

func (g *Generator) generateBuildValue() (string, error) {
	var buf = new(bytes.Buffer)
	err := buildValueTmpl.Execute(buf, struct {
		StructType string
	}{
		StructType: g.typeName,
	})

	return buf.String(), err
}

func (g *Generator) generateBuildPointer() (string, error) {
	var buf = new(bytes.Buffer)
	err := buildPointerTmpl.Execute(buf, struct {
		StructType string
	}{
		StructType: g.typeName,
	})

	return buf.String(), err
}
