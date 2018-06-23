// Copyright (c) 2015, b3log.org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/b3log/wide/util"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

const (
	sessionStateActive = iota
	sessionStateClosed // (not used so far)
)

var (
	// sessionWS holds all session channels. <sid, *util.WSChannel>
	sessionWS = map[string]*util.WSChannel{}

	// editorWS holds all editor channels. <sid, *util.WSChannel>
	editorWS = map[string]*util.WSChannel{}

	// notificationWS holds all notification channels. <sid, *util.WSChannel>
	notificationWS = map[string]*util.WSChannel{}
)

// HTTP session store.
var httpSessionStore = sessions.NewCookieStore([]byte("BEYOND"))

// CoditorSession represents a session associated with a browser tab.
type coditorSession struct {
	ID          string            // id
	Username    string            // username
	Color       string            // authorship color
	HTTPSession *sessions.Session // HTTP session related
	EventQueue  *userEventQueue   // event queue
	State       int               // state
	Created     time.Time         // create time
	Updated     time.Time         // the latest use time
}

// Type of Coditor sessions.
type cSessions []*coditorSession

// Coditor sessions.
var coditorSessions cSessions

// Exclusive lock.
var mutex sync.Mutex

func fixedTimeRelease() {
	go func() {
		for _ = range time.Tick(time.Hour) {
			hour, _ := time.ParseDuration("-30m")
			threshold := time.Now().Add(hour)

			for _, s := range coditorSessions {
				if s.Updated.Before(threshold) {
					logger.Debugf("Removes a invalid session [%s], user [%s]", s.ID, s.Username)

					coditorSessions.remove(s.ID)
				}
			}
		}
	}()
}

// Online user statistic report.
type userReport struct {
	username   string
	sessionCnt int
	updated    time.Time
}

// report returns a online user statistics in pretty format.
func (u *userReport) report() string {
	return "[" + u.username + "] has [" + strconv.Itoa(u.sessionCnt) + "] sessions, latest activity [" + u.updated.Format("2006-01-02 15:04:05") + "]"
}

func fixedTimeReport() {
	go func() {
		for _ = range time.Tick(10 * time.Minute) {
			users := userReports{}

			for _, s := range coditorSessions {
				if report, exists := contains(users, s.Username); exists {
					if s.Updated.After(report.updated) {
						report.updated = s.Updated
					}

					report.sessionCnt++
				} else {
					users = append(users, &userReport{username: s.Username, sessionCnt: 1, updated: s.Updated})
				}
			}

			var buf bytes.Buffer
			buf.WriteString("\n  [" + strconv.Itoa(len(users)) + "] users, [" +
				strconv.Itoa(len(coditorSessions)) + "] sessions currently\n")

			for _, t := range users {
				buf.WriteString("    " + t.report() + "\n")
			}

			logger.Info(buf.String())
		}
	}()
}

func contains(reports []*userReport, username string) (*userReport, bool) {
	for _, ur := range reports {
		if username == ur.username {
			return ur, true
		}
	}

	return nil, false
}

type userReports []*userReport

func coditorSessionWSHandler(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query()["sid"][0]
	cSession := coditorSessions.get(sid)
	if nil == cSession {
		httpSession, _ := httpSessionStore.Get(r, "coditor-session")

		if httpSession.IsNew {
			return
		}

		httpSession.Options.MaxAge = conf.HTTPSessionMaxAge
		httpSession.Save(r, w)

		cSession = coditorSessions.new(httpSession, sid)

		logger.Tracef("Created a Coditor session [%s] for websocket reconnecting, user [%s]", sid, cSession.Username)
	}

	conn, _ := websocket.Upgrade(w, r, nil, 1024, 1024)
	wsChan := util.WSChannel{Sid: sid, Conn: conn, Request: r, Time: time.Now()}

	ret := map[string]interface{}{"output": "Session initialized", "cmd": "init-session"}
	err := wsChan.WriteJSON(&ret)
	if nil != err {
		return
	}

	sessionWS[sid] = &wsChan

	logger.Tracef("Open a new [Session Channel] with session [%s], %d", sid, len(sessionWS))

	input := map[string]interface{}{}

	for {
		if err := wsChan.ReadJSON(&input); err != nil {
			logger.Tracef("[Session Channel] of session [%s] disconnected, releases all resources with it, user [%s]", sid, cSession.Username)

			coditorSessions.remove(sid)

			return
		}

		ret = map[string]interface{}{"output": "", "cmd": "session-output"}

		if err := wsChan.WriteJSON(&ret); err != nil {
			logger.Error("[Session Channel] ERROR: " + err.Error())

			return
		}

		wsChan.Time = time.Now()
	}
}

// Refresh refreshes the channel by updating its use time.
func (s *coditorSession) refresh() {
	s.Updated = time.Now()
}

// New creates a Coditor session.
func (sessions *cSessions) new(httpSession *sessions.Session, sid string) *coditorSession {
	mutex.Lock()
	defer mutex.Unlock()

	now := time.Now()

	// create user event queue
	userEventQueue := userEventQueues.new(sid)

	r := util.Rand.Int(0, 255)
	g := util.Rand.Int(0, 255)
	b := util.Rand.Int(0, 255)
	rgb := "rgb(" + strconv.Itoa(r) + "," + strconv.Itoa(g) + "," + strconv.Itoa(b) + ")"

	ret := &coditorSession{
		ID:          sid,
		Username:    httpSession.Values["username"].(string),
		Color:       rgb,
		HTTPSession: httpSession,
		EventQueue:  userEventQueue,
		State:       sessionStateActive,
		Created:     now,
		Updated:     now,
	}

	*sessions = append(*sessions, ret)

	return ret
}

// Get gets a Coditor session with the specified session id.
func (sessions *cSessions) get(sid string) *coditorSession {
	mutex.Lock()
	defer mutex.Unlock()

	for _, s := range *sessions {
		if s.ID == sid {
			return s
		}
	}

	return nil
}

func (sessions *cSessions) remove(sid string) {
	mutex.Lock()
	defer mutex.Unlock()

	for i, s := range *sessions {
		if s.ID == sid {
			// remove from session set
			*sessions = append((*sessions)[:i], (*sessions)[i+1:]...)

			// close user event queue
			userEventQueues.close(sid)

			// close websocket channels
			if ws, ok := notificationWS[sid]; ok {
				ws.Close()
				delete(notificationWS, sid)
			}

			if ws, ok := sessionWS[sid]; ok {
				ws.Close()
				delete(sessionWS, sid)
			}

			// release cursors in document
			for _, doc := range documentHolder.docs {
				doc.removeCursor(sid)
			}

			cnt := 0 // count Coditor sessions associated with HTTP session
			for _, ses := range *sessions {
				if ses.HTTPSession.Values["id"] == s.HTTPSession.Values["id"] {
					cnt++
				}
			}

			logger.Debugf("Removed a session [%s] of user [%s], it has [%d] sessions currently", sid, s.Username, cnt)

			return
		}
	}
}

// GetByUsername gets Coditor sessions.
func (sessions *cSessions) getByUsername(username string) []*coditorSession {
	mutex.Lock()
	defer mutex.Unlock()

	ret := []*coditorSession{}

	for _, s := range *sessions {
		if s.Username == username {
			ret = append(ret, s)
		}
	}

	return ret
}
