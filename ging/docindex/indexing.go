package docindex

import (
	"bytes"
	"go/doc"

	"github.com/blevesearch/bleve"
)

// DocKind ...
// TODO(alvivi): doc this
type DocKind string

const (
	// PackageKind is the kind of a package.
	PackageKind DocKind = "p"
	// FuncKind is the kind of a function.
	FuncKind DocKind = "f"
	// MethodKind is the kind of a function.
	MethodKind DocKind = "m"
	// ConstKind is the kind of a const.
	ConstKind DocKind = "c"
	// VarKind is the kind of a variable.
	VarKind DocKind = "v"
	// TypeKind is the kind of a variable.
	TypeKind DocKind = "t"
)

// Package ...
// TODO(alvivi): doc this
type Package struct {
	Doc        string  `json:"doc"`
	Name       string  `json:"name"`
	ImportPath string  `json:"import"`
	Kind       DocKind `json:"kind"`

	Funcs  []*Func  `json:"funcs"`
	Consts []*Value `json:"const"`
	Vars   []*Value `json:"vars"`
	Types  []*Type  `json:"types"`
}

// NewPackage ...
// TODO(alvivi): doc this
func NewPackage(pkgDoc *doc.Package) *Package {
	pkg := new(Package)
	pkg.Kind = PackageKind
	pkg.Name = pkgDoc.Name
	pkg.ImportPath = pkgDoc.ImportPath
	buf := new(bytes.Buffer)
	doc.ToHTML(buf, pkgDoc.Doc, nil)
	pkg.Doc = removeDocSourcecode(buf.String())
	// Top level functions
	funcs := make([]*Func, len(pkgDoc.Funcs))
	for i, fn := range pkgDoc.Funcs {
		funcs[i] = NewFunction(pkg, fn)
	}
	pkg.Funcs = funcs
	// Top level constants
	consts := []*Value{}
	for _, c := range pkgDoc.Consts {
		consts = append(consts, NewConsts(pkg, c)...)
	}
	pkg.Consts = consts
	// Top level variables
	vars := []*Value{}
	for _, v := range pkgDoc.Vars {
		vars = append(vars, NewVars(pkg, v)...)
	}
	pkg.Vars = vars
	// Type declarations
	ts := make([]*Type, len(pkgDoc.Types))
	for i, t := range pkgDoc.Types {
		tt, fs := NewType(pkg, t)
		ts[i] = tt
		pkg.Funcs = append(pkg.Funcs, fs...)
	}
	pkg.Types = ts
	return pkg
}

// Type ...
// TODO(alvivi): doc this
func (pkg Package) Type() string {
	return "package"
}

// Func ...
// TODO(alvivi): doc this
type Func struct {
	Doc        string  `json:"doc"`
	Name       string  `json:"name"`
	ImportPath string  `json:"import"`
	Kind       DocKind `json:"kind"`
}

// NewFunction ...
// TODO(alvivi): doc this
func NewFunction(pkg *Package, fn *doc.Func) *Func {
	return &Func{
		Doc:        fn.Doc,
		Name:       fn.Name,
		ImportPath: pkg.ImportPath,
		Kind:       FuncKind,
	}
}

// NewMethod ...
// TODO(alvivi): doc this
func NewMethod(pkg *Package, fn *doc.Func) *Func {
	return &Func{
		Doc:        fn.Doc,
		Name:       fn.Name,
		ImportPath: pkg.ImportPath,
		Kind:       MethodKind,
	}
}

// Type ...
// TODO(alvivi): doc this
func (fn Func) Type() string {
	return "func"
}

// Value represents top level constants and variables.
type Value struct {
	Doc        string  `json:"doc"`
	Name       string  `json:"name"`
	ImportPath string  `json:"import"`
	Kind       DocKind `json:"kind"`
}

// NewConsts ...
// TODO(alvivi): doc this
func NewConsts(pkg *Package, v *doc.Value) []*Value {
	return newValues(pkg, v, ConstKind)
}

// NewVars ...
// TODO(alvivi): doc this
func NewVars(pkg *Package, v *doc.Value) []*Value {
	return newValues(pkg, v, VarKind)
}

// Type ...
// TODO(alvivi): doc this
func (v Value) Type() string {
	return "value"
}

