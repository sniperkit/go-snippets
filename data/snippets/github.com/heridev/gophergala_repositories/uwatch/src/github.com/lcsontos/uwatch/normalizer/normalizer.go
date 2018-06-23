package normalizer

import (
	"strings"
)

func Normalize(text string) string {
	if text == "" {
		return text
	}

	var normalized = make([]rune, len(text))

	text = strings.ToLower(text)

	for index, runeValue := range text {
		char := runeValue

		switch {
		case runeValue >= '0' && runeValue <= '9':
			normalized[index] = char
		case runeValue >= 'a' && runeValue <= 'z':
			normalized[index] = char
		case runeValue >= 224 && runeValue <= 229:
			normalized[index] = 'a'
		case runeValue >= 232 && runeValue <= 235:
			normalized[index] = 'e'
		case runeValue >= 236 && runeValue <= 239:
			normalized[index] = 'i'
		case runeValue == 241:
			normalized[index] = 'n'
		case (runeValue >= 242 && runeValue <= 246) || (runeValue == 337):
			normalized[index] = 'o'
		case (runeValue >= 249 && runeValue <= 252) || (runeValue == 369):
			normalized[index] = 'u'
		case runeValue == 253:
			normalized[index] = 'y'
		case runeValue == 255:
			normalized[index] = 'y'
		default:
			normalized[index] = '-'
		}

	}

	return string(normalized)
}
