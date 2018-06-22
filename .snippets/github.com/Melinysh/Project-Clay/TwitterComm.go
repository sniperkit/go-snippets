package main

import (
	"encoding/json"
	"fmt"
	"github.com/kurrik/oauth1a"
	"github.com/kurrik/twittergo"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type twitterComm struct {
	client *twittergo.Client
}

func switchClient(twComm *twitterComm) {
	// Code removed as it contained my API keys
}

func NewTwitterAppClient() *twitterComm {
	config := &oauth1a.ClientConfig{
		ConsumerKey:    "my-api-key",
		ConsumerSecret: "my-api-secret",
	}
	user := oauth1a.NewAuthorizedConfig("oauth", "sucks")
	client := twittergo.NewClient(config, user)
	if err := client.FetchAppToken(); err != nil {
		if err2 := client.FetchAppToken(); err2 != nil {
			fmt.Println("Error occured creating the client.", err)
			panic("ERROR CREATING USER CLAY CLIENT")
		}
	}

	return &twitterComm{client}
}

func NewTwitterUserClient() *twitterComm {
	config := &oauth1a.ClientConfig{
		ConsumerKey:    "my-api-key",
		ConsumerSecret: "my-api-secret",
	}

	user := oauth1a.NewAuthorizedConfig("oauth", "sucks")
	client := twittergo.NewClient(config, user)
	if err := client.FetchAppToken(); err != nil {
		fmt.Println("Error occured creating the client.", err)
		panic("ERROR CREATING USER CLIENT")
	}
	return &twitterComm{client}
}

func (twComm *twitterComm) Tweets(userID int64 /* Max 200 */) []map[string]interface{} {
	url := fmt.Sprintf("https://api.twitter.com/1.1/statuses/user_timeline.json?user_id=%d&count=200", userID)
	data, err := twComm.responseData(url)
	tweets := []map[string]interface{}{}
	if err != nil {
		return tweets
	}
	if jsonErr := json.Unmarshal(data, &tweets); jsonErr != nil {
		fmt.Println("Error in unmarshalling response", jsonErr)
		fmt.Println("String is:", string(data))
	}
	return tweets
}

func (twComm *twitterComm) User(id int64) map[string]interface{} {
	url := fmt.Sprintf("https://api.twitter.com/1.1/users/show.json?user_id=%d", id)
	data, err := twComm.responseData(url)
	responseMap := map[string]interface{}{}

	if err != nil {
		fmt.Println("Error:", err, ". Returning empty user.")
	}

	if jsonErr := json.Unmarshal(data, &responseMap); jsonErr != nil {
		fmt.Println("Error in unmarshalling response data", jsonErr)
		fmt.Println("String is:", string(data))
		return twComm.User(id)
	}
	return responseMap
}

func (twComm *twitterComm) Users(ids []int64) []map[string]interface{} {
	if len(ids) == 0 {
		return []map[string]interface{}{}
	}
	url := "https://api.twitter.com/1.1/users/lookup.json?user_id=" // MAX 100
	for index, id := range ids {
		if index == 100 {
			panic("Too many ids passed.")
		}
		url = url + strconv.FormatInt(id, 10) + ","
	}
	url = url[:len(url)-1]
	data, err := twComm.responseData(url)
	users := []map[string]interface{}{}
	if err != nil {
		println("Error occured on fetching users", ids, err)
		return users
	}

	if jsonErr := json.Unmarshal(data, &users); jsonErr != nil {
		fmt.Println("Error unmarshalling response data", jsonErr, url)
		fmt.Println("String is:", string(data))
	}
	return users
}

func (twComm *twitterComm) Favourites(userID int64 /* Max 200 */) []map[string]interface{} {
	url := fmt.Sprintf("https://api.twitter.com/1.1/favorites/list.json?user_id=%d&count=200", userID)
	data, err := twComm.responseData(url)
	favourites := []map[string]interface{}{}
	if err != nil {
		return favourites
	}

	if jsonErr := json.Unmarshal(data, &favourites); jsonErr != nil {
		fmt.Println("Error in unmarshalling response data", jsonErr)
	}
	return favourites
}

func (twComm *twitterComm) FollowingIDS(userID int64) []interface{} {
	/* Fetches following ids and assigns them to user map*/
	url := fmt.Sprintf("https://api.twitter.com/1.1/friends/ids.json?user_id=%d&count=1000", userID)
	data, err := twComm.responseData(url)
	responseMap := map[string]interface{}{}
	if err != nil {
		return []interface{}{}
	}

	if jsonErr := json.Unmarshal(data, &responseMap); jsonErr != nil {
		fmt.Println("Error in unmarshalling response data", jsonErr)
		fmt.Println("String is:", string(data))
	}

	if ids, ok := responseMap["ids"]; ok {
		idsToReturn := ids.([]interface{})
		return idsToReturn
	}
	fmt.Println("No ids found in response map for following ids. Map:", responseMap)
	return []interface{}{}
}

func (twComm *twitterComm) FollowerIDS(userID int64) []interface{} {
	url := fmt.Sprintf("https://api.twitter.com/1.1/followers/ids.json?user_id=%d&count=1000", userID)
	data, err := twComm.responseData(url)

	responseMap := map[string]interface{}{}
	if err != nil {
		return []interface{}{}
	}

	if jsonErr := json.Unmarshal(data, &responseMap); jsonErr != nil {
		fmt.Println("Error in unmarshalling response data", jsonErr)
		fmt.Println("String is:", string(data))
	}

	if ids, ok := responseMap["ids"]; ok {
		idsToReturn := ids.([]interface{})
		return idsToReturn
	}
	fmt.Println("No ids found in response map for follower ids. Map:", responseMap)
	return []interface{}{}
}

func (twComm *twitterComm) responseData(url string) ([]byte, error) {
	fmt.Println("Getting response for", url)
	var req *http.Request
	var err error
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		fmt.Println("Error occured on creating request url", url, err)
		return nil, err
	}

	var resp *twittergo.APIResponse
	if resp, err = twComm.client.SendRequest(req); err != nil {
		fmt.Println("Error occured on sending request", url, err)
		return nil, err
	}
	defer resp.Body.Close()

	var body = []byte{}                                    //need to parse own response for JSON
	if body, err = ioutil.ReadAll(resp.Body); err != nil { //read all the data as []byte
		fmt.Println("Error in ioutil Reading of response", url, err)
		return nil, err
	}

	if strings.Contains(string(body), "Rate limit exceeded") {
		fmt.Println("---------------------------------------------")
		fmt.Println("URL: ", url)
		fmt.Printf("Rate limit:           %v\n", resp.RateLimit())
		fmt.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
		fmt.Printf("Rate limit reset:     %v\n---------------------------------------------\n", resp.RateLimitReset())
		twComm.handleRateLimit()
		return twComm.responseData(url)
	}

	return body, nil
}

func (twComm *twitterComm) handleRateLimit() {
	fmt.Println("Rate limited. Switching twitter clients after 30 seconds...")
	time.Sleep(time.Duration(30 * time.Second))
	switchClient(twComm)
}
