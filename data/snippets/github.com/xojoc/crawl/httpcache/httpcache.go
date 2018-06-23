/*  Copyright (C) 2018 Alexandru Cojocaru

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>. */

package httpcache

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"xojoc.pw/must"
)

func init() {
	http.DefaultTransport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		DisableKeepAlives:     true,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
}

type Cache interface {
	Fetch(string) (*http.Response, bool)
	IsCached(string) bool
}

type DiskCache struct {
	Base string
}

func NewDiskCache(path string) *DiskCache {
	d := &DiskCache{}
	d.Base = path
	return d
}

func (d *DiskCache) md5path(u string) string {
	m := fmt.Sprintf("%x", md5.Sum([]byte(u)))
	// TODO:
	return d.Base + "/" + m[:2] + "/" + m[2:]
}

type myBody struct {
	responseBody io.ReadCloser
	file         *os.File

	closed bool
}

func (b *myBody) Close() error {
	if b.closed {
		return nil
	}
	b.closed = true

	//	berr := b.responseBody.Close()
	ferr := b.file.Close()
	if ferr != nil {
		//		log.Println("ferr: ", ferr)
		return ferr
	}
	/*
		if berr != nil {
			log.Println("berr: ", berr)
			return berr
		}
	*/
	return nil
}
func (r *myBody) Read(p []byte) (int, error) {
	return r.responseBody.Read(p)
}

func (d *DiskCache) Fetch(u string) (*http.Response, error) {
	p := d.md5path(u)
	_, err := os.Stat(p)
	if err == nil {
		f := must.Open(p)
		r, err := http.ReadResponse(bufio.NewReader(f), nil)
		if r != nil {
			r.Body = &myBody{
				responseBody: r.Body,
				file:         f,
			}
			r.Header.Add("X-From-Cache", "true")
		}
		return r, err
	}
	// remove
	http.DefaultClient.Timeout = 5 * time.Second

	r, err := http.Get(u)
	if err != nil {
		return nil, err
	}
	must.OK(os.MkdirAll(path.Dir(p), 0766))
	w := must.Create(p)
	defer must.Close(w)
	err = r.Write(w)
	if err != nil {
		return nil, err
	}
	f := must.Open(p)
	r, err = http.ReadResponse(bufio.NewReader(f), nil)
	r.Body = &myBody{
		responseBody: r.Body,
		file:         f,
	}
	return r, err
}
