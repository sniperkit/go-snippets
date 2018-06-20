package wiki

import (
	"net/http"
	"strings"

	"github.com/gliderlabs/stdcom/web/console"
)

// ConsoleMenuItems provides web/console with menu items for this wiki
func (c *Component) ConsoleMenuItems() []console.MenuItem {
	return []console.MenuItem{
		{Title: "Wiki", Link: "/wiki"},
	}
}

// MatchHTTP implements the web/console "sub" request handler
func (c *Component) MatchHTTP(r *http.Request) bool {
	return strings.HasPrefix(r.URL.Path, "/wiki")
}

// ServeHTTP implements the web/console "sub" request handler
func (c *Component) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the page name from the path
	pageName := strings.TrimPrefix(strings.TrimSuffix(r.URL.Path, "/edit"), "/wiki")
	pageName = strings.TrimLeft(pageName, "/")

	// if no page, redirect to the Home page
	if pageName == "" {
		http.Redirect(w, r, r.URL.Path+"/Home", http.StatusTemporaryRedirect)
		return
	}

	// try to fetch the page from the store
	page, _ := c.Store.Page(pageName)

	// if we're editing we take a different path
	if strings.HasSuffix(r.URL.Path, "/edit") {
		// make an empty page if the page doesn't exist
		if page == nil {
			page = &Page{
				Title: pageName,
			}
		}
		// put a page into store based on form fields
		// then redirect to the page, if POST
		if r.Method == "POST" {
			r.ParseForm()
			page.Title = r.FormValue("title")
			page.Content = r.FormValue("content")
			page.Touch(c.Console.Auth.CurrentUser(r).NickName)
			c.Store.PutPage(pageName, page)
			http.Redirect(w, r, "/wiki/"+pageName, http.StatusFound)
			return
		}

		// if not POST, render the edit template
		c.Console.RenderTemplate(w, r, "wikiedit.jet", nil, page)
		return
	}

	// if we're not editing, if the page does not exist
	// redirect to the edit page
	if page == nil {
		http.Redirect(w, r, r.URL.Path+"/edit", http.StatusTemporaryRedirect)
		return
	}

	// otherwise render the wiki page
	c.Console.RenderTemplate(w, r, "wikipage.jet", nil, page)
}
