package common

import (
	"strings"
)

const SEPARATOR string = "0192837465"

type Message struct {
	From string
	To   string
	Body string
}

func (m *Message) FromBytes(b []byte) Message {
	raw := string(b)
	fields := strings.Split(raw, SEPARATOR)
	return Message{fields[0], fields[1], fields[2]}
}

func (m *Message) ToBytes() []byte {
	raw := []string{m.From, m.To, m.Body}
	return []byte(strings.Join(raw, SEPARATOR))
}
