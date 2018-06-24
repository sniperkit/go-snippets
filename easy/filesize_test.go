package main

import (
	"log"
	"os"
	"testing"
)

func TestFilesize(t *testing.T) {
	for k, v := range map[string]int64{
		"/usr/src/linux/COPYING": 18693} {
		if r := filesize(k); r != v {
			t.Errorf("failed: filesize %s is %d, got %d",
				k, v, r)
		}
	}
}

func filesize(q string) int64 {
	data, err := os.Open(q)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	stat, err := data.Stat()
	if err != nil {
		log.Fatal(err)
	}
	return stat.Size()
}
