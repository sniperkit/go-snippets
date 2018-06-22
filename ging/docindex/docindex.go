package docindex

import (
	"fmt"
	"net/http"

	"go/doc"

	"github.com/blevesearch/bleve"
	"github.com/golang/gddo/gosrc"
)

// SetLocalDevMode sets the package to local development mode.
func SetLocalDevMode(path string) {
	gosrc.SetLocalDevMode(path)
}

// IndexPackage ...
// TODO(alvivi): doc this
func IndexPackage(client *http.Client, index bleve.Index, pkgPath string) error {
	pkg, err := fetchPackage(client, pkgPath)
	if err != nil {
		return err
	}
	pkgDesc := NewPackage(doc.New(pkg, pkgPath, 0))
	// Functions
	for _, fnDesc := range pkgDesc.Funcs {
		fnName := fmt.Sprintf("%s.%s", pkgDesc.ImportPath, fnDesc.Name)
		err := index.Index(fnName, fnDesc)
		if err != nil {
			return err
		}
	}
	// Constants
	for _, constDesc := range pkgDesc.Consts {
		constName := fmt.Sprintf("%s.%s", pkgDesc.ImportPath, constDesc.Name)
		err := index.Index(constName, constDesc)
		if err != nil {
			return err
		}
	}
	// Variables
	for _, varDesc := range pkgDesc.Vars {
		varName := fmt.Sprintf("%s.%s", pkgDesc.ImportPath, varDesc.Name)
		err := index.Index(varName, varDesc)
		if err != nil {
			return err
		}
	}
	// Types
	for _, typeDesc := range pkgDesc.Types {
		typeName := fmt.Sprintf("%s.%s", pkgDesc.ImportPath, typeDesc.Name)
		err := index.Index(typeName, typeDesc)
		if err != nil {
			return err
		}
	}
	return index.Index(pkgDesc.ImportPath, pkgDesc)
}
