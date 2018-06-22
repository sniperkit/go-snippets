package stackoverflow

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	//APIBase is the base URL for the StackOverflow API
	APIBase = "https://api.stackexchange.com/2.2"
)

type CommonResponse struct {
	ErrorID      int    `json:"error_id"`
	ErrorMessage string `json:"error_message"`
	ErrorName    string `json:"error_name"`

	QuotaMax       int `json:"quota_max"`
	QuotaRemaining int `json:"quota_remaining"`

	HasMore bool `json:"has_more"`
}

type Owner struct {
	Reputation   int    `json:"reputation"`
	UserID       int    `json:"user_id"`
	UserType     string `json:"user_type"`
	ProfileImage string `json:"profile_image"`
	DisplayName  string `json:"display_name"`
	Link         string `json:"link"`
}

type SearchResponse struct {
	CommonResponse
	Items []struct {
		Tags             []string `json:"tags"`
		Owner            `json:"owner"`
		IsAnswered       bool   `json:"is_answered"`
		ViewCount        int    `json:"view_count"`
		AnswerCount      int    `json:"answer_count"`
		AcceptedAnswerID int    `json:"accepted_answer_id"`
		Score            int    `json:"score"`
		LastActivityDate int    `json:"last_activity_date"`
		CreationDate     int    `json:"creation_date"`
		QuestionID       int    `json:"question_id"`
		Link             string `json:"link"`
		Title            string `json:"title"`
	} `json:"items"`
}

type AnswerResponse struct {
	CommonResponse
	Items []struct {
		Owner            `json:"owner"`
		IsAccepted       bool   `json:"is_accepted"`
		Score            int    `json:"score"`
		LastActivityDate int    `json:"last_activity_date"`
		CreationDate     int    `json:"creation_date"`
		AnswerID         int    `json:"answer_id"`
		QuestionID       int    `json:"question_id"`
		Body             string `json:"body"`
	} `json:"items"`
}

type Validator interface {
	IsValid() bool
	Error() error
}

func (res CommonResponse) IsValid() bool {
	if res.ErrorID > 0 || res.QuotaRemaining == 0 {
		return false
	}

	return true
}

func (res CommonResponse) Error() error {
	return fmt.Errorf("Error Id: %v, %v: %v", res.ErrorID, res.ErrorName, res.ErrorMessage)
}

func get(url string) (body []byte, err error) {
	res, err := http.Get(url)
	defer res.Body.Close()

	if err != nil {
		return
	}

	return ioutil.ReadAll(res.Body)
}

func load(url string, result Validator) (err error) {
	body, err := get(url)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)

	if !result.IsValid() {
		err = result.Error()
	}

	return
}

func Search(request *SearchRequest) (*SearchResponse, error) {
	url := makeSearchRequest(request)
	result := new(SearchResponse)
	err := load(url, result)

	return result, err
}

func GetAnswers(request *AnswerRequest) (*AnswerResponse, error) {
	url := makeAnswerRequest(request)
	result := new(AnswerResponse)
	err := load(url, result)

	return result, err
}
