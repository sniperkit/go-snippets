package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/gophergala/vebomm/core"
	"github.com/phaikawl/wade/app"
	wadeserv "github.com/phaikawl/wade/platform/serverside"

	"github.com/gophergala/vebomm/client"
)

const (
	ServersidePrerender = false
	ConfigFilePath      = "data/config.toml"
	IndexFile           = "public/index.html"
)

type Config struct {
	Db DbConfig `toml:"database"`
}

func main() {
	println("Starting...")
	conf := &Config{}

	_, err := toml.DecodeFile(ConfigFilePath, conf)
	if err != nil {
		fmt.Println("Cannot load config file. " + err.Error())
	}
	InitDb(conf.Db)

	r := gin.Default()
	gopath := os.Getenv("GOPATH")
	if gopath != "" {
		r.ServeFiles("/gopath/*filepath", http.Dir(gopath))
	}
	
	go h.run()

	staticServ := gzipStatic{http.StripPrefix("/public", http.FileServer(http.Dir("public/")))}
	r.GET("/public/*filepath", staticServ.serve)
	r.GET("/handshake/:id", func(c *gin.Context) {
		idStr := c.Params.ByName("id")
		var id int64
		n, err := fmt.Sscan(idStr, &id)
		if n == 0 || err != nil {
			c.Abort(http.StatusUnauthorized)
			return
		}

		var user core.User
		err = Db().SelectOne(&user, "SELECT * FROM users WHERE id=$1", id)
		if err == sql.ErrNoRows {
			c.Abort(http.StatusUnauthorized)
			return
		}

		if err != nil {
			panic(err)
		}

		wsHandler(&user, c.Writer, c.Request)
	})

	r.GET("/web/*sub", func(c *gin.Context) {
		if !ServersidePrerender {
			c.File(IndexFile)
		} else {
			indexBytes, err := ioutil.ReadFile(IndexFile)
			if err != nil {
				panic(err)
			}
			httpBkn := wadeserv.NewHttpBackend(r, c.Request, "/api")
			app := wadeserv.NewApp(app.Config{BasePath: "/web"},
				bytes.NewReader(indexBytes), c.Request.URL.Path, httpBkn)
			err = wadeserv.StartRender(app, client.AppMain{Application: app}, c.Writer)

			if err != nil {
				panic(err)
			}
		}
	})

	api(r.Group("/api/"))
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	r.Run(":" + port)
}
