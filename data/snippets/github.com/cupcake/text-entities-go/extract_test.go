package text

import (
	"io/ioutil"
	"reflect"
	"testing"

	"launchpad.net/goyaml"
)

type conformanceSuite struct {
	Tests struct {
		Hashtags []struct {
			Description string
			Text        string
			Expected    []string
		}
		HashtagIndices []struct {
			Description string
			Text        string
			Expected    []struct {
				Hashtag string
				Indices []int
			}
		} `yaml:"hashtags_with_indices"`

		URLs []struct {
			Description string
			Text        string
			Expected    []string
		}
		URLIndices []struct {
			Description string
			Text        string
			Expected    []struct {
				URL     string
				Indices []int
			}
		} `yaml:"urls_with_indices"`
	}
}

var conformance conformanceSuite

func init() {
	data, err := ioutil.ReadFile("conformance/extract.yml")
	if err != nil {
		panic(err)
	}
	if err := goyaml.Unmarshal(data, &conformance); err != nil {
		panic(err)
	}
}

func TestExtractHashtags(t *testing.T) {
	for _, test := range conformance.Tests.Hashtags {
		res := Extract(test.Text, URLsAndHashtags).Hashtags // extract urls too so that overlapping entities are ignored
		if len(test.Expected) == 0 && len(res) == 0 {
			continue
		}
		hashtags := make([]string, len(res))
		for i, m := range res {
			hashtags[i] = m.Text
		}
		if !reflect.DeepEqual(hashtags, test.Expected) {
			t.Errorf("%s: want %v, got %v", test.Description, test.Expected, res)
		}
	}
}

func TestExtractHashtagIndices(t *testing.T) {
	for _, test := range conformance.Tests.HashtagIndices {
		res := ExtractHashtagMatches(test.Text)
		if len(test.Expected) != len(res) {
			t.Errorf("%s: want %v, got %v", test.Description, test.Expected, res)
			continue
		}
		for i, expected := range test.Expected {
			ei := unicodeToByteOffset(test.Text, [2]int{expected.Indices[0], expected.Indices[1]})
			if res[i].Text != expected.Hashtag || res[i].Indices[0] != ei[0] || res[i].Indices[1] != ei[1] {
				t.Errorf("%s: [%d] want %v, got {%s %v}", test.Description, i, expected, res[i].Text, ei)
			}
		}
	}
}

var skipURLTests = map[int]bool{
	32: true, 33: true, 34: true, // CJK surrounded without protocol
	65: true, 66: true, // special t.co extraction
	26: true, // broken: https://github.com/twitter/twitter-text-conformance/pull/73
}

func TestExtractURLs(t *testing.T) {
	for i, test := range conformance.Tests.URLs {
		if skipURLTests[i] {
			continue
		}
		res := ExtractURLs(test.Text)
		if len(test.Expected) == 0 && len(res) == 0 {
			continue
		}
		if !reflect.DeepEqual(res, test.Expected) {
			t.Errorf("%s: [%d] got %#v, want %#v", test.Description, i, res, test.Expected)
		}
	}
}

var skipURLIndexTests = map[int]bool{
	3: true, 4: true, // CJK surrounded without protocol
	5: true, // special t.co extraction
	8: true, // contains unassigned idn tld
}

func TestExtractURLIndices(t *testing.T) {
	for i, test := range conformance.Tests.URLIndices {
		if skipURLIndexTests[i] {
			continue
		}
		res := ExtractURLMatches(test.Text)
		if len(test.Expected) != len(res) {
			t.Errorf("%s: [%d] want %v, got %v", test.Description, i, test.Expected, res)
			continue
		}
		for j, expected := range test.Expected {
			ei := unicodeToByteOffset(test.Text, [2]int{expected.Indices[0], expected.Indices[1]})
			if res[j].Text != expected.URL || res[j].Indices[0] != ei[0] || res[j].Indices[1] != ei[1] {
				t.Errorf("%s: [%d-%d] want %v, got {%s %v}", test.Description, i, j, expected, res[j].Text, ei)
			}
		}
	}
}

func BenchmarkExtractHashtags(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractHashtagMatches("Getting my Oktoberfest on #münchen")
	}
}

func BenchmarkExtractURLs(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ExtractURLMatches("Extract valid URL: http://google.com/#search?q=iphone%20-filter%3Alinks")
	}
}

func unicodeToByteOffset(s string, charIdx [2]int) (byteIdx [2]int) {
	var j int
	for i := range s {
		j++
		if j == charIdx[0]+1 {
			byteIdx[0] = i
		} else if j > charIdx[1] {
			byteIdx[1] = i
			break
		}
	}
	if byteIdx[1] == 0 {
		byteIdx[1] = len(s)
	}
	return
}
