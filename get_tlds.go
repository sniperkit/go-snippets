// +build ignore

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"code.google.com/p/go.net/idna"
)

func main() {
	res, err := http.Get("https://data.iana.org/TLD/tlds-alpha-by-domain.txt")
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	lines := bytes.Split(bytes.ToLower(data), []byte("\n"))
	var idntld []string
	var cctld, gtld [][]byte
	for _, line := range lines[1:] {
		if len(line) == 2 {
			cctld = append(cctld, line)
			continue
		}
		if bytes.HasPrefix(line, []byte("xn--")) {
			tld, err := idna.ToUnicode(string(line))
			if err != nil {
				log.Fatal(err)
			}
			idntld = append(idntld, tld)
			continue
		}
		if len(line) > 0 {
			gtld = append(gtld, line)
		}
	}

	fmt.Printf("regexen[\"GTLD\"] = \"(?:%s)\"\n", bytes.Join(gtld, []byte("|")))
	fmt.Printf("regexen[\"CCTLD\"] = \"(?:%s)\"\n", bytes.Join(cctld, []byte("|")))
	fmt.Printf("regexen[\"IDNTLD\"] = \"(?:%s)\"\n", strings.Join(idntld, "|"))
}
