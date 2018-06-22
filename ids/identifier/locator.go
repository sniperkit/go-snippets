package identifier

import (
	"github.com/distributed-vision/go-resources/ids"
	lru "github.com/hashicorp/golang-lru"
)

type locator struct {
	ids.Identifier
	entity interface{}
}

func (loc *locator) Get() (interface{}, error) {
	return nil, nil
}

func (loc *locator) GetAs(typeId ids.TypeIdentifier) (interface{}, error) {
	return nil, nil
}

func (loc *locator) Resolve() (chan interface{}, chan error) {
	return nil, nil
}

func (loc *locator) ResolveAs(typeId ids.TypeIdentifier) (chan interface{}, chan error) {
	return nil, nil
}

func NewLocator(id ids.Identifier) ids.Locator {
	return &locator{id, nil}
}

var locators, _ = lru.NewARC(500)

func getLocator(id ids.Identifier) ids.Locator {

	locatorKey := string(id.Value())
	locator, ok := locators.Get(locatorKey)

	if ok {
		return locator.(ids.Locator)
	}

	locator = NewLocator(id)
	locators.Add(locatorKey, locator)
	return locator.(ids.Locator)
}
