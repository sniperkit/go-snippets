package main

import (
	"regexp"
	"testing"
)

func TestEmail(t *testing.T) {
	for k, v := range map[string]bool{
		"a\"b(c)d,e:f;gi[j\\k]l@example.com": false,
		"":                                           false,
		"this is not true":                           false,
		"Abc.example.com":                            false,
		"this\\ still\\\"not\\\\allowed@example.com": false,
		"just\"not\"right@example.com":               false,
		"niceandsimple@example.com":                  true,
		"a.little.lengthy.but.fine@dept.example.com": true,
		"bob123@alice123.com":                        true,
		"this is\"not\\allowed@example.com":          false,
		"b@d.eu": true,
		"disposable.style.email.with+156@example.com": true,
		"very.common@example.com":                     true,
		"\"very.unusual.@.unusual.com\"@example.com":  true,
		"A@b@c@example.com":                           false,
		"1@d.eu":                                      true,
		"hfij#kjdfvkl":                                false,
		"one_char_domain@foo.c":                       false,
		"b@domain.eu":                                 true,
		"@foo.com":                                    false,
		"disposable.style.email.with+symbol@example.com": true} {
		if r := email(k); r != v {
			t.Errorf("failed: email %s is %t, got %t",
				k, v, r)
		}
	}
}

var emailRegex *regexp.Regexp

func init() {
	emailRegex = regexp.MustCompile(`^((\"[0-9A-Za-z@.+-=]+\")|([0-9A-Za-z.+-=]+))@\w*(\w+\.)+\w{2,4}$`)
}

func email(q string) bool {
	return emailRegex.MatchString(q)
}
