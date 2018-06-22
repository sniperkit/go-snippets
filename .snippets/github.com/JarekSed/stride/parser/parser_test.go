package parser

import "testing"

func linksEqualIgnoringOrder(first []Link, second []Link) bool {
	firstinterface := make([]interface{}, len(first))
	for i, v := range first {
		firstinterface[i] = v
	}
	secondinterface := make([]interface{}, len(second))
	for i, v := range second {
		secondinterface[i] = v
	}
	return itemsEqualIgnoringOrder(firstinterface, secondinterface)
}

func stringsEqualIgnoringOrder(first []string, second []string) bool {
	firstinterface := make([]interface{}, len(first))
	for i, v := range first {
		firstinterface[i] = v
	}
	secondinterface := make([]interface{}, len(second))
	for i, v := range second {
		secondinterface[i] = v
	}
	return itemsEqualIgnoringOrder(firstinterface, secondinterface)
}

// itemsEqualIgnoringOrder return true if the two slices contain identical elements,
// even if they may be in different orders. Duplicates matter, ["a", "a"] != ["a"]
func itemsEqualIgnoringOrder(first []interface{}, second []interface{}) bool {
	// if lengths are equal, one slice must have elements not in the other
	if len(first) != len(second) {
		return false
	}

	items := map[interface{}]int{}
	// Count occurences in first
	for _, item := range first {
		items[item]++
	}
	// Check everything in second
	for _, item := range second {
		// if this wasn't in first (or we've already found a value in second for every value in first),
		// the slices aren't equal
		if items[item] == 0 {
			return false
		}
		items[item]--
	}
	// If the count isn't zero, item was incremented more than decremented (found more in first)
	// or decremented more than incremented (found more in second)
	for _, count := range items {
		if count != 0 {
			return false
		}
	}
	return true
}

func TestMentions(t *testing.T) {
	input := "@chris are you around"
	expected := []string{"chris"}
	found := Mentions(input)
	if !stringsEqualIgnoringOrder(found, expected) {
		t.Errorf("Mentions(%q): expected %q got %q", input, expected, found)
	}
}

func TestMentionsNoMatches(t *testing.T) {
	input := "Yo, are you around"
	expected := []string{}
	found := Mentions(input)
	if !stringsEqualIgnoringOrder(found, expected) {
		t.Errorf("Mentions(%q): expected %q got %q", input, expected, found)
	}
}

func TestMultipleMatches(t *testing.T) {
	input := "@bob @john (success) such a cool feature; https://twitter.com/jdorfman/status/430511497475670016"
	expected := []string{"bob", "john"}
	found := Mentions(input)
	if !stringsEqualIgnoringOrder(found, expected) {
		t.Errorf("Mentions(%q): expected %q got %q", input, expected, found)
	}
}

func TestEmoticons(t *testing.T) {
	input := "Good morning! (megusta) (coffee)"
	expected := []string{"megusta", "coffee"}
	found := Emoticons(input)
	if !stringsEqualIgnoringOrder(found, expected) {
		t.Errorf("Mentions(%q): expected %q got %q", input, expected, found)
	}
}

func TestEmoticonsNoMatches(t *testing.T) {
	input := "Good morning! (this isn't valid) (thisalsoiswaytoolongtobeavalidemoticon)"
	expected := []string{}
	found := Emoticons(input)
	if !stringsEqualIgnoringOrder(found, expected) {
		t.Errorf("Mentions(%q): expected %q got %q", input, expected, found)
	}
}

// TODO(jarek): mock out http calls, so tests are hermetic
func TestLinks(t *testing.T) {
	input := "@bob @john (success) such a cool feature; https://twitter.com/jdorfman/status/430511497475670016"
	expected := []Link{
		{
			URL:   "https://twitter.com/jdorfman/status/430511497475670016",
			Title: "Justin Dorfman; on Twitter: \"nice @littlebigdetail from @HipChat (shows hex colors when pasted in chat). http://t.co/7cI6Gjy5pq\"",
		},
	}
	found, err := Links(input)
	if err != nil {
		t.Errorf("Unexpected error from Links(%q): %v", input, err)
	}
	if !linksEqualIgnoringOrder(found, expected) {
		t.Errorf("Links(%q): expected %q got %q", input, expected, found)
	}
	input = "Olympics are starting soon; http://www.nbcolympics.com"
	expected = []Link{
		{
			URL:   "http://www.nbcolympics.com",
			Title: "2018 PyeongChang Olympic Games | NBC Olympics",
		},
	}
	found, err = Links(input)
	if err != nil {
		t.Errorf("Unexpected error from Links(%q): %v", input, err)
	}
	if !linksEqualIgnoringOrder(found, expected) {
		t.Errorf("Links(%q): expected %q got %q", input, expected, found)
	}
}

// TODO(jarek): clarify expected behavior on overlapping conditions, and test. for example "@(random)"
