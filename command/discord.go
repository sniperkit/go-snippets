package command

import "fmt"

func makeString(message string, sender string, sign string) string {
	var output string

	if len(sender) > 0 {
		output = fmt.Sprintf("<@%s> ", sender)
	}
	output = fmt.Sprintf("%s %s %s", output, sign, message)

	return output
}

func SendSuccess(message string, sender ...string) string {
	var s string

	if len(sender) == 0 {
		s = ""
	} else {
		s = sender[0]
	}
	return makeString(message, s, ":white_check_mark:")
}

func SendError(message string, sender ...string) string {
	var s string

	if len(sender) == 0 {
		s = ""
	} else {
		s = sender[0]
	}
	return makeString(message, s, ":warning:")
}

func SendFatal(message string, sender ...string) string {
	var s string

	if len(sender) == 0 {
		s = ""
	} else {
		s = sender[0]
	}
	return makeString(message, s, ":octagonal_sign:")
}
