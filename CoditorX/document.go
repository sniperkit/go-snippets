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
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sergi/go-diff/diffmatchpatch"
)

var (
	documentHolder *DocumentHolder
)

// the holder to hold document.
type DocumentHolder struct {
	docs map[string]*Document
	lock sync.RWMutex
}

func newDocumentHolder() *DocumentHolder {
	dh := &DocumentHolder{}
	dh.docs = make(map[string]*Document, 0)
	return dh
}

func (dh *DocumentHolder) getDoc(docName string) *Document {
	dh.lock.RLock()
	defer func() {
		dh.lock.RUnlock()
	}()

	return dh.docs[docName]
}

func (dh *DocumentHolder) setDoc(docName string, doc *Document) error {
	dh.lock.Lock()
	defer func() {
		dh.lock.Unlock()
	}()

	oldDoc := dh.docs[docName]
	if oldDoc != nil {
		return errors.New("document is opened!")
	}
	dh.docs[docName] = doc
	return nil
}

func (dh *DocumentHolder) delDoc(docName string) {
	dh.lock.Lock()
	defer func() {
		dh.lock.Unlock()
	}()
	delete(dh.docs, docName)
}

// Document's MetaData
type DocumentMetaData struct {
	fileName string     // need not save to MetaData file
	Owner    string     `json:"owner"`
	Editors  []string   `json:"editors"`
	Viewers  []string   `json:"viewers"`
	IsPublic int        `json:"isPublic"`
	Version  DocVersion `json:"version"`
}

// newDocumentMetaData creates a new DocumentMetaData with the specified fileName full path.
func newDocumentMetaData(fileName string) (*DocumentMetaData, error) {
	metaDataFileName := fileName + ".json.coditor"
	absMetaDataFileName := filepath.Clean(metaDataFileName)
	file, err := openOrCreateFile(absMetaDataFileName)
	if err != nil {
		return nil, err
	}
	defer func() {
		file.Close()
	}()
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var dmd DocumentMetaData
	if len(fileData) < 10 {
		logger.Infof("the MetaData file is empty. FileName: %s", fileName)
		// is empty or error!
		dmd = DocumentMetaData{}
		dmd.Editors = []string{}
		dmd.Viewers = []string{}
		dmd.fileName = fileName
	} else {
		err := json.Unmarshal(fileData, &dmd)
		dmd.fileName = fileName
		return &dmd, err
	}
	return &dmd, nil
}

