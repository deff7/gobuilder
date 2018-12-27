package generator

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	t.Run("with empty typeName expect error", func(t *testing.T) {
		_, err := NewGenerator("", "", "")

		if err == nil {
			t.Error("expect error")
		}
	})

	t.Run("with all params expect ok", func(t *testing.T) {
		var (
			typeName    = "Foo"
			packageName = "foo"
			filter      = ".+"
		)

		g, err := NewGenerator(typeName, packageName, filter)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if g.typeName != typeName {
			t.Errorf("want %q, got %q", typeName, g.typeName)
		}

		if g.packageName != packageName {
			t.Errorf("want %q, got %q", packageName, g.packageName)
		}

		if g.filterRE == nil {
			t.Error("filter regexp is not initializaed")
		}
	})
}

func TestAddField(t *testing.T) {
	t.Run("without field filter", func(t *testing.T) {
		var (
			g         = newGenerator()
			fieldName = "foo"
			fieldType = "string"
		)

		g.AddField(fieldName, fieldType)

		if len(g.fields) == 0 {
			t.Fatal("field is not added")
		}

		if got := g.fields[0].name; got != fieldName {
			t.Errorf("expect %q, got %q", fieldName, got)
		}

		if got := g.fields[0].typeName; got != fieldType {
			t.Errorf("expect %q, got %q", fieldType, got)
		}
	})

	t.Run("with field filter as regexp and field name that satisfy one expect skip field", func(t *testing.T) {
		var (
			g         = newGenerator()
			fieldName = "foo"
			fieldType = "string"
		)
		g.filterRE = regexp.MustCompile(".oo")

		g.AddField(fieldName, fieldType)

		if len(g.fields) != 0 {
			t.Error("field is not skipped")
		}
	})

	t.Run("with field filter and field name that not satisfy one expect add field", func(t *testing.T) {
		var (
			g         = newGenerator()
			fieldName = "far"
			fieldType = "string"
		)
		g.filterRE = regexp.MustCompile(".oo")

		g.AddField(fieldName, fieldType)

		if len(g.fields) == 0 {
			t.Fatal("field is not added")
		}
	})
}

func TestGenerate(t *testing.T) {
	g := newGenerator(
		withFields([]field{newField("Foo", "string")}),
	)

	got, err := g.Generate()

	assertEqualFromFile(t, "builder.golden", got, err, false)
}

func assertEqualFromFile(t *testing.T, wantFile, got string, err error, ignoreComments bool) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	wantBytes, err := ioutil.ReadFile(filepath.Join("testdata", wantFile))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := string(wantBytes)

	if ignoreComments {
		want = removeComments(want)
		got = removeComments(got)
	}

	if want != got {
		t.Errorf("expect %s, got %s", want, got)
	}
}

func removeComments(src string) string {
	var buf bytes.Buffer
	for _, line := range strings.Split(src, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "//") {
			continue
		}
		buf.WriteString(line + "\n")
	}
	s := buf.String()
	s = s[:len(s)-1]
	return s
}

func TestGenerateBuildValue(t *testing.T) {
	g := newGenerator()

	got, err := g.generateBuildValue()

	assertEqualFromFile(t, "build_value.golden", got, err, true)
}

func TestGenerateBuildValueSamePackage(t *testing.T) {
	g := newGenerator(withPackageName(""))

	got, err := g.generateBuildValue()

	assertEqualFromFile(t, "build_value_same_package.golden", got, err, true)
}

func TestGenerateBuildPointer(t *testing.T) {
	g := newGenerator()

	got, err := g.generateBuildPointer()

	assertEqualFromFile(t, "build_pointer.golden", got, err, true)
}

func TestGenerateBuildPointerSamePackage(t *testing.T) {
	g := newGenerator(withPackageName(""))

	got, err := g.generateBuildPointer()

	assertEqualFromFile(t, "build_pointer_same_package.golden", got, err, true)
}

func TestGenerateDeclaration(t *testing.T) {
	g := newGenerator()

	got, err := g.generateDeclaration()

	assertEqualFromFile(t, "declaration.golden", got, err, true)
}

func TestGenerateDeclarationSamePackage(t *testing.T) {
	g := newGenerator(withPackageName(""))

	got, err := g.generateDeclaration()

	assertEqualFromFile(t, "declaration_same_package.golden", got, err, true)
}

func TestGenerateSetMethod(t *testing.T) {
	for _, tc := range []struct {
		name     string
		field    field
		err      error
		wantFile string
	}{
		{
			name:  "with empty field name expect error",
			field: newField("", "string"),
			err:   ErrEmptyFieldName,
		},
		{
			name:  "with empty field type expect error",
			field: newField("Foo", ""),
			err:   ErrEmptyFieldType,
		},
		{
			name:     "with valid field expect returns method source code",
			field:    newField("Foo", "string"),
			wantFile: "method.golden",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			g := newGenerator()

			got, err := g.generateSetMethod(tc.field)

			if tc.err != nil {
				if tc.err != err {
					t.Errorf("want error: %v, got: %v", tc.err, err)
				}
				return
			}

			assertEqualFromFile(t, tc.wantFile, got, err, true)
		})
	}
}

type generatorOption func(*Generator)

func withPackageName(packageName string) generatorOption {
	return func(g *Generator) {
		g.packageName = packageName
	}
}

func withFields(fields []field) generatorOption {
	return func(g *Generator) {
		g.fields = fields
	}
}

func newGenerator(opts ...generatorOption) *Generator {
	g := &Generator{
		typeName:    "SampleType",
		packageName: "domain",
	}
	for _, apply := range opts {
		apply(g)
	}
	return g
}

func newField(name, typeName string) field {
	return field{
		name:     name,
		typeName: typeName,
	}
}
