package gopherpaint

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/user"
	"time"
)

type Image struct {
	OwnerID      string
	Blobkey      appengine.BlobKey
	Style        string
	CreationTime time.Time
	MD5          string
	Size         int64
}

func (m *Image) GenerateID() string {
	return (string)(m.Blobkey) + "_" + m.OwnerID
}

func GenID(blobkey, oid string) string {
	return blobkey + "_" + oid
}

func ImagesPOST(c appengine.Context,
	usr *user.User,
	blobinfo *blobstore.BlobInfo,
	style string) (*datastore.Key, error) {
	data := &Image{
		OwnerID:      usr.ID,
		Blobkey:      blobinfo.BlobKey,
		Style:        style,
		CreationTime: time.Now(),
		MD5:          blobinfo.MD5,
		Size:         blobinfo.Size,
	}

	mcKey := data.GenerateID()
	key := datastore.NewKey(c, "Images", mcKey, 0, nil)

	mcItem := &memcache.Item{
		Key:    mcKey,
		Object: data,
	}
	memcache.Gob.Set(c, mcItem)

	return datastore.Put(c, key, data)
}

func Images_OfUser_GET(c appengine.Context, usr *user.User) ([]Image, error) {
	mcKey := "pics_" + usr.ID
	var items []Image
	_, err := memcache.Gob.Get(c, mcKey, &items)
	if err == nil {
		return items, nil
	}

	q := datastore.NewQuery("Images").
		Filter("OwnerID =", usr.ID).
		Order("-CreationTime")
	_, err = q.GetAll(c, &items)

	// Add list of images to memcache
	mcItem := &memcache.Item{
		Key:    mcKey,
		Object: items,
	}
	memcache.Gob.Set(c, mcItem)

	return items, err
}

func Images_GetOne(c appengine.Context, usr *user.User, blobkey string) (*Image, error) {
	// Try to get from memcache
	mcKey := GenID(blobkey, usr.ID)
	var item Image
	_, err := memcache.Gob.Get(c, mcKey, &item)
	if err == nil {
		return &item, nil
	}

	// Get from datastore
	q := datastore.NewQuery("Images").
		Filter("OwnerID =", usr.ID).
		Filter("Blobkey =", blobkey)
	var images []Image
	_, err = q.GetAll(c, &images)
	if err != nil && len(images) > 0 {
		return nil, err
	}
	img := images[0]

	// Saves in memcache so we don't hit datastore
	mcItem := &memcache.Item{
		Key:    mcKey,
		Object: img,
	}
	memcache.Gob.Set(c, mcItem)

	return &img, nil
}

func Images_UpdateStyle(c appengine.Context,
	usr *user.User,
	blobkey string,
	newstyle string) (*datastore.Key, error) {
	// Retrieve key
	m, err := Images_GetOne(c, usr, blobkey)
	if err != nil {
		c.Infof("Images_GetOne: %v", err)
		return nil, err
	}

	// Updates the value
	if m.Style == newstyle || m.OwnerID != usr.ID {
		c.Infof("Old style %v new %v | owner %v current %v", m.Style, newstyle, m.OwnerID, usr.ID)
		return nil, nil
	}

	m.Style = newstyle
	key := datastore.NewKey(c, "Images", m.GenerateID(), 0, nil)
	nk, err := datastore.Put(c, key, m)

	if err != nil {
		c.Infof("datastore put: %v", err)
		return nk, err
	}
	memcache.Delete(c, "pics_"+usr.ID)
	mcKey := GenID(blobkey, usr.ID)
	mcItem := &memcache.Item{
		Key:    mcKey,
		Object: m,
	}
	memcache.Gob.Set(c, mcItem)

	return nk, err
}

func Images_Delete(c appengine.Context,
	usr *user.User,
	blobkey string) error {
	m, err := Images_GetOne(c, usr, blobkey)
	if err != nil {
		return err
	}
	if m.OwnerID != usr.ID {
		return nil
	}
	itemKey := GenID(blobkey, usr.ID)
	key := datastore.NewKey(c, "Images", itemKey, 0, nil)
	err = datastore.Delete(c, key)
	if err != nil {
		return err
	}

	return memcache.Delete(c, itemKey)
}
