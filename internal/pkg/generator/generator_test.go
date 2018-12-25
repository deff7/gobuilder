package generator

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	t.Run("with empty typeName expect error", func(t *testing.T) {
		_, err := NewGenerator("", "")

		if err == nil {
			t.Error("expect error")
		}
	})

	t.Run("with all params expect ok", func(t *testing.T) {
		var (
			typeName    = "Foo"
			packageName = "foo"
		)

		g, err := NewGenerator(typeName, packageName)

		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}

		if g.typeName != typeName {
			t.Errorf("want %q, got %q", typeName, g.typeName)
		}

		if g.packageName != packageName {
			t.Errorf("want %q, got %q", packageName, g.packageName)
		}
	})
}

func TestAddField(t *testing.T) {
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
}

func TestGenerate(t *testing.T) {
	g := newGenerator(
		withFields([]field{newField("Foo", "string")}),
	)

	got, err := g.Generate()

	assertEqualFromFile(t, "builder.golden", got, err)
}

func assertEqualFromFile(t *testing.T, wantFile, got string, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want, err := ioutil.ReadFile(filepath.Join("testdata", wantFile))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(want) != got {
		t.Errorf("expect %s, got %s", want, got)
	}
}

func TestGenerateBuildValue(t *testing.T) {
	g := newGenerator()

	got, err := g.generateBuildValue()

	assertEqualFromFile(t, "build_value.golden", got, err)
}

func TestGenerateBuildValueSamePackage(t *testing.T) {
	g := newGenerator(withPackageName(""))

	got, err := g.generateBuildValue()

	assertEqualFromFile(t, "build_value_same_package.golden", got, err)
}

func TestGenerateBuildPointer(t *testing.T) {
	g := newGenerator()

	got, err := g.generateBuildPointer()

	assertEqualFromFile(t, "build_pointer.golden", got, err)
}

func TestGenerateBuildPointerSamePackage(t *testing.T) {
	g := newGenerator(withPackageName(""))

	got, err := g.generateBuildPointer()

	assertEqualFromFile(t, "build_pointer_same_package.golden", got, err)
}

func TestGenerateDeclaration(t *testing.T) {
	g := newGenerator()

	got, err := g.generateDeclaration()

	assertEqualFromFile(t, "declaration.golden", got, err)
}

func TestGenerateDeclarationSamePackage(t *testing.T) {
	g := newGenerator(withPackageName(""))

	got, err := g.generateDeclaration()

	assertEqualFromFile(t, "declaration_same_package.golden", got, err)
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

			want, err := ioutil.ReadFile(filepath.Join("./testdata", tc.wantFile))
			if err != nil {
				t.Fatal(err)
			}

			if got != string(want) {
				t.Errorf("expect %q got %q", want, got)
			}
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
