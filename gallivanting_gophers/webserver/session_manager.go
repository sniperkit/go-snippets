package webserver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/rand"
	"time"

	"github.com/gophergala/gallivanting_gophers/data"
	"gopkg.in/mgo.v2/bson"
)

// SessionManager is setup to handle session management for users. I wish I could
// have used the system I've already implemented for work, but it isn't currently
// open source. This did not get developed yet but will be a stub/placeholder for
// the real system which will drop in.
type SessionManager struct {
	Sessions []*Session

	idChannel chan string
	database  *data.DB
	running   bool
}

// NewSessionManager creates a new session manager.
func NewSessionManager(db *data.DB) *SessionManager {
	sm := &SessionManager{
		Sessions:  make([]*Session, 0),
		idChannel: make(chan string, 100),
		database:  db,
	}

	return sm
}

// Start turns on the session manager threads (cleanup and id generation).
func (sm *SessionManager) Start() {
	sm.running = true
}

// GetID returns an ID from the buffered channel.
func (sm *SessionManager) GetID() string {
	return <-sm.idChannel
}

// Since IDs are verified against the database for duplicates, prepopulate
// sessionids in a buffered channel to handle spike loads.
func (sm *SessionManager) produceIDs() {
	for sm.running {
		sm.idChannel <- sm.createID()
	}
}

func (sm *SessionManager) createID() string {
	var unique = false
	var key string
	for !unique {
		buffer := new(bytes.Buffer)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		binary.Write(buffer, binary.LittleEndian, r.Int63())
		key = string(buffer.Bytes())

		count, _ := sm.database.Sessions.Find(bson.M{"key": key}).Count()
		if count == 0 {
			unique = true
		}
	}

	fmt.Printf("Creating sessionid: %s\n", key)

	return key
}
