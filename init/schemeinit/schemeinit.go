package schemeinit

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/scheme"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/resolvers/fileresolver"
	"github.com/distributed-vision/go-resources/types/gotypeid"
)

func Init() {
	schemePath := os.Getenv("DV_DOMAIN_SCOPE_PATH")

	if schemePath == "" {
		schemePath = "../../../id-schemes"
	}

	entityType := gotypeid.IdOf(reflect.TypeOf((*ids.Scheme)(nil)).Elem())

	resolverFactory, err := fileresolver.NewResolverFactory(
		resolvers.NewResolverInfo(fileresolver.PublicType,
			[]ids.TypeIdentifier{entityType},
			[]ids.Domain{},
			scheme.KeyExtractor,
			map[interface{}]interface{}{
				"location": "schemeinfo.json",
				"paths":    filepath.SplitList(schemePath)}))

	if err != nil {
		panic(fmt.Sprintf("Unexpected error creating resolver factory: %s", err))
	}

	scheme.RegisterResolverFactory(resolverFactory)
}
