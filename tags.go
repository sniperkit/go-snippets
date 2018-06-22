package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	port = flag.String("port", "8082", "Listening HTTP port")
	ssl = flag.Bool("ssl", true, "Use SSL")

	db *Database
	loginForm []byte

	ltmpl = template.Must(
		template.New("login.html").ParseFiles("templates/login.html"))

	utmpl = template.Must(
		template.New("user.html").Funcs(template.FuncMap{
			"GetTags": func(tags []string) string {
				return strings.Join(tags, ", ")
			},
			// from rsc'eq
			"eq": func(x, y interface{}) bool {
				switch x := x.(type) {
				case string, int, int64, int32:
					return x == y
				}
				return false
			},
			// XXX allows comment to be inserted on subsequent lines.
			"GetURL" : func(url string) string {
				return strings.SplitN(url, "\n", 2)[0]
			},
			"GetComment" : func(url string) string {
				ret := strings.SplitN(url, "\n", 2)
				if len(ret) == 2 {
					return ret[1]
				} else {
					return ""
				}
			},
		}).ParseFiles("templates/user.html"))

	ntmpl = template.Must(
		template.New("navbar.html").ParseFiles("templates/navbar.html"))
)

func splitTags(tags string) []string {
	return strings.FieldsFunc(tags, func(r rune) bool {
		switch r {
		case []rune(TagSep)[0], ' ', ',', '\n', '\t':
			return true
		}
		return false
	})
}

func getType(c string) string {
	// XXX may regexp
	// "^((f|ht)tps?:// | [A-Z]:\\ | /)"
	switch {
	case strings.HasPrefix(c, "http://"),
		strings.HasPrefix(c, "https?://"),
		strings.HasPrefix(c, "ftp://"),
		strings.HasPrefix(c, "ftps://"),
		strings.HasPrefix(c, "C:"),
		strings.HasPrefix(c, "/"):
		return "url"
	default:
		return "text"
	}
}

func index(w http.ResponseWriter, r *http.Request, _ int32) {
	writeFiles(w, "templates/index.html")
}

