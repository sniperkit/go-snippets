package main

// objStringVal returns a string value from a map[string]interface{} or ""
func objStringVal(obj map[string]interface{}, key string) string {
	if val, ok := obj[key].(string); ok {
		return val
	}
	return ""
}
