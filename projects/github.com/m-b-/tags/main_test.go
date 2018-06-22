package main

// Keep this for dev purpose:
// DROP TABLE docs CASCADE; DROP TABLE tagsdocs CASCADE; DROP TABLE tags CASCADE; DROP TYPE dtype;

// SELECT * FROM (SELECT docs.id, docs.name, docs.type, docs.content, string_agg(tags.name, U&'\001F') AS ctags FROM tags, tagsdocs, docs WHERE tagsdocs.idtag = tags.id AND tagsdocs.iddoc = docs.id GROUP BY docs.id) AS foo WHERE ctags LIKE '%bookmarks%physics%';
// SELECT docs.id FROM tags, tagsdocs, docs WHERE tagsdocs.idtag = tags.id AND tagsdocs.iddoc = docs.id AND tags.name IS IN ('bookmarks', 'physics') GROUP BY docs.id
// SELECT docs.id FROM tags, tagsdocs, docs WHERE tagsdocs.idtag = tags.id AND tagsdocs.iddoc = docs.id AND tags.name IN ('bookmarks', 'physics') GROUP BY docs.id HAVING COUNT(docs.id) = 2;
import (
	"encoding/json"
	"strings"
	"testing"
)

var testdb *Database

type testdocs struct {
	json		[]byte		// JSON-ified doc
	inserted	bool		// previous doc should have been inserted?
}

var docs = []testdocs {
	{
		[]byte(`{
			"Id"		:	-1,
			"Name"		:	"nice website",
			"Type"		:	"url",
			"Content"	:	"http://awesom.eu",
			"Uid"		:	1,
			"Tags"		:	["irc", "awesom", "index"]
		}`),
		true,
	},
	{
		[]byte(`{
			"Id"		:	-1,
			"Name"		:	"nice website",
			"Type"		:	"badtype",
			"Content"	:	"http://awesom.eu",
			"Uid"		:	1,
			"Tags"		:	["irc", "awesom", "index"]
		}`),
		false,	// bad type
	},
	{
		[]byte(`{
			"Id"		:	-1,
			"Name"		:	"bad tag",
			"Type"		:	"text",
			"Content"	:	"whatever",
			"Uid"		:	1,
			"Tags"		:	["bad\u001Ftag"]
		}`),
		true,
	},
	{
		[]byte(`{
			"Id"		:	-1,
			"Name"		:	"The ArXiV",
			"Type"		:	"url",
			"Content"	:	"http://arxiv.org/",
			"Uid"		:	1,
			"Tags"		:	["papers", "maths", "physics"]
		}`),
		true,
	},	
	{
		[]byte(`{
			"Id"		:	-1,
			"Name"		:	"/r/physics",
			"Type"		:	"url",
			"Content"	:	"http://www.reddit.com/r/physics/",
			"Uid"		:	1,
			"Tags"		:	["bookmarks", "physics", "/r/"]
		}`),
		true,
	},
	{
		[]byte(`{
			"Id"		:	-1,
			"Name"		:	"/r/programming",
			"Type"		:	"url",
			"Content"	:	"http://www.reddit.com/r/programming/",
			"Uid"		:	1,
			"Tags"		:	["bookmarks", "programming", "/r/"]
		}`),
		true,
	},
	{
		[]byte(`{
			"Id"		:	-1,
			"Name"		:	"Hacker News",
			"Type"		:	"url",
			"Content"	:	"https://news.ycombinator.com/",
			"Uid"		:	1,
			"Tags"		:	["bookmarks", "programming" ]
		}`),
		true,
	},
	{
		[]byte(`{
			"Id"		:	-1,
			"Name"		:	"Slashdot",
			"Type"		:	"url",
			"Content"	:	"http://beta.slashdot.org/",
			"Uid"		:	1,
			"Tags"		:	["bookmarks", "programming", "physics", "science" ]
		}`),
		true,
	},
	{
		[]byte(`{
			"Id"		:	-1,
			"Name"		:	"Slackware",
			"Type"		:	"url",
			"Content"	:	"http://www.slackware.com/",
			"Uid"		:	1,
			"Tags"		:	["bookmarks", "linux", "slackware" ]
		}`),
		true,
	},
}

// XXX for now, we only support tag1 AND tag2 AND tag3, etc.
type testqueries struct {
	query		[]string	// list of tags
	results		[]string	// Name of the docs; unique in testdocs
}

var queries = []testqueries {
	{
		[]string{ "bookmarks" },
		[]string{
			"/r/physics", "/r/programming",
			"Hacker News", "Slashdot", "Slackware",
		},
	},
	{
		[]string{ "physics" },
		[]string{ "The ArXiV","/r/physics","Slashdot" },
	},
	{
		[]string{ "physics", "bookmarks" },
		[]string{ "/r/physics", "Slashdot" },
	},
}

func TestDocs(t *testing.T) {
	// Database Creation
//	testdb = NewDB("test")
	testdb = NewDB()
	if testdb == nil {
		t.Error("Database creation failure")
	}

	var ids []int32

	// Add every test documents
	for _, doc := range docs {
		var d Doc
		err := json.Unmarshal(doc.json, &d)
		if err != nil {
			t.Log(string(doc.json))
			t.Error("Cannot retrieve doc:", err)
		}
		id := testdb.AddDoc(&d)
		if id != -1 && !doc.inserted ||
		   id == -1 && doc.inserted {
			t.Error("Wrong expectation about doc insertion")
		}

		// Fetch previously added document
		if id != -1 {
			ids = append(ids, id)
			d.Id = id
			d2 := testdb.GetDoc(id)
			if d.Name != d2.Name {
				t.Error("Bad name:", d.Name, d2.Name)
			} else if d.Type != d2.Type {
				t.Error("Bad types:", d.Type, d2.Type)
			} else if d.Content != d2.Content {
				t.Error("Bad Content:", d.Content, d2.Content)
			} else if len(d.Tags) != len(d2.Tags) {
				t.Error("Number of tags mismatch")
			}
			for _, t1 := range d.Tags {
				// Remove separator in tag
				t1 = strings.Replace(t1, TagSep, "", -1)
				found := false
				for _, t2 := range d2.Tags {
					if t2 == t1 { found = true }
				}
				if !found {
					t.Error("Tag was not added:", t1)
				}
			}
		}
	}

	// Launch every query
	for _, query := range queries {
		ds := testdb.GetDocs(1, query.query)
		if len(ds) != len(query.results) {
			t.Log(len(ds), ds)
			t.Log(len(query.results), query.results)
			t.Error("Wrong number of results")
		}
		for _, d := range ds {
			found := false
			for _, r := range query.results {
				if d.Name == r { found = true }
			}
			if !found {
				t.Error("Doc was not retrieved")
			}
		}
	}

	// Delete everything.

	// Drop everything
//	testdb.Query("DROP TABLE docs CASCADE; DROP TABLE tagsdocs CASCADE; DROP TABLE tags CASCADE; DROP TYPE dtype;")
}
