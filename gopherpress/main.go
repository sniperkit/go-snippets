package main

// built in
import (
	"html/template"
	"log"
)

// our st00f
import (
	"github.com/gophergala/gopherpress/handlers"
	"github.com/gophergala/gopherpress/wpdb"
)

// Martini
import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

// Gorm
import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	db, err := gorm.Open("mysql", "root:henkje5@/wordpress")

	if err != nil {
		log.Fatal(err)
	}

	db.LogMode(true)

	m := martini.Classic()
	m.Map(&db)

	var wpoptions []wpdb.Option
	db.Find(&wpoptions)

	m.Use(render.Renderer(render.Options{
		Funcs: []template.FuncMap{{"bloginfo": func(k string) string {
			for i, option := range wpoptions {
				if option.Name == k {
					return option.Value
				}
				i++
			}

			return ""
		}}},
	}))

	store := sessions.NewCookieStore([]byte("secret1337"))
	m.Use(sessions.Sessions("gopherpress", store))

	m.Get("/", handlers.Homepage)

	m.Get("/posts/all.json", handlers.GetAllPosts)
	m.Get("/posts/:id.json", handlers.GetPostById)
	m.Get("/posts/comments/:id.json", handlers.GetCommentsForArticle)
	m.Post("/posts/comments", handlers.AddComment)

	m.Get("/dashboard/login", handlers.DashboardLogin)
	m.Post("/dashboard/login/post", handlers.DashboardLoginPost)
	m.Get("/dashboard", handlers.Dashboard)
	m.Put("/dashboard/post", handlers.DashboardPostPut)

	m.Get("/search", handlers.Search)
	//	})

	/*m.Get("/hashtest/:hash", func(params martini.Params) string {
	                // why didn't they just use sha1 or php's built in thingie *sigh*
			pw := params["hash"]
			user := structures.User{}

			db.Model(&user).Where("user_login = 'emansom'")

			pwhash := user.Password
			var h = phpass.New(&phpass.Config{Count:8, Itoa: "./ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"})
			return strconv.FormatBool(h.Check([]byte(pw), []byte(pwhash)))
		})*/

	m.Run()
}
