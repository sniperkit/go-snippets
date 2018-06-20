package base

import bmsg "mdts/protocols/brokermsg"

type TransToResult struct {
	Method bmsg.EnumMethod
	Head   map[string]string
	Body   []byte
	URL    string
}

type TransFromResult struct {
	Head map[string]string
	Body []byte
}

// Transformer need to be implementec by broker entity.
type Transformer interface {
	ID() string
	TransTo(APICODE string, Data []byte) (*TransToResult, error)
	TransFrom(APICODE string, Data []byte) (*TransFromResult, error)
}
