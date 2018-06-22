package dshash

import (
	"appengine"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
)

func getHandler(c appengine.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	person := &Person{}
	person.Handler = "chischaschos"
	bytes, err := person.Marshal()

	if err != nil {
		panic(err)
	}

	s := PersonsService{
		PersonRepository: PersonRepository{c},
	}

	locations, err := s.GetAll(person)

	if err != nil {
		panic(err)
	}

	person.Locations = locations

	_, err = w.Write(bytes)

	if err != nil {
		panic(err)
	}
}

func postHandler(c appengine.Context, w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	defer r.Body.Close()
	bodyBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		panic(err)
	}

	person := &Person{}

	if err := person.Unmarshal(bodyBytes); err != nil {
		panic(err)
	}

	s := PersonsService{
		PersonRepository: PersonRepository{c},
	}

	if err := s.Save(person); err != nil {
		panic(err)
	}

	bytes, err := person.Marshal()

	if err != nil {
		panic(err)
	}

	if _, err := w.Write(bytes); err != nil {
		panic(err)
	}
}
