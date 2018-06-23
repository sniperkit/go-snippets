package types

import (
	"errors"
	"fmt"
	"sync"

	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/identifier"
	"github.com/distributed-vision/go-resources/version"
)

var typeIds = make(map[string]ids.TypeIdentifier)
var typeIdsMutex = &sync.Mutex{}

func MustNewId(domain interface{}, id []byte, idVersion version.Version) ids.TypeIdentifier {
	typeid, err := NewId(domain, id, idVersion)

	if err != nil {
		panic(fmt.Sprintf("Can't create new type id: %s", err))
	}

	return typeid
}

func NewId(domain interface{}, id []byte, idVersion version.Version) (ids.TypeIdentifier, error) {

	base, err := identifier.New(domain, id, idVersion)

	if err != nil {
		return nil, err
	}

	typeId := &typeIdentifier{base}

	typeIdsMutex.Lock()
	defer typeIdsMutex.Unlock()
	if registeredId, ok := typeIds[string(typeId.Value())]; ok {
		return registeredId, nil
	}

	typeIds[string(typeId.Value())] = typeId
	return typeId, nil
}

type typeIdentifier struct {
	ids.Identifier
}

func (this *typeIdentifier) Elements() (ids.SignatureElements, error) {
	return nil, errors.New("TODO")
}

func (this *typeIdentifier) IsAssignableFrom(typeId ids.TypeIdentifier) bool {
	return false
}

func RegisterId(typeid ids.TypeIdentifier) {
	typeIdsMutex.Lock()
	typeIds[string(typeid.Value())] = typeid
	typeIdsMutex.Unlock()
}
