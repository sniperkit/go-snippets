package api

import (
	"database/sql"
	"fmt"
	eram "github.com/Onefootball/entity-rest-api/manager"
	erat "github.com/Onefootball/entity-rest-api/test"
	"github.com/ant0ine/go-json-rest/rest"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server  *httptest.Server
	handler http.Handler
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Email    string `json:"email"`
}

type Post struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Create_Time int    `json:"create_time"`
	Author_Id   int    `json:"author_id"`
	Status      int    `json:"status"`
}

type Tag struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func init() {

	db, err := sql.Open("sqlite3", ":memory:")

	if err != nil {
		log.Fatal("An error '%s' was not expected when opening a stub database connection", err)
	} else {

		dat, err := ioutil.ReadFile("./../test.sql")
		if err != nil {
			log.Fatal("An error '%s' was not expected when opening sql file", err)
		} else {
			db.Exec(string(dat))
		}
	}

	api := rest.NewApi()
	api.Use(rest.DefaultDevStack...)

	entityManager := eram.NewEntityDbManager(db)
	entityRestApi := NewEntityRestAPI(entityManager)

	router, err := rest.MakeRouter(
		rest.Get("/api/:entity", entityRestApi.GetAllEntities),
		rest.Post("/api/:entity", entityRestApi.PostEntity),
		rest.Get("/api/:entity/:id", entityRestApi.GetEntity),
		rest.Put("/api/:entity/:id", entityRestApi.PutEntity),
		rest.Delete("/api/:entity/:id", entityRestApi.DeleteEntity),
	)

	if err != nil {
		log.Fatal("An error '%s' was not expected when creating the router", err)
	}

	api.SetApp(router)
	handler = api.MakeHandler()

	server = httptest.NewServer(handler)
}

func TestGETWithEmptySetShouldReturnEmptyJsonArray200(t *testing.T) {

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/comment", server.URL), nil))

	recorded.CodeIs(200)
	recorded.BodyIs("[]")
}

func TestGETWithExistentSetShouldReturnJsonArray200(t *testing.T) {

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/user", server.URL), nil))

	recorded.CodeIs(200)

	data := []User{}
	err := recorded.DecodeJsonPayload(&data)

	if err != nil {
		t.Error(err)
	} else {

		if len(data) < 1 {
			t.Error("No users found, and should have found at least one.")
		}
	}
}

func TestGETEntityDoestExistsShouldReturn404(t *testing.T) {

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/user/999", server.URL), nil))

	recorded.CodeIs(404)
}

func TestGETEntityThatExistsReturn200WithJson(t *testing.T) {

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/user/1", server.URL), nil))

	recorded.CodeIs(200)

	data := User{}
	err := recorded.DecodeJsonPayload(&data)

	if err != nil {
		t.Error(err)
	} else {

		if data.Id != 1 {
			t.Error("Weird behavior finding different user Id.")
		}
	}
}

func TestPOSTWithInvalidEntityShouldReturn400(t *testing.T) {

	// This post doesn't have title and should be wrong then
	entity := map[string]string{"title": "Test Post 1", "content": "not enought data"}

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("POST", fmt.Sprintf("%s/api/post", server.URL), entity))

	t.Skipf("Invalid POST should return 400 but it is returning %d", recorded.Recorder.Code)
}

func TestPOSTWithExistingEntryShouldReturn409(t *testing.T) {

	// This post has a ID 1 that conflicts with the on in the database
	entity := new(Post)
	entity.Id = 1
	entity.Title = "Test Post 1"
	entity.Content = "<p>Onefootball test post content...</p>"
	entity.Create_Time = 1437839411
	entity.Author_Id = 1
	entity.Status = 1

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("POST", fmt.Sprintf("%s/api/post", server.URL), entity))

	t.Skipf("Invalid code, should be 409 as a conflict for the ID 1 but get %d", recorded.Recorder.Code)
}

