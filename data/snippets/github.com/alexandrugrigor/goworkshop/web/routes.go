package web

import "net/http"

const (
	authorBaseUrl   = "/authors"
	authorByUuidUrl = authorBaseUrl + "/{uuid}"
	booksBaseUrl    = "/books"
	bookByUuidUrl   = booksBaseUrl + "/{uuid}"
)

type Route struct {
	//method type (GET, POST, PUT, DELETE, etc)
	Method string

	//the method path
	Pattern string

	//the method that the endpoint should call
	HandlerFunc RouteFunc
}

type RouteFunc func(http.ResponseWriter, *http.Request) error

type Routes []Route

func (server *RestServer) initRoutes(){
	server.routes = Routes{
		//book_handlers
		Route{
			Method:      "GET",
			Pattern:     "/",
			HandlerFunc: server.Index,
		},
		Route{
			Method:      "GET",
			Pattern:     booksBaseUrl,
			HandlerFunc: server.GetAllBooks,
		},
		Route{
			Method:      "POST",
			Pattern:     booksBaseUrl,
			HandlerFunc: server.AddBook,
		},
		Route{
			Method:      "GET",
			Pattern:     bookByUuidUrl,
			HandlerFunc: server.GetBookByUUID,
		},
		Route{
			Method:      "DELETE",
			Pattern:     bookByUuidUrl,
			HandlerFunc: server.DeleteBookByUUID,
		},
		Route{
			Method:      "PUT",
			Pattern:     bookByUuidUrl,
			HandlerFunc: server.UpdateBook,
		},

		//author_handlers
		Route{
			Method:      "GET",
			Pattern:     authorBaseUrl,
			HandlerFunc: server.GetAllAuthors,
		},
		Route{
			Method:      "POST",
			Pattern:     authorBaseUrl,
			HandlerFunc: server.AddAuthor,
		},
		Route{
			Method:      "GET",
			Pattern:     authorByUuidUrl,
			HandlerFunc: server.GetAuthorByUUID,
		},
		Route{
			Method:      "DELETE",
			Pattern:     authorByUuidUrl,
			HandlerFunc: server.DeleteAuthorByUUID,
		},
		Route{
			Method:      "PUT",
			Pattern:     authorByUuidUrl,
			HandlerFunc: server.UpdateAuthor,
		},
	}
}

