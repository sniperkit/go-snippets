package chillingeffects

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Notice struct {
	ID             int       `bson:"_id" json:"id"`
	Type           string    `bson:"type,omitempty" json:"type"`
	Title          string    `bson:"title,omitempty" json:"title"`
	Body           string    `bson:"body,omitempty" json:"body"`
	DateSent       time.Time `bson:"date_send,omitempty" json:"date_send"`
	DateReceived   time.Time `bson:"date_received,omitempty" json:"date_received"`
	Topics         []string  `bson:"topics,omitempty" json:"topics"`
	Tags           []string  `bson:"tags,omitempty" json:"tags"`
	Jurisdiction   []string  `bson:"jurisdiction,omitempty" json:"jurisdiction"`
	ActionTaken    string    `bson:"action_taken,omitempty" json:"action_taken"`
	SenderName     string    `bson:"sender_name,omitempty" json:"sender_name"`
	RecipientName  string    `bson:"recipient_name,omitempty" json:"recipient_name"`
	PrincipalName  string    `bson:"principal_name,omitempty" json:"principal_name"`
	Works          []Work    `bson:"works,omitempty" json:"works`
	Marks          []Work    `bson:"marks,omitempty" json:"marks"`
	Language       string    `bson:"language,omitempty" json:"language"`
	LegalComplaint string    `bson:"legal_complaint,omitempty" json:"legal_complaint"`
}

type Work struct {
	CopyrightedURLs []URL  `bson:"copyrighted_urls,omitempty" json:"copyrighted_urls"`
	DefamatoryURLs  []URL  `bson:"defamatory_urls,omitempty" json:"defamatory_urls"`
	Description     string `bson:"description,omitempty" json:"description"`
	InfringingURLs  []URL  `bson:"infringing_urls,omitempty" json:"infringing_urls"`
}

type URL struct {
	URL string `bson:"url,omitempty" json:"url"`
}

func RequestNotice(id int) (*Notice, error) {
	url := fmt.Sprintf("https://chillingeffects.org/notices/%d.json", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("StatusCode: %s", resp.Status)
	}

	dec := json.NewDecoder(resp.Body)
	value := make(map[string]*Notice)
	err = dec.Decode(&value)
	if err != nil {
		return nil, err
	}

	if len(value) > 1 {
		panic(fmt.Sprintf("Returned more than one notice on a single request with id %d.", id))
	}

	var notice *Notice
	for k, v := range value {
		v.Type = k
		notice = v
	}
	return notice, nil
}
