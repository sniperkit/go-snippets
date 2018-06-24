package dcat

import (
	"encoding/json"
	//"github.com/datatogether/linked_data"
	"bytes"
	//"fmt"
	"testing"
	"time"
)

// This is an example setup for testing MarshalJSON for the dcat structs
func TestDcatMarshalJSON(t *testing.T) {

	cases := []struct {
		in  *Catalog
		out []byte
		err error
	}{
		{&Catalog{Title: "myTitle",
			Description: "myDescription",
			Issued:      time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			Modified:    time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			Language:    []string{"en", "fr", "abc"},
			Homepage:    "http://epa.gov",
			Publisher:   "epa",
			Spatial:     "abc",
			Themes:      []string{"solarized-dark", "climate", "#tbt"},
			License:     "GPL",
			Rights:      "abc",
			Dataset:     &Dataset{Title: "abc"},
			Record:      &CatalogRecord{Title: "abc", Description: "description", Issued: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC), Modified: time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC), PrimaryTopic: "abc"},
		}, []byte(`{"dct:title":"myTitle","dct:description":"myDescription","dct:issued":"2009-11-17T20:34:58.651387237Z","dct:modified":"2009-11-17T20:34:58.651387237Z","dct:language":["en","fr","abc"],"foaf:homepage":"http://epa.gov","dct:publisher":"epa","dct:spatial":"abc","dcat:themeTaxonomy":["solarized-dark","climate","#tbt"],"dct:license":"GPL","dct:rights":"abc","dcat:dataset":{"dct:title":"abc","dct:issued":"0001-01-01T00:00:00Z","dct:modified":"0001-01-01T00:00:00Z"},"dcat:record":{"dct:title":"abc","dct:description":"description","dct:issued":"2009-11-17T20:34:58.651387237Z","dct:modified":"2009-11-17T20:34:58.651387237Z","foaf:primaryTopic":"abc"}}`), nil},
		//case 1: omit empty title
		{&Catalog{Title: "",
			Description: "",
			Issued:      time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			Modified:    time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC),
			Language:    []string{""},
			Homepage:    "",
			Publisher:   "",
			Spatial:     "",
			Themes:      []string{""},
			License:     "",
			Rights:      "",
			Dataset:     nil,
			Record:      nil,
		}, []byte(`{"dct:issued":"2009-11-17T20:34:58.651387237Z","dct:modified":"2009-11-17T20:34:58.651387237Z","dct:language":[""],"dcat:themeTaxonomy":[""]}`), nil},
	}

	for i, c := range cases {
		got, err := json.Marshal(c.in)
		//fmt.Println(got)
		if err != c.err {
			t.Errorf("case %d error mismatch. expected: '%s', got: '%s'", i, c.err, err)
			continue
		}

		if !bytes.Equal(c.out, got) {
			t.Errorf("case %d error mismatch. %s != %s", i, string(c.out), string(got))
			continue
		}
	}

	// strbytes, err := json.Marshal(&Dataset{path: datastore.NewKey("/path/to/dataset")})
	// if err != nil {
	// 	t.Errorf("unexpected string marshal error: %s", err.Error())
	// 	return
	// }

	// if !bytes.Equal(strbytes, []byte("\"/path/to/dataset\"")) {
	// 	t.Errorf("marshal strbyte interface byte mismatch: %s != %s", string(strbytes), "\"/path/to/dataset\"")
	// }
}
