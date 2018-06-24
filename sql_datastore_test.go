package sql_datastore

import (
	"database/sql"
	//datastore "github.com/ipfs/go-datastore"
	"testing"
)

func TestNewDatastore(t *testing.T) {
	cases := []struct {
		in  *sql.DB
		out *Datastore
		err error
	}{
		//case 0
		{nil, &Datastore{DB: nil}, nil},
	}

	for i, c := range cases {
		got := NewDatastore(c.in)

		// TODO: NewDatastore should also return nil/error
		// if err != c.err {
		// 	t.Errorf("case %d error mismatch. expected: '%s', got: '%s'", i, c.err, err)
		// 	continue
		// }

		if c.out.DB != got.DB {
			t.Errorf("case %d error mismatch. %s != %s", i, c.out.DB, got.DB)
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
