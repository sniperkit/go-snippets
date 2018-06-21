package session

import (
	"context"
	"newcache"
	"result"
	"rut"
	"util"
)

type Session struct {
	Username   string
	Password   string
	ID         string
	RUT        *rut.RUT             //For run script
	Result     <-chan result.Result //For run script
	CaseResult map[string]chan *result.Result
	NewCache   *newcache.NewCache
	Cancel     context.CancelFunc
	Admin      bool
}

func New(name, pass string) *Session {
	id := util.GenerateSessionIDByUserNameAndPassword(name, pass)
	cache := newcache.New(id)
	cache.Restore()

	adminid := util.GenerateSessionIDByUserNameAndPassword("mra", "")

	return &Session{
		Username:   name,
		Password:   pass,
		ID:         id,
		NewCache:   cache,
		Result:     make(<-chan result.Result),
		CaseResult: make(map[string]chan *result.Result),
		Admin:      adminid == id,
	}
}

func GenerateSessionIDByUserNameAndPassword(name, pass string) string {
	return util.GenerateSessionIDByUserNameAndPassword(name, pass)
}

func (s *Session) Destroy() {
	s.NewCache.Save()
	s.NewCache = nil
}

func IsAdmin(id string) bool {
	return id == util.GenerateSessionIDByUserNameAndPassword("mra", "")
}
