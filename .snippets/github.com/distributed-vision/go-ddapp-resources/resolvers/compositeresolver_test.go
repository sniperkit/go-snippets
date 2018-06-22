package resolvers_test

import (
	"context"
	"os"
	"testing"

	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/resolvers/localresolver"
	"github.com/distributed-vision/go-resources/types"
	"github.com/distributed-vision/go-resources/types/publictypeid"
	"github.com/distributed-vision/go-resources/util"
	"github.com/distributed-vision/go-resources/util/random"
	"github.com/distributed-vision/go-resources/version"
)

func TestMain(m *testing.M) {
	resolvers.ResisterNewFactoryFunction(testResolverType, NewResolverFactory)
	os.Exit(m.Run())
}

var testResolverType = types.MustNewId(publictypeid.ResolverDomain, []byte("TestResolver"), version.New(1, 0, 0))

type factory struct {
	resolverInfo resolvers.ResolverInfo
}

func NewResolverFactory(resolverInfo resolvers.ResolverInfo) (resolvers.ResolverFactory, error) {
	return &factory{resolverInfo}, nil
}

func (this *factory) New(resolutionContext context.Context) (resolvers.Resolver, error) {

	resolver, err := localresolver.New(this.resolverInfo)

	if err == nil {
		componentResolvers = append(componentResolvers, resolver)
	}

	return resolver, err
}

func (this *factory) ResolverType() ids.TypeIdentifier {
	return testResolverType
}

func (this *factory) ResolverInfo() resolvers.ResolverInfo {
	return this.resolverInfo
}

var componentResolvers = []*localresolver.LocalResolver{}
var testResolvableTypes = []ids.TypeIdentifier{testEntityType}

var resolverInfos = []resolvers.ResolverInfo{
	resolvers.NewResolverInfo(testResolverType, testResolvableTypes, nil, testExtractor, nil),
	resolvers.NewResolverInfo(testResolverType, testResolvableTypes, nil, testExtractor, nil),
	resolvers.NewResolverInfo(testResolverType, testResolvableTypes, nil, testExtractor, nil),
	resolvers.NewResolverInfo(testResolverType, testResolvableTypes, nil, testExtractor, nil),
}

func TestCompositeGet(t *testing.T) {
	resolver, err := resolvers.NewCompositeResolver(nil)

	if err == nil {
		t.Error("TestCompositeGet: NewLocalResolver with nil info should fail")
	}

	resolver, err = resolvers.NewCompositeResolver(testInfo)

	if err != nil {
		t.Fatal("TestCompositeGet: NewLocalResolver failed:", err)
	}

	for _, resolverInfo := range resolverInfos {
		factory, err := resolvers.NewResolverFactory(resolverInfo)

		if err == nil {
			err = resolver.RegisterComponentFactory(factory, true)

			if err != nil {
				t.Fatal("TestCompositeGet: Factory registration failed: err:", err)
			}
		} else {
			t.Fatal("TestCompositeGet: Factory creation failed: err:", err)
		}
	}

	if len(componentResolvers) < 4 {
		t.Fatal("TestCompositeGet: Expected 4 component resolvers got:", len(componentResolvers))
	}

	keys := make([]string, 256)
	values := make([]string, 256)

	for i := 0; i < 256; i++ {
		keys[i] = random.RandomString(20)
		values[i] = random.RandomString(20)

		componentResolvers[i%4].Put(testContext, entity{keys[i], values[i]})
	}

	for i := 0; i < 256; i++ {
		resolved, err := resolver.Get(testContext, &untypedSelector{key: keys[i]})

		if err != nil {
			t.Fatal("TestLocalResolverGet: LocalResolver.Get failed:", err)
		}

		if keys[i] != resolved.(entity).key ||
			values[i] != resolved.(entity).value {
			t.Fatalf("TestLocalResolverGet: Key %v get unexpectedly failed", i)
		}

		resolved, err = util.Await(resolver.Resolve(testContext, &untypedSelector{key: keys[i]}))

		if err != nil {
			t.Fatal("TestLocalResolverGet: LocalResolver.Get failed:", err)
		}

		if keys[i] != resolved.(entity).key ||
			values[i] != resolved.(entity).value {
			t.Fatalf("TestLocalResolverGet: Key %v get unexpectedly failed", i)
		}
	}

	if resolver.Cache().Len() != 256 {
		t.Fatal("Invalid cache length:", resolver.Cache().Len())
	}

	resolver.Cache().Purge()

	for i := 0; i < 256; i++ {
		resolved, err := resolver.Get(testContext, &typedSelector{key: keys[i]})

		if err != nil {
			t.Fatal("TestLocalResolverGet: LocalResolver.Get failed:", err)
		}

		if keys[i] != resolved.(entity).key ||
			values[i] != resolved.(entity).value {
			t.Fatalf("TestLocalResolverGet: Key %v get unexpectedly failed", i)
		}

		resolved, err = util.Await(resolver.Resolve(testContext, &typedSelector{key: keys[i]}))

		if err != nil {
			t.Fatal("TestLocalResolverGet: LocalResolver.Get failed:", err)
		}

		if keys[i] != resolved.(entity).key ||
			values[i] != resolved.(entity).value {
			t.Fatalf("TestLocalResolverGet: Key %v get unexpectedly failed", i)
		}
	}

}
