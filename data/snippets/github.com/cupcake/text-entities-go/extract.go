package text

import (
	"sort"
)

type Match struct {
	Text    string
	Indices [2]int
}

func ExtractHashtags(s string) []string {
	hashtags := ExtractHashtagMatches(s)
	res := make([]string, len(hashtags))
	for i, tag := range hashtags {
		res[i] = tag.Text
	}
	return res
}

func ExtractHashtagMatches(s string) []Match {
	matches := hashtagPattern.FindAllStringSubmatchIndex(s, -1)
	res := make([]Match, 0, len(matches))
	for _, m := range matches {
		if invalidHashtagEnd.MatchString(s[m[3]:]) {
			continue
		}
		res = append(res, Match{s[m[4]:m[5]], [2]int{m[2], m[3]}})
	}
	return res
}

func ExtractURLs(s string) []string {
	urls := ExtractURLMatches(s)
	res := make([]string, len(urls))
	for i, url := range urls {
		res[i] = url.Text
	}
	return res
}

func ExtractURLMatches(s string) []Match {
	matches := urlPattern.FindAllStringSubmatchIndex(s, -1)
	res := make([]Match, 0, len(matches))
	for _, m := range matches {
		var before, protocol, path string
		if m[2] != -1 {
			before = s[m[2]:m[3]]
		}
		if m[6] != -1 {
			protocol = s[m[6]:m[7]]
		}
		domain := s[m[8]:m[9]]
		afterDomain := s[m[9]:]
		if m[10] != -1 {
			path = s[m[10]:m[11]]
		}

		// We don't have lookahead, so manually check that the tld doesn't end in the middle of a word
		if afterDomain != "" && ((afterDomain[0] >= '0' && afterDomain[0] <= '9') ||
			(afterDomain[0] >= 'A' && afterDomain[0] <= 'Z') ||
			(afterDomain[0] >= 'a' && afterDomain[0] <= 'z') ||
			afterDomain[0] == '@') {
			continue
		}

		if protocol == "" {
			if invalidBeforeDomain.MatchString(before) {
				continue
			}
			ad := asciiDomain.FindString(domain)
			if len(domain) != len(ad) || invalidWithoutPath.MatchString(ad) && path == "" {
				continue
			}
		}
		res = append(res, Match{s[m[4]:m[5]], [2]int{m[4], m[5]}})
	}
	return res
}

type Entities struct {
	Hashtags []Match
	URLs     []Match
}

const (
	FlagURLs Flag = 1 << iota
	FlagHashtags
	FlagOverlapping
)

var URLsAndHashtags = FlagURLs | FlagHashtags

type Flag int

type matchInfo struct {
	Indices [2]int
	Type    Flag
	Index   int
}

type matchInfos []matchInfo

func (m matchInfos) Len() int           { return len(m) }
func (m matchInfos) Less(i, j int) bool { return m[i].Indices[0] < m[j].Indices[0] }
func (m matchInfos) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }

func Extract(s string, flags Flag) *Entities {
	var res Entities
	if flags&FlagURLs != 0 {
		res.URLs = ExtractURLMatches(s)
	}
	if flags&FlagHashtags != 0 {
		res.Hashtags = ExtractHashtagMatches(s)
	}
	if flags&FlagOverlapping == 0 && len(res.URLs) > 0 && len(res.Hashtags) > 0 {
		matches := make(matchInfos, len(res.Hashtags)+len(res.URLs))
		for i, m := range res.Hashtags {
			matches[i] = matchInfo{m.Indices, FlagHashtags, i}
		}
		for i, m := range res.URLs {
			matches[i+len(res.Hashtags)] = matchInfo{m.Indices, FlagURLs, i}

		}

		sort.Sort(matches)
		for i, m := range matches {
			if i > 0 && matches[i-1].Indices[1] > m.Indices[0] {
				switch m.Type {
				case FlagURLs:
					res.URLs = deleteMatch(res.URLs, m.Index)
				case FlagHashtags:
					res.Hashtags = deleteMatch(res.Hashtags, m.Index)
				}
			}
		}
	}
	return &res
}

func deleteMatch(m []Match, i int) []Match {
	if i == len(m) {
		return m[:i-1]
	}
	return append(m[:i], m[i+1:]...)
}
