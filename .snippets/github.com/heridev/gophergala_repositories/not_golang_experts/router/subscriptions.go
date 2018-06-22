package router

import (
	"encoding/json"
	"github.com/gophergala/not_golang_experts/model"
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

type SubscriptionParams struct {
	Url string `json:"url"`
}

func SubscriptionsIndex(res http.ResponseWriter, req *http.Request) {
	token := getToken(req)

	model.GetSubscriptionsForUser(token, func(subscriptions []model.Subscription) {
		respondWith(subscriptions, 200, res)
	}, func(message string) {
		respondWith(map[string]string{"error": message}, 401, res)
	})
}

func SubscriptionsCreate(res http.ResponseWriter, req *http.Request) {
	token := getToken(req)
	url, err := parseSubscriptionsRequest(req.Body)
	PanicIf(err, res)

	model.SubscribeUser(url, token, func(subscription model.Subscription) {
		respondWith(map[string]interface{}{"Subscription": subscription}, 201, res)
	}, func(message string) {
		respondWith(map[string]interface{}{"error": message}, 401, res)
	})
}

func SubscriptionsDestroy(res http.ResponseWriter, req *http.Request) {
	token := getToken(req)
	vars := mux.Vars(req)

	model.UnsubscribeUser(vars["id"], token, func(message string) {
		respondWith(map[string]interface{}{"message": message}, 200, res)
	}, func(message string) {
		respondWith(map[string]interface{}{"error": message}, 401, res)
	})
}

func parseSubscriptionsRequest(body io.Reader) (string, error) {
	subscription := SubscriptionParams{}

	decoder := json.NewDecoder(body)
	err := decoder.Decode(&subscription)

	return subscription.Url, err
}

func getToken(req *http.Request) string {
	params := req.URL.Query()
	if len(params["token"]) > 0 {
		return params["token"][0]
	} else {
		return ""
	}
}
