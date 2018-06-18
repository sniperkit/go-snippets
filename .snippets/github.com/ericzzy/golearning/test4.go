package main

import (
    "fmt"
    "crypto/hmac"
    "crypto/sha1"
    "encoding/base64"
    "time"
)

func generateAuth(accessID, accessKey, action string) (string, error) {
    date := time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
    toSign := "GET\n"
    toSign += "\n"
    toSign += "application/json\n"
    toSign += date + "\n"
    toSign += "/" + action
    hash := hmac.New(sha1.New, []byte(accessKey))
    hash.Write([]byte(toSign))
    signed := base64.StdEncoding.EncodeToString(hash.Sum(nil))
    auth := "Authorization: LWS " + accessID + ":" + signed
    return auth, nil
}

func main() {
    accessID := "CWXYQX8Z1LDPZ5DJ14NG"
    accessKey := "k/MOjnPk2yrjNecndc212kNocLX1yPVaq4PTNg=="
    auth, err := generateAuth(accessID, accessKey, "PutMetricData")
    fmt.Println(err)
    fmt.Println(auth)
}
