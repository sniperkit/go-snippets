package storage

type Language struct {
	Code   string
	Native string
}

var languageMap map[string]string = map[string]string{
	"de": "Deutsch",
	"eo": "Esperanto",
	"en": "English",
	"fr": "Français",
}

var Languages map[string]Language

func init() {
	Languages = make(map[string]Language, len(languageMap))

	for code, native := range languageMap {
		Languages[code] = Language{
			Code:   code,
			Native: native,
		}
	}
}
