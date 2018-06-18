package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// Using Google prediction API to detect the language of
// given words (or sentences) and translate it to english.
func main() {

	s := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	// create a loop of s.Scan to act as interpreter
	for s.Scan() {
		res, err := post(s.Text())
		if err != nil {
			fmt.Println("failed:", err)
			continue
		}

		// predictedLang := res.OutputMulti[0].Label
		// label, ok := labels[res.OutputLabel]
		// if !ok {
		// 	label = res.OutputLabel
		// }
		fmt.Println(res.OutputMulti)
		fmt.Print("> ")
	}
}

func detectLang(t string) string {
	return ""
}

func translate(t string, fromLang string, toLang string) string {
	return ""
}

// response struct to handle
// response from Google Prediction API
type response struct {
	Kind        string
	ID          string
	OutputLabel string
	OutputMulti []struct {
		Label string
		Score string
	}
	OutputValue float64
}

func post(q string) (*response, error) {
	v := url.Values{
		"model":  []string{"Language Detection"},
		"Phrase": []string{q},
	}
	const endpoint = "http://try-prediction.appspot.com/predict"
	res, err := http.PostForm(endpoint, v)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// decode response body to result struct
	var result response
	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
