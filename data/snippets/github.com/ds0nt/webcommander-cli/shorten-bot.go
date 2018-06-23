package main

import (
  "log"
	"fmt"
	"golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "google.golang.org/api/urlshortener/v1"
	"github.com/mvdan/xurls"
  "strings"
)

type shortenBot struct {
  processor *processor
  name string
  urlService *urlshortener.Service
}

func newShortenBot(processor *processor, name string) (*shortenBot, error) {
    bot := shortenBot{ processor, name, nil }
  	// Use oauth2.NoContext if there isn't a good context to pass in.
    client, err := google.DefaultClient(oauth2.NoContext, urlshortener.UrlshortenerScope)
    if err != nil {
    	return nil, err
    }

  	bot.urlService, err = urlshortener.New(client)
    if err != nil {
    	return nil, err
    }
  return &bot, nil
}

func (bot *shortenBot) run() {
  bot.processor.Out <-newMessage("nick", bot.name)
	for msg := range bot.processor.In {
  	if msg.Type == "chat" {
      message := msg.Payload.(string)
      if strings.HasPrefix(message, bot.name) {
        continue
      }
      urls := xurls.Relaxed.FindAllString(message, -1)

      for _, v := range urls {
        bot.processor.Out <-newMessage("chat", fmt.Sprintf("Shortening URL: %s", v))
        go func(u string) {
          url, err := bot.urlService.Url.Insert(&urlshortener.Url{LongUrl: u}).Do()
          if err != nil {
            log.Printf("Translation List Error: %v", err)
            bot.processor.Out <-newMessage("chat", fmt.Sprintf("Sorry, on cooldown.", u))
          } else {
            log.Printf("Short url: %v", url)
            bot.processor.Out <-newMessage("chat", fmt.Sprintf("Short URL %s", url.Id))
          }
        }(v)
      }
    }
  }
}
