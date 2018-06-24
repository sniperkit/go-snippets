package core

// Example is an interface representing an example instance.
type Example interface {
	// GetExampleValue must return a value for a feature name.
	GetExampleValue(featureName string) interface{}
}

// MapExample represents and Example as key-value pairs.
type MapExample map[string]interface{}

// GetExampleValue implements Example
func (m MapExample) GetExampleValue(featureName string) interface{} { return m[featureName] }
