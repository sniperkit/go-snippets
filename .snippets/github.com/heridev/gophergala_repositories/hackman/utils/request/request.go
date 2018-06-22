package request

import "net/http"
import "fmt"
import "io"
import "io/ioutil"

func Get(url, accessToken string) []byte {
	req, _ := http.NewRequest("GET", url, nil)

	authHeader := "token " + accessToken
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	res, _ := client.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func Post(url, accessToken string, payloadReader io.Reader, media bool) []byte {
	req, _ := http.NewRequest("POST", url, payloadReader)

        if media {
                req.Header.Set("Accept", "application/vnd.github.sersi-preview+json")
        } else {
                req.Header.Set("Accept", "application/json")
        }

	if accessToken != "" {
		authHeader := "token " + accessToken
		req.Header.Set("Authorization", authHeader)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, _ := client.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}

func Put(url, accessToken string, payloadReader io.Reader) []byte {
	req, _ := http.NewRequest("PUT", url, payloadReader)

	if accessToken != "" {
		fmt.Println("including")
		authHeader := "token " + accessToken
		req.Header.Set("Authorization", authHeader)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	res, _ := client.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return body
}
