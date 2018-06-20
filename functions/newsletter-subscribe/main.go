package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"

	apex "github.com/apex/go-apex"
)

type subscriber struct {
	Email       string                 `json:"email_address"`
	Status      string                 `json:"status"`
	MergeFields map[string]interface{} `json:"merge_fields"`
}

func getURL() string {
	var url bytes.Buffer
	url.WriteString("https://us12.api.mailchimp.com/3.0/lists/")
	url.WriteString("lists")
	url.WriteString(os.Getenv("MAILCHIMP_LIST_ID"))
	url.WriteString("/members")
	return url.String()
}

func subscription(subs *subscriber) error {

	subscriber := &subscriber{
		Email:       subs.Email,
		Status:      "subscribed",
		MergeFields: subs.MergeFields,
	}

	json, err := json.Marshal(subscriber)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", getURL(), bytes.NewBuffer(json))
	req.SetBasicAuth("anystring", os.Getenv("MAILCHIMP_API_KEY"))
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		var s subscriber
		if err := json.Unmarshal(event, &s); err != nil {
			return nil, err
		}
		if err := subscription(&s); err != nil {
			return nil, err
		}
		return map[string]interface{}{"status": "success"}, nil
	})
}
