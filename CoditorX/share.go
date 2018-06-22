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
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/b3log/wide/util"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
	"html/template"
)

const (
	EDITABLE = 1 // editable
	VIEWABLE = 0 // viewable
)

type Share struct {
	Owner     string `json:"owner"`
	DocName   string `json:"docName"`
	ShareType int    `json:"shareType"` // 0 - view, 1 - edit
}

func shareHandler(w http.ResponseWriter, r *http.Request) {
	if "GET" == r.Method {
		model := map[string]interface{}{}
		toHtml(w, "share.html", model, conf.Locale)
		return
	} else if "POST" == r.Method {
		data := map[string]interface{}{"succ": true}
		defer util.RetJSON(w, r, data)

		httpSession, _ := httpSessionStore.Get(r, "coditor-session")
		userSession := httpSession.Values[user_session]
		if nil == userSession {
			data["succ"] = false
			data["msg"] = "permission denied"
			return
		}

		user := userSession.(*User)

		var args map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
			logger.Error(err)
			data["succ"] = false
			data["msg"] = "args decode error!"
			return
		}

		fileName := args["fileName"].(string)
		editorsStr := ""
		if args["editors"] != nil {
			editorsStr = args["editors"].(string)
		}
		viewersStr := ""
		if args["viewers"] != nil {
			viewersStr = args["viewers"].(string)
		}
		isPublic := 0
		if args["isPublic"] != nil {
			isPublic = int(args["isPublic"].(float64))
		}

		docOpen := true

		doc := documentHolder.getDoc(fileName)
		if doc == nil {
			docOpen = false
			metaData, err := newDocumentMetaData(fileName)
			if err != nil {
				data["succ"] = false
				data["msg"] = "open document error!"
				return
			}
			doc, err = newDocument(metaData, 10)
			if err != nil {
				data["succ"] = false
				data["msg"] = err.Error()
				return
			}
		}
		// get old editors and old viewers.To del the invalids.
		// check permission
		oldEditors, err := doc.getEditors(user.Username)
		if err != nil {
			data["succ"] = false
			data["msg"] = err.Error()
			return
		}
		oldViewers, err := doc.getViewers(user.Username)
		if err != nil {
			data["succ"] = false
			data["msg"] = err.Error()
			return
		}

		// only get the file name
		index := strings.LastIndex(fileName, string(os.PathSeparator))
		if index > -1 {
			fileName = fileName[index+1:]
		}

		tempEditors := strings.Split(editorsStr, ",")
		tempViewers := strings.Split(viewersStr, ",")

		// check user exist and save to share.json
		editors := []string{}
		viewers := []string{}
		for _, editor := range tempEditors {
			if editor == "" {
				continue
			}
			share := &Share{}
			share.Owner = user.Username
			share.DocName = fileName
			share.ShareType = EDITABLE
			editors, _ = checkAndSave(editor, editors, share)
		}
		for _, viewer := range tempViewers {
			if viewer == "" {
				continue
			}
			share := &Share{}
			share.Owner = user.Username
			share.DocName = fileName
			share.ShareType = VIEWABLE
			viewers, _ = checkAndSave(viewer, viewers, share)
		}

		delEditors := []string{}
		delViewers := []string{}
		for _, oldEditor := range oldEditors {
			delAble := true
			for _, tempEditor := range tempEditors {
				if oldEditor == tempEditor {
					delAble = false
					break
				}
			}
			if delAble {
				delEditors = append(delEditors, oldEditor)
			}
		}
		for _, oldViewer := range oldViewers {
			delAble := true
			for _, tempViewer := range tempViewers {
				if oldViewer == tempViewer {
					delAble = false
					break
				}
			}
			if delAble {
				delViewers = append(delViewers, oldViewer)
			}
		}
		share := &Share{}
		share.Owner = user.Username
		share.DocName = fileName
		share.ShareType = VIEWABLE
		logger.Debugf("%s cancle share file %s to %v", user.Username, fileName, delEditors)
		for _, delEditor := range delEditors {
			checkAndDel(delEditor, share)
		}
		logger.Debugf("%s cancle share file %s to %v", user.Username, fileName, delViewers)
		for _, delViewer := range delViewers {
			checkAndDel(delViewer, share)
		}

		err = doc.setIsPublic(isPublic, user.Username)
		if err != nil {
			data["succ"] = false
			data["msg"] = err.Error()
			return
		}
		logger.Debugf("%s share file %s to %v", user.Username, fileName, editors)
		err = doc.setEditors(editors, user.Username)
		if err != nil {
			data["succ"] = false
			data["msg"] = err.Error()
			return
		}
		logger.Debugf("%s share file %s to %v", user.Username, fileName, viewers)
		err = doc.setViewers(viewers, user.Username)
		if err != nil {
			data["succ"] = false
			data["msg"] = err.Error()
			return
		}

		if !docOpen {
			doc.close(1)
		}
	}
}

