package vars

import (
	"regexp"
	"strings"
)

// Parse parses a vars string i.e "var_name=value"
func Parse(str string) (map[string]string, error) {
	vars := make(map[string]string)
	regex, err := regexp.Compile(`\s?[a-zA-Z_.-]+\s*=\s*[^\s]+`)

	if err != nil {
		return vars, err
	}

	matches := regex.FindAllString(str, -1)
	for _, match := range matches {
		s := strings.Split(match, "=")
		key, value := strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
		vars[key] = value
	}
	return vars, nil
}
