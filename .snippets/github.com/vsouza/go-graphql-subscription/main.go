package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/graphql-go/graphql"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

var commentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Comment",
		Fields: graphql.Fields{
			"commentTitle": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return randSeq(10), nil
				},
			},
			"commentDescription": &graphql.Field{
				Type: graphql.String,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return randSeq(10), nil
				},
			},
		},
	},
)

var subscriptionType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Subscription",
		Fields: graphql.Fields{
			"newComments": &graphql.Field{
				Type: commentType,
				Args: graphql.FieldConfigArgument{
					"postId": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return randSeq(10), nil
				},
			},
		},
	})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Subscription: subscriptionType,
		Query:        commentType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
	}
	return result
}

func subscriptionHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()["query"]
	if len(query) > 0 {
		result := executeQuery(query[0], schema)
		json.NewEncoder(w).Encode(result)
		return
	}
	http.Error(w, "query cannot be empty", http.StatusBadRequest)
}

func main() {
	http.HandleFunc("/subscription", subscriptionHandler)
	log.Println("Listening...")
	http.ListenAndServe(":8000", nil)

}
