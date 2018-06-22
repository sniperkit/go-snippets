package main

import (
  "log"
  // "github.com/PuerkitoBio/goquery"
  "strings"
  "os/exec"
)

type insightBot struct {
  processor *processor
  name string
}

func newInsightBot(processor *processor, name string) *insightBot {
    return &insightBot{processor, name}
}

func (bot *insightBot) run() {
  bot.processor.Out <-newMessage("nick", bot.name)

	for msg := range bot.processor.In {
  	if msg.Type == "chat" {
      message := msg.Payload.(string)
      if strings.HasPrefix(message, bot.name) {
        continue
      }

      if !strings.Contains(message, bot.name) {
        continue
      }

      go func() {
        insight, err := exec.Command("./respond.js").Output()
        if err != nil {
          log.Println("Not very insightful:", err)
          return
        }
        log.Println("Received Insight:", string(insight))
        bot.processor.Out <-newMessage("chat", string(insight))
      }()

      // go func() {
      //   doc, err := goquery.NewDocument("https://datahexagon.com/fddata/Documents/insight.htm")
      //   if err != nil {
      //     log.Println("Not very insightful:", err)
      //     return
      //   }
      //   insight, err := doc.Find("body").Html()
      //   if err != nil {
      //     return
      //   }
      //
      // }()
    }
  }
}
