package generator

import "text/template"

var setMethodTmpl = template.Must(template.New("").Parse(`func (b *{{ .StructType }}Builder) {{ .FieldName }}(v {{ .FieldType }}) *{{ .StructType }}Builder {
	b.instance.{{ .FieldName }} = v
	return b
}
`))

var buildValueTmpl = template.Must(template.New("").Parse(`func (b *{{ .StructType }}Builder) V() {{ .PackageName }}{{ .StructType }} {
	return *b.instance
}
`))

var buildPointerTmpl = template.Must(template.New("").Parse(`func (b *{{ .StructType }}Builder) P() *{{ .PackageName }}{{ .StructType }} {
	return b.instance
}
`))

var declarationTmpl = template.Must(template.New("").Parse(`type {{ .StructType }}Builder struct {
	instance *{{ .PackageName }}{{ .StructType }}
}

func {{ .NewFuncPrefix }}{{ .StructType }}() *{{ .StructType }}Builder {
	return &{{ .StructType }}Builder{
		instance: &{{ .PackageName }}{{ .StructType }}{},
	}
}
`))
