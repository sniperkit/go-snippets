package elasticsearch

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"reflect"
	"strings"
	"time"

	elastigo "github.com/mattbaird/elastigo/lib"
	"gopkg.in/errgo.v1"
)

var (
	defaultIndex = "iloveindexes"
	defaultConn  *elastigo.Conn
)

func init() {
	initDefaultConn()
}

type Document interface {
	GetId() string
	SetId(id string)
	Type() string
}

func initDefaultConn() {
	u := os.Getenv("ELASTICSEARCH_URL")
	if u != "" {
		elasticsearchURL, err := url.Parse(u)
		if err != nil {
			log.Fatalln("Invlid URL:", u)
		}
		user, password := "", ""
		if elasticsearchURL.User != nil {
			user = elasticsearchURL.User.Username()
			password, _ = elasticsearchURL.User.Password()
		}
		splittedHost := strings.Split(elasticsearchURL.Host, ":")
		defaultConn = &elastigo.Conn{
			Protocol:       elastigo.DefaultProtocol,
			Username:       user,
			Password:       password,
			Domain:         splittedHost[0],
			ClusterDomains: []string{splittedHost[0]},
			Port:           splittedHost[1],
			DecayDuration:  time.Duration(elastigo.DefaultDecayDuration * time.Second),
		}
	} else {
		defaultConn = elastigo.NewConn()
	}
}

func index(_type string, id string, args map[string]interface{}, data interface{}) (elastigo.BaseResponse, error) {
	return defaultConn.Index(defaultIndex, _type, id, args, data)
}

func Index(d Document) error {
	res, err := index(d.Type(), d.GetId(), nil, &d)
	if err != nil {
		errgo.Mask(err)
	}
	if res.Created {
		d.SetId(res.Id)
	}
	return nil
}

func ListDocuments(_struct Document) ([]Document, error) {
	res, err := list(_struct.Type())
	if err != nil {
		return nil, errgo.Mask(err)
	}

	t := reflect.TypeOf(_struct).Elem()

	var out []Document
	for _, h := range res.Hits.Hits {
		doc := reflect.New(t).Interface().(Document)
		err := json.Unmarshal(*h.Source, &doc)
		if err != nil {
			return nil, errgo.Mask(err)
		}
		doc.SetId(h.Id)
		out = append(out, doc)
	}
	return out, nil
}

func list(_type string) (*elastigo.SearchResult, error) {
	res, err := elastigo.Search(defaultIndex).Type(_type).Size("100").Result(defaultConn)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return res, nil
}
