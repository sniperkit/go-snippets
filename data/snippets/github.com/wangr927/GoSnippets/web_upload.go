package main

import (
	"net/http"
	"fmt"
	"path"
	"time"
	"strconv"
	"os"
	"io"
)

const index_page = `<html><head><title>上传文件</title></head>
<body><form enctype="multipart/form-data" action="/upload" method="post">
<input type="file" name="uploadfile">
<input type="hidden" name="token" value="{...{.}...}">
<input type="submit" value="upload">
</form></body></html>`

func main() {
	http.HandleFunc("/",index)
	http.HandleFunc("/upload",upload)
	err := http.ListenAndServe(":6767",nil)
	if err != nil {
		fmt.Println("Server Boot Failed:", err.Error())
		return
	}
}

func index(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte(index_page))
	// 使用http.ResponseWriter 渲染数据，需要将固定的html内容以byte的形式发出
}

func upload(writer http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(32<<20)
	file, handler, err := request.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	// 注意x.File 在使用完后一定要shut down gracefully
	ext := path.Ext(handler.Filename)
	fileNewName := string(time.Now().Format("20060102150405")) + strconv.Itoa(time.Now().Nanosecond()) + ext

	f, err := os.OpenFile("./ppp"+fileNewName,os.O_WRONLY|os.O_CREATE,0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	io.Copy(f, file)
	fmt.Fprintln(writer, "upload done!" + fileNewName)
}
