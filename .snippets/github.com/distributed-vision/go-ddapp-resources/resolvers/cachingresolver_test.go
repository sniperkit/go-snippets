package resolvers_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/types/gotypeid"
	"github.com/distributed-vision/go-resources/util"
	"github.com/distributed-vision/go-resources/util/random"
)

type untypedSelector struct {
	key   string
	value string
}

func (this *untypedSelector) Type() ids.TypeIdentifier {
	return nil
}

func (this *untypedSelector) Key() interface{} {
	return this.key
}

func (this *untypedSelector) Test(i interface{}) bool {
	if e, ok := i.(entity); ok {
		if this.key != "" && e.key == this.key {
			return true
		}

		if this.value != "" && e.value == this.value {
			return true
		}
	}

	return false
}

type typedSelector struct {
	key   string
	value string
}

func (this *typedSelector) Type() ids.TypeIdentifier {
	return testEntityType
}

func (this *typedSelector) Key() interface{} {
	return this.key
}

func (this *typedSelector) Test(i interface{}) bool {
	if e, ok := i.(entity); ok {
		if this.key != "" && e.key == this.key {
			return true
		}

		if this.value != "" && e.value == this.value {
			return true
		}
	}

	return false
}

type entity struct {
	key   string
	value string
}

func testExtractor(e ...interface{}) (interface{}, bool) {
	//fmt.Printf("e=%v", e)
	return e[0].(entity).key, true
}

var testContext = context.Background()
var testInfo = resolvers.RootInfo.WithExtractor(testExtractor)
var testEntityType ids.TypeIdentifier = gotypeid.IdOf(reflect.TypeOf(entity{}))

func TestResolverCache(t *testing.T) {
	resolver, err := resolvers.NewCachingResolver(nil, 128)

	if err == nil {
		t.Error("TestLocalResolver: NewLocalResolver with nil info should fail")
	}

	resolver, err = resolvers.NewCachingResolver(testInfo, 128)

	if err != nil {
		t.Fatal("TestLocalResolver: NewLocalResolver failed:", err)
	}

	if (!resolver.ResolverInfo().Matches(&untypedSelector{})) {
		t.Fatal("TestLocalResolver: ResolverInfo .Match failed")
	}

	keys := make([]string, 256)
	values := make([]string, 256)

	for i := 0; i < 256; i++ {
		keys[i] = random.RandomString(20)
		values[i] = random.RandomString(20)
		resolver.Cache().Add(keys[i], entity{keys[i], values[i]})
	}

	if resolver.Cache().Len() != 128 {
		t.Fatalf("bad cache size: %v", resolver.Cache().Len())
	}

	for i, k := range resolver.Cache().Keys() {
		if e, ok := resolver.Cache().Get(k); !ok ||
			e.(entity).key != keys[i+128] || e.(entity).value != values[i+128] {
			t.Fatalf("bad key: %v", k)
		}
	}

	for i := 0; i < 128; i++ {
		_, err := resolver.Get(testContext, &untypedSelector{key: keys[i]})
		if err == nil {
			t.Fatalf("should be evicted")
		}
		_, err = util.Await(resolver.Resolve(testContext, &untypedSelector{key: keys[i]}))
		if err == nil {
			t.Fatalf("should be evicted")
		}
	}

	for i := 128; i < 256; i++ {
		_, err := resolver.Get(testContext, &untypedSelector{key: keys[i]})
		if err != nil {
			t.Fatalf("should not be evicted")
		}
		_, err = util.Await(resolver.Resolve(testContext, &untypedSelector{key: keys[i]}))
		if err != nil {
			t.Fatalf("should not be evicted")
		}
	}

	for i := 128; i < 192; i++ {
		resolver.Cache().Remove(keys[i])
		_, err := resolver.Get(testContext, &untypedSelector{key: keys[i]})
		if err == nil {
			t.Fatalf("should be deleted")
		}
		_, err = util.Await(resolver.Resolve(testContext, &untypedSelector{key: keys[i]}))
		if err == nil {
			t.Fatalf("should be deleted")
		}
	}

	resolver.Cache().Purge()
	if resolver.Cache().Len() != 0 {
		t.Fatalf("bad len: %v", resolver.Cache().Len())
	}
	if _, err := resolver.Get(testContext, &untypedSelector{key: keys[200]}); err == nil {
		t.Fatalf("should contain nothing")
	}

}
