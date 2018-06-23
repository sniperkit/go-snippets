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
	"net/http"
	"strconv"
	"time"

	"github.com/b3log/wide/util"
	"github.com/gorilla/websocket"
)

const (
	severityError = "ERROR" // notification.severity: ERROR
	severityWarn  = "WARN"  // notification.severity: WARN
	severityInfo  = "INFO"  // notification.severity: INFO

	typeServer = "server"
)

// notification represents a notification.
type notification struct {
	event    *event
	Type     string `json:"type"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

// event2Notification processes user event by converting the specified event to a notification, and then push it to front
// browser with notification channel.
func event2Notification(e *event) {
	if nil == notificationWS[e.Sid] {
		return
	}

	wsChannel := notificationWS[e.Sid]
	if nil == wsChannel {
		return
	}

	//httpSession, _ := HTTPSession.Get(wsChannel.Request, "coditor-session")
	// username := httpSession.Values["username"].(string)
	locale := "en_US" // TODO: user locale

	var noti *notification

	switch e.Code {
	case EvtCodeServerInternalError:
		noti = &notification{event: e, Type: typeServer, Severity: severityError,
			Message: getMsg(locale, "notification_"+strconv.Itoa(e.Code)).(string) + " [" + e.Data.(string) + "]"}
	default:
		logger.Warnf("Can't handle event[code=%d]", e.Code)

		return
	}

	wsChannel.WriteJSON(noti)

	wsChannel.Refresh()
}

func notificationWSHandler(w http.ResponseWriter, r *http.Request) {
	sid := r.URL.Query()["sid"][0]

	cSession := coditorSessions.get(sid)

	if nil == cSession {
		return
	}

	conn, _ := websocket.Upgrade(w, r, nil, 1024, 1024)
	wsChan := util.WSChannel{Sid: sid, Conn: conn, Request: r, Time: time.Now()}

	ret := map[string]interface{}{"notification": "Notification initialized", "cmd": "init-notification"}
	err := wsChan.WriteJSON(&ret)
	if nil != err {
		return
	}

	notificationWS[sid] = &wsChan

	logger.Tracef("Open a new [Notification] with session [%s], %d", sid, len(notificationWS))

	// add user event handler
	cSession.EventQueue.addHandler(eventHandleFunc(event2Notification))

	input := map[string]interface{}{}

	for {
		if err := wsChan.ReadJSON(&input); err != nil {
			return
		}
	}
}
