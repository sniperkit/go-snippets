package lib

import (
	"fmt"
	"sync"
	"time"
)

const (
	outdatedTime = 3600 //sec
)

//ImageStore keep the images which can be reused
type ImageStore struct {
	lock   *sync.RWMutex
	images map[string]*Image
}

//Image ...
type Image struct {
	Name       string
	Tag        string
	ActiveTime int64
}

//NewImageStore ...
func NewImageStore() *ImageStore {
	return &ImageStore{
		lock:   new(sync.RWMutex),
		images: make(map[string]*Image),
	}
}

//Put image in with the key
func (is *ImageStore) Put(img, tag string) {
	if len(img) == 0 || len(tag) == 0 {
		return
	}

	key := fmt.Sprintf("%s:%s", img, tag)
	if image, ok := is.images[key]; ok {
		image.ActiveTime = time.Now().Unix()
		return
	}

	image := &Image{
		Name:       img,
		Tag:        tag,
		ActiveTime: time.Now().Unix(),
	}

	is.images[key] = image
}

//Get what you want
func (is *ImageStore) Get(key string) (*Image, bool) {
	img, ok := is.images[key]
	if ok {
		img.ActiveTime += 5 //for safe calling
	}
	return img, ok
}

//Garbage collection
func (is *ImageStore) Garbage() []*Image {
	outdatedOnes := make([]*Image, 0)
	now := time.Now().Unix()
	for k, v := range is.images {
		if now > v.ActiveTime+outdatedTime {
			outdatedOnes = append(outdatedOnes, v)
			delete(is.images, k)
		}
	}

	return outdatedOnes
}
