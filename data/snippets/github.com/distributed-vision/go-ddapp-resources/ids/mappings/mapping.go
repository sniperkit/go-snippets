package mappings

import (
	"fmt"
	"time"

	"github.com/distributed-vision/go-resources/ids"
)

func KeyExtractor(entity ...interface{}) (interface{}, bool) {
	if len(entity) > 0 {
		if mappings, ok := entity[0].(Mappings); ok {
			return mappings.fromId.String() + "->" + mappings.toDomain.String(), true
		}
	}
	return nil, false
}

func AwaitMapping(cres chan ids.Mapping, cerr chan error) (result ids.Mapping, err error) {
	if cres == nil || cerr == nil {
		return nil, fmt.Errorf("Await Failed: channels are undefined")
	}

	resolved := false
	for !resolved {
		select {
		case res, ok := <-cres:
			if ok {
				result = res
				resolved = true
			}
		case error, ok := <-cerr:
			if ok {
				err = error
				resolved = true
			}
		}
	}

	return result, err
}

var MaxTime = time.Unix(1<<63-62135596801, 999999999)
var MinTime = time.Time{}

type mappedId struct {
	from time.Time
	to   time.Time
	id   ids.Identifier
}

type Mappings struct {
	fromId    ids.Identifier
	toDomain  ids.IdentityDomain
	mappedIds []mappedId
}

type mapping struct {
	mappings *Mappings
	index    int
}

func (this *mapping) FromId() ids.Identifier {
	return this.mappings.fromId
}

func (this *mapping) ToDomain() ids.IdentityDomain {
	return this.mappings.toDomain
}

func (this *mapping) From() time.Time {
	return this.mappings.mappedIds[this.index].from
}

func (this *mapping) To() time.Time {
	return this.mappings.mappedIds[this.index].to
}

func (this *mapping) ToId() ids.Identifier {
	return this.mappings.mappedIds[this.index].id
}
