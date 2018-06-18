package main

import (
	"fmt"
	// "io/ioutil"
	"net/http"
)

func Test() {
	req, err := http.NewRequest("GET", "http://www.hao123.com", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.SetBasicAuth("test", "123456")
	fmt.Println(req.Proto)
	cookie := &http.Cookie{
		Name:  "test",
		Value: "12",
	}
	req.AddCookie(cookie)
	fmt.Println(req.Cookie("test")) //获取cookie
	fmt.Println(req.Cookies())
	req.Header.Set("User-Agent", "user-agent")
	fmt.Println(req.UserAgent())
	fmt.Println(req.URL)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		content, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(content))
	}
}

func main() {
	Test()
}
