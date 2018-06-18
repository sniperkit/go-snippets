package main

import (
	"log"
	"net/http"
	"os"

	"github.com/smaxwellstewart/golang-decorator-pattern/client"
)

// We can construct a client with all the functionality we need.
func main() {
	c := client.Decorate(client.MockClient{},
		client.Logging(log.New(os.Stdout, "[logging decorator] ", 0)),
		client.Audit(log.New(os.Stdout, "[audit decorator] ", 0)),
	)
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	c.Do(req)
}