func login(w http.ResponseWriter, r *http.Request, _ int32) {
	switch r.Method {
	case "GET":
		w.Write(loginForm)
	case "POST":
		// login may be email, name or token
		login := r.FormValue("login")

		resp, err := login2(login)
		if err != nil {
			LogHttp(w, err)
			return
		}

		switch resp {
		// received a valid token, effectively login the user
		case "ok":
			// fetch id from server
			udata, err := info(login)
			if err != nil {
				LogHttp(w, err)
				return
			}
			uid, _ := strconv.ParseInt(strings.Split(udata, "\n")[0], 10, 32)

			// generate a new token
			token, _ := chain(login)
			setToken(w, token, int32(uid))

			// everything went well, redirect
			http.Redirect(w, r, "/user/", http.StatusFound)
		// wrong data.
		case "ko":
			SetError(w, errors.New("Wrong token/email/name"))
			http.Redirect(w, r, "/login", http.StatusFound)
		// new token has been generated
		case "new":
			SetInfo(w, "Check your AAS account!")
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request, _ int32) {
	// XXX can safely getToken() here as
	// logout is not in mustauth
	if token, _, err := getToken(r); err == nil {
		logout2(token)
		unsetToken(w)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func user(w http.ResponseWriter, r *http.Request, uid int32) {
	// fetch docs
	docs := db.GetDocs(uid, splitTags(r.FormValue("search")))

	d := struct {
		Empty Doc
		Docs  []Doc
		Uid   int32
	}{Doc{uid, "", "", "Some content", -1, []string{""}}, docs, uid}

	if err := utmpl.Execute(w, &d); err != nil {
		LogHttp(w, err)
		return
	}
}

func add(w http.ResponseWriter, r *http.Request, uid int32) {
	name := r.FormValue("name")
	tags := splitTags(r.FormValue("tags"))
	// XXX element not added (can't be retrieved)
	if len(tags) == 0 {
		SetError(w, errors.New("At least one tag is required"))
		http.Redirect(w, r, "/user/", http.StatusFound)
		return
	}

	// XXX  remove  html content.
	// (it should have been removed by js, but...)
	content := strings.TrimSpace(r.FormValue("content"))
	typ := getType(content)

	id := db.AddDoc(&Doc{-1, name, typ, content, uid, tags})
	if id == -1 {
		SetError(w, errors.New("Can't add that (weird)"))
		http.Redirect(w, r, "/user/", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/user/", http.StatusFound)
}

// edit document (delete too)
func edit(w http.ResponseWriter, r *http.Request, uid int32) {
	i, _ := strconv.ParseInt(r.FormValue("id"), 10, 32)
	id := int32(i)

	if !db.HasOwner(id, uid) {
		SetError(w, errors.New("You don't own this."))
		http.Redirect(w, r, "/user/", http.StatusFound)
		return
	}

	log.Println(strings.TrimSpace(r.FormValue("content")))

	switch r.FormValue("action") {
	case "edit":
		name := r.FormValue("name")
		tags := splitTags(r.FormValue("tags"))
		if len(tags) == 0 {
			return
		}
		content := strings.TrimSpace(r.FormValue("content"))
		typ := getType(content)
		db.UpdateDoc(&Doc{id, name, typ, content, uid, tags})
	case "delete":
		db.DelDoc(id)
	}

	http.Redirect(w, r, "/user/", http.StatusFound)
}

var tagsfuncs = map[string]func(http.ResponseWriter, *http.Request, int32){
	"":       index,
	"login":  login,
	"logout": logout,
	"user":   user,
	"add":    add,
	"edit":   edit,
//	"settings"	:	settings,
}

var mustauth = map[string]bool{
	"user": true,
	"add":  true,
	"edit": true,
//	"settings"	:	true,
}

func tags(w http.ResponseWriter, r *http.Request) {
	var uid int32

	f := r.URL.Path[1:]
	if len(f) != 0 && f[len(f)-1] == '/' {
		f = f[:len(f)-1]
	}

	if tagsfuncs[f] == nil {
		http.NotFound(w, r)
		return
	}

	// pages requiring to be connected
	if mustauth[f] {
		var err error
		uid, err = ChainToken(w, r)
		if err != nil {
			SetError(w, errors.New("Invalid token"))
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
	} else {
		// XXX try getting uid nevertheless (navbar display...)
		_, uid, _ = getToken(r)
	}

	if (r.Method == "GET" && f != "logout") || f == "user" {
		writeFiles(w, "templates/header.html")
		d := struct{ Connected bool }{Connected: uid > 0}
		if err := ntmpl.Execute(w, &d); err != nil {
			log.Println(err)
		}
	}

	tagsfuncs[f](w, r, uid)

	if (r.Method == "GET" && f != "logout") || f == "user" {
		writeFiles(w, "templates/footer.html")
	}
}

func main() {
	var err error

	loginForm, err = ioutil.ReadFile("templates/login.html")
	if err != nil {
		log.Fatal(err)
	}

	// load auth certificate
	pem, err := ioutil.ReadFile("auth-cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	
	certs := x509.NewCertPool()
	if !certs.AppendCertsFromPEM(pem) {
		log.Fatal(errors.New("can't add auth-cert.pem"))
	}

	// create Client for auth requests
	Client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: certs,
			},
		},
	}

	// Load Database
	db = NewDB()

	http.HandleFunc("/", tags)

	// TODO automatically tag :bookmark for type=url

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static"))))

	if *ssl {
		log.Print("Launching on https://localhost:" + *port)
		log.Fatal(http.ListenAndServeTLS(":"+*port, "cert.pem", "key.pem", nil))
	} else {
		log.Print("Launching on http://localhost:" + *port)
		log.Fatal(http.ListenAndServe(":"+*port, nil))
	}
}
