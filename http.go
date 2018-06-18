package toolkit

/*
Http related
*/
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	httpurl "net/url"
	"strings"
)

func HttpCall(url string, callType string,
	datas []byte,
	config M) (*http.Response, error) {
	var err error
	if config == nil {
		config = M{}
	}
	config.Set("calltype", callType)

	var req *http.Request

	//fmt.Println(callType)
	//return nil, fmt.Errorf("ERROR: "+callType+" %v", config.Has("formvalues"))
	/*
		if callType == "POST" && config.Has("formvalues") {
			fmt.Println("Parsing values")
			fvs := config["formvalues"].(M)
			vs := httpurl.Values{}
			for k, v := range fvs {
				fmt.Printf("Add formvalue %s = %v \n", k, v)
				//q += k + "="
				//q += v.(string)
				vs.Set(k, v.(string))
			}
			//rdr := bytes.NewBuffer([]byte(q))
			//req, err = http.NewRequest(callType, url, rdr)
			return http.PostForm(url, vs)
		} else */
	if datas == nil || len(datas) == 0 {
		req, err = http.NewRequest(callType, url, nil)
	} else {
		//fmt.Printf("Datas: %v\n%s \n", datas, string(datas))
		rdr := bytes.NewBuffer(datas)
		req, err = http.NewRequest(callType, url, rdr)
	}
	if err != nil {
		return nil, err
	}

	if config.Has("auth") {
		authType := strings.ToLower(config.Get("auth", "").(string))
		if authType == "basic" {
			username := config.Get("user", "").(string)
			pass := config.Get("password", "").(string)
			req.SetBasicAuth(username, pass)
		}
	}

	return httpcall(req, config)
}

func httpcall(req *http.Request, config M) (*http.Response, error) {
	var client *http.Client

	//-- handling cookie
	if config.Has("cookie") == false {
		client = new(http.Client)
	} else {
		//-- preparing cookie jar and http client
		jar, err := cookiejar.New(nil)
		if err != nil {
			return nil, fmt.Errorf("Unable to initialize cookie jar: %s", err.Error())
		}

		tjar := config.Get("cookie", nil).(*cookiejar.Jar)

		if tjar != nil {
			jar = tjar
		}

		client = &http.Client{
			Jar: jar,
		}
	}

	//--- handling header
	if headers, hasHeaders := config["headers"]; hasHeaders {
		mHeaders := headers.(M)
		for k, v := range mHeaders {
			req.Header.Add(k, v.(string))
		}
	}

	var resp *http.Response
	var errCall error

	if config.Has("formvalues") {
		fvs := config["formvalues"].(M)
		vs := httpurl.Values{}
		for k, v := range fvs {
			// fmt.Printf("Add formvalue %s = %v \n", k, v)
			//q += k + "="
			//q += v.(string)
			vs.Set(k, v.(string))
		}
		resp, errCall = client.PostForm(req.URL.String(), vs)
	} else {
		resp, errCall = client.Do(req)
	}
	if errCall == nil {
		if expectedStatus := config.Get("expectedstatus", 0).(int); expectedStatus != 0 && resp.StatusCode != expectedStatus {
			return nil, fmt.Errorf("Code error: " + resp.Status)
		}
	}
	return resp, errCall
}

func HttpContent(r *http.Response) []byte {
	defer r.Body.Close()
	bytes, _ := ioutil.ReadAll(r.Body)
	return bytes
}

func HttpContentM(r *http.Response) M {
	bytes := HttpContent(r)
	obj := M{}
	_ = json.Unmarshal(bytes, &obj)
	return obj
}

func HttpContentString(r *http.Response) string {
	bytes := HttpContent(r)
	return string(bytes)
}

func HttpGetCookieJar(url string, callType string,
	config M) (*cookiejar.Jar, error) {

	var resp *http.Response
	var errCall error
	var client *http.Client

	jar, e := cookiejar.New(nil)
	if e != nil {
		return nil, fmt.Errorf("Unable to initialize cookie jar: %s", e.Error())
	}

	client = &http.Client{
		Jar: jar,
	}

	if callType == "POST" {
		if config.Has("loginvalues") {
			fvs := config["loginvalues"].(M)
			vs := httpurl.Values{}
			for k, v := range fvs {
				vs.Set(k, v.(string))
			}

			resp, errCall = client.PostForm(url, vs)
			if errCall == nil {
				resp.Body.Close()
			}
		}
	} else {
		_, errCall = client.Get(url)
	}

	return jar, errCall
}

type builder struct {
	handler func(http.ResponseWriter, *http.Request)
}

func (b builder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.handler(w, r)
}

func ToHttpHandler(handleFunc func(http.ResponseWriter, *http.Request)) http.Handler {
	return builder{handler: handleFunc}
}
