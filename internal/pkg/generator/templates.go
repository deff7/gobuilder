package generator

import "text/template"

var setMethodTmpl = template.Must(template.New("").Parse(`func (b *{{ .StructType }}Builder) {{ .FieldName }}(v {{ .FieldType }}) *{{ .StructType }}Builder {
	b.instance.{{ .FieldName }} = v
	return b
}
`))
