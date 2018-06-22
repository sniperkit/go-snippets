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
	"path/filepath"
	"time"

	"github.com/b3log/wide/util"
	"github.com/gorilla/websocket"
)

func editorWSHandler(w http.ResponseWriter, r *http.Request) {
	httpSession, _ := httpSessionStore.Get(r, "coditor-session")
	userSession := httpSession.Values[user_session]

	if nil == userSession {
		http.Error(w, "Forbidden", http.StatusForbidden)

		return
	}

	sid := r.URL.Query()["sid"][0]

	cSession := coditorSessions.get(sid)
	if nil == cSession {
		return
	}

	conn, _ := websocket.Upgrade(w, r, nil, 1024, 1024)
	wsChan := util.WSChannel{Sid: sid, Conn: conn, Request: r, Time: time.Now()}

	ret := map[string]interface{}{"editor": "Editor initialized", "cmd": "init-editor"}
	err := wsChan.WriteJSON(&ret)
	if nil != err {
		return
	}

	editorWS[sid] = &wsChan

	logger.Tracef("Open a new [Editor] with session [%s], %d", sid, len(editorWS))

	input := map[string]interface{}{}

	userName := coditorSessions.get(sid).Username

	for {
		if err := wsChan.ReadJSON(&input); err != nil {
			return
		}

		//logger.Trace(input)

		docName := input["docName"].(string)
		docName = filepath.Clean(docName)

		doc := documentHolder.getDoc(docName)
		err := doc.setContent(input["content"].(string), userName)
		if err != nil {
			// TODO maybe should send error here
			logger.Errorf("set document error %v", err)
			continue
		}

		for _, cursor := range doc.cursors {
			if cursor.Sid == sid { // skip the current session itself
				continue
			}

			content := input["content"].(string)

			ret = map[string]interface{}{"content": content, "cmd": "changes",
				"docName": docName, "user": input["user"], "cursor": input["cursor"], "color": input["color"]}

			if err := editorWS[cursor.Sid].WriteJSON(&ret); err != nil {
				logger.Error("[Editor Channel] ERROR: " + err.Error())

				return
			}
		}

		wsChan.Time = time.Now()
	}
}
