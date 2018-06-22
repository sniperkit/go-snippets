# golang workshop

# Part 1 - Write REST services with data loaded from files
### 1. hello world in GO
1. create main directory
1. create main.go file
1. declare package main
1. declare main function
1. import fmt package
1. make a call to fmt.Println and print "Hello world, my name is <your_name_here>"
### 2. defining the data model
1. create datamodel.go file
1. create AuthorDto struct with the following fields:
 - UUID string
 - FirstName string
 - LastName string
 - Birthday string
 - Death string
1. create BookDto struct with the following fields:
- UUID string
- Title string
- NoPages int
- ReleaseDate string
- Author AuthorDto
4. create a slice of AuthorDto which will hold all authors in the system
5. create a slice of BookDto which will hold all books in the system
### 3. add JSON marshalling
#### AuthorDto JSON mappings:
* map **AuthorDto.UUID** field to **uuid** field in the resulting JSON
* map **AuthorDto.FirstName** field to **firstName** field in the resulting JSON
* map **AuthorDto.LastName** field to **lastName** field in the resulting JSON
* map **AuthorDto.Birthday** field to **birthday** field in the resulting JSON
* map **AuthorDto.Death** field to **death** field in the resulting JSON
#### BookDto JSON mappings
* map **BookDto.UUID** field to **uuid** field in the resulting JSON
* map **BookDto.Title** field to **title** field in the resulting JSON
* map **BookDto.NoPages** field to **noPages** field in the resulting JSON
* map **BookDto.ReleaseDate** field to **releaseDate** field in the resulting JSON
* map **BookDto.Author** field to **author** field in the resulting JSON
### 4. read sample data from files
1. create model package
1. move datamodel.go in model folder
1. create importer package
1. create authors.json file in importer folder
1. add an array of authors with sample data in JSON format
1. create books.json file in importer folder
1. add an array of books with sample data in JSON format
1. create dataImporter.go file in importer folder
1. import encoding/json package
1. import io/ioutil
1. create a function named ImportAuthors that reads authors.json file and return a slice of authors
1. create a function named ImportBooks that reads the books.json file and returns a slice of books
1. print all authors to the STDOUT
1. print all books to the STDOUT
### 5. start a web server, listening on a configured port
1. create web package
1. create webserver.go
1. import net/http
1. read listening port from API_PORT environment variable, defaulting to 8000
1. define a http.Handler function that will handle the incoming HTTP requests to the web server
1. launch the web server on API_PORT port by calling **ListenAndServe**
1. if the web server returns an error, exit the application
### 6. writing a simple REST endpoint
1. change the path of the REST endpoint from /test to /books
1. configure the handler to respond only to GET methods
1. marshal the books array to JSON bytes
1. convert the JSON bytes to string
1. set the "Content-Type" response header to "application/json"
1. when /books endpoint is hit, return all books previously loaded from the JSON file
### 7. implement all REST endpoints with data loaded from files
1. install github.com/gorilla/mux package using 'go get -u' command
1. create util.go file in web package and implement a generic function that should write any array of data as JSON to the http.ResponseWriter
1. create routes.go file
1. define authorBaseUrl, authorByUuidUrl, booksBaseUrl, bookByUuidUrl constants with proper values
1. define a Route struct with Method, Pattern, ExpectedCode and HandlerFunc fields
1. define a Routes struct which should contain an array of Route structs
1. declare a variable named routes of type Route and define all the routes mappings
1. create author_handlers.go file in web package
1. create book_handlers.go file in web package
1. in each file, create HTTP handlers for:
 * get all entities
 * get entity by id
 * delete entity by id
 * add new entity
 * update entity
11. use mux.Vars(r)["..."] to extract path variables
12. in webserver.StartServer function, create a Gorilla mux router, register the routes array with the router and then start the http server with the router as handler
# Part 2 - Add unit tests for all REST endpoints and create the persistence layer
### 8. initialize database persistence
1. use go get to retrieve github.com/jinzhu/gorm package
1. use go get to retrieve github.com/jinzhu/gorm/dialects/postgres
1. create persistence package
1. create config.go file in persistence package
1. import for side effects github.com/jinzhu/gorm/dialects/postgres
1. create an InitDB function that uses gorm.Open to open a connection to a PostgreSQL db
1. gorm.Open should receive the following parameters: dialect, host, port, user, password, dbname, sslmode
1. on the DBInstance set DB().SetMaxOpenConns, LogMode, SingularTable, AutoMigrate for Book and Author entities, add foreign constraint from Book to Author and add unique constraint on book title
### 9. add persistence mappings
1. delete importer package
1. remove importer references from main.go
1. remove importer references from book_handlers.go
1. remove importer references from author_handlers.go
1. rename BookDto to Book and AuthroDto to Author
1. create an Entity struct in datamodel.go file with UUID field
1. add `gorm:"primary_key"` primary key annotation to the Entity.UUID field
1. embed Entity struct in Book struct
1. add `gorm:"ForeignKey:AuthorUUID"` foreign key annotation to the Book.Author field
1. add `sql:"type:text REFERENCES author(uuid) ON DELETE CASCADE"` cascade delete annotation to the Book.AuthorUUID field
1. embed Entity struct in the Author struct
### 10. add persistence services to retrieve data from db
1. create datastore.go file in persistence package
1. create GormDataStore struct with a field named DBInstance of a type representing a pointer to a gorm.DB
1. create DataStore interface in datastore.go with all books and authors operations
1. create a Store variable of DataStore type
1. in persistence/config.go file initialize the store variable before returning the DBInstance
1. create books_datastore.go file in persistence package
1. create authors_datastore.go file in persistence package
1. implement DataStore interface in books_datastore.go and authors_datastore.go files by creating functions with a pointer receiver to a gorm.DB type
### 11. switch from loading data from db instead of files
1. in main.go init the persistence
1. in web/book_handlers remove Books type and books variable
1. in web/author_handlers remove Authors type and authors variable
1. use methods on persistence.Store to work with Book and Author entities

#Part 3 - Add error handling and logging
### 12 write unit tests for All REST endpoints using mocks for all external dependencies
