package signature

import (
	"github.com/distributed-vision/go-resources/ids"
	"github.com/distributed-vision/go-resources/ids/identifier"
)

type signature struct {
	ids.Identifier
	elements ids.SignatureElements
}

func NewSignature(domain ids.SignatureDomain, id []byte) (ids.Signature, error) {
	base, err := identifier.New(domain, id, nil)

	if err != nil {
		return nil, err
	}

	return &signature{base, nil}, nil
}

func (this *signature) Elements() (ids.SignatureElements, error) {
	if this.elements == nil {
		elements, err := GetElements(this)

		if err != nil {
			return nil, err
		}

		this.elements = elements
	}

	return this.elements, nil
}
