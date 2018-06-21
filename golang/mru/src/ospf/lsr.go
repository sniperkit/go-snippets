package ospf

type LSR struct {
	Header Header
	LSAs []LSAIdentity
}

func UnMarshalLSR(b []byte, length int) (*LSR, error) {
	return nil, nil
}
