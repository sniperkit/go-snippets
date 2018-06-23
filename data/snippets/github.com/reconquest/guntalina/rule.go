package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Rule struct {
	Group     string   `yaml:"group"`
	Masks     []string `yaml:"masks"`
	Workflow  []string `yaml:"workflow"`
	regexps   map[string]*regexp.Regexp
	processed bool
}

type Rules []*Rule

// SetProcessed sets original rule processed value to true, if rule has a
// group, than all rules with same rule will be setted to processed too.
func (rules *Rules) SetProcessed(rule *Rule) {
	rule.setProcessed()

	if rule.Group != "" {
		for _, item := range *rules {
			if item.Group == item.Group {
				item.setProcessed()
			}
		}
	}
}

func (rule *Rule) setProcessed() {
	rule.processed = true
}

func (rules *Rules) Compile() error {
	for _, rule := range *rules {
		for _, mask := range rule.Masks {
			glob := strings.Replace(mask, "*", "_GROD_", -1)
			pattern := strings.Replace(
				regexp.QuoteMeta(glob), "_GROD_", "([^/]*)", -1,
			)

			re, err := regexp.Compile(pattern)
			if err != nil {
				return fmt.Errorf(
					"can't create regexp for mask '%s': %s", mask, err,
				)
			}

			rule.regexps[mask] = re
		}
	}

	return nil
}

func (rule *Rule) Match(modification string) bool {
	for _, re := range rule.regexps {
		match := re.MatchString(modification)
		if match {
			return true
		}
	}

	return false
}
