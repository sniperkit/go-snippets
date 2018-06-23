package urls

import (
	"net/http"
)

type HomeHandle struct{}

type TestHandle struct{}

func (t *HomeHandle) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r := Render()
	r.Text(w, http.StatusOK, "hello world")
}

func (t *TestHandle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	render := Render()
	rows, _ := DB().Query("SELECT * FROM b_story LIMIT 20")
	result := FetchAll(rows)
	render.JSON(w, http.StatusOK, result)
}
