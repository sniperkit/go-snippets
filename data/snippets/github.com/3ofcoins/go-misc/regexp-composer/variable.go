package composer

import "fmt"

import "bytes"
import "regexp/syntax"
import "text/template"

// A variable is a named regexp. It is translated to a
// `*regexp.Regexp` variable and constants enumerating named captures.
type Variable struct {
	Name string
	*syntax.Regexp
}

// Returns new variable object or error
func NewVariable(name string, elements ...RxElement) (*Variable, error) {
	rx_str := string(Sequence(elements...))
	rx, err := syntax.Parse(rx_str, syntax.Perl)
	if err != nil {
		return nil, fmt.Errorf("Error parsing %q: %v", rx_str, err)
	}
	return &Variable{name, rx.Simplify()}, nil
}

// Returns new variable object, panics on error
func MustVariable(name string, elements ...RxElement) *Variable {
	if vrbl, err := NewVariable(name, elements...); err != nil {
		panic(err)
	} else {
		return vrbl
	}
}

var vrblTemplate = template.Must(template.New("Variable").Parse(`var {{.Name}} = regexp.MustCompile({{.Regexp|printf "%q"}})
{{if gt .MaxCap 0}}
const ({{$prefix := .Name}}{{range .CapNames}}
  {{if eq . ""}}_{{else}}{{$prefix}}_{{.}}{{end}} = iota{{end}}
){{end}}`))

func (vrbl *Variable) String() string {
	buf := bytes.NewBuffer(nil)
	if err := vrblTemplate.Execute(buf, vrbl); err != nil {
		panic(err)
	}
	return buf.String()
}
