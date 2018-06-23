package goconfigure

import (
	"github.com/coreos/go-etcd/etcd"
	"testing"
)

func TestNestedEtcdToInterface(t *testing.T) {
	// Setup
	client := etcd.NewClient([]string{"http://127.0.0.1:4001"})
	client.CreateDir("/goconfigure", 200)
	client.CreateDir("/goconfigure/blah", 200)
	client.Create("/goconfigure/blah/asd", "1", 200)
	client.Create("/goconfigure/blah/meh", "2", 200)
	client.CreateDir("/goconfigure/blah/yoyo", 200)
	client.Create("/goconfigure/blah/yoyo/boo", "3", 200)
	client.CreateDir("/goconfigure/hay", 200)
	client.AddChild("/goconfigure/hay", "4", 200)
	client.AddChild("/goconfigure/hay", "5", 200)
	client.AddChild("/goconfigure/hay", "6", 200)
	client.Create("/goconfigure/yo", "7", 200)

	// type GoConfigure struct {
	// 	Blah struct {
	// 		Asd  string `goconfigure:"asd"`
	// 		Meh  string `goconfigure:"meh"`
	// 		Yoyo struct {
	// 			Boo string `goconfigure:"boo"`
	// 		} `goconfigure:"yoyo"`
	// 	} `goconfigure:"blah"`

	// 	Hay map[string]string `goconfigure:"hay"`
	// 	Yo  string            `goconfigure:"yo"`
	// }

	// goconfig := &GoConfigure{}

	resp, _ := client.Get("/goconfigure", true, true)

	r := nestedEtcdToMap(resp.Node)

	if r["yo"] != "7" {
		t.Fatal("yo is not set to 7")
	}

	if r["blah"].(map[string]interface{})["asd"] != "1" ||
		r["blah"].(map[string]interface{})["meh"] != "2" ||
		r["blah"].(map[string]interface{})["yoyo"].(map[string]interface{})["boo"] != "3" {
		t.Fatal("blah is not set correctly")
	}

	if len(r["hay"].(map[string]interface{})) < 3 {
		t.Fatal("hay is not set correctly")
	}

	// Clean up
	client.Delete("/goconfigure", true)
}

// struct Blah {
//   Asd string `goconfigure:"asd"`
//   Meh string `goconfigure:"meh"`
//   Hay map[string]string `goconfigure:"hay"`
// }

// struct Blah {
//   Asd string `goconfigure:"asd"`
//   Meh string `goconfigure:"-"`
//   Yoyo struct{
//     Boo string `goconfigure:"boo"`
//   } `goconfigure:"yoyo"`
// }
