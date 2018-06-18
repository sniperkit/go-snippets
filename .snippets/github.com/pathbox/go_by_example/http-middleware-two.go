package httpmiddleware

import (
	"net/http"
)

type AppendMiddleware struct {
	handler http.Handler
}

// ServeHTTP 是http.Handler接口默认使用的方法
func (this *AppendMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// this.handler.ServeHTTP(w, r) // 这是底层的　ServeHTTP, w r传入到了myHandler中
	w.Write([]byte("Hey, this is middleware!"))
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!"))
}

func main() {
	mid := &AppendMiddleware{http.HandlerFunc(myHandler)} // myHandler is a 闭包。传入了　ｗ和ｒ以及它的执行代码
	http.ListenAndServe(":8080", mid)
}
