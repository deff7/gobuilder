package generator

import "text/template"

var setMethodTmpl = template.Must(template.New("").Parse(`func (b *{{ .StructType }}Builder) {{ .FieldName }}(v {{ .FieldType }}) *{{ .StructType }}Builder {
	b.instance.{{ .FieldName }} = v
	return b
}
`))

var buildValueTmpl = template.Must(template.New("").Parse(`func (b *{{ .StructType }}Builder) V() {{ .StructType }} {
	return *b.instance
}
`))

var buildPointerTmpl = template.Must(template.New("").Parse(`func (b *{{ .StructType }}Builder) P() *{{ .StructType }} {
	return b.instance
}
`))

var declarationTmpl = template.Must(template.New("").Parse(`type {{ .StructType }}Builder struct {
	instance *{{ .StructType }}
}

func {{ .StructType }}() *{{ .StructType }}Builder {
	return &{{ .StructType }}Builder{
		instance: &{{ .StructType }}{},
	}
}
`))
