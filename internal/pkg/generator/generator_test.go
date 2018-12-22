package generator

import (
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	t.Run("with empty typeName expect error", func(t *testing.T) {
		_, err := NewGenerator("")

		if err == nil {
			t.Error("expect error")
		}
	})
}

func TestGenerate(t *testing.T) {
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

func newGenerator() *Generator {
	return &Generator{
		typeName: "SampleType",
	}
}

func newField(name, typeName string) field {
	return field{
		name:     name,
		typeName: typeName,
	}
}
