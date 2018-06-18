// go-swagger examples.
//
// The purpose of this application is to provide some
// use cases describing how to generate docs for your API
//
//     Schemes: http, https
//     Host: localhost
//     BasePath: /
//     Version: 0.0.1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import "net/http"

// An OrderBodyParams model.
//
// This is used for operations that want an Order as body of the request
// swagger:parameters createOrder
type OrderBodyParams struct {
	// The order to submit.
	//
	// in: body
	// required: true
	Order *Order `json:"order"`
}

// An OrderQueryParams model.
//
// swagger:parameters listOrder
type OrderQueryParams struct {
	// List order limitations.
	//
	// required: true
	// minimum: 1
	// maximum: 12
	Limit int32 `json:"limit"`
}

// An OrderResponse response model
//
// This is used for returning a response with a single order as body
//
// swagger:response orderResponse
type OrderResponse struct {
	//Get number of rate limted requests remaining
	//
	RateLimitRemaining string `json:"Rate-Limit-Remaining"`
	//
	// in: body
	Payload *Order `json:"order"`
}

// A UserNameParams parameter model.
//
// This is used for operations that want the ID of a user in the path
// swagger:parameters getUserByNickname
type UserNameParams struct {
	// The nickname of user
	//
	// in: path
	// required: true
	Nickname string `json:"nickname"`
}

// A ValidationError is a swagger response to represent error
//
// swagger:response validationError
type ValidationError struct {
	// in: body
	Body struct {
		Code    int32  `json:"code"`
		Message string `json:"message"`
		Field   string `json:"field"`
	} `json:"body"`
}

// CreateOrder swagger:route POST /orders orders createOrder
//
// Handler to create an order.
//
// Responses:
//        200: orderResponse
//        422: validationError
func CreateOrder(w http.ResponseWriter, req *http.Request) {
	// your code here
}

// GetOrders swagger:route GET /orders orders listOrder
//
// Handler returning list of orders.
//
// List of orders made by current user
//
// Responses:
// 		    200: []order
func ListOrder(w http.ResponseWriter, r *http.Request) {
	// your code here
}

// Me swagger:route GET /users{nickname} users getUserByNickname
//
// Handler returning information about user.
//
// Information about user
//
// Responses:
// 		    200: user
func GetUserByNickname(w http.ResponseWriter, r *http.Request) {
	// your code here
}

func main() {
}
