package web

import (
	"net/http"
	"goworkshop/model"
	"io/ioutil"
	"encoding/json"
)

func (server *RestServer) GetAllAuthors(w http.ResponseWriter, _ *http.Request) error {
	authors, err := (*server.Store).GetAuthors()
	if err != nil {
		return err
	}
	WriteJson(w, authors)
	return nil
}

func (server *RestServer) GetAuthorByUUID(w http.ResponseWriter, r *http.Request) error {
	authorUUID := ExtractUuid(r)
	if author, err := (*server.Store).GetAuthor(authorUUID); err != nil {
		return err
	} else {
		WriteJson(w, author)
		return nil
	}
}

func (server *RestServer) DeleteAuthorByUUID(w http.ResponseWriter, r *http.Request) error {
	authorUUID := ExtractUuid(r)
	if err := (*server.Store).DeleteAuthor(authorUUID); err != nil {
		return err
	} else {
		WriteJson(w, struct{ Message string }{Message: "Deleted"})
		return nil
	}
}

func (server *RestServer) AddAuthor(w http.ResponseWriter, r *http.Request) error {
	var author model.Author
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &author); err != nil {
		return err
	} else if err := (*server.Store).CreateAuthor(&author); err != nil {
		return err
	} else {
		WriteJson(w, author)
		return nil
	}

}

func (server *RestServer) UpdateAuthor(w http.ResponseWriter, r *http.Request) error {
	authorUUID := ExtractUuid(r)
	var author model.Author
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &author); err != nil {
		return err
	}

	if err := (*server.Store).UpdateAuthor(authorUUID, &author); err != nil {
		return err
	} else {
		WriteJson(w, author)
	}
	return nil
}
