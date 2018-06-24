// transactional email handled by postmark
package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"github.com/datatogether/task_mgmt/tasks"
// 	"io"
// 	"io/ioutil"
// 	"net/http"
// 	"strings"
// )

// // SendTaskRequestEmail sends an email to cfg.EmailNotificationRecipients
// // with details for a newly requested task
// func SendTaskRequestEmail(t *tasks.Task) error {
// 	if len(cfg.EmailNotificationRecipients) == 0 {
// 		return fmt.Errorf("no recipients are set to send email to")
// 	}

// 	body := fmt.Sprintf(`{
//     "From" : "brendan@qri.io",
//     "To" : "%s",
//     "Tag" : "feedback",
//     "Subject" : "Injest Request: %s",
//     "TextBody" : "requested: %s\nsource url: %s\n"
//   }`,
// 		strings.Join(cfg.EmailNotificationRecipients, ","),
// 		t.Title,
// 		t.Request,
// 		t.SourceUrl,
// 	)

// 	return sendEmail(strings.NewReader(body))
// }

// // SendTaskRequestEmail sends an email to cfg.EmailNotificationRecipients
// // notifying them of a cancelled request
// func SendTaskCancelEmail(t *tasks.Task) error {
// 	if len(cfg.EmailNotificationRecipients) == 0 {
// 		return fmt.Errorf("no recipients are set to send email to")
// 	}

// 	body := fmt.Sprintf(`{
//     "From" : "brendan@qri.io",
//     "To" : "%s",
//     "Tag" : "feedback",
//     "Subject" : "Request Cancelled: %s",
//     "TextBody" : "requested: %s\nsource url: %s\ncancelled: %s"
//   }`,
// 		strings.Join(cfg.EmailNotificationRecipients, ","),
// 		t.Title,
// 		t.Request,
// 		t.SourceUrl,
// 		t.Fail,
// 	)

// 	return sendEmail(strings.NewReader(body))
// }

// // send an email using postmark transactional email service
// // postmarkapp.com
// func sendEmail(jsonBody io.Reader) error {
// 	if cfg.PostmarkKey == "" {
// 		return fmt.Errorf("missing postmark key for sending email")
// 	}

// 	url := "https://api.postmarkapp.com/email/"

// 	req, err := http.NewRequest("POST", url, jsonBody)
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Add("X-Postmark-Server-Token", cfg.PostmarkKey)
// 	req.Header.Add("Accept", "application/json")
// 	req.Header.Add("Content-Type", "application/json")
// 	req.Body = ioutil.NopCloser(jsonBody)

// 	res, err := http.DefaultClient.Do(req)
// 	// if the server responds with an error, process & log out
// 	if res.StatusCode == 422 {
// 		responseBody := map[string]interface{}{}
// 		json.NewDecoder(res.Body).Decode(&responseBody)
// 		log.Info(responseBody)
// 	}

// 	return err
// }
