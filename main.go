package main

import (
	"fmt"

	"github.com/golang-plus/i18n"
	_ "github.com/golang-plus/i18n-resources/en"
	_ "github.com/golang-plus/i18n-resources/zh-hans"
)

func languages() {
	for _, lang := range i18n.AllLanguages() {
		fmt.Printf("language: %s\n", lang.Code)
		fmt.Println("==> names")
		langs := lang.Name.SupportedLanguages()
		for _, lang2 := range langs {
			name := lang.Name.Value(lang2)
			fmt.Printf("  %s (language: %s)\n", name, lang2.NativeName)
		}
	}
}

func currencies() {
	for _, curr := range i18n.AllCurrencies() {
		fmt.Printf("currency: %s\n", curr.Code)
		fmt.Println("==> names")
		langs := curr.Name.SupportedLanguages()
		for _, lang := range langs {
			name := curr.Name.Value(lang)
			fmt.Printf("  %s (language: %s)\n", name, lang.NativeName)
		}
	}
}

func countries() {
	for _, coun := range i18n.AllCountries() {
		fmt.Printf("country: %s %s %s\n", coun.Alpha2Code, coun.Alpha3Code, coun.NumericCode)
		fmt.Println("==> names")
		langs := coun.Name.SupportedLanguages()
		for _, lang := range langs {
			name := coun.Name.Value(lang)
			fmt.Printf("  %s (language: %s)\n", name, lang.NativeName)
		}

		fmt.Println("==> aliases")
		langs = coun.Aliases.SupportedLanguages()
		for _, lang := range langs {
			aliases := coun.Aliases.Values(lang)
			fmt.Printf("  %v (language: %s)\n", aliases, lang.NativeName)
		}
	}
}

func cultures() {
	for _, cul := range i18n.AllCultures() {
		fmt.Printf("culture: %s, native name: %s\n", cul.Code, cul.NativeName)
		fmt.Println("==> names")
		langs := cul.Name.SupportedLanguages()
		for _, lang := range langs {
			name := cul.Name.Value(lang)
			fmt.Printf("  %s (language: %s)\n", name, lang.NativeName)
		}
	}
}

func main() {
	languages()
	currencies()
	countries()
	cultures()
}
