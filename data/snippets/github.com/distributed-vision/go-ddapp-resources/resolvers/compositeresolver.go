package resolvers

import (
	"context"
	"fmt"
	"sync"

	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/util"
)

type componentEntry struct {
	entityType ids.TypeIdentifier
	factory    ResolverFactory
	resolver   Resolver
}

func (this *componentEntry) ResolverInfo() ResolverInfo {
	if this.resolver != nil {
		return this.resolver.ResolverInfo()
	}

	return this.factory.ResolverInfo()
}

type CompositeResolver struct {
	*CachingResolver
	componentMap      map[string][]*componentEntry
	componentMapMutex *sync.Mutex
}

var DefaultCacheSize = 300

func NewCompositeResolver(baseInfo ResolverInfo) (*CompositeResolver, error) {

	base, err := NewCachingResolver(baseInfo, DefaultCacheSize)

	if err != nil {
		return nil, err
	}

	return &CompositeResolver{base,
		make(map[string][]*componentEntry),
		&sync.Mutex{}}, nil
}

func (this *CompositeResolver) RegisterComponentFactory(resolverFactory ResolverFactory, initialiseResolver bool) error {
	this.componentMapMutex.Lock()
	defer this.componentMapMutex.Unlock()

	resolvableTypes := resolverFactory.ResolverInfo().ResolvableTypes()

	if len(resolvableTypes) == 0 {
		return fmt.Errorf("Resolver has no resolvable types")
	}

	var resolver Resolver
	var err error

	if initialiseResolver {
		resolver, err = resolverFactory.New(context.Background())

		if err != nil {
			return err
		}
	}

	for _, resolvableType := range resolvableTypes {
		entry := componentEntry{resolvableType, resolverFactory, resolver}

		if entries, ok := this.componentMap[string(resolvableType.Value())]; ok {
			/*for _, entry := range entries {
				if entry.factory.ResolverType().Equals(resolverFactory.ResolverType()) {
					// ignore duplicate entries
					// TODO look at resolver info to determine duplication
					continue
				}
			}*/

			this.componentMap[string(resolvableType.Value())] = append(entries, &entry)
		} else {
			this.componentMap[string(resolvableType.Value())] = []*componentEntry{&entry}
		}
	}

	return nil
}

func (this *CompositeResolver) RegisterComponent(resolver Resolver) error {
	this.componentMapMutex.Lock()
	defer this.componentMapMutex.Unlock()

	resolvableTypes := resolver.ResolverInfo().ResolvableTypes()

	if len(resolvableTypes) == 0 {
		return fmt.Errorf("Resolver has no resolvable types")
	}

	for _, resolvableType := range resolvableTypes {
		entry := componentEntry{resolvableType, nil, resolver}

		if entries, ok := this.componentMap[string(resolvableType.Value())]; ok {
			/*for _, entry := range entries {
				if entry.factory.ResolverType().Equals(resolverFactory.ResolverType()) {
					// ignore duplicate entries
					// TODO look at resolver info to determine duplication
					continue
				}
			}*/

			this.componentMap[string(resolvableType.Value())] = append(entries, &entry)
		} else {
			this.componentMap[string(resolvableType.Value())] = []*componentEntry{&entry}
		}
	}

	return nil
}
func (this *CompositeResolver) GetMutableComponents(getContext context.Context, selector Selector) []Resolver {
	mutableComponents := []Resolver{}

	if selector != nil {
		this.componentMapMutex.Lock()
		if selector.Type() == nil {
			for _, entries := range this.componentMap {
				for _, entry := range entries {
					if entry.ResolverInfo().Matches(selector) &&
						entry.ResolverInfo().IsMutable() {
						if entry.resolver == nil {
							resolver, err := entry.factory.New(getContext)
							if err != nil {
								entry.resolver = resolver
								mutableComponents = append(mutableComponents, entry.resolver)
							}
						}
						mutableComponents = append(mutableComponents, entry.resolver)
					}
				}
			}
		} else {
			if entries, ok := this.componentMap[string(selector.Type().Value())]; ok {
				for _, entry := range entries {
					if entry.ResolverInfo().Matches(selector) && entry.ResolverInfo().IsMutable() {
						if entry.resolver == nil {
							resolver, err := entry.factory.New(getContext)
							if err != nil {
								entry.resolver = resolver
								mutableComponents = append(mutableComponents, entry.resolver)
							}
						}
						mutableComponents = append(mutableComponents, entry.resolver)
					}
				}
			}
		}
		this.componentMapMutex.Unlock()
	}

	return mutableComponents
}

