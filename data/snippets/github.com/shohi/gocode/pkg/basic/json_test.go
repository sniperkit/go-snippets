package basic

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

type hello struct {
	A int    `json:"fieldA"`
	B string `json:"fieldB"`
}

func TestMarshal(t *testing.T) {
	bs, _ := json.Marshal(hello{10, "hello"})
	log.Println(string(bs))

	bs, _ = json.MarshalIndent(hello{10, "hello"}, "", "  ")
	log.Println(string(bs))
}

func TestUnmarshal(t *testing.T) {
	bs, _ := json.Marshal(hello{10, "hello"})
	var data hello
	err := json.Unmarshal(bs, &data)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(data)
	}

	// use raw message
	// Unmarshall will
	var dd hello
	info := json.RawMessage(`{"fieldA":100}`)
	json.Unmarshal(info, &dd)

	log.Println(dd.B == "")
}

func TestRawMessag(t *testing.T) {
	aa := json.RawMessage("hello")
	log.Println(aa)
	log.Println(string(aa))
}

func TestNestedMarshal(t *testing.T) {

	type Inner struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	type Outer struct {
		Inner
		Type string `json:"type"`
	}

	var data = Outer{
		Inner: Inner{
			Name:  "apple",
			Value: 100,
		},
		Type: "Fruit",
	}

	bs, _ := json.MarshalIndent(data, "", "  ")

	log.Println(string(bs))

	var data2 Outer

	json.Unmarshal(bs, &data2)

	log.Println(data2)
	log.Println(data2.Name)

}

func TestJsonArray(t *testing.T) {
	type Data struct {
		Value int `json:"value"`
	}
	type DataArray struct {
		DataList []Data `json:"dataList"`
	}

	data := Data{10}
	dataList := DataArray{
		DataList: []Data{data},
	}

	bs, _ := json.MarshalIndent(dataList, "", "  ")
	// log.Println(string(bs))

	var dataList2 DataArray
	json.Unmarshal(bs, &dataList2)
	log.Println(dataList2.DataList)
}

func TestJsonList(t *testing.T) {
	type Data struct {
		Index int `json:"index"`
		Value int `json:"value"`
	}

	ds := []*Data{
		&Data{Index: 1, Value: 1},
		&Data{Index: 2, Value: 2},
	}

	fn := func(d *Data) string {
		return fmt.Sprintf("index: %d, value: %d", d.Index, d.Value)
	}

	bs, _ := json.MarshalIndent(ds, "", "  ")
	log.Printf("json: %s\n", string(bs))

	var dds []*Data
	json.Unmarshal(bs, &dds)

	for _, v := range dds {
		fmt.Println(fn(v))
	}
}
