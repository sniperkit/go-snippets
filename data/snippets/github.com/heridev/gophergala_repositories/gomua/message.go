package gomua

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"strings"
)

type Message struct {
	mail.Message
	Content  string
	Filename string
}

func (m *Message) store() {
	b, err := ioutil.ReadAll(m.Message.Body)
	if err != nil {
		log.Fatal(err)
	}
	m.Content = string(b)
}

// ReadMessage embeds a mail.Message inside the gomua.Message and stores the Body content
func ReadMessage(msg *mail.Message) *Message {
	m := new(Message)
	m.Message = *msg
	m.store()

	return m
}

// Flag sets a filename flag on the message.
func (m *Message) Flag(flag string) {
	s := strings.Split(m.Filename, ":2")
	if len(s) != 2 {
		log.Fatal(fmt.Errorf("filename %s does not contain ':2'", m.Filename))
	}
	name := s[0]
	flags := s[1]

	if !strings.Contains(flags, flag) {
		if flags[len(flags)-1] != ',' {
			flags += ","
		}
		flags += flag

		newname := name + ":2" + flags
		os.Rename(m.Filename, newname)
		m.Filename = newname
	}
}

// IsFlagged checks if the filename of the message is flagged with the given flag.
func (m *Message) IsFlagged(flag string) bool {
	s := strings.Split(m.Filename, ":2")
	if len(s) != 2 {
		log.Fatal(fmt.Errorf("filename %s does not contain ':2'", m.Filename))
	}
	flags := s[1]

	return strings.Contains(flags, flag)
}

// Unread returns true if the message is unread.
func (m *Message) Unread() bool {
	return !m.IsFlagged("S")
}

// SanitizeContent returns only text/plain content portions of the email.
// --Probably not fully working.--
func (m *Message) SanitizeContent() string {
	var bound string
	var boundB bool = false
	t := m.Header.Get("Content-Type")
	if strings.Contains(t, "boundary=") {
		bs := strings.Split(t, "boundary=")
		bound = "--" + bs[1]
		boundB = true
	}

	raw := bufio.NewScanner(strings.NewReader(m.Content))
	buf := new(bytes.Buffer)
	var write bool = true
	for raw.Scan() {
		line := raw.Text()
		if strings.Contains(line, "Content-Type:") {
			write = false
		}

		if !strings.Contains(line, "Content-Transfer-Encoding") {
			if !boundB || boundB && line != bound {
				if write {
					buf.WriteString(line + "\r\n")
				}
			}
		}

		if strings.Contains(line, "Content-Type: text/plain") {
			write = true
		}
	}
	return string(buf.Bytes())
}

// WriteMessage interactively prompts the user for an email to send.
func WriteMessage(r io.Reader) *mail.Message {
	cli := bufio.NewScanner(r)

	fmt.Print("To: ")
	cli.Scan()
	to := cli.Text()
	fmt.Print("From: ")
	cli.Scan()
	from := cli.Text()
	fmt.Print("Subject: ")
	cli.Scan()
	subject := cli.Text()
	content := WriteContent(r)

	msg := "Content-Type: text/plain; charset=UTF-8\r\n"
	msg += fmt.Sprintf(
		"To: %v\r\nFrom: %v\r\nSubject: %v\r\n\r\n%v",
		to, from, subject, content)

	m, err := mail.ReadMessage(strings.NewReader(msg))
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func WriteContent(r io.Reader) string {
	cli := bufio.NewScanner(r)
	fmt.Print("Content: (Enter SEND to finish adding content and send the email.\n")

	var content string
	for {
		cli.Scan()
		line := cli.Text()
		if line == "SEND" {
			break
		} else {
			content += line + "\n"
		}
	}
	return content
}

func Save(file string, m string) error {
	b := bytes.NewBufferString(m).Bytes()
	ioutil.WriteFile(file, b, 0600)
	return nil
}
