/*
* @Author: wuhailin
* @Date:   2017-10-20 09:37:29
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-10-20 10:34:48
 */

package main

import (
	"compress/gzip"
	"io/ioutil"
	"os"
)

func main() {
	filename := "story.zip"
	f, _ := os.Create("story.gz")
	defer f.Close()
	r := gzip.NewWriter(f)
	defer r.Close()
	content, _ := ioutil.ReadFile(filename)
	r.Name = filename
	r.Write(content)
}
