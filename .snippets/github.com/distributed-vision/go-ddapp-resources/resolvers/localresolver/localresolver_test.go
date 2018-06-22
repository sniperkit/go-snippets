package localresolver_test

import (
	"context"
	"testing"

	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/resolvers"
	"github.com/distributed-vision/go-resources/resolvers/localresolver"
	"github.com/distributed-vision/go-resources/util"
	"github.com/distributed-vision/go-resources/util/random"
)

type selector struct {
	key   string
	value string
}

func (this *selector) Type() ids.TypeIdentifier {
	return nil
}

func (this *selector) Key() interface{} {
	return this.key
}

func (this *selector) Test(i interface{}) bool {
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

var testContext = context.Background()
var testInfo = resolvers.RootInfo.WithExtractor(func(e ...interface{}) (interface{}, bool) {
	//fmt.Printf("e=%v", e)
	return e[0].(entity).key, true
})

func TestLocalResolverGet(t *testing.T) {
	resolver, err := localresolver.New(nil)

	if err == nil {
		t.Error("TestLocalResolverGet: NewLocalResolver with nil info should fail")
	}

	resolver, err = localresolver.New(testInfo)

	if err != nil {
		t.Fatal("TestLocalResolverGet: NewLocalResolver failed:", err)
	}

	if (!resolver.ResolverInfo().Matches(&selector{})) {
		t.Fatal("TestLocalResolver: ResolverInfo .Match failed")
	}

	keys := make([]string, 256)
	values := make([]string, 256)

	for i := 0; i < 256; i++ {
		keys[i] = random.RandomString(20)
		values[i] = random.RandomString(20)

		resolver.Put(testContext, entity{keys[i], values[i]})
	}

	for i := 0; i < 256; i++ {
		resolved, err := resolver.Get(testContext, &selector{key: keys[i]})

		if err != nil {
			t.Fatal("TestLocalResolverGet: LocalResolver.Get failed:", err)
		}

		if keys[i] != resolved.(entity).key ||
			values[i] != resolved.(entity).value {
			t.Fatalf("TestLocalResolverGet: Key %v get unexpectedly failed", i)
		}
	}

	for i := 0; i < 256; i++ {
		resolved, err := util.Await(resolver.Resolve(testContext, &selector{key: keys[i]}))

		if err != nil {
			t.Fatal("TestLocalResolverGet: LocalResolver.Get failed:", err)
		}

		if keys[i] != resolved.(entity).key ||
			values[i] != resolved.(entity).value {
			t.Fatalf("TestLocalResolverGet: Key %v get unexpectedly failed", i)
		}
	}

	for i := 128; i < 256; i++ {
		err := resolver.Delete(testContext, &selector{key: keys[i]})

		if err != nil {
			t.Fatal("TestLocalResolverGet: LocalResolver.Delete failed:", err)
		}
	}

	for i := 0; i < 128; i++ {
		resolved, err := resolver.Get(testContext, &selector{key: keys[i]})

		if err != nil {
			t.Fatal("TestLocalResolverGet: LocalResolver.Get failed:", err)
		}

		if keys[i] != resolved.(entity).key ||
			values[i] != resolved.(entity).value {
			t.Fatalf("TestLocalResolverGet: Key %v get unexpectedly failed", i)
		}
	}

	for i := 0; i < 128; i++ {
		resolved, err := util.Await(resolver.Resolve(testContext, &selector{key: keys[i]}))

		if err != nil {
			t.Fatal("TestLocalResolverGet: LocalResolver.Get failed:", err)
		}

		if keys[i] != resolved.(entity).key ||
			values[i] != resolved.(entity).value {
			t.Fatalf("TestLocalResolverGet: Key %v get unexpectedly failed", i)
		}
	}

	for i := 128; i < 256; i++ {
		_, err := resolver.Get(testContext, &selector{key: keys[i]})

		if err == nil {
			t.Fatal("TestLocalResolverGet: LocalResolver.Get after delete succeded at:", i)
		}
	}

	for i := 128; i < 256; i++ {
		_, err := util.Await(resolver.Resolve(testContext, &selector{key: keys[i]}))

		if err == nil {
			t.Fatal("TestLocalResolverGet: LocalResolver.Get after delete succeded at:", i)
		}
	}

	for i := 0; i < 128; i++ {
		err := resolver.Post(testContext, entity{keys[i], values[i+128]})

		if err != nil {
			t.Fatal("TestLocalResolverGet: LocalResolver.Post failed:", err)
		}
	}

	for i := 128; i < 256; i++ {
		err := resolver.Post(testContext, entity{keys[i], values[i-128]})

		if err == nil {
			t.Fatal("TestLocalResolverGet: LocalResolver.Post unexpectedly succeded at:", i)
		}
	}

	for i := 0; i < 128; i++ {
		resolved, err := resolver.Get(testContext, &selector{key: keys[i]})

		if err != nil {
			t.Fatal("TestLocalResolverGet: LocalResolver.Get failed:", err)
		}

		if keys[i] != resolved.(entity).key ||
			values[i+128] != resolved.(entity).value {
			t.Fatalf("TestLocalResolverGet: Key %v get unexpectedly failed", i)
		}
	}
}
