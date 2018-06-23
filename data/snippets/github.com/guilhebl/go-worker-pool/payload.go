package job

// Data a Job receives as Param
type Payload struct {
	Params map[string]string
}

func NewPayload(params map[string]string) Payload {
	return Payload{
		Params: params,
	}
}
