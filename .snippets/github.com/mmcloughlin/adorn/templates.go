package adorn

import "text/template"

var packageTemplate = template.Must(template.New("package").Parse(
	"package {{ .Package }}\n",
))

var interfaceTemplate = template.Must(template.New("interface").Parse(`
{{if .Documentation}}// {{ .Documentation }}{{end}}
type {{ .TypeName }} interface {
	{{ .MethodName }}({{ .ArgumentsUnnamed }}){{if .ReturnSignature}} {{ .ReturnSignature }}{{end}}
}
`))

var funcTemplate = template.Must(template.New("func").Parse(`
// {{ .FuncTypeName }} is an adapter to allow ordinary functions to be used as {{ .TypeName }} implementations.
type {{ .FuncTypeName }} func({{ .ArgumentsUnnamed }}){{if .ReturnSignature}} {{ .ReturnSignature }}{{end}}

// {{ .MethodName }} calls f.
func (f {{ .FuncTypeName }}) {{ .MethodName }}({{ .ArgumentsNamed }}){{if .ReturnSignature}} {{ .ReturnSignature }}{{end}} {
	{{if .ReturnSignature}}return {{end}}f({{ .ArgumentsCalling }})
}
`))

var assertionTemplate = template.Must(template.New("assertion").Parse(`
// Compile time assertion that {{ .FuncTypeName }} implements {{ .TypeName }}.
var _ {{ .TypeName }} = new({{ .FuncTypeName }})
`))
