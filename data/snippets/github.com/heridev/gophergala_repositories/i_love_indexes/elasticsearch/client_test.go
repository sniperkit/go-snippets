package elasticsearch

import (
	"log"
	"os"
	"testing"
	"time"
)

type D struct {
	Id    string `json:"_id,omitempty"`
	Value string `json:"value"`
}

func (d *D) Type() string {
	return "d"
}

func (d *D) SetId(id string) {
	d.Id = id
}

func (d *D) GetId() string {
	return d.Id
}

func TestMain(m *testing.M) {
	defaultIndex = "testing"
	ret := m.Run()
	os.Exit(ret)
}

func teardown() {
	res, err := defaultConn.DeleteIndex(defaultIndex)
	if err != nil {
		log.Fatalln(res, err)
	}
}

func TestIndex(t *testing.T) {
	d := &D{Value: "helloworld"}
	err := Index(d)
	if err != nil {
		t.Fatal("expected nil, got", err)
	}
	if d.Id == "" {
		t.Error("expected not nil ID")
	}
	teardown()
}

func TestList(t *testing.T) {
	err := Index(&D{Value: "helloworld"})
	if err != nil {
		t.Fatal("err should be nil", err)
	}
	err = Index(&D{Value: "helloworld2"})
	if err != nil {
		t.Fatal("err should be nil", err)
	}

	time.Sleep(2 * time.Second)

	ds, err := ListDocuments((*D)(nil))
	if err != nil {
		t.Fatal("err should be nil", err)
	}

	if len(ds) != 2 {
		t.Fatal("len should be 2, got", len(ds))
	}
	teardown()
}
