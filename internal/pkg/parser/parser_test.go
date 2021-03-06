package parser

import (
	"bytes"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	p := newParser()
	buf := bytes.NewBufferString(exampleSrc)

	packageName := "foo"
	structsMap, err := p.Parse(buf, []string{})
	structs, ok := structsMap[packageName]

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if !ok {
		t.Fatalf("structs for package %q not found", packageName)
	}

	if want := 3; len(structs) != want {
		t.Errorf("len(structs) = %d, want %d", len(structs), want)
	}

}

func TestParseDir(t *testing.T) {
	path := "testdata"
	allowedStructs := []string{}
	p := newParser()

	structsMap, err := p.ParseDir(path, allowedStructs)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	want := map[string][]StructDecl{
		"foo": []StructDecl{
			{
				Name:   "Foo",
				Fields: []Field{{"Num", "int"}},
			},
		},
		"bar": []StructDecl{
			{
				Name:   "Bar",
				Fields: []Field{{"String", "string"}},
			},
		},
	}

	if !reflect.DeepEqual(structsMap, want) {
		t.Errorf("want %v but got %v", want, structsMap)
	}
}

func TestParseStructs(t *testing.T) {
	t.Run("unexported structs and fields", func(t *testing.T) {
		for _, tc := range []struct {
			name       string
			unexported bool
			want       []StructDecl
		}{
			{
				name:       "disallowed",
				unexported: false,
				want: []StructDecl{
					newStructDecl("Foo", []Field{newField("Bar", "int")}),
				},
			},
			{
				name:       "allowed",
				unexported: true,
				want: []StructDecl{
					newStructDecl("Foo", []Field{newField("foo", "string"), newField("Bar", "int")}),
					newStructDecl("bar", []Field{}),
				},
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				var (
					p    = newParser()
					file = newASTFile(exampleWithUnexported)
				)
				p.unexported = tc.unexported
				got, err := p.parseStructs(file, []string{})

				if err != nil {
					t.Errorf("unexpected error: %s", err)
				}

				if !reflect.DeepEqual(tc.want, got) {
					t.Fatalf("want %v got %v", tc.want, got)
				}
			})
		}
	})

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
				newStructDecl("Third", []Field{newField("Floats", "[2]*foreign.Float64")}),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				p    = newParser()
				file = newASTFile(exampleSrc)
			)
			got, err := p.parseStructs(file, tc.allowedStructs)

			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("want %v got %v", tc.want, got)
			}
		})
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
	Floats [2]*foreign.Float64
}`

var exampleWithUnexported = `package foo

type Foo struct {
	foo string
	Bar int
}

type bar struct {
}`

func newASTFile(src string) *ast.File {
	fset := token.NewFileSet()
	buf := bytes.NewBufferString(src)
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