// Save MetaData to File.
func (dmd *DocumentMetaData) save() error {
	metaDataFileName := dmd.fileName + ".json.coditor"
	absFileName := filepath.Clean(metaDataFileName)
	metaDataBytes, err := json.MarshalIndent(dmd, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(absFileName, metaDataBytes, 0644)
	return err
}

func (dmd *DocumentMetaData) checkOwner(userName string) bool {
	if dmd.Owner == userName {
		return true
	}
	return false
}

func (dmd *DocumentMetaData) checkEditAble(userName string) bool {
	if dmd.Owner == userName {
		return true
	}
	for _, editor := range dmd.Editors {
		if editor == userName {
			return true
		}
	}
	return false
}

func (dmd *DocumentMetaData) checkViewAble(userName string) bool {
	if dmd.Owner == userName {
		return true
	}
	for _, viewer := range dmd.Viewers {
		if viewer == userName {
			return true
		}
	}
	for _, editor := range dmd.Editors {
		if editor == userName {
			return true
		}
	}
	return false
}

// The version of Document.
type DocVersion uint32

// Document.
type Document struct {
	content    string
	cursors    []*Cursor
	cursorLock sync.Mutex
	ringCache  *RingCache
	lock       sync.RWMutex
	metaData   *DocumentMetaData
	docFile    *os.File
	binLog     *BinLog

	flushTicker *time.Ticker
}

func newDocument(metaData *DocumentMetaData, flushTime int64) (*Document, error) {
	doc := &Document{}
	// TODO should check metaData first!
	doc.metaData = metaData
	absDocFileName := filepath.Clean(metaData.fileName)
	docFile, err := openOrCreateFile(absDocFileName)
	if err != nil {
		return nil, err
	}
	doc.docFile = docFile
	binLogFileName := metaData.fileName + ".logs.coditor"
	absBinLogFileName := filepath.Clean(binLogFileName)
	binLog, err := openBinLog(absBinLogFileName)
	if err != nil {
		return nil, err
	}
	doc.binLog = binLog

	contentBytes, err := ioutil.ReadAll(docFile)
	if err != nil {
		return nil, err
	}
	doc.content = string(contentBytes)
	rc, err := newRingCache(20)
	if err != nil {
		return nil, err
	}
	doc.ringCache = rc

	// TODO If there are many document.Ticker should by be changed to timingwheel.
	ticker := time.NewTicker(time.Duration((flushTime * int64(time.Second))))
	doc.flushTicker = ticker
	go func(doc *Document) {
		for {
			select {
			case <-doc.flushTicker.C:
				doc.flush()
			}
		}
	}(doc)
	return doc, nil
}

func (doc *Document) merge(content string, version DocVersion, userName string) (string, DocVersion, error) {
	doc.lock.Lock()
	defer func() {
		doc.lock.Unlock()
	}()

	if !doc.metaData.checkEditAble(userName) {
		return "", doc.metaData.Version, errors.New(userName + " can not editor file.")
	}

	if doc.metaData.Version != version {
		return "", doc.metaData.Version, errors.New("version must equest with doc's version.")
	}

	dmp := diffmatchpatch.New()
	patchs := dmp.PatchMake(doc.content, content)
	content, results := dmp.PatchApply(patchs, doc.content)
	// maybe the results check can be remove from here.
	for _, result := range results {
		if !result {
			return "", doc.metaData.Version, errors.New("merge to doc error.")
		}
	}

	patchsStr := dmp.PatchToText(patchs)

	// merge success!
	doc.content = content
	doc.metaData.Version++
	doc.ringCache.put(patchsStr)

	// save patchs to BinLog.
	doc.binLog.append(uint32(doc.metaData.Version), []byte(patchsStr))

	return patchsStr, doc.metaData.Version, nil
}

func (doc *Document) tail(version DocVersion, userName string) ([]string, error) {
	doc.lock.RLock()
	defer func() {
		doc.lock.RUnlock()
	}()

	if doc.metaData.IsPublic <= 0 && !doc.metaData.checkViewAble(userName) {
		return nil, errors.New(userName + " can not access file.")
	}

	if version > doc.metaData.Version || version < 0 {
		return nil, errors.New("version must less than document's version and great than 0.")
	}
	tempPatchss, err := doc.ringCache.tail(int(doc.metaData.Version - version))
	if err != nil {
		return nil, err
	}
	output := make([]string, len(tempPatchss))
	for i, tempPatchs := range tempPatchss {
		output[i] = tempPatchs.(string)
	}
	return output, nil
}

func (doc *Document) getContents(userName string) (string, error) {
	doc.lock.RLock()
	defer func() {
		doc.lock.RUnlock()
	}()

	if doc.metaData.IsPublic <= 0 && !doc.metaData.checkViewAble(userName) {
		return "", errors.New(userName + " can not access file.")
	}

	return doc.content, nil
}

func (doc *Document) getVersion(userName string) (DocVersion, error) {
	doc.lock.RLock()
	defer func() {
		doc.lock.RUnlock()
	}()

	if doc.metaData.IsPublic <= 0 && !doc.metaData.checkViewAble(userName) {
		return 0, errors.New(userName + " can not access file.")
	}

	return doc.metaData.Version, nil
}

func (doc *Document) setIsPublic(isPublic int, userName string) error {
	doc.lock.RLock()
	defer func() {
		doc.lock.RUnlock()
	}()

	if !doc.metaData.checkOwner(userName) {
		return errors.New(userName + " is not the owner.")
	}

	doc.metaData.IsPublic = isPublic
	return nil
}

func (doc *Document) setEditors(editors []string, userName string) error {
	doc.lock.RLock()
	defer func() {
		doc.lock.RUnlock()
	}()

	if !doc.metaData.checkOwner(userName) {
		return errors.New(userName + " is not the owner.")
	}
	doc.metaData.Editors = editors
	return nil
}

func (doc *Document) setViewers(viewers []string, userName string) error {
	doc.lock.RLock()
	defer func() {
		doc.lock.RUnlock()
	}()

	if !doc.metaData.checkOwner(userName) {
		return errors.New(userName + " is not the owner.")
	}
	doc.metaData.Viewers = viewers
	return nil
}

func (doc *Document) getEditors(userName string) ([]string, error) {
	doc.lock.RLock()
	defer func() {
		doc.lock.RUnlock()
	}()

	if !doc.metaData.checkOwner(userName) {
		return nil, errors.New(userName + " is not the owner.")
	}
	return doc.metaData.Editors, nil
}

func (doc *Document) getViewers(userName string) ([]string, error) {
	doc.lock.RLock()
	defer func() {
		doc.lock.RUnlock()
	}()

	if !doc.metaData.checkOwner(userName) {
		return nil, errors.New(userName + " is not the owner.")
	}
	return doc.metaData.Viewers, nil
}

func (doc *Document) setContent(content, userName string) error {
	doc.lock.RLock()
	defer func() {
		doc.lock.RUnlock()
	}()

	if !doc.metaData.checkEditAble(userName) {
		return errors.New(userName + " can not edit this document.")
	}

	doc.content = content
	return nil
}

func (doc *Document) flush() {
	contentBytes := []byte(doc.content)
	absDocFileName := filepath.Clean(doc.metaData.fileName)
	err := ioutil.WriteFile(absDocFileName, contentBytes, 0644)
	if err != nil {
		logger.Errorf("file flush error! FileName: %s.", absDocFileName)
		return
	}
	err = doc.metaData.save()
	if err != nil {
		logger.Errorf("MetaData file flush error! FileName: %s.", absDocFileName+".json.coditor")
	}
}

// to close document.
func (doc *Document) close(stopFlag int) {
	doc.flushTicker.Stop()
	doc.binLog.close()
	doc.docFile.Close()
	if stopFlag > 0 {
		doc.flush()
	}
}

// Ring Cache.
type RingCache struct {
	caches []interface{}
	index  int
	cap    int // capacity
	lock   sync.RWMutex
}

// New Ring Cache.Cap is the capacity of Ring Cache.
func newRingCache(cap int) (*RingCache, error) {
	if cap <= 0 {
		return nil, errors.New("capacity must great than 0.")
	}
	rc := &RingCache{}
	caches := make([]interface{}, cap)
	rc.caches = caches
	rc.index = 0
	rc.cap = cap
	return rc, nil
}

func (rc *RingCache) put(value interface{}) {
	if value == nil {
		return
	}

	rc.lock.Lock()
	defer func() {
		rc.lock.Unlock()
	}()

	rc.index = rc.index + 1
	rc.index = rc.index % rc.cap

	rc.caches[rc.index] = value
}

// return the last 'count' values.
func (rc *RingCache) tail(count int) ([]interface{}, error) {
	if count <= 0 || count > rc.cap {
		return nil, errors.New("count must great than 0 and less than capacity.")
	}

	rc.lock.RLock()
	defer func() {
		rc.lock.RUnlock()
	}()

	tempIndex := rc.index
	values := make([]interface{}, count)
	for i := 0; i < count; i++ {
		tempVal := rc.caches[tempIndex]
		if tempVal != nil {
			values[i] = tempVal
			tempIndex--
			if tempIndex < 0 {
				tempIndex = len(rc.caches) - 1
			}
		} else {
			return values[:i], nil
		}
	}
	return values, nil
}