// Type represents top level type declaration.
type Type struct {
	Doc        string  `json:"doc"`
	Name       string  `json:"name"`
	ImportPath string  `json:"import"`
	Kind       DocKind `json:"kind"`

	Methods []Func `json:"methods"`
}

// NewType ...
// TODO(alvivi): doc this
func NewType(pkg *Package, docType *doc.Type) (*Type, []*Func) {
	t := new(Type)
	t.Doc = docType.Doc
	t.Name = docType.Name
	t.ImportPath = pkg.ImportPath
	t.Kind = TypeKind
	fns := make([]*Func, len(docType.Funcs)+len(docType.Methods))
	for i, f := range docType.Funcs {
		fns[i] = NewFunction(pkg, f)
	}
	for i, m := range docType.Methods {
		fns[i+len(docType.Funcs)] = NewMethod(pkg, m)
	}
	return t, fns
}

// Type is the most recurisve method out there.
func (v Type) Type() string {
	return "type"
}

// OpenOrCreateIndex ...
// TODO(alvivi): doc this
func OpenOrCreateIndex(indexPath string) (bleve.Index, error) {
	idx, err := bleve.Open(indexPath)
	if err == nil {
		return idx, nil
	}
	mapping, err := buildDefaultMapping()
	if err != nil {
		return nil, err
	}
	return bleve.New(indexPath, mapping)
}

func buildDefaultMapping() (*bleve.IndexMapping, error) {
	// a generic reusable mapping for docucmentation content
	docFieldMapping := bleve.NewTextFieldMapping()
	docFieldMapping.Analyzer = "doc"

	// a generic reusable mapping for keyword text
	keywordFieldMapping := bleve.NewTextFieldMapping()
	keywordFieldMapping.Analyzer = "keyword"

	// a generic reusable mapping which only stores (but no index) a text
	noindexTextFieldMapping := bleve.NewTextFieldMapping()
	noindexTextFieldMapping.Store = true
	noindexTextFieldMapping.Index = false

	// a generinc reusable entry mapping
	entryMapping := bleve.NewDocumentStaticMapping()
	entryMapping.AddFieldMappingsAt("name", keywordFieldMapping)
	entryMapping.AddFieldMappingsAt("doc", docFieldMapping)
	entryMapping.AddFieldMappingsAt("kind", noindexTextFieldMapping)

	// Package Mapping
	packageMapping := bleve.NewDocumentStaticMapping()
	packageMapping.AddFieldMappingsAt("name", keywordFieldMapping)
	// TODO(alvivi): ImportPath must be indexed, but requires a custom analayzer
	// that removes the host.
	packageMapping.AddFieldMappingsAt("import", noindexTextFieldMapping)
	packageMapping.AddFieldMappingsAt("doc", docFieldMapping)
	packageMapping.AddFieldMappingsAt("kind", noindexTextFieldMapping)
	packageMapping.AddSubDocumentMapping("funcs", entryMapping)
	packageMapping.AddSubDocumentMapping("consts", entryMapping)
	packageMapping.AddSubDocumentMapping("vars", entryMapping)
	packageMapping.AddSubDocumentMapping("types", entryMapping)

	// Index Mapping
	indexMapping := bleve.NewIndexMapping()
	err := indexMapping.AddCustomAnalyzer("doc",
		map[string]interface{}{
			"type":          "custom",
			"char_filters":  []string{"html"},
			"tokenizer":     "whitespace",
			"token_filters": []string{"to_lower", "stop_en"},
		})
	if err != nil {
		return nil, err
	}
	indexMapping.AddDocumentMapping("package", packageMapping)
	indexMapping.AddDocumentMapping("func", entryMapping)
	indexMapping.AddDocumentMapping("const", entryMapping)
	indexMapping.AddDocumentMapping("vars", entryMapping)
	indexMapping.AddDocumentMapping("types", entryMapping)
	return indexMapping, nil
}

func newValues(pkg *Package, value *doc.Value, t DocKind) []*Value {
	vs := make([]*Value, len(value.Names))
	for i, n := range value.Names {
		vs[i] = &Value{
			Doc:        value.Doc,
			Name:       n,
			ImportPath: pkg.ImportPath,
			Kind:       t,
		}
	}
	return vs
}
