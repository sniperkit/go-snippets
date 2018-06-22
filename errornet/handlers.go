package errornet

import (
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/HouzuoGuo/tiedot/dberr"

	"fmt"
	"net/http"
)

type TieDot struct {
	Collection *db.Col
}

func (td *TieDot) Handle(e Error) error {
	// insert error into collection
	if _, err := td.Collection.Insert(e); err != nil {
		return err
	}

	return nil
}

func (td *TieDot) ServeHTTP(addr string) {
	http.Handle("/", func(w http.ResponseWriter, r *http.Request) {
		td.Collection.ForEachDoc(func(id int, doc []byte) {
			fmt.Fprintf(w, "id: %v:\n%#v\n\n", id, doc)
		})
	})

	http.ListenAndServe(addr, nil)
}

func NewTieDot(col *db.Col) *TieDot {
	return &TieDot{col}
}
