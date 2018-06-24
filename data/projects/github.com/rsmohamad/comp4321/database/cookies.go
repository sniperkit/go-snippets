package database

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/rsmohamad/comp4321/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	UserHistory = iota
)

type CookieDb struct {
	db *bolt.DB
}

var cookieDb *CookieDb

func GetCookieInstance() *CookieDb {
	if cookieDb == nil {
		cookieDb, _ = loadCookieDb("user.db")
	}
	return cookieDb
}

func loadCookieDb(filename string) (*CookieDb, error) {
	var cookieDb CookieDb
	var err error
	cookieDb.db, err = bolt.Open(filename, 0666, nil)
	if err != nil {
		return nil, err
	}

	cookieDb.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists(intToByte(UserHistory))
		return nil
	})

	return &cookieDb, nil
}

func (c *CookieDb) containsUserId(userId uint64) bool {
	var data []byte
	c.db.View(func(tx *bolt.Tx) error {
		users := tx.Bucket(intToByte(UserHistory))
		data = users.Get(uint64ToByte(userId))
		return nil
	})
	return data != nil
}

func (c *CookieDb) GetCookieId(r *http.Request) uint64 {
	var userId uint64 = 0
	if len(r.Cookies()) == 0 {
		userId = c.generateNewId()
	} else {
		for _, c := range r.Cookies() {
			if c.Name == "GoSearchID" {
				userId, _ = strconv.ParseUint(c.Value, 10, 64)
				break
			}
		}
		if !c.containsUserId(userId) {
			userId = c.generateNewId()
		}
	}
	return userId
}

func (c *CookieDb) generateNewId() uint64 {
	var uniqueId uint64
	c.db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket(intToByte(UserHistory))
		uniqueId, _ = users.NextSequence()
		return nil
	})
	c.ClearSearchHistory(uniqueId)
	return uniqueId
}

func (c *CookieDb) SetCookieResponse(userId uint64, w http.ResponseWriter) {
	cookieVal := fmt.Sprint(userId)
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "GoSearchID", Value: cookieVal, Expires: expiration}
	http.SetCookie(w, &cookie)
}

func (c *CookieDb) AddQuery(userId uint64, query string) {
	if !c.containsUserId(userId) {
		log.Println("user ID not found", userId)
		return
	}

	q := models.NewSearchHistory(query)
	c.db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket(intToByte(UserHistory))
		history := byteToHistory(users.Get(uint64ToByte(userId)))
		history = append([]models.SearchHistory{q}, history...)
		users.Put(uint64ToByte(userId), historyToByte(history))
		return nil
	})
}

func (c *CookieDb) GetSearchHistory(userId uint64) []models.SearchHistory {
	rv := make([]models.SearchHistory, 0)
	c.db.View(func(tx *bolt.Tx) error {
		users := tx.Bucket(intToByte(UserHistory))
		rv = byteToHistory(users.Get(uint64ToByte(userId)))
		return nil
	})
	return rv
}

func (c *CookieDb) ClearSearchHistory(userId uint64) {
	c.db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket(intToByte(UserHistory))
		interfaceArr := make([]models.SearchHistory, 0)
		users.Put(uint64ToByte(userId), historyToByte(interfaceArr))
		return nil
	})
}

func (c *CookieDb) Close() {
	c.db.Close()
}
