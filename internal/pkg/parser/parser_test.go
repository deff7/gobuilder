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

func TestParseStructs(t *testing.T) {
	var (
		p    = newParser()
		file = newASTFile()
	)

	t.Run("with specified struct name expect parse only this structure", func(t *testing.T) {
		structs, err := p.parseStructs(file, []string{"Foo"})

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if len(structs) != 1 {
			t.Fatalf("len(structs) = %d, expect %d", len(structs), 1)
		}

		got := structs[0]

		want := structDecl{
			name: "Foo",
			fields: []field{
				{
					name:     "Bar",
					typeName: "string",
				},
			},
		}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("want %v got %v", want, got)
		}
	})

	t.Run("with specified several struct names expect parse these structs", func(t *testing.T) {
		got, err := p.parseStructs(file, []string{"Foo", "FooBar"})

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if len(got) != 2 {
			t.Fatalf("len(structs) = %d, expect %d", len(got), 2)
		}

		want := []structDecl{
			{
				name: "FooBar",
				fields: []field{
					{
						name:     "Number",
						typeName: "int",
					},
				},
			},
			{
				name: "Foo",
				fields: []field{
					{
						name:     "Bar",
						typeName: "string",
					},
				},
			},
		}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("want %v got %v", want, got)
		}
	})

	t.Run("without specified struct names expect parse all structs", func(t *testing.T) {
		got, err := p.parseStructs(file, []string{})

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if want := 3; len(got) != want {
			t.Fatalf("len(structs) = %d, expect %d", len(got), want)
		}

		want := []structDecl{
			{
				name: "FooBar",
				fields: []field{
					{
						name:     "Number",
						typeName: "int",
					},
				},
			},
			{
				name: "Foo",
				fields: []field{
					{
						name:     "Bar",
						typeName: "string",
					},
				},
			},
			{
				name: "Third",
				fields: []field{
					{
						name:     "Float",
						typeName: "float64",
					},
				},
			},
		}
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("want %v got %v", want, got)
		}
	})
}

func newParser() *Parser {
	return &Parser{}
}

var exampleSrc = `package foo

type FooBar struct {
	Number int
}

type Foo struct {
	Bar string
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
