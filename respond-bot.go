package main

import (
  "log"
  "strings"
)

type respondBot struct {
  processor *processor
  name string
  reply string
}

func newRespondBot(processor *processor, name, reply string) (*respondBot) {
    return &respondBot{processor, name, reply}
}

func (bot *respondBot) run() {
  bot.processor.Out <-newMessage("nick", bot.name)
	for msg := range bot.processor.In {
  	if msg.Type == "chat" {
      message := msg.Payload.(string)
      if strings.HasPrefix(message, bot.name) {
        continue
      }

      if strings.Contains(message, bot.name) {
        log.Printf("found %s in %s\n", bot.name, message)
        bot.processor.Out <-newMessage("chat", bot.reply)
        log.Printf("sent %s\n", bot.reply)
      }
    }
  }
}
