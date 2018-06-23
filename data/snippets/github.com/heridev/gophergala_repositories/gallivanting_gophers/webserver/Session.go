package webserver

import "time"

// Session ...
type Session struct {
	Created   time.Time `bson:"created"`
	ID        string    `bson:"id"`
	IPAddress string    `bson:"ipaddress"`
}

// NewSession ...
func NewSession(id string, ip string) *Session {
	s := &Session{
		Created:   time.Now(),
		ID:        id,
		IPAddress: ip,
	}

	return s
}
