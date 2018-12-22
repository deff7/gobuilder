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
	if err != nil {
		return "", nil
	}

	return buf.String(), nil
}
