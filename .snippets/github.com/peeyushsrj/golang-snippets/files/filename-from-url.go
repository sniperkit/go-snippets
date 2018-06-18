package golangSnippets

import (
	"log"
	"net/url"
	"path/filepath"
)

//FilenameFromUrl takes input as a escaped url & outputs filename from it (unescaped - normal one)
func FilenameFromUrl(urlstr string) string {
	u, err := url.Parse(urlstr)
	if err != nil {
		log.Fatal("Error due to parsing url: ", err)
	}
	x, _ := url.QueryUnescape(u.EscapedPath())
	return filepath.Base(x)
}
