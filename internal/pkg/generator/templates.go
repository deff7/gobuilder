package generator

import "text/template"

var setMethodTmpl = template.Must(template.New("").Parse(`// {{ .FieldName }} sets field with type {{ .FieldType }}
func (b *{{ .StructType }}Builder) {{ .FieldName }}(v {{ .FieldType }}) *{{ .StructType }}Builder {
	b.instance.{{ .FieldName }} = v
	return b
}
`))

var buildValueTmpl = template.Must(template.New("").Parse(`// V returns value of {{ .StructType }} instance
func (b *{{ .StructType }}Builder) V() {{ .PackageName }}{{ .StructType }} {
	return *b.instance
}
`))

var buildPointerTmpl = template.Must(template.New("").Parse(`// P returns pointer to {{ .StructType }} instance
func (b *{{ .StructType }}Builder) P() *{{ .PackageName }}{{ .StructType }} {
	return b.instance
}
`))

var declarationTmpl = template.Must(template.New("").Parse(`// {{ .StructType }}Builder is builder for type {{ .StructType }}
type {{ .StructType }}Builder struct {
	instance *{{ .PackageName }}{{ .StructType }}
}

// {{ .NewFuncPrefix }}{{ .StructType }} creates new builder
func {{ .NewFuncPrefix }}{{ .StructType }}() *{{ .StructType }}Builder {
	return &{{ .StructType }}Builder{
		instance: &{{ .PackageName }}{{ .StructType }}{},
	}
}
`))
