package pairs

//-----------------------------------------------------
type Marshaler interface {
	Marshal() (data []byte, err error)
}

type Unmarshaler interface {
	Unmarshal(data []byte) (err error)
}

//-----------------------------------------------------
type Serializer interface {
	Serialize() (data []byte, err error)
}

type Deserializer interface {
	Deserialize(data []byte) (err error)
}

//-----------------------------------------------------
type Encoder interface {
	Encode() ([]byte, error)
}

type Decoder interface {
	Decode([]byte) error
}

//-----------------------------------------------------
type Encrypter interface {
	Encrypt() (data []byte, err error)
}

type Decrypter interface {
	Decrypt(data []byte) (err error)
}

//-----------------------------------------------------
type Filler interface {
	Fill(bs []byte) error
}

type Cleaner interface {
	Clean()
}

//-----------------------------------------------------
