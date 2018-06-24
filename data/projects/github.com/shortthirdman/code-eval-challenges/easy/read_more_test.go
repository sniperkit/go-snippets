package main

import "testing"

func TestReadMore(t *testing.T) {
	for k, v := range map[string]string{
		"Tom exhibited.": "Tom exhibited.",
		"Amy Lawrence was proud and glad, and she tried to make Tom see it in her face - but he wouldn't look.": "Amy Lawrence was proud and glad, and... <Read More>",
		"Tom was tugging at a button-hole and looking sheepish.":                                                "Tom was tugging at a button-hole and looking sheepish.",
		"Two thousand verses is a great many - very, very great many.":                                          "Two thousand verses is a great many -... <Read More>",
		"Tom's mouth watered for the apple, but he stuck to his work.":                                          "Tom's mouth watered for the apple, but... <Read More>",
		"AnnaElsaOlafKristoffSvenHansDukeOfWeseltonOakenMarshmallow":                                            "AnnaElsaOlafKristoffSvenHansDukeOfWeselt... <Read More>"} {
		if r := readMore(k); r != v {
			t.Errorf("failed: readMore %s is %s, got %s",
				k, v, r)
		}
	}
}

func readMore(q string) string {
	if len(q) > 55 {
		for i := 39; i >= 0; i-- {
			if q[i] == ' ' {
				q = q[:i]
				break
			}
		}
		if len(q) > 55 {
			q = q[:40]
		}
		q += "... <Read More>"
	}
	return q
}
