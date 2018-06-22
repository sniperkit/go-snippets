package mp3fetcher

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/GopherGala/i_love_indexes/conn_throttler"
	"github.com/mikkyang/id3-go"
)

func ArtisteAndAlbum(songURL string) string {
	u, err := url.Parse(songURL)
	if err != nil {
		log.Println("Wrong URL", err)
		return ""
	}
	sem := conn_throttler.Acquire(u.Host)
	defer sem.Release()

	res, err := http.Get(songURL)
	if err != nil {
		log.Println("fail to request", songURL, ":", err)
		return ""
	}

	var buffer [1024]byte
	_, err = res.Body.Read(buffer[:])
	if err != nil {
		log.Println("fail to read first KB of", songURL, ":", err)
		return ""
	}
	res.Body.Close()

	f, err := ioutil.TempFile("/tmp", "crawlmp3")
	if err != nil {
		log.Println("fail to create tmp file", songURL, ":", err)
		return ""
	}
	b := bytes.NewBuffer(buffer[:])
	_, err = io.Copy(f, b)
	if err != nil {
		log.Println("fail to copy to tmp file", songURL, ":", err)
		return ""
	}
	err = f.Close()
	if err != nil {
		log.Println("fail to close tmp file", songURL, ":", err)
	}

	mp3, err := id3.Open(f.Name())
	if err != nil {
		log.Println("fail to open mp3 file", songURL, ":", err)
		return ""
	}

	return mp3.Artist() + " " + mp3.Album()
}
