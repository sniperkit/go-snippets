package core

// NewModel initializes a new model
func NewModel(features ...*Feature) *Model {
	m := &Model{Features: make(map[string]*Feature, len(features))}
	for _, f := range features {
		m.Features[f.Name] = f
	}
	return m
}

// Feature returns a feature by name (or nil).
func (m *Model) Feature(name string) *Feature {
	return m.Features[name]
}
