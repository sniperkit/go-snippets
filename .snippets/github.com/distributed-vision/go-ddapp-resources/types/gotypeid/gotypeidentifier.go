package gotypeid

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/distributed-vision/go-resources/encoding/encodertype"
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/domain"
	"github.com/distributed-vision/go-resources/ids/identifier"
	"github.com/distributed-vision/go-resources/types"
)

var identifiers = make(map[reflect.Type]ids.TypeIdentifier)
var mapMutex = &sync.Mutex{}

func init() {
	ids.LocalTypeInit(IdOf, types.NewId)
}

func IdOf(gotype reflect.Type) ids.TypeIdentifier {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	var err error

	id, ok := identifiers[gotype]

	if !ok {
		id, err = newIdFor(gotype)

		if err != nil {
			panic(fmt.Sprintf("Can't create type id for:%+v\n", err))
		}
		identifiers[gotype] = id
	}
	//fmt.Printf("type(%s):%+v\n", gotype, id)
	return id
}

type GoTypeIdentifier struct {
	ids.Identifier
	TypeOf reflect.Type
}

func (this *GoTypeIdentifier) Elements() (ids.SignatureElements, error) {
	return nil, errors.New("TODO")
}

func (this *GoTypeIdentifier) IsAssignableFrom(typeId ids.TypeIdentifier) bool {
	return false
}

var goTypeDomain []byte = domain.MustDecodeId(encodertype.BASE62, "2", "0")

func newIdFor(gotype reflect.Type) (ids.TypeIdentifier, error) {

	base, err := identifier.New(goTypeDomain, []byte(gotype.String()), nil)

	if err != nil {
		return nil, err
	}

	id := &GoTypeIdentifier{base, gotype}

	types.RegisterId(id)

	return id, nil
}
