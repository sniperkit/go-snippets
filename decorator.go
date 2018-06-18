package mapdecor

// Decor - decorator interface function
type Decor func(map[string]interface{}) map[string]interface{}

// Decorators - gel all Decorators
func Decorators() []Decor {
	registry.mutex.Lock()
	defer registry.mutex.Unlock()
	return registry.decorators
}

// DecorateFromRegistry - decorates input from registry
func DecorateFromRegistry(input map[string]interface{}) map[string]interface{} {
	registry.mutex.Lock()
	decorators := registry.decorators
	registry.mutex.Unlock()

	if len(decorators) == 0 {
		return input
	}
	for _, decorate := range decorators {
		input = decorate(input)
	}
	return input
}

// Decorate - decorate from given decorators
func Decorate(input map[string]interface{}, decorators ...Decor) map[string]interface{} {
	if len(decorators) == 0 {
		return input
	}
	for _, decorate := range decorators {
		input = decorate(input)
	}
	return input
}
