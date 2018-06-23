package composer

import "bytes"
import "text/template"

type File struct {
	Package   string
	Variables []*Variable
}

func NewFile(pkg string, vars ...*Variable) *File {
	return &File{pkg, vars}
}

func (f *File) Add(v *Variable) {
	f.Variables = append(f.Variables, v)
}

func (f *File) Get(name string) *Variable {
	for _, v := range f.Variables {
		if v.Name == name {
			return v
		}
	}
	return nil
}

var fileTemplate = template.Must(template.New("Variable").Parse(`package {{.Package}}

import "regexp"

{{range .Variables}}
{{.}}
{{end}}`))

func (f *File) String() string {
	buf := bytes.NewBuffer(nil)
	if err := fileTemplate.Execute(buf, f); err != nil {
		panic(err)
	}
	return buf.String()
}
