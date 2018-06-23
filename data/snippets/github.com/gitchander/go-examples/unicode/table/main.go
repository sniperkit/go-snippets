package main

import (
	"bytes"
	"fmt"
	"os"
)

func GenerateUnicodeTable(fileName string, runeMin, runeMax rune) error {

	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	buffer := new(bytes.Buffer)

	r := runeMin
	for {

		buffer.Reset()
		for j := 0; j < 8; j++ {
			buffer.WriteString(fmt.Sprintf("[%c:U+%x], ", r, r))
			r++

			if r > runeMax {
				break
			}
		}
		buffer.WriteRune('\n')

		f.Write(buffer.Bytes())

		if r > runeMax {
			break
		}
	}

	return nil
}

func main() {

	s := "\u1699 \u03A8 \u03CF \u090B \u0DA3 \u0FD5 \u2230"
	fmt.Println(s)

	GenerateUnicodeTable("unicode_chars.txt", 64, 10074)
}
