package handlers

import (
	"log"
	"net/http"
	"strconv"
)

import (
	"github.com/gophergala/gopherpress/wpdb"
)

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

import (
	"github.com/jinzhu/gorm"
)

func GetAllPosts(r render.Render, db *gorm.DB) {
	var posts []wpdb.Post
	db.Where("post_type='post' AND post_status='publish'").Find(&posts)

	// Hacky workaround.
	// We do not want to give the Content of all posts to the client.
	var posts2 []wpdb.Post
	for _, p := range posts {
		p.Content = ""
		posts2 = append(posts2, p)
	}

	r.JSON(200, posts2)
}

func GetPostById(r render.Render, args martini.Params, db *gorm.DB) {
	id, err := strconv.Atoi(args["id"])
	if err != nil {
		// TO DO: Give an error
		return
	}

	var post wpdb.Post
	db.Where("id=?", id).Find(&post)
	if post.Id == 0 {
		// TO DO: Give an error
		return
	}
	r.JSON(200, post)
}

func GetCommentsForArticle(r render.Render, args martini.Params, db *gorm.DB) {
	id, err := strconv.Atoi(args["id"])
	if err != nil {
		// TO DO: Give an error
		return
	}

	log.Println("I am here")

	var cmts []wpdb.Comment
	db.Where("comment_post_ID=?", id).Find(&cmts) //  AND (comment_approved = '1' OR (user_id = 1 AND comment_approved = '0' ) )

	r.JSON(200, cmts)
}

func AddComment(r render.Render, s sessions.Session, db *gorm.DB, req *http.Request) {
	var user wpdb.User
	if s.Get("login") != nil {
		db.Where("ID=?", s.Get("login").(int)).Find(&user)
	}

	id, err := strconv.Atoi(req.FormValue("comment_post_ID"))
	if err != nil {
		r.JSON(400, map[string]string{"error": "invalid id"})
		return
	}

	comment := req.FormValue("comment")
	author := req.FormValue("author")
	email := req.FormValue("email")
	url := req.FormValue("url")

	if len(comment) == 0 || ((len(author) == 0 || len(email) == 0) && user.Id != 0) {
		r.JSON(400, map[string]string{"error": "please fill in all required fields"})
		return
	}

	cmt := wpdb.Comment{
		Author:         author,
		AuthorMail:     email,
		AuthorHomepage: url,
		Content:        comment,
		UserId:         user.Id,
		PostId:         int64(id),
	}

	db.Save(&cmt)
	r.JSON(200, map[string]int64{"id": cmt.Id})
}