func checkAndSave(user string, users []string, share *Share) ([]string, error) {
	shareFilePath := filepath.Join(conf.Workspace, user, "share.json")
	data := func(shareFilePath string) (data []byte) {
		file, err := os.Open(shareFilePath)
		if err != nil {
			data = []byte{}
			logger.Errorf("share file error, %v", err)
			return data
		}
		data, err = ioutil.ReadAll(file)
		defer file.Close()
		if err != nil {
			data = []byte{}
			logger.Errorf("share file error, %v", err)
		}
		return data
	}(shareFilePath)
	shareList := []*Share{}
	if len(data) > 2 {
		// not empty!
		err := json.Unmarshal(data, &shareList)
		if err != nil {
			logger.Errorf("share file error, %v", err)
			return nil, err
		}
	}
	// check if this share is exist!
	index := -1
	for i, oShare := range shareList {
		if share.Owner == oShare.Owner && share.DocName == oShare.DocName {
			index = i
		}
	}
	if index == -1 {
		shareList = append(shareList, share)
	} else {
		shareList[index] = share
	}
	data, err := json.MarshalIndent(shareList, "", "    ")
	if err != nil {
		logger.Errorf("share file error, %v", err)
		return nil, err
	}
	err = ioutil.WriteFile(shareFilePath, data, 0644)
	if err != nil {
		logger.Errorf("share file error, %v", err)
		return nil, err
	}
	users = append(users, user)
	return users, nil
}

func checkAndDel(user string, share *Share) error {
	shareFilePath := filepath.Join(conf.Workspace, user, "share.json")
	data := func(shareFilePath string) (data []byte) {
		file, err := os.Open(shareFilePath)
		if err != nil {
			data = []byte{}
			logger.Errorf("share file error, %v", err)
			return data
		}
		data, err = ioutil.ReadAll(file)
		defer file.Close()
		if err != nil {
			data = []byte{}
			logger.Errorf("share file error, %v", err)
		}
		return data
	}(shareFilePath)
	shareList := []*Share{}
	if len(data) < 2 {
		// empty!
		return nil
	} else {
		// not empty!
		err := json.Unmarshal(data, &shareList)
		if err != nil {
			logger.Errorf("share file error, %v", err)
			return err
		}
	}
	// check if this share is exist!
	index := -1
	for i, oShare := range shareList {
		if share.Owner == oShare.Owner && share.DocName == oShare.DocName {
			index = i
		}
	}
	if index == -1 {
		return nil
	} else {
		shareList = append(shareList[:index], shareList[index+1:]...)
	}
	data, err := json.Marshal(shareList)
	if err != nil {
		logger.Errorf("share file error, %v", err)
		return err
	}
	err = ioutil.WriteFile(shareFilePath, data, 0644)
	if err != nil {
		logger.Errorf("share file error, %v", err)
	}
	return err
}

func checkAndUpdate(user string, oldShare, newShare *Share) error {
	shareFilePath := filepath.Join(conf.Workspace, user, "share.json")
	data := func(shareFilePath string) (data []byte) {
		file, err := os.Open(shareFilePath)
		if err != nil {
			data = []byte{}
			logger.Errorf("share file error, %v", err)
			return data
		}
		data, err = ioutil.ReadAll(file)
		defer file.Close()
		if err != nil {
			data = []byte{}
			logger.Errorf("share file error, %v", err)
		}
		return data
	}(shareFilePath)
	shareList := []*Share{}
	if len(data) < 2 {
		// empty!
		return nil
	} else {
		// not empty!
		err := json.Unmarshal(data, &shareList)
		if err != nil {
			logger.Errorf("share file error, %v", err)
			return err
		}
	}
	// check if this share is exist!
	index := -1
	for i, oShare := range shareList {
		if oldShare.Owner == oShare.Owner && oldShare.DocName == oShare.DocName {
			index = i
		}
	}
	if index == -1 {
		return nil
	} else {
		newShare.ShareType = oldShare.ShareType
		shareList[index] = newShare
	}
	data, err := json.Marshal(shareList)
	if err != nil {
		logger.Errorf("share file error, %v", err)
		return err
	}
	err = ioutil.WriteFile(shareFilePath, data, 0644)
	if err != nil {
		logger.Errorf("share file error, %v", err)
	}
	return err
}

func shareListHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"succ": true}
	defer util.RetJSON(w, r, data)

	httpSession, _ := httpSessionStore.Get(r, "coditor-session")
	userSession := httpSession.Values[user_session]
	if nil == userSession {
		data["succ"] = false
		data["msg"] = "permission denied"
		return
	}

	user := userSession.(*User)
	shareList, err := getOrInitShareFiles(user)
	if err != nil {
		data["succ"] = false
		data["msg"] = err.Error()
		return
	}

	data["shares"] = shareList
}

func getShareInfoHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"succ": true}
	defer util.RetJSON(w, r, data)

	httpSession, _ := httpSessionStore.Get(r, "coditor-session")
	userSession := httpSession.Values[user_session]
	if nil == userSession {
		data["succ"] = false
		data["msg"] = "permission denied"
		return
	}

	var args map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&args); err != nil {
		logger.Error(err)
		data["succ"] = false
		data["msg"] = "args decode error!"
		return
	}

	docName := args["docName"]
	if docName == nil || len(docName.(string)) == 0 {
		data["succ"] = false
		data["msg"] = "docName can not be null!"
		return
	}
	filePath := docName.(string)

	// check file first
	metaDataFileName := filePath + ".json.coditor"
	absMetaDataFileName := filepath.Clean(metaDataFileName)
	file, err := os.Open(absMetaDataFileName)
	if err != nil {
		logger.Error(err)
		data["succ"] = false
		data["msg"] = err.Error()
		return
	}
	file.Close()

	dmd, err := newDocumentMetaData(filePath)
	if err != nil {
		logger.Error(err)
		data["succ"] = false
		data["msg"] = err.Error()
		return
	}

	data["shareInfo"] = dmd
}

func publicViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	owner := vars["owner"]
	fileName := vars["fileName"]

	model := make(map[string]interface{}, 0)
	model["owner"] = owner
	model["fileName"] = fileName

	docName := filepath.Join(conf.Workspace, owner, "workspace", fileName)
	metaDataFileName := docName + ".json.coditor"
	absMetaDataFileName := filepath.Clean(metaDataFileName)
	file, err := os.Open(absMetaDataFileName)
	if err != nil {
		logger.Error(err)
		toHtml(w, "public_view_error.html", model, "")
		return
	}
	file.Close()

	metaData, err := newDocumentMetaData(docName)
	if err != nil {
		logger.Error(err)
		toHtml(w, "public_view_error.html", model, "")
		return
	}

	if metaData.IsPublic != 1 {
		logger.Error(err)
		toHtml(w, "public_view_error.html", model, "")
		return
	} else {
		contentBytes, err := ioutil.ReadFile(docName)
		if err != nil {
			logger.Error(err)
			toHtml(w, "public_view_error.html", model, "")
			return
		}
		contentBytes = blackfriday.MarkdownCommon(contentBytes)
		content := string(contentBytes)
		model["content"] = template.HTML(content)
		toHtml(w, "public_view.html", model, "")
	}
}

func getOrInitShareFiles(u *User) ([]*Share, error) {
	shareFilePath := u.getFileBasePath("share.json")
	file, err := os.Open(shareFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// new and init share.json
			shareList := []*Share{}
			data, err := json.Marshal(shareList)
			if err != nil {
				return nil, err
			}
			err = ioutil.WriteFile(shareFilePath, data, 0644)
			if err != nil {
				return nil, err
			}
			file, err = os.Open(shareFilePath)
			if err != nil {
				return nil, err
			}
		}
	}
	data, err := ioutil.ReadAll(file)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	shareList := []*Share{}
	if len(data) > 2 {
		// empty!
		err = json.Unmarshal(data, &shareList)
		if err != nil {
			return nil, err
		}
	}
	return shareList, nil
}
