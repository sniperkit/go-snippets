package responses

//BuilderFactory just creates new Builder's
type BuilderFactory interface {
	Builder() Builder
}

//NewBuilderFactory creates a new BuilderFactory
func NewBuilderFactory(format Format) BuilderFactory {
	return &builderFactory{
		format: format,
	}
}

type builderFactory struct {
	format Format
}

func (r *builderFactory) Builder() Builder {
	return &builder{
		format: r.format,
	}
}
