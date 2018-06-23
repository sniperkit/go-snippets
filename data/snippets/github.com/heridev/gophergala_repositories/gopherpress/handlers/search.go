package handlers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/go-martini/martini"
	"github.com/gophergala/gopherpress/wpdb"
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
)

func Search(r render.Render, db *gorm.DB, args martini.Params, req *http.Request) {
	s, err := url.QueryUnescape(req.FormValue("s"))
	if err != nil || len(s) == 0 {
		r.Redirect("/")
		return
	}

	log.Println(s)
	var posts []wpdb.Post
	db.Where("(((wp_posts.post_title LIKE '%'?'%') OR (wp_posts.post_content LIKE '%'?'%')))  AND wp_posts.post_type IN ('post', 'page', 'attachment')", s, s).Find(&posts)

	r.XML(200, posts)
}
