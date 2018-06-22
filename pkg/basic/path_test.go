package basic

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"testing"
)

func TestPathJoin(t *testing.T) {
	log.Println(filepath.Join("", "a"))
	log.Println(filepath.Join("", "a", ".dat"))

	log.Println(filepath.Join("/a/b/c", "/b/c"))
	log.Println(filepath.Join("http://a.b.c.d/", "/a/b/c/"))
}

func TestDirPath(t *testing.T) {
	var ex, _ = os.Executable()
	var basePath = path.Dir(ex)
	log.Println(basePath)
	log.Println(filepath.Dir(basePath))
}
