package main

// An Order for one or more products by a user.
// swagger:model order
type Order struct {
	// ID of the order
	//
	// required: true
	ID int64 `json:"id"`

	// the id of the user who placed the order.
	//
	// required: true
	UserID int64 `json:"user_id"`

	// items for this order
	// mininum items: 1
	OrderItems []struct {

		// the id of the product to order
		//
		// required: true
		ProductID int64 `json:"product_id"`

		// the quantity of this product to order
		//
		// required: true
		// minimum: 1
		Quantity int32 `json:"qty"`
	} `json:"items"`
}
