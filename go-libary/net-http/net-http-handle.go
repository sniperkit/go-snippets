package main

import (
	// "fmt"
	"net/http"
)

func main() {
	http.Handle("/test", http.FileServer(http.Dir("./")))
	http.Handle("/download/", http.StripPrefix("/download/", http.FileServer(http.Dir("./"))))
	http.Handle("/code/", http.StripPrefix("/code/", http.FileServer(http.Dir("/Users/path/code")))) // 会自动去找目录下的index.html文件并且渲染,如果没有的话，就显示目录文件
	http.ListenAndServe(":9998", nil)

}
