package check

import (
	"github.com/fatih/color"
	"log"
)

func Error(err error) {
	if err != nil {

		// other symbols: ⚠️  ❌  ✘

		errMsg := `
		▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
		▓▓           ✘ ЄɌɌϴɌ ✘           ▓▓
		▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓
		`

		color.Red(errMsg)
		log.Fatal(err)
	}
}
