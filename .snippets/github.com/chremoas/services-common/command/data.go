package command

import (
	"regexp"
)

func IsDiscordUser(user string) bool {
	var discordUser = regexp.MustCompile(`<@!?\d*>`)
	return discordUser.MatchString(user)
}

func ExtractUserId(user string) string {
	var discordUser = regexp.MustCompile(`^<@!?(\d+)>.*$`)
	return discordUser.FindStringSubmatch(user)[1]
}
