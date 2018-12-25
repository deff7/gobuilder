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

func TestParse(t *testing.T) {
	p := newParser()
	buf := bytes.NewBufferString(exampleSrc)

	structs, err := p.Parse(buf, []string{})
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if want := 3; len(structs) != want {
		t.Errorf("len(structs) = %d, want %d", len(structs), want)
	}
}

func TestParseStructs(t *testing.T) {
	for _, tc := range []struct {
		name           string
		allowedStructs []string
		want           []StructDecl
	}{
		{
			name:           "with specified struct name expect parse only this structure",
			allowedStructs: []string{"Second"},
			want: []StructDecl{
				newStructDecl("Second", []Field{newField("String", "string")}),
			},
		},
		{
			name:           "with specified several struct names expect parse these structs",
			allowedStructs: []string{"First", "Second"},
			want: []StructDecl{
				newStructDecl("First", []Field{newField("Number", "int")}),
				newStructDecl("Second", []Field{newField("String", "string")}),
			},
		},
		{
			name:           "without specified struct names expect parse all structs",
			allowedStructs: []string{},
			want: []StructDecl{
				newStructDecl("First", []Field{newField("Number", "int")}),
				newStructDecl("Second", []Field{newField("String", "string")}),
				newStructDecl("Third", []Field{newField("Floats", "[2]*float64")}),
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
	Floats [2]*float64
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

func newField(name, typeName string) Field {
	return Field{Name: name, TypeName: typeName}
}

func newStructDecl(name string, fields []Field) StructDecl {
	return StructDecl{
		Name:   name,
		Fields: fields,
	}
}
