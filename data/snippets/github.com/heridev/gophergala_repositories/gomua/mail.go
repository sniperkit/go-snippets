package gomua

import "fmt"

// Mail interface defines mail to be read.
// This can be a threaded list of messages, or a single message.
//
// String() should return a string representing what the user needs to see
// to further interact with the Mail. For a single message, this should be
// simply the message. For a threaded list of messages, this might be a list
// of all messages in the thread, or possibly a just newest unread message?
//
// Summary() should return a short string that can fit on one line, suitable
// for printing as 1 item in a list on a screen.
// For instance, for a single message:
//     Test Message from Frenata <mr.k.frenata@gmail.com>
type Mail interface {
	String() string
	Summary() string
}

func (m *Message) String() string {
	var output string = fmt.Sprintf("From: %v\n", m.Header.Get("From")) +
		fmt.Sprintf("To: %v\n", m.Header.Get("To")) +
		fmt.Sprintf("Date: %v\n", m.Header.Get("Date")) +
		fmt.Sprintf("Subject: %v\n", m.Header.Get("Subject"))

		//output += fmt.Sprintf("\n%s\n", m.Content)

	output += fmt.Sprintf("\n%s\n", m.SanitizeContent())

	return output
}

func (m *Message) Summary() string {
	subject := m.Header.Get("Subject")
	from := m.Header.Get("From")
	return fmt.Sprintf("%s from %s", color(subject, "31"), color(from, "33"))
}

// adds ANSI color to text
func color(s string, color string) string {
	return "\033[" + color + "m" + s + "\033[0m"
}