func (this *CompositeResolver) Get(resolutionContext context.Context, selector Selector) (entity interface{}, err error) {
	return util.Await(this.Resolve(resolutionContext, selector))
}

func (this *CompositeResolver) Resolve(resolutionContext context.Context, selector Selector) (chan interface{}, chan error) {
	cResOut := make(chan interface{}, 1)
	cErrOut := make(chan error, 1)

	if result, err := this.CachingResolver.Get(resolutionContext, selector); err == nil {
		cResOut <- result
		close(cResOut)
		close(cErrOut)
		return cResOut, cErrOut
	}

	var resolverEntries = []*componentEntry{}
	var err error

	if selector != nil {
		this.componentMapMutex.Lock()
		if selector.Type() == nil {
			for _, entries := range this.componentMap {
				for _, entry := range entries {
					//fmt.Println("resolverinfo=", entry.resolver.ResolverInfo())
					if entry.ResolverInfo().Matches(selector) {
						resolverEntries = append(resolverEntries, entry)
					}
				}
			}
		} else {
			if entries, ok := this.componentMap[string(selector.Type().Value())]; ok {
				for _, entry := range entries {
					if entry.ResolverInfo().Matches(selector) {
						resolverEntries = append(resolverEntries, entry)
					}
				}
			} else {
				err = fmt.Errorf("Resolve Failed: no resolver entries for type=%v", selector.Type())
			}
		}
		this.componentMapMutex.Unlock()
	} else {
		err = fmt.Errorf("Resolve Failed: selector cannot be nil")
	}

	if err != nil {
		cErrOut <- err
		close(cResOut)
		close(cErrOut)
		return cResOut, cErrOut
	}

	if len(resolverEntries) == 0 {
		cErrOut <- fmt.Errorf("Resolve Failed: no resolver for entity type=%s", selector.Type())
		close(cResOut)
		close(cErrOut)
		return cResOut, cErrOut
	}

	var wg sync.WaitGroup
	mergeContext, cancel := context.WithCancel(resolutionContext)
	errors := make([]error, 0)
	resultMutex := &sync.Mutex{}
	var result interface{}

	resolve := func(componentEntry *componentEntry) {
		defer wg.Done()

		if componentEntry.resolver == nil {
			resolver, err := componentEntry.factory.New(mergeContext)

			if err != nil {
				errors = append(errors, err)
				return
			}

			componentEntry.resolver = resolver
		}

		cres, cerr := componentEntry.resolver.Resolve(mergeContext, selector)

		if cres == nil || cerr == nil {
			errors = append(errors, fmt.Errorf("Resolve Failed: Resolve channels are undefined"))
		}

		for result == nil {
			select {
			case res, ok := <-cres:
				if ok {
					if result == nil {
						resultMutex.Lock()
						if result == nil {
							result = res
							if key, ok := componentEntry.ResolverInfo().KeyExtractor()(result); ok {
								this.Cache().Add(key, result)
							}
							cancel()
						}
						resultMutex.Unlock()
					}
					return
				}
			case err, ok := <-cerr:
				if ok {
					errors = append(errors, err)
					return
				}
			case <-mergeContext.Done():
				return
			}
		}
	}

	wg.Add(len(resolverEntries))
	for _, r := range resolverEntries {
		go resolve(r)
	}

	go func() {
		wg.Wait()
		if result != nil {
			cResOut <- result
		} else if len(errors) > 0 {
			cErrOut <- fmt.Errorf("Resolve failed with the following errors %v", errors)
		}
		close(cResOut)
		close(cErrOut)
	}()

	return cResOut, cErrOut
}
