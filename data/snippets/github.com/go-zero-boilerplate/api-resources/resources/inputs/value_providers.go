package inputs

//QueryValueProvider provides a Value from the given Query
type QueryValueProvider interface {
	Value(input Query) Value
}

type queryValueProvider struct {
	key string
}

func (q *queryValueProvider) Value(input Query) Value {
	return input.Query(q.key)
}

//ParamValueProvider provides a Value from the given Param
type ParamValueProvider interface {
	Value(input Param) Value
}

type paramValueProvider struct {
	key string
}

func (p *paramValueProvider) Value(input Param) Value {
	return input.Param(p.key)
}
