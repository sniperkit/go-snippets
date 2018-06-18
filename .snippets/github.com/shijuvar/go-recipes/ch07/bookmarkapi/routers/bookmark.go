package routers

import (
	"github.com/gorilla/mux"

	"github.com/shijuvar/go-recipes/ch07/bookmarkapi/common"
	"github.com/shijuvar/go-recipes/ch07/bookmarkapi/controllers"
)

// SetBookmarkRoutes registers routes for bookmark entity.
func SetBookmarkRoutes(router *mux.Router) *mux.Router {
	bookmarkRouter := mux.NewRouter()
	bookmarkRouter.HandleFunc("/bookmarks", controllers.CreateBookmark).Methods("POST")
	bookmarkRouter.HandleFunc("/bookmarks/{id}", controllers.UpdateBookmark).Methods("PUT")
	bookmarkRouter.HandleFunc("/bookmarks", controllers.GetBookmarks).Methods("GET")
	bookmarkRouter.HandleFunc("/bookmarks/{id}", controllers.GetBookmarkByID).Methods("GET")
	bookmarkRouter.HandleFunc("/bookmarks/users/{id}", controllers.GetBookmarksByUser).Methods("GET")
	bookmarkRouter.HandleFunc("/bookmarks/{id}", controllers.DeleteBookmark).Methods("DELETE")
	router.PathPrefix("/bookmarks").Handler(common.AuthorizeRequest(bookmarkRouter))
	return router
}
