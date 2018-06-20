package conf

import (
	"os"
	"strings"
)

// InitConfMapFromEnv get value of key of confMap from environment,
// if value is not empty, write value to confMap.
func InitConfMapFromEnv(confMap map[string]string) {
	for k, _ := range confMap {
		if v := os.Getenv(strings.ToUpper(k)); v != "" {
			confMap[k] = v
		}
	}

}
