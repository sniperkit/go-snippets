package zkutils

import (
	"sort"
	"testing"
)

func AssertParseSequenceNodeNotMatch(t *testing.T, name, prefix string) {
	_, err := ParseSequenceNode(name, prefix)

	if err != ErrNodeNotMatch {
		t.Errorf("Expected error parsing sequence node `%s` with prefix `%s` to be ErrNodeNotMatch, but it is: %v", name, prefix, err)
	}
}

func AssertParseSequenceNodeSequenceNumber(t *testing.T, name, prefix string, sequenceNumber int32) {
	sn, err := ParseSequenceNode(name, prefix)
	if err != nil {
		t.Errorf("Unexpected error parsing sequence node `%s` with prefix `%s`: %v", name, prefix, err)
	} else {
		if sn.Name != name {
			t.Errorf("Parsed sequence node name `%s` does not match original name `%s`", sn.Name, name)
		}
		if sn.SequenceNumber != sequenceNumber {
			t.Errorf("Expected sequence number of parsed sequence name `%s` to be %d but it is %d", name, sequenceNumber, sn.SequenceNumber)
		}
	}
}

func TestParseSequenceNode(t *testing.T) {
	// Test without prefix.
	AssertParseSequenceNodeSequenceNumber(t, "0000000000", "", 0)
	AssertParseSequenceNodeSequenceNumber(t, "0000000001", "", 1)
	AssertParseSequenceNodeSequenceNumber(t, "-2147483648", "", -2147483648)
	AssertParseSequenceNodeSequenceNumber(t, "2147483647", "", 2147483647)
	AssertParseSequenceNodeSequenceNumber(t, "prefix_a_0000000000", "", 0)
	AssertParseSequenceNodeSequenceNumber(t, "prefix_a_0000000001", "", 1)
	AssertParseSequenceNodeSequenceNumber(t, "prefix_a_-2147483648", "", -2147483648)
	AssertParseSequenceNodeSequenceNumber(t, "prefix_a_2147483647", "", 2147483647)

	AssertParseSequenceNodeNotMatch(t, "", "")
	AssertParseSequenceNodeNotMatch(t, "000000000a", "")

	// Test with prefix.
	for _, prefix := range []string{
		"",
		"_guid_-",
		"/my/sub/dir/",
		"/my/sub/dir/_guid_",
	} {
		AssertParseSequenceNodeSequenceNumber(t, prefix+"prefix_a_0000000000", "prefix_a_", 0)
		AssertParseSequenceNodeSequenceNumber(t, prefix+"prefix_a_0000000001", "prefix_a_", 1)
		AssertParseSequenceNodeSequenceNumber(t, prefix+"prefix_a_-2147483648", "prefix_a_", -2147483648)
		AssertParseSequenceNodeSequenceNumber(t, prefix+"prefix_a_2147483647", "prefix_a_", 2147483647)

		AssertParseSequenceNodeNotMatch(t, prefix+"prefix_a_0000000000", "prefix_b_")
		AssertParseSequenceNodeNotMatch(t, prefix+"prefix_a_0000000001", "prefix_b_")
		AssertParseSequenceNodeNotMatch(t, prefix+"prefix_a_-2147483648", "prefix_b_")
		AssertParseSequenceNodeNotMatch(t, prefix+"prefix_a_2147483647", "prefix_b_")
	}

	AssertParseSequenceNodeNotMatch(t, "0000000000", "prefix_a_")
	AssertParseSequenceNodeNotMatch(t, "0000000001", "prefix_a_")
	AssertParseSequenceNodeNotMatch(t, "-2147483648", "prefix_a_")
	AssertParseSequenceNodeNotMatch(t, "2147483647", "prefix_a_")
}

func AssertParseSequenceNodes(t *testing.T, names []string, prefix string, expected []SequenceNode) {
	actual := ParseSequenceNodes(names, prefix)

	if len(expected) != len(actual) {
		t.Errorf("Expected %d parsed node(s) for %v, but got %d", len(expected), names, len(actual))
	}

	compareCount := len(expected)
	if compareCount > len(actual) {
		compareCount = len(actual)
	}

	for i := 0; i < compareCount; i++ {
		a := actual[i]
		e := expected[i]

		if !a.Equals(e) {
			t.Errorf("Expected parsed node #%d of %v to be %v, but it is %v", i, names, e, a)
		}
	}
}

