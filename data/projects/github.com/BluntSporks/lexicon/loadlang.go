package lex

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// LoadLang loads a language file.
func LoadLang(langFile string) map[string]bool {
	// Open file.
	handle, err := os.Open(langFile)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Scan file line by line.
	words := make(map[string]bool)
	scanner := bufio.NewScanner(handle)
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		word := strings.TrimSpace(strings.ToLower(line))
		words[word] = true
	}
	return words
}
