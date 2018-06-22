package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"
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

func checkSession(s sessions.Session, db *gorm.DB) (wpdb.User, bool) {
	// Check if the user is logged in.
	if s.Get("login") == nil {
		return wpdb.User{}, false
	}

	uid, ok := s.Get("login").(int64)
	if ok == false {
		s.Delete("login")
		return wpdb.User{}, false
	}

	var user wpdb.User
	db.Where("ID=?", uid).Find(&user)

	if user.Id == 0 {
		s.Delete("login")
		return wpdb.User{}, false
	}

	return user, true
}

func DashboardLogin(r render.Render, s sessions.Session, db *gorm.DB) {
	if _, ok := checkSession(s, db); ok == true {
		r.Redirect("/dashboard")
		return
	}

	type page struct {
		Title string
		Login string
	}

	// 'incorrect-login' holds a value for the login input field.
	// It is set before redirecting to this handler when invalid user/pass is tried in /dashboard/login/post.
	var login string
	if tmp := s.Get("incorrect-login"); tmp != nil {
		login = tmp.(string)
		s.Delete("incorrect-login")
	}

	r.HTML(200, "dashboard-login", page{
		Title: "GopherPress",
		Login: login,
	})
}

func DashboardLoginPost(r render.Render, s sessions.Session, db *gorm.DB, args martini.Params, req *http.Request) {
	if _, ok := checkSession(s, db); ok == true {
		r.Redirect("/dashboard")
		return
	}

	login := req.FormValue("login")

	if login == "" {
		r.Redirect("/dashboard/login")
		return
	}

	var user wpdb.User
	db.Where("user_login=?", login).Find(&user)

	if user.Id == 0 {
		s.Set("incorrect-login", login)
		r.Redirect("/dashboard/login")
		return
	}

	s.Set("login", user.Id)
	r.Redirect("/dashboard")

}

func Dashboard(r render.Render, s sessions.Session, db *gorm.DB) {
	user, ok := checkSession(s, db)
	if ok == false {
		r.Redirect("/dashboard/login")
		return
	}

	type page struct {
		Title string
		User  wpdb.User
	}

	r.HTML(200, "dashboard", page{
		Title: "GopherPress",
		User:  user,
	})
}

// Update a post, or create one if id is 0
// /dashboard/post
func DashboardPostPut(r render.Render, db *gorm.DB, s sessions.Session, req *http.Request, args martini.Params) {
	user, ok := checkSession(s, db)
	if ok == false {
		r.JSON(403, map[string]string{"error": "not logged in"})
		return
	}

	id, err := strconv.Atoi(req.FormValue("id"))
	if err != nil {
		r.JSON(400, map[string]string{"error": "invalid id"})
		return
	}
	title := req.FormValue("title")
	content := req.FormValue("content")

	if id > 0 {
		var post wpdb.Post
		db.Where("ID=?", id).Find(&post)
		if post.Id == 0 {
			r.JSON(400, map[string]string{"error": "post not found"})
			return
		}

		// TO DO: Revisioning stuff for last edited by
		post.Title = title
		post.Content = content
		db.Save(&post)
		return
	}

	post := wpdb.Post{
		UserId:   user.Id,
		PostDate: time.Now().Format("2006-01-02 15:04:05"),
		Content:  content,
		Title:    title,
		Name:     strings.ToLower(strings.Replace(title, " ", "-", -1)),
	}

	db.Save(&post)
	r.JSON(200, map[string]int64{"postid": post.Id})
}

/*
INSERT INTO `wp_posts`
	(`post_author`,`post_date`,`post_date_gmt`,`post_content`,`post_content_filtered`,`post_title`,`post_excerpt`,`post_status`,`post_type`,`comment_status`,`ping_status`,`post_password`,`post_name`,`to_ping`,`pinged`,`post_modified`,`post_modified_gmt`,`post_parent`,`menu_order`,`guid`)
	VALUES (1,'2015-01-25 14:48:57','2015-01-25 14:48:57','POST CONTENT','','deathsknot09','','inherit','revision','open','open','','11-revision-v1','','','2015-01-25 14:48:57','2015-01-25 14:48:57',11,0,'')
*/
