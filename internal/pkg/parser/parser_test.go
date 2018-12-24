package parser

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

func TestNewParser(t *testing.T) {
	p, err := NewParser()

	if p == nil {
		t.Error("must create new parser instance")
	}

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestParseStructsNew(t *testing.T) {
	for _, tc := range []struct {
		name           string
		allowedStructs []string
		want           []structDecl
	}{
		{
			name:           "with specified struct name expect parse only this structure",
			allowedStructs: []string{"Second"},
			want: []structDecl{
				{
					name:   "Second",
					fields: []field{newField("String", "string")},
				},
			},
		},
		{
			name:           "with specified several struct names expect parse these structs",
			allowedStructs: []string{"First", "Second"},
			want: []structDecl{
				{
					name:   "First",
					fields: []field{newField("Number", "int")},
				},
				{
					name:   "Second",
					fields: []field{newField("String", "string")},
				},
			},
		},
		{
			name:           "without specified struct names expect parse all structs",
			allowedStructs: []string{},
			want: []structDecl{
				{
					name:   "First",
					fields: []field{newField("Number", "int")},
				},
				{
					name:   "Second",
					fields: []field{newField("String", "string")},
				},
				{
					name:   "Third",
					fields: []field{newField("Float", "float64")},
				},
			},
		},
	} {
		var (
			p    = newParser()
			file = newASTFile()
		)
		got, err := p.parseStructs(file, tc.allowedStructs)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
		if !reflect.DeepEqual(tc.want, got) {
			t.Fatalf("want %v got %v", tc.want, got)
		}
	}
}

func newParser() *Parser {
	return &Parser{}
}

var exampleSrc = `package foo

type First struct {
	Number int
}

type Second struct {
	String string
	Complicated struct {
		A int
	}
}

type Third struct {
	Float float64
}`

func newASTFile() *ast.File {
	fset := token.NewFileSet()
	buf := bytes.NewBufferString(exampleSrc)
	file, err := parser.ParseFile(fset, "", buf, 0)
	if err != nil {
		panic(err)
	}
	return file
}

func newField(name, typeName string) field {
	return field{name: name, typeName: typeName}
}
