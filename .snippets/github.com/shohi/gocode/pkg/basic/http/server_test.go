package http

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"testing"
	"time"
)

func readContent(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, errors.New("resp is nil")
	}

	contents, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return contents, err
}

func TestServer(t *testing.T) {
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello"))
	})

	s := &http.Server{Handler: handler}
	go s.Serve(l)

	resp, err := http.Get("http://" + l.Addr().String())
	log.Println(err, resp.Body)

	bs, _ := readContent(resp)
	log.Printf("content: %s", string(bs))
}

func TestListenAndServe(t *testing.T) {
	server := &http.Server{Addr: ":10010"}
	errC := make(chan error)
	go func() { errC <- server.ListenAndServe() }()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		log.Println("Hello world")
	case <-errC:
		log.Println("Hello error")
	}
}

// this blocks
/*
func TestListenAndServeToHang(t *testing.T) {
	server := &http.Server{Addr: ":10010"}
	err := server.ListenAndServe()

	log.Println(err)
}
*/
