package main

import (
	"os"
	"fmt"
	"strings"
	"crypto/md5"
	"io"
	"encoding/hex"
	http "net/http"
)

func checkErr(err error) {
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
}

func DownloadHttp(url string) error {

	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]

	h := md5.New()
	io.WriteString(h, fileName)
	s := hex.EncodeToString(h.Sum(nil))

	output, err := os.Create(fmt.Sprintf("/tmp/%s", s))
	if err != nil {
		checkErr(err)
		return err
	}

	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		checkErr(err)
		return err
	}

	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		checkErr(err)
		return err
	}

	fmt.Println(n, "bytes downloaded.")

	return nil
}

func DownloadFtp(url string) error {



	return nil
}

func Download(url string) error {

	if strings.Contains(url, "ftp://") {
		return DownloadFtp(url)
	} else if strings.Contains(url, "http://") {
		return DownloadHttp(url)
	} else {
		return fmt.Errorf("url sting not valid")
	}

	return nil
}

func main() {
//	err := Download("ftp://mds:mds@mds.kallisto.ru/pionerfm/Dzherom_K._Dzherom_-_Kot_Dika_Dankermana.mp3")
	err := Download("http://mds.kallisto.ru/pionerfm/Dzherom_K._Dzherom_-_Kot_Dika_Dankermana.mp3")
	if err != nil {
		checkErr(err)
	}
}