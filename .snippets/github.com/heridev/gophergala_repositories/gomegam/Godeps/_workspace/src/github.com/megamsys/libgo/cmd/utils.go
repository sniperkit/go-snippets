package cmd

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"crypto/hmac"
	"crypto/sha1"
	"github.com/tsuru/config"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	X_Megam_EMAIL       = "X-Megam-EMAIL"
	X_Megam_APIKEY      = "X-Megam-APIKEY"
	X_Megam_DATE        = "X-Megam-DATE"
	X_Megam_HMAC        = "X-Megam-HMAC"
	API_GATEWAY_VERSION = "/v1"
	Content_Type               = "Content-Type"
	application_json           = "application/json"
	Accept                     = "Accept"
	application_vnd_megam_json = "application/vnd.megam+json"
)

type Authly struct {
	UrlSuffix string
	Date      string
	Email     string
	APIKey    string
	JSONBody  []byte
	AuthMap   map[string]string
}

func NewAuthly(urlsuffix string, jsonbody []byte) (*Authly, error) {
	email, err := config.GetString("api:email")

	if err != nil {
		return nil, fmt.Errorf("Failed to find the email (%s).", err)
	}

	api_key, err := config.GetString("api:api_key")

	if err != nil {
		return nil, fmt.Errorf("Failed to find the api_key (%s).", err)
	}

	return &Authly{
		UrlSuffix: urlsuffix,
		Date:      time.Now().Format(time.RFC850),
		Email:     email,
		APIKey:    api_key,
		JSONBody:  jsonbody,
		AuthMap:   map[string]string{},
	}, nil
}

func (authly *Authly) AuthHeader() error {
	timeStampedPath := authly.Date + "\n" + API_GATEWAY_VERSION + authly.UrlSuffix
	md5Body := ""
	if len(authly.JSONBody) > 0 {
		md5Body = GetMD5Hash(authly.JSONBody)
	}
	headMap := make(map[string]string)
	headMap[X_Megam_DATE] = authly.Date
	headMap[X_Megam_EMAIL] = authly.Email
	headMap[X_Megam_APIKEY] = authly.APIKey
	headMap[Accept] = application_vnd_megam_json
	headMap[X_Megam_HMAC] = authly.Email + ":" + CalcHMAC(authly.APIKey, (timeStampedPath+"\n"+md5Body))
	headMap["Content-Type"] = "application/json"
	authly.AuthMap = headMap
	return nil
}

func CalcHMAC(secret string, message string) string {
    key := []byte(secret)
    h := hmac.New(sha1.New, key)
    h.Write([]byte(message))    
    sumh := h.Sum(nil)

    sumi := make([]string, len(sumh))
    for i := 0; i < len(sumh); i++ {
    	sumi[i] = ("00" + fmt.Sprintf("%x",(sumh[i] & 0xff)))
    	sumi[i]=sumi[i][len(sumi[i])-2: len(sumi[i])]
    }
    outs := strings.Join(sumi, "")
    return outs
  }
 

func GetMD5Hash(text []byte) string {
	hasher := md5.New()
	hasher.Write(text)
	return base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	
}

func CalcBase64(data string) string {
	var buf bytes.Buffer
	encoder := base64.NewEncoder(base64.StdEncoding, &buf)
	encoder.Write([]byte(data))
	encoder.Close()
	return buf.String()
}

func GetURL(path string) (string, error) {
	var prefix string
	target, err := config.GetString("api:server")
	if err != nil {
		return "", err
	}
	if m, _ := regexp.MatchString("^https?://", target); !m {
		prefix = "http://"
	}
	return prefix + strings.TrimRight(target, "/") + strings.TrimRight(API_GATEWAY_VERSION, "/") + path, nil
}

/*
type ServiceModel struct {
	Service   string
	Instances []string
}

func ShowServicesInstancesList(b []byte) ([]byte, error) {
	var services []ServiceModel
	err := json.Unmarshal(b, &services)
	if err != nil {
		return []byte{}, err
	}
	if len(services) == 0 {
		return []byte{}, nil
	}
	table := NewTable()
	table.Headers = Row([]string{"Services", "Instances"})
	for _, s := range services {
		insts := strings.Join(s.Instances, ", ")
		r := Row([]string{s.Service, insts})
		table.AddRow(r)
	}
	return table.Bytes(), nil
}*/
