package model

import (
	"errors"
	"reflect"
	"strings"
	"time"
)

func Call(name string, params ...interface{}) (result []reflect.Value, err error) {
	m := map[string]interface{}{}
	m[TypeCategory] = insertCategory
	m[TypeTag] = insertTag
	m[TypeAuthor] = insertAuthor
	m[TypeUrl] = insertUrl

	f := reflect.ValueOf(m[name])
	paramLen := f.Type().NumIn()
	if len(params) < paramLen {
		err = errors.New("The number of params is not adapted.")
		return
	}
	in := make([]reflect.Value, paramLen)
	for k, param := range params {
		if k >= paramLen {
			break
		}
		in[k] = reflect.ValueOf(param)
	}
	result = f.Call(in)
	return
}

func insertCategory(source Source, types Types, name string, url string) int {
	m := Category{
		Name:     name,
		State:    SwitchOn,
		CreateAt: time.Now().Unix(),
	}
	_, err := o.InsertOrUpdate(&m, "name", "create_at = excluded.create_at")
	if err != nil {
		panic(err)
	}
	return m.Id
}

func insertTag(source Source, types Types, name string, url string) int {
	m := Tag{
		Name:     name,
		State:    SwitchOn,
		CreateAt: time.Now().Unix(),
	}
	_, err := o.InsertOrUpdate(&m, "name", "create_at = excluded.create_at")
	if err != nil {
		panic(err)
	}
	return m.Id
}

func insertAuthor(source Source, types Types, name string, url string) int {
	m := Author{
		Name:     name,
		State:    SwitchOn,
		CreateAt: time.Now().Unix(),
	}
	_, err := o.InsertOrUpdate(&m, "name", "create_at = excluded.create_at")
	if err != nil {
		panic(err)
	}
	return m.Id
}

func insertUrl(source Source, types Types, name string, url string) int {
	url = strings.TrimSpace(url)
	if url == "" {
		return 0
	}

	basePrefixHas := strings.HasPrefix(url, "http://")
	sslPrefixHas := strings.HasPrefix(url, "https://")
	if !basePrefixHas && !sslPrefixHas {
		url = source.Url + url
	}

	m := Url{
		SourceId: source.Id,
		Url:      url,
		State:    SwitchOff,
		CreateAt: time.Now().Unix(),
	}

	if qs := o.QueryTable(&Url{}).Filter("url", url).Filter("source_id", source.Id); qs.Exist() {
		if err := qs.One(&m); err != nil {
			panic(err)
		}
		return m.Id
	}
	_, err := o.InsertOrUpdate(&m, "url", "create_at = excluded.create_at")
	if err != nil {
		panic(err)
	}
	return m.Id
}
