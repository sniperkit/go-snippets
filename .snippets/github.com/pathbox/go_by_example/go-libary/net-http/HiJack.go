package main

import (
	"fmt"
	"net/http"
)

func HiJack(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "webservr doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	bufrw.WriteString("Now we're speaking raw TCP. Say hi: \n")
	bufrw.Flush()

	fmt.Fprintf(bufrw, "You said: %s Bye.\n", "Good")
	bufrw.Flush()
}

func main() {
	http.HandleFunc("/hijack", HiJack)
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Println(err)
	}
}
