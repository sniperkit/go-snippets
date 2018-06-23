package hello

import (
	"log"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func extendMap(a map[string]Gif, b map[string]Gif) {
	for k, v := range b {
		a[k] = v
	}
}

func extractKeys(a map[string]Gif) []string {
	keys := make([]string, 0, len(a))
	for i, _ := range a {
		keys = append(keys, i)
	}
	return keys
}
