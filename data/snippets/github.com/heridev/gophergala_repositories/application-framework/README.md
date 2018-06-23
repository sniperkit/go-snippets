##Application Framework

This project is a full modular web application framework.
Useful when we want to put many specific tasks or services in one application with authentication.

The modularity is based on some go specific characteristics and modules are plugable at compile time.

Everithing but main.go is a module and have the same structure. You can add or remove any of mod_*.go file and program compile and run flawless (wow!).

Basically, application is a puzzle of modules linked after a rule to give a speed in developing applications.

Note that included modules do not some useful things but show some techniques to write
real things and interact with framework and other modules.

![8](https://cloud.githubusercontent.com/assets/6298396/5888340/ae14f154-a403-11e4-9c65-a0ab748f6d6c.png)

[![last-version-blue](https://cloud.githubusercontent.com/assets/6298396/5602522/8967405e-935b-11e4-8777-de3623ed6ad7.png)] (https://github.com/gophergala/application-framework/archive/master.zip)

**How to use**

Compile program. Se here (https://golang.org/doc/code.html) how.

Run and open http://localhost:8080 in your favorite browser (default user is george without password).

Back button is disabled in browser and is nice to run with Google Chrome in app mode

         google-chrome --app=http://localhost:8080
		
Note that using a HTML5 browser like Chrome you have automatic forms validation.

**Tools used in this project**

   * compiler http://golang.org
   * ide      https://github.com/visualfc/liteide
   * gopei shell    https://golang.org/geosoft1/tools for faster development

**How it works**

Basicaly, you have a module template (see [mod_template.go](https://github.com/gophergala/application-framework/blob/master/mod_template.go)) and a technique to plugin or plugout into main application.

init() function make go module plugable. Here we put only the web handler

         http.HandleFunc("/ModuleName", ModuleName)
	
handler structure

         func ModuleName(w http.ResponseWriter, r *http.Request) {
         	//this must add at begin of every session code
         	c, err := r.Cookie("session")
         	if err != nil || c.Value == "" {
         		http.Error(w, "Session expired", 401)
         		return
         	}
         
         	//build page content
         	b := `<pre> Page content`
         	
         	//finally show the page
         	p := Page{
         		Title:    "Module Title",
         		Status:   c.Value,		// e.g connected user
         		Body:     template.HTML(b),
         	}
         	t.ExecuteTemplate(w, "index.html", p)
         }

as you see the structure are fixed

   * cookie checker
   * build page content in b variable
   * show page with go template

Note that a simple cookie mecanism are used to implement sessions in modules (see [mod_login.go](https://github.com/gophergala/application-framework/blob/master/mod_login.go)).
Also, module has access to a global logfile with

         log.Println("message")

Thats all folks about modules. Now, adding a menu line in [templates/index.html](https://github.com/gophergala/application-framework/blob/master/templates/index.html) file like

         <a href="/ModuleName" >ModuleName</a> 

make visible ModuleName to application. Remove this line and coresponding module,recompile application and module are removed. Nothing else to do.

Database used is sqlite (see https://github.com/mattn/go-sqlite3)

I used a nice preformatted text to make templates more readable (see [templates/style.html](https://github.com/gophergala/application-framework/blob/master/templates/style.html))

**Included modules**

Minimal modules list (but compile even are removed too):
- mod_login autenticate user,init session, set cookie, launch mod_index
- mod_index make main page and show menus

For demo puposes i put some simple modules:
- mod_checkUser show forms using,search a user in database
- mod_showPersons show some persons informations from a database
- mod_addPerson add a person into a database and call showPerson module
- mod_ajax show integration of a simple ajax module who show a clock

**Enjoy!**