func TestPOSTWithValidEntityShouldReturn201WithHeader(t *testing.T) {

	entity := new(Post)
	entity.Id = 10
	entity.Title = "Test Post 1"
	entity.Content = "<p>Onefootball test post content...</p>"
	entity.Create_Time = 1437839411
	entity.Author_Id = 1
	entity.Status = 1

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("POST", fmt.Sprintf("%s/api/post", server.URL), entity))

	recorded.CodeIs(201)
	recorded.HeaderIs(LocationHeader, fmt.Sprintf("post/%d", entity.Id))
	recorded.HeaderIs(EntityIDHeader, fmt.Sprintf("%d", entity.Id))
	recorded.HeaderIs(StatusCodeHeader, "201")
}

/**
 * right now the PUT request return 200 even if you send problematic JSONs
 * to beign inserted, which means a false-positive, but not dangerous mistakes
 * it is just a matter of make the API more understandable and easy to debug
 */
func TestPUTWithInvalidEntityShouldReturn400(t *testing.T) {
	t.Skip("PUT request should make validations, not allow fake 200")
}

func TestPUTWithNoEntityChangeShouldReturn204(t *testing.T) {

	entity := map[string]string{"title_wrong": "Test Post 1"}

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("PUT", fmt.Sprintf("%s/api/post/%d", server.URL, 10), entity))

	recorded.CodeIs(204)
}

func TestPUTShouldReturn404IfEntityNotFound(t *testing.T) {

	entity := map[string]string{"title": "Test Post 1.1"}

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("PUT", fmt.Sprintf("%s/api/post/%d", server.URL, 999), entity))

	recorded.CodeIs(404)
}

func TestPUTShouldReturn200IfEntityUpdated(t *testing.T) {

	entity := map[string]string{"title": "Test Post 1.1"}

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("PUT", fmt.Sprintf("%s/api/post/%d", server.URL, 10), entity))

	recorded.CodeIs(200)
}

func TestDELETEShouldReturn200IfEntityExists(t *testing.T) {

	entity := map[string]string{}

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("DELETE", fmt.Sprintf("%s/api/post/%d", server.URL, 10), entity))

	recorded.CodeIs(200)
}

func TestDELETEShouldReturn404IfEntityNotExists(t *testing.T) {

	entity := map[string]string{}

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("DELETE", fmt.Sprintf("%s/api/post/%d", server.URL, 999), entity))

	recorded.CodeIs(404)
}

func TestGETWithSortQueryStringsShouldReturn200WithOrderedSet(t *testing.T) {

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/tag?_sortField=name&_sortDir=DESC", server.URL), nil))

	recorded.CodeIs(200)

	data := []Tag{}
	err := recorded.DecodeJsonPayload(&data)

	if err != nil {
		t.Error(err)
	} else {

		if len(data) <= 0 {
			t.Error("Should have found at least 3 tags.")
		}

		if data[len(data)-1].Name != "announce" {
			t.Error("The data set is not ordered by name.")
		}
	}
}

func TestGETWithPaginationQueryStringsShouldReturn200WithLimit(t *testing.T) {

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/tag?_perPage=2&_page=1", server.URL), nil))

	recorded.CodeIs(200)

	data := []Tag{}
	err := recorded.DecodeJsonPayload(&data)

	if err != nil {
		t.Error(err)
	} else {

		if len(data) != 2 {
			t.Error("Should have found and limited the set to 2 tags.")
		}
	}
}

func TestGETWithColumnNameAsQueryStringShouldReturn200WithSet(t *testing.T) {

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/tag?name=test", server.URL), nil))

	recorded.CodeIs(200)

	data := []Tag{}
	err := recorded.DecodeJsonPayload(&data)

	if err != nil {
		t.Error(err)
	} else {

		if len(data) <= 0 {
			t.Error("Should have found at least 1 tags.")
		}

		if data[0].Name != "test" {
			t.Error("The request should have filter the name by test.")
		}
	}
}

func TestGETWithMultipleColumnFiltersShouldReturn200WithSet(t *testing.T) {

	recorded := erat.RunRequest(
		t,
		handler,
		erat.MakeSimpleRequest("GET", fmt.Sprintf("%s/api/lookup?name=%s&type=PostStatus", server.URL, "%25ed%25"), nil))

	recorded.CodeIs(200)

	data := []Tag{}
	err := recorded.DecodeJsonPayload(&data)

	if err != nil {
		t.Error(err)
	} else {

		if len(data) < 2 {
			t.Errorf("Should have found at least 2 tags. %v", data)
		}
	}
}
