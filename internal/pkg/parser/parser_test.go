package parser

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
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
	t.Run("with specified struct name expect parse only this structure", func(t *testing.T) {
		p := newParser()
		file := newASTFile()

		structs, err := p.parseStructs(file, []string{"Foo"})

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if len(structs) != 1 {
			t.Fatalf("len(structs) = %d, expect %d", len(structs), 1)
		}

		s := structs[0]
		want := "Foo"
		if s.name != want {
			t.Errorf("expect %q, got %q", want, s.name)
		}

		if len(s.fields) != 1 {
			t.Fatalf("len(fields) = %d, expect %d", len(s.fields), 1)
		}

		f := s.fields[0]

		want = "Bar"
		if f.name != want {
			t.Errorf("expect %q, got %q", want, f.name)
		}

		want = "string"
		if f.typeName != want {
			t.Errorf("expect %q, got %q", want, f.typeName)
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
