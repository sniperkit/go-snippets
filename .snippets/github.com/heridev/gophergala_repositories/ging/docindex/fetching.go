package docindex

import (
	"errors"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"net/http"
	"regexp"
	"strings"

	"github.com/golang/gddo/gosrc"
)

func fetchPackage(client *http.Client, path string) (*ast.Package, error) {
	dir, err := gosrc.Get(client, path, "")
	if err != nil {
		return nil, err
	}
	// Fisrt, we have to find a proper build context for the package
	bctx := build.Context{
		CgoEnabled:  true,
		ReleaseTags: build.Default.ReleaseTags,
		BuildTags:   build.Default.BuildTags,
		Compiler:    "gc",
	}
	var bpkg *build.Package
	for _, env := range buildEnvs {
		bctx.GOOS = env.GOOS
		bctx.GOARCH = env.GOARCH
		bpkg, err = dir.Import(&bctx, build.ImportComment)
		if _, ok := err.(*build.NoGoError); !ok {
			break
		}
	}
	// Parse all package's sourcecode
	fileSet := token.NewFileSet()
	filesData := map[string][]byte{}
	pkgFiles := map[string]*ast.File{}
	for _, file := range dir.Files {
		if strings.HasSuffix(file.Name, ".go") {
			gosrc.OverwriteLineComments(file.Data)
		}
		filesData[file.Name] = file.Data
		// TODO(alvivi): else { addReferences(references, file.Data) }
	}
	fileNames := append(bpkg.GoFiles, bpkg.CgoFiles...)
	for _, fname := range fileNames {
		pfile, err :=
			parser.ParseFile(fileSet, fname, filesData[fname], parser.ParseComments)
		if err != nil {
			return nil, err
		}
		pkgFiles[fname] = pfile
	}
	// Actually, we don't care about building the package. Only the parser have
	// to succeed to read its documentation.
	pkg, _ := ast.NewPackage(fileSet, pkgFiles, simpleImporter, nil)
	return pkg, nil
}

var buildEnvs = []struct{ GOOS, GOARCH string }{
	{"linux", "amd64"},
	{"darwin", "amd64"},
	{"windows", "amd64"},
}

/*
A stub importer implementation
*/

// simpleImporter is a importert which actually does not import anything.
// From github.com/golang/gddo/doc/builder.go, at line 335.
func simpleImporter(imports map[string]*ast.Object, path string) (*ast.Object, error) {
	pkg := imports[path]
	if pkg != nil {
		return pkg, nil
	}
	// Guess the package name without importing it.
	for _, pat := range packageNamePats {
		m := pat.FindStringSubmatch(path)
		if m != nil {
			pkg = ast.NewObj(ast.Pkg, m[1])
			pkg.Data = ast.NewScope(nil)
			imports[path] = pkg
			return pkg, nil
		}
	}
	return nil, errors.New("package not found")
}

// From github.com/golang/gddo/doc/builder.go, at line 315.
var packageNamePats = []*regexp.Regexp{
	regexp.MustCompile(`/([^-./]+)[-.](?:git|svn|hg|bzr|v\d+)$`),
	regexp.MustCompile(`/([^-./]+)[-.]go$`),
	regexp.MustCompile(`/go[-.]([^-./]+)$`),
	regexp.MustCompile(`^code\.google\.com/p/google-api-go-client/([^/]+)/v[^/]+$`),
	regexp.MustCompile(`^code\.google\.com/p/biogo\.([^/]+)$`),
	regexp.MustCompile(`([^/]+)$`),
}
