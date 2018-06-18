package slackwebhook

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/haisum/recaptcha"
	"github.com/spf13/viper"
)

type webhook struct {
	webhook  string
	server   string
	channels map[string]bool
	origins  map[string]bool
	re       *recaptcha.R
}

func NewHandler(server, configDir string) func(w http.ResponseWriter, r *http.Request) {
	v := viper.New()

	v.AddConfigPath(configDir)
	if err := v.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	wh := &webhook{
		v.GetString("webhook"),
		server,
		make(map[string]bool),
		make(map[string]bool),
		nil,
	}

	if s := v.GetString("recaptchasecret"); s != "" {
		wh.re = &recaptcha.R{
			Secret: s,
		}
	}

	for _, o := range v.GetStringSlice("origins") {
		wh.origins[o] = true
	}
	for _, c := range v.GetStringSlice("channels") {
		wh.channels[c] = true
	}

	if v.GetString("webhook") == "" || len(wh.origins) == 0 || len(wh.channels) == 0 {
		log.Fatal("Invalid configuration")
	}

	return wh.handlePayload
}

func (wh *webhook) Send(payload slack.Payload) error {
	if payload.Channel == "" {
		return fmt.Errorf("Channel cannot be empty")
	}

	if !wh.channels[payload.Channel] {
		return fmt.Errorf("Invalid channel: %s", payload.Channel)
	}

	err := slack.Send(wh.webhook, "", payload)
	if len(err) > 0 {
		return err[0]
	}
	return nil
}

// test:
// curl -d "channel=#test&text=hello+world&author=yolo" -X POST http://localhost:9090
func (wh *webhook) handlePayload(w http.ResponseWriter, r *http.Request) {
	var err error
	var site, redirect string
	var skipThankyou, wantJson bool

	defer func() {
		if wantJson {
			if err != nil {
				http.Error(w, `{"ok":false,"msg":"`+err.Error()+`"}`, http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, `{"ok":true,"msg":"Payload sent"}`)
		} else {
			if err != nil {
				w.Write([]byte(templateFormError(err.Error(), site)))
				return
			}

			if redirect == "" {
				redirect = site
			}
			if skipThankyou {
				http.Redirect(w, r, redirect, http.StatusSeeOther)
				return
			}
			w.Write([]byte(templateFormSubmitted(redirect)))
		}
	}()

	// cors
	origin := strings.SplitAfter(r.Header.Get("Origin"), "://")
	referer := strings.SplitAfter(r.Header.Get("Referer"), "://")

	if len(origin) < 2 {
		if len(referer) > 1 && referer[1] != wh.server {
			err = errors.New("Bad referer")
			return
		}
		err = errors.New("Invalid origin")
		return
	} else if !wh.origins[origin[1]] {
		err = errors.New("Bad origin")
		return
	}

	// prepare slack payload
	attachment := &slack.Attachment{}
	payload := slack.Payload{
		Username:  "WebFormBot",
		IconEmoji: ":robot_face:",
	}

	r.ParseForm()

	// honeypot
	if strings.Join(r.Form["_gotcha"], "") != "" {
		return
	}

	site = strings.Join(r.Form["_site"], "")

	wantJson = strings.Join(r.Form["_json"], "") == "true"
	if err = wh.checkParams(r); err != nil {
		return
	}

	for k, v := range r.Form {
		switch strings.ToLower(k) {
		// app specific (underscore prefix)
		case "_site", "g-recaptcha-response", "_gotcha":
			continue
		case "_redirect":
			redirect = strings.Join(v, "")
		case "_skipthankyou":
			if strings.Join(v, "") == "true" {
				skipThankyou = true
			}

		// slack message
		case "text":
			payload.Text = strings.Join(v, "")
		case "channel":
			payload.Channel = strings.Join(v, "")

		// slack attachment
		case "pretext":
			pretext := strings.Join(v, "")
			attachment.PreText = &pretext
		case "attachmentcolor":
			color := strings.Join(v, "")
			attachment.Color = &color

		// slack fields
		default:
			attachment = attachment.AddField(slack.Field{Title: k, Value: strings.Join(v, "")})
		}
	}

	payload.Attachments = []slack.Attachment{*attachment}
	err = wh.Send(payload)
	if err != nil && strings.Contains(err.Error(), "no such host") {
		err = errors.New("Error connecting to slack")
	}
}

func (wh *webhook) checkParams(r *http.Request) error {
	// check `_site`
	if strings.Join(r.Form["_site"], "") == "" {
		return errors.New("`_site` parameter is required")
	}

	// verify google recaptcha
	if wh.re != nil && !wh.re.Verify(*r) {
		return errors.New("Invalid captcha")
	}

	return nil
}
