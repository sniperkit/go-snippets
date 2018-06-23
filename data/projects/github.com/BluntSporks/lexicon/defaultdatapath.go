package lex

import (
	"log"
	"os"
)

// DefaultDataPath gets the lexicon word list file location from the LEX_DATA environment variable.
func DefaultDataPath() string {
	dir := os.Getenv("LEX_DATA")
	if dir == "" {
		log.Fatal("Set LEX_DATA variable to directory of lexicon data files")
	}
	return dir
}
