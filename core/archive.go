package core

import (
	"io"
	"net/http"

	"github.com/datatogether/archive"
	"github.com/datatogether/rewrite"
	"github.com/datatogether/warc"
	"github.com/ipfs/go-datastore"
	// "github.com/ipfs/go-ipfs/commands/files"
	"github.com/qri-io/cafs"
)

type ArchiveRequests struct {
	Store cafs.Filestore
}

type ArchiveUrlsParams struct {
	Urls []string
}

func (r ArchiveRequests) ArchiveUrls(p *ArchiveUrlsParams, path *datastore.Key) error {
	records := warc.Records{}

	for _, rawurl := range p.Urls {
		httpreq, err := http.NewRequest("GET", rawurl, nil)
		if err != nil {
			return err
		}

		rw := rewrite.NewWarcRecordRewriter(rawurl)
		if records, err = archive.ArchiveUrl(httpreq, rw, records); err != nil {
			return err
		}
	}

	pkg, err := archive.PackageRecords(p.Urls, records)
	if err != nil {
		return err
	}

	key, err := r.Store.Put(pkg, false)
	if err != nil {
		return err
	}

	*path = key
	return nil
}

type ArchiveUrlParams struct {
	Method string
	Url    string
	Body   io.Reader
}

func (r ArchiveRequests) ArchiveUrl(p *ArchiveUrlParams, path *datastore.Key) error {
	if p.Method == "" {
		p.Method = "GET"
	}
	httpreq, err := http.NewRequest(p.Method, p.Url, p.Body)
	if err != nil {
		return err
	}

	records := warc.Records{}
	rw := rewrite.NewWarcRecordRewriter(p.Url)
	if records, err = archive.ArchiveUrl(httpreq, rw, records); err != nil {
		return err
	}

	// // paths := []string{}
	// cafs.Walk(pkg, 0, func(f files.File, depth int) error {
	// 	// paths = append(paths, f.FullPath())
	// 	// fmt.Println(f.FullPath())
	// 	return nil
	// })
	pkg, err := archive.PackageRecords([]string{p.Url}, records)
	if err != nil {
		return err
	}

	key, err := r.Store.Put(pkg, false)
	if err != nil {
		return err
	}

	*path = key
	return nil
}
