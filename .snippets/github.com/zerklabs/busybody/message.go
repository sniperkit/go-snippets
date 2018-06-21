package busybody

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"

	"github.com/zerklabs/auburn/log"
	"github.com/zerklabs/busybody/protocol"
)

func (m *BusyMember) defaultMessage() *protocol.Message {
	ct := protocol.NoCompression

	if m.config.SnappyCompression {
		ct = protocol.SnappyCompression
	}

	if m.config.DeflateCompression {
		ct = protocol.DeflateCompression
	}

	return protocol.NewMessage(protocol.StandardMessage, ct, m.id)
}

func (m *BusyMember) hellomsg() *protocol.Message {
	ct := protocol.NoCompression

	if m.config.SnappyCompression {
		ct = protocol.SnappyCompression
	}

	if m.config.DeflateCompression {
		ct = protocol.DeflateCompression
	}

	return protocol.NewMessage(protocol.HelloMessage, ct, m.id)
}

func UnmarshalIntroduction(p *protocol.Message) (*Introduction, error) {
	var intro Introduction

	if p.MessageType() != protocol.HelloMessage {
		return nil, fmt.Errorf("not a hello message")
	}

	body, err := p.Body()
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(body)
	decoder := gob.NewDecoder(buffer)

	if err := decoder.Decode(&intro); err != nil {
		if err != io.EOF {
			return nil, fmt.Errorf("error gob decoding introduction: %v", err)
		} else {
			log.Warnf("error during gob decoding: %v", err)
		}
	}

	if intro.Id == "" {
		return nil, fmt.Errorf("invalid introduction message: id missing")
	}

	if intro.Uri == "" {
		return nil, fmt.Errorf("invalid introduction message: uri missing")
	}

	if intro.Key == "" {
		log.Warnf("shared key is empty")
	}

	return &intro, nil
}

func (m *BusyMember) hello() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	buffer := bytes.NewBuffer(nil)
	encoder := gob.NewEncoder(buffer)

	msg := m.hellomsg()

	if err := encoder.Encode(m.Introduction()); err != nil {
		return fmt.Errorf("error gob encoding message: %v", err)
	}

	if _, err := msg.Write(buffer.Bytes()); err != nil {
		return fmt.Errorf("error writing content to message: %v", err)
	}

	return m.send(msg)
}

// Send writes the given byte slice to the underlying protocol message
func (m *BusyMember) Send(content []byte) error {
	msg := m.defaultMessage()

	if _, err := msg.Write(content); err != nil {
		return fmt.Errorf("error writing content to message: %v", err)
	}

	return m.send(msg)
}

func (m *BusyMember) send(msg *protocol.Message) error {
	sendbuf := bytes.NewBuffer(nil)
	n, err := sendbuf.ReadFrom(msg)
	if err != nil {
		return err
	}

	if m.config.LogLevel >= log.DEBUG && msg.MessageType() == protocol.StandardMessage {
		log.Debugf("sending %d bytes", n)
	}

	if err := m.bussock.Send(sendbuf.Bytes()); err != nil {
		return fmt.Errorf("error sending message: %v", err)
	}

	return nil
}

func (m *BusyMember) share() error {
	peers := m.Members()

	for _, peer := range peers {
		if peer.Id != "" && peer.state == HealthyState {
			buffer := bytes.NewBuffer(nil)
			encoder := gob.NewEncoder(buffer)
			msg := m.hellomsg()

			if err := encoder.Encode(peer); err != nil {
				return fmt.Errorf("error gob encoding message: %v", err)
			}

			msg.Write(buffer.Bytes())

			if err := m.send(msg); err != nil {
				return err
			}
		}
	}

	return nil
}
