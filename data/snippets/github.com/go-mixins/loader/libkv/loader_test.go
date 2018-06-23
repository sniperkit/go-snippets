package libkv_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/docker/libkv/store"
	"github.com/go-test/deep"

	"github.com/go-mixins/loader/libkv"
)

type kvMock []*store.KVPair

func (kvm kvMock) Close() {}

func (kvm kvMock) Get(key string) (res *store.KVPair, err error) {
	for _, kv := range kvm {
		if kv.Key == key {
			res = kv
			return
		}
	}
	err = store.ErrKeyNotFound
	return
}

func (kvm kvMock) List(directory string) (res []*store.KVPair, err error) {
	for _, kv := range kvm {
		if strings.HasPrefix(kv.Key, directory) && kv.Key != directory {
			res = append(res, kv)
		}
	}
	return
}

func (kvm kvMock) WatchTree(directory string, stopCh <-chan struct{}) (<-chan []*store.KVPair, error) {
	return nil, nil
}

type testStruct struct {
	B struct {
		C int
	}
	D string
	E []struct {
		X float64
		Y string
	}
}

func TestLoad(t *testing.T) {
	var src, dest testStruct
	kv := kvMock{
		{Key: "a/b/c", Value: []byte("1")},
		{Key: "a/d", Value: []byte("string")},
		{Key: "a/e/1/x", Value: []byte("0.1")},
		{Key: "a/e/2/x", Value: []byte("0.2")},
		{Key: "a/e/1/y", Value: []byte("y1")},
		{Key: "a/e/2/y", Value: []byte("y2")},
	}
	json.Unmarshal([]byte(`{
			"B": {
				"C": 1
			},
			"D": "string",
			"E": [
				{
					"X": 0.1,
					"Y": "y1"
				},
				{
					"X": 0.2,
					"Y": "y2"
				}
			]
		}
	`), &src)
	loader, err := libkv.New("a", kv)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer loader.Close()
	err = loader.Load(&dest)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	if diff := deep.Equal(src, dest); diff != nil {
		t.Errorf("%+v", diff)
	}
}
