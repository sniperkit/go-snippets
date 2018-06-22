package handlers

import (
	"github.com/gophergala/gopherpress/wpdb"
)

import (
	"github.com/jinzhu/gorm"
	"github.com/martini-contrib/render"
)

type page struct {
	PageTitle string
	BodyClass string
	Articles  *[]wpdb.Post
}

func Homepage(r render.Render, db *gorm.DB) {
	var posts []wpdb.Post
	db.Where("post_type='post' AND post_status='publish'").Find(&posts)

	r.HTML(200, "home", page{
		PageTitle: "Home - Gopherpress",
		BodyClass: "",
		Articles:  &posts,
	}, render.HTMLOptions{
		Layout: "_wpbase",
	})
}