func TestParseSequenceNodes(t *testing.T) {
	names := []string{
		"0000000000",
		"0000000001",
		"-2147483648",
		"2147483647",
		"prefix_a_0000000000",
		"prefix_a_0000000001",
		"prefix_a_-2147483648",
		"prefix_a_2147483647",
		"prefix_b_0000000000",
		"prefix_b_0000000001",
		"prefix_b_-2147483648",
		"prefix_b_2147483647",
		"",
		"000000000a",
	}

	for _, fixture := range []struct {
		Prefix   string
		Expected []SequenceNode
	}{
		{"", []SequenceNode{
			{"0000000000", 0},
			{"0000000001", 1},
			{"-2147483648", -2147483648},
			{"2147483647", 2147483647},
			{"prefix_a_0000000000", 0},
			{"prefix_a_0000000001", 1},
			{"prefix_a_-2147483648", -2147483648},
			{"prefix_a_2147483647", 2147483647},
			{"prefix_b_0000000000", 0},
			{"prefix_b_0000000001", 1},
			{"prefix_b_-2147483648", -2147483648},
			{"prefix_b_2147483647", 2147483647},
		}},
		{"prefix_a", []SequenceNode{}},
		{"prefix_a_", []SequenceNode{
			{"prefix_a_0000000000", 0},
			{"prefix_a_0000000001", 1},
			{"prefix_a_-2147483648", -2147483648},
			{"prefix_a_2147483647", 2147483647},
		}},
	} {
		for _, prefix := range []string{
			"",
			"_guid_",
			"/my/sub/dir/",
			"/my/sub/dir/_guid_",
		} {
			prefixedNames := make([]string, len(names))
			for i, n := range names {
				prefixedNames[i] = prefix + n
			}

			prefixedExpected := make([]SequenceNode, len(fixture.Expected))
			for i, e := range fixture.Expected {
				prefixedExpected[i] = SequenceNode{
					Name:           prefix + e.Name,
					SequenceNumber: e.SequenceNumber,
				}
			}

			AssertParseSequenceNodes(t, prefixedNames, fixture.Prefix, prefixedExpected)
		}
	}
}

func TestSequenceNodesAscendinglyNegativeLast(t *testing.T) {
	expected := []SequenceNode{
		{"", 0},
		{"", 1},
		{"", 2},
		{"", 3},
		{"", -19},
		{"", -18},
		{"", -17},
	}

	actual := []SequenceNode{
		{"", -19},
		{"", 3},
		{"", 0},
		{"", 2},
		{"", -17},
		{"", 1},
		{"", -18},
	}
	sort.Sort(sequenceNodesAscendinglyNegativeLast(actual))

	for i, a := range actual {
		e := expected[i]

		if !a.Equals(e) {
			t.Errorf("Expected sorted node #%d to be %v but is %v", i, e, a)
		}
	}
}

func TestSortSequenceNodes(t *testing.T) {
	for _, fixture := range []struct {
		Nodes    []SequenceNode
		Expected []SequenceNode
	}{
		{
			[]SequenceNode{
				{"", 3},
				{"", 0},
				{"", 2},
				{"", 1},
			},
			[]SequenceNode{
				{"", 0},
				{"", 1},
				{"", 2},
				{"", 3},
			},
		},

		{
			[]SequenceNode{
				{"", -2147483646},
				{"", -2147483645},
				{"", -2147483648},
				{"", -2147483647},
			},
			[]SequenceNode{
				{"", -2147483648},
				{"", -2147483647},
				{"", -2147483646},
				{"", -2147483645},
			},
		},

		{
			[]SequenceNode{
				{"", -2147483646},
				{"", -2147483645},
				{"", 2147483647},
				{"", -2147483648},
				{"", -2147483647},
				{"", 2147483646},
				{"", -1073741824},
			},
			[]SequenceNode{
				{"", 2147483646},
				{"", 2147483647},
				{"", -2147483648},
				{"", -2147483647},
				{"", -2147483646},
				{"", -2147483645},
				{"", -1073741824},
			},
		},

		{
			[]SequenceNode{
				{"", -1073741822},
				{"", -1073741823},
				{"", -1073741821},
			},
			[]SequenceNode{
				{"", -1073741823},
				{"", -1073741822},
				{"", -1073741821},
			},
		},

		{
			[]SequenceNode{
				{"", 2147483647},
				{"", 0},
				{"", -1073741822},
				{"", -1073741823},
				{"", -1073741821},
			},
			[]SequenceNode{
				{"", -1073741823},
				{"", -1073741822},
				{"", -1073741821},
				{"", 0},
				{"", 2147483647},
			},
		},
	} {
		SortSequenceNodes(fixture.Nodes)

		for i, a := range fixture.Nodes {
			e := fixture.Expected[i]

			if !a.Equals(e) {
				t.Errorf("Expected sorted node #%d to be %v but is %v", i, e, a)
			}
		}
	}
}
