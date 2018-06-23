package tools

import "testing"

func TestParseGitUrl(t *testing.T) {

	if parseGitUrl("https://github.com/ChappIO/dashfiles") != "https://github.com/ChappIO/dashfiles" {
		t.Error("A full https url did not remain unchanged")
	}

	if parseGitUrl("ChappIO/dashfiles") != "git@github.com:ChappIO/dashfiles" {
		t.Error("A github shorthand did not expand to an ssh url")
	}

	if parseGitUrl("git@bitbucket.org:Test/test") != "git@bitbucket.org:Test/test" {
		t.Error("A bitbucket ssh url did not remain unchanged")
	}
}
