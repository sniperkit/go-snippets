package main

import "time"

// TimeStamp struct for the Model Interface
type TimeStamp struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateTime Method for setting/updating the time
// struct elements
func (t *TimeStamp) UpdateTime() {

	currentTime := time.Now().UTC()
	if !t.CreatedAt.IsZero() {
		t.UpdatedAt = currentTime
		return
	}
	t.CreatedAt = currentTime
	t.UpdatedAt = currentTime
	return
}
