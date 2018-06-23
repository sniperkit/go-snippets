/*
* @Author: wuhailin
* @Date:   2017-10-19 18:44:21
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-10-19 18:50:14
 */

package main

import (
	"archive/zip"
	"io/ioutil"
	"os"
)

func main() {
	filename := "story.csv"
	fzip, _ := os.Create("story.zip")
	defer fzip.Close()
	w := zip.NewWriter(fzip)
	defer w.Close()
	fw, _ := w.Create(filename)
	fileContent, err := ioutil.ReadFile(filename)
	if nil != err {
		panic(err)
	}
	if _, err = fw.Write(fileContent); nil != err {
		panic(err)
	}
}
