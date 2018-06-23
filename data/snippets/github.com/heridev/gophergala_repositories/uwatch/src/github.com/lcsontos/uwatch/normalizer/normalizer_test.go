package normalizer

import (
	"testing"
)

type testCase struct {
	want string
	text string
}

var testCaseData = []testCase{
	{
		"zen---musica-de-relajacion-y-balance-espiritual",
		"ZEN - MÚSICA DE RELAJACIÓN Y BALANCE ESPIRITUAL",
	},
	{
		"meditations-with-sri-chinmoy-vol--1",
		"Meditations with Sri Chinmoy Vol. 1",
	},
	{
		"ouooueaui",
		"öüóőúéáűí",
	},
}

func TestNormalize(t *testing.T) {
	for _, data := range testCaseData {
		text := data.text

		got, want := Normalize(text), data.want

		// TODO manual test succeeds, but this fails, check later why?
		if got != want {
			t.Fatalf("Normalize(%s) = %s, but wanted %s", text, got, want)
		}
	}
}
