package ospf

type LSAck struct {
	Header Header
}

func UnMarshalLSAck(b []byte, length int) (*LSAck, error) {
	return nil, nil
}
