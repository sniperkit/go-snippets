package core

import (
	"fmt"
	// "fmt"
	"github.com/datatogether/sql_datastore"
	"github.com/ipfs/go-datastore"
	"testing"
)

func TestUrlStorage(t *testing.T) {
	store := datastore.NewMapDatastore()

	u := &Url{Url: "http://youtube.com"}
	if err := u.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	u.ContentLength = 10
	if err := u.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	u2 := &Url{Id: u.Id}
	if err := u2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if !u2.Created.Equal(u.Created) {
		t.Errorf("created doesn't match: %s != %s", u2.Created.String(), u.Created.String())
	}

	if !u2.Updated.Equal(u.Updated) {
		t.Errorf("updated doesn't match: %s != %s", u2.Updated.String(), u.Updated.String())
	}

	if err := u.Delete(store); err != nil {
		t.Error(err.Error())
		return
	}
}

func TestUrlSQLStorage(t *testing.T) {
	defer resetTestData(appDB, "urls", "links")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Url{}); err != nil {
		t.Error(err)
		return
	}

	u := &Url{Url: "http://youtube.com"}
	if err := u.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	u.ContentLength = 10
	if err := u.Save(store); err != nil {
		t.Error(err.Error())
		return
	}

	u2 := &Url{Url: "http://youtube.com"}
	if err := u2.Read(store); err != nil {
		t.Error(err.Error())
		return
	}

	if !u2.Created.Equal(u.Created) {
		t.Errorf("created doesn't match: %s != %s", u2.Created.String(), u.Created.String())
	}

	if !u2.Updated.Equal(u.Updated) {
		t.Errorf("updated doesn't match: %s != %s", u2.Updated.String(), u.Updated.String())
	}

	if err := u.Delete(store); err != nil {
		t.Error(err.Error())
		return
	}
}

func TestShouldEnqueue(t *testing.T) {
	defer resetTestData(appDB, "urls")

	store := sql_datastore.NewDatastore(appDB)
	if err := store.Register(&Url{}); err != nil {
		t.Error(err)
		return
	}

	epa := &Url{Url: "http://www.epa.gov"}
	if err := epa.Read(store); err != nil {
		t.Fatal(err.Error())
	}

	cases := []struct {
		url       *Url
		get, head bool
	}{
		// TODO - this test isn't working the func properly. Should enhance with DB interaction
		{&Url{Url: "https://youtube.com"}, true, true},
		{&Url{Url: "http://www.fda.gov"}, true, true},
		{&Url{Url: "http://epa.gov/new"}, true, true},
		{epa, true, true},
	}

	for _, c := range cases {
		u := c.url
		head := u.ShouldEnqueueHead()
		if head != c.head {
			t.Errorf("shouldEnqueueHead: %s error. expected %t, got %t", u.Url, c.head, head)
		}

		get := u.ShouldEnqueueGet()
		if get != c.get {
			t.Errorf("shouldEnqueueGet: %s expected %t, got %t", u.Url, c.get, get)
		}
	}
}

// func TestUrlGet(t *testing.T) {
// 	u := &Url{Url: "https://www.apple.com"}
// 	done := make(chan bool, 0)
// 	links, err := u.Get(appDB, func(err error) {
// 		if err != nil {
// 			fmt.Println(err.Error())
// 		}
// 		done <- true
// 	})

// 	if err != nil {
// 		t.Error(err.Error())
// 	}

// 	if len(links) == 0 {
// 		t.Error("didn't find any links?")
// 	}
// 	<-done
// }

func TestUrlSuspectedContentUrl(t *testing.T) {
	cases := []struct {
		url    *Url
		expect bool
	}{
		{&Url{
			Url:           "https://opendata.epa.gov/home.xhtml?view",
			Status:        200,
			ContentType:   "text/html;charset=UTF-8",
			ContentLength: 18327,
			ContentSniff:  "text/plain; charset=utf-8",
		}, false},
	}

	for i, c := range cases {
		if got := c.url.SuspectedContentUrl(); got != c.expect {
			t.Errorf("case %d fail: %t != %t", i, c.expect, got)
		}
	}
}

func CompareUrls(a, b *Url) error {
	if a == nil && b != nil || a != nil && b == nil {
		return fmt.Errorf("nil mismatch %s != %s", a, b)
	} else if a == nil && b == nil {
		return nil
	}

	if a.Id != b.Id {
		return fmt.Errorf("id mismatch: %s != %s", a.Id, b.Id)
	}

	if !a.Created.Equal(b.Created) {
		return fmt.Errorf("created mismatch: %s != %s", a.Created, b.Created)
	}

	if !a.Updated.Equal(b.Updated) {
		return fmt.Errorf("updated mismatch: %s != %s", a.Updated, b.Updated)
	}

	if a.Url != b.Url {
		return fmt.Errorf("url mismatch: %s != %s ", a.Url, b.Url)
	}
	if a.Hash != b.Hash {
		return fmt.Errorf("Hash mistmatch: %s != %s", a.Hash, b.Hash)
	}

	if a.LastGet != b.LastGet {
		return fmt.Errorf("LastGet mistmatch: %s != %s", a.LastGet, b.LastGet)
	}
	if a.LastHead != b.LastHead {
		return fmt.Errorf("LastHead mistmatch: %s != %s", a.LastHead, b.LastHead)
	}
	if a.Status != b.Status {
		return fmt.Errorf("Status mistmatch: %s != %s", a.Status, b.Status)
	}
	if a.ContentType != b.ContentType {
		return fmt.Errorf("ContentType mistmatch: %s != %s", a.ContentType, b.ContentType)
	}
	if a.ContentSniff != b.ContentSniff {
		return fmt.Errorf("ContentSniff mistmatch: %s != %s", a.ContentSniff, b.ContentSniff)
	}
	if a.ContentLength != b.ContentLength {
		return fmt.Errorf("ContentLength mistmatch: %s != %s", a.ContentLength, b.ContentLength)
	}
	if a.FileName != b.FileName {
		return fmt.Errorf("FileName mistmatch: %s != %s", a.FileName, b.FileName)
	}
	if a.Title != b.Title {
		return fmt.Errorf("Title mistmatch: %s != %s", a.Title, b.Title)
	}
	if a.DownloadTook != b.DownloadTook {
		return fmt.Errorf("DownloadTook mistmatch: %s != %s", a.DownloadTook, b.DownloadTook)
	}
	if a.HeadersTook != b.HeadersTook {
		return fmt.Errorf("HeadersTook mistmatch: %s != %s", a.HeadersTook, b.HeadersTook)
	}
	// TODO - proper comparison
	// if a.Headers != b.Headers {
	// 	return fmt.Errorf("Headers mistmatch: %s != %s", a.Headers, b.Headers)
	// }
	// if a.Meta != b.Meta {
	// 	return fmt.Errorf("Meta mistmatch: %s != %s", a.Meta, b.Meta)
	// }
	if a.ContentUrl != b.ContentUrl {
		return fmt.Errorf("ContentUrl mistmatch: %s != %s", a.ContentUrl, b.ContentUrl)
	}
	if a.Uncrawlable != b.Uncrawlable {
		return fmt.Errorf("Uncrawlable mistmatch: %s != %s", a.Uncrawlable, b.Uncrawlable)
	}

	return nil
}
