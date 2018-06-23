package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// resp, err := http.Get("http://www.hao123.com")
	// defer resp.Body.Close()
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("Error!")
	// }
	// fmt.Println(body)

	// resp, err := http.Post("http://example.com/upload", "image/jpeg", &buf)
	// resp, err := http.PostForm("http://example.com/form", url.Values{"key": {"Value"}, "id": {"123"}})
	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}
	resp, _ := client.Get("http://www.hao123.com")

	req, err := client.NewRequest("GET", "http://example.com", nil)
	req.Header.Add("If-None-Match", `W/"wyzzy"`)
	resp, err := client.Do(req)
	tr := &http.Transport{
		TLSClientConfig:    &tls.Config{RootCAs: pool},
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get("http://example.com")

	http.Handle("/foo", fooHandler)
	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

	s := &http.Server{
		Addr:           ":8080",
		Handler:        myHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())


	func HelloServer(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "hello, world!\n")
	}

	func main() {
		http.HandleFunc("/hello", HelloServer)
		log.Fatal(http.ListenAndServe(":8080", nil))
	}

	func handler(w http.ResponseWriter, r *http.Request)
	{
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("This is an example server.\n"))
	}

	func main() {
		http.HandleFunc("/", handler)
		log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
		err := http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", nil)
		log.Fatal(err)
	}

	type Cookie struct {
		Name  string
		Value  string
		Path  string
		Domain  string
		Expires  time.Time
		RawExpires string
		MaxAge  int
		Secure  bool
		HttpOnly  bool
		Raw  string
		Unparsed  []string
	}

log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("/usr/share/doc"))))

















}
