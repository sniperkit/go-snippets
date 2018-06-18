package pubsub

import "regexp"

type Topic string

type TopicMatcher interface {
	Match(Topic) bool
}

func (t Topic) Match(topic Topic) bool {
	return t == topic
}

//TODO: regex match
type RegexpMatcher regexp.Regexp

func MatchRegexp(expression string) TopicMatcher {
	return (*RegexpMatcher)(regexp.MustCompile(expression))
}

func (m *RegexpMatcher) Match(topic Topic) bool {
	r := (*regexp.Regexp)(m)
	return r.MatchString(string(topic))
}

type allMatcher struct{}

func (*allMatcher) Match(topic Topic) bool {
	return true
}

var MatchAll TopicMatcher = (*allMatcher)(nil)
