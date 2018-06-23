package main

// User available to create an order.
// swagger:model user
type User struct {
	// ID of user
	//
	// required: true
	ID int64 `json:"id"`

	// nickname of user.
	//
	// required: true
	Nickname string `json:"nickname"`
}
