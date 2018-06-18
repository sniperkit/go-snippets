package stringutil

import (
	"github.com/ghetzel/uuid"
)

type Uuid struct {
	uuid.UUID
}

func UuidFromBytes(b []byte) (*Uuid, error) {
	if uuid, err := uuid.FromBytes(b); err == nil {
		return &Uuid{uuid}, nil
	} else {
		return nil, err
	}
}

func ParseUUID(in string) (*Uuid, error) {
	if uuid, err := uuid.Parse(in); err == nil {
		return &Uuid{
			UUID: uuid,
		}, nil
	} else {
		return nil, err
	}
}

func MustUUID(in string) *Uuid {
	if uuid, err := uuid.Parse(in); err == nil {
		return &Uuid{
			UUID: uuid,
		}
	} else {
		panic(err)
	}
}

func UUID() *Uuid {
	return &Uuid{
		UUID: uuid.New(),
	}
}

func (self *Uuid) Bytes() []byte {
	return []byte(self.UUID[:])
}
