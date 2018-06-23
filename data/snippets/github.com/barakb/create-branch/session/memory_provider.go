package session

import (
	"time"
)

func init() {
	Register("memory", memoryProvider{make(map[string]memorySession)})
}

type memoryProvider struct {
	sessions map[string]memorySession
}

func (p memoryProvider) SessionInit(sid string) (Session, error) {
	return p.getSessionCreated(sid)
}

func (p memoryProvider) SessionRead(sid string) (Session, error) {
	return p.getSessionCreated(sid)
}

func (p memoryProvider) getSessionCreated(sid string) (Session, error) {
	res, ok := p.sessions[sid]
	if !ok {
		time := time.Now()
		res = memorySession{m: make(map[interface{}]interface{}), sid: sid, lmt: &time}
		p.sessions[sid] = res
	}
	res.touch()
	return res, nil
}

func (p memoryProvider) SessionDestroy(sid string) error {
	delete(p.sessions, sid)
	return nil
}

func (p memoryProvider) SessionGC(maxLifeTime int64) {
	for key, session := range p.sessions {
		if session.lmt.Add(30 * time.Minute).Before(time.Now()) {
			delete(p.sessions, key)
		}
	}
}

type memorySession struct {
	m   map[interface{}]interface{}
	sid string
	lmt *time.Time
}

func (s memorySession) Set(key, value interface{}) error {
	s.touch()
	s.m[key] = value
	return nil
}

func (s memorySession) Get(key interface{}) interface{} {
	s.touch()
	return s.m[key]
}

func (s memorySession) Delete(key interface{}) error {
	s.touch()
	delete(s.m, key)
	return nil

}

func (s memorySession) SessionID() string {
	return s.sid
}

func (s memorySession) touch() {
	time := time.Now()
	s.lmt = &time
}
