package util

import (
	"fmt"
	"math/rand"
	"time"
)

// format keys for User
// example: roc => roc:usrid
func FormatUserKey(userID string) string {
	return fmt.Sprintf("%s:usrid", userID)
}

// format key to associate with a user's subscription list
// example roc => roc:sublist
func FormatSubListKey(userID string) string {
	return fmt.Sprintf("%s:sublist", userID)
}

// format key for a tribble post
// example roc make a post => roc:post_time_srvId (time and srvId in %x)
// srvId is a random number to break ties for post id, not perfect but will work with very high probability.
// If it turns out to be not unique, call this function again to generate a new one.
func FormatPostKey(userID string, postTime int64) string {
	return fmt.Sprintf("%s:post_%x_%x", userID, postTime,
		rand.New(rand.NewSource(time.Now().Unix())))
}

// format key to associate with a user's list for tribble keys
// example roc => roc:triblist
func FormatTribListKey(userID string) string {
	return fmt.Sprintf("%s:triblist", userID)
}
