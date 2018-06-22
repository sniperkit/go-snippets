package web

import (
	"net/http"
	"goworkshop/model"
	"io/ioutil"
	"encoding/json"
)

//Demonstrates the basic functionality of private and public modifiers in GO
func (server *RestServer) Index(w http.ResponseWriter, r *http.Request) error {
	helloWorkshop := struct {
		Message        string `json:"message"`
		privateMessage string `json:"privateMessage"`
		NoTagField     string `json:"-"`
	}{
		Message:        "Hello workshop!",
		privateMessage: "Message that does not appear in response :).",
		NoTagField:     "This message won't appear either",
	}
	WriteJson(w, helloWorkshop)
	return nil
}

func (server *RestServer) GetAllBooks(w http.ResponseWriter, _ *http.Request) error {
	books, err := (*server.Store).GetBooks()
	if err != nil {
		return err
	}
	WriteJson(w, books)
	return nil
}

func (server *RestServer) GetBookByUUID(w http.ResponseWriter, r *http.Request) error {
	bookUUID := ExtractUuid(r)
	if book, err := (*server.Store).GetBook(bookUUID); err != nil {
		return err
	} else {
		WriteJson(w, book)
		return nil
	}
}

func (server *RestServer) DeleteBookByUUID(w http.ResponseWriter, r *http.Request) error {
	bookUUID := ExtractUuid(r)
	if err := (*server.Store).DeleteBook(bookUUID); err != nil {
		return err
	} else {
		WriteJson(w, struct{ Message string }{Message: "Deleted"})
		return nil
	}
}

func (server *RestServer) AddBook(w http.ResponseWriter, r *http.Request) error {
	var book model.Book
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &book); err != nil {
		return err
	} else if err := (*server.Store).CreateBook(&book); err != nil {
		return err
	} else {
		WriteJson(w, book)
		return nil
	}
}

func (server *RestServer) UpdateBook(w http.ResponseWriter, r *http.Request) error {
	var book model.Book
	bookUUID := ExtractUuid(r)
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &book); err != nil {
		return err
	}

	if err := (*server.Store).UpdateBook(bookUUID, &book); err != nil {
		return err
	} else {
		WriteJson(w, book)
	}
	return nil
}
