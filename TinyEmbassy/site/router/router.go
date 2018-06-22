/*
* @Author: souravray
* @Date:   2015-01-24 10:47:03
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-25 23:41:14
 */

package router

import (
	"github.com/gophergala/tinyembassy/site/controllers"
	"github.com/gorilla/mux"
	"net/http"
)

func Routes(rtr *mux.Router) {
	//static assets
	rtr.PathPrefix("/script/").Handler(http.StripPrefix("/script/", http.FileServer(http.Dir("./static/script"))))
	rtr.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img"))))
	rtr.PathPrefix("/style/").Handler(http.StripPrefix("/style/", http.FileServer(http.Dir("./static/style"))))
	rtr.PathPrefix("/template/").Handler(http.StripPrefix("/template/", http.FileServer(http.Dir("./static/template"))))

	//basic navigations
	rtr.HandleFunc("/", controllers.LandingPage).Methods("GET").Name("Homepage")
	rtr.HandleFunc("/signin", controllers.LoginPage).Methods("GET").Name("LoginPage")
	rtr.HandleFunc("/signin", controllers.Login).Methods("POST").Name("LoginRequest")
	rtr.HandleFunc("/signup", controllers.SignupPage).Methods("GET").Name("SignupPage")
	rtr.HandleFunc("/signup", controllers.Signup).Methods("POST").Name("SignupRequest")
	rtr.HandleFunc("/logout", controllers.Logout).Methods("GET").Name("LogoutRequest")

	campSubrtr := rtr.PathPrefix("/campaign").Subrouter()
	campSubrtr.HandleFunc("/create", controllers.CreateCampaignPage).Methods("GET").Name("CreateCampaignPage")
	campSubrtr.HandleFunc("/create", controllers.CreateCampaign).Methods("POST").Name("CreateCampaignRequest")

	// badge board
	badgeSubrtr := rtr.PathPrefix("/b").Subrouter()
	badgeSubrtr.HandleFunc("/{camp_uri}", controllers.BoardPage).Methods("GET").Name("CreateCampaignPage")

	// groupSubrtr := rtr.PathPrefix("/group").Subrouter()
	// groupSubrtr.HandleFunc("/create", controllers.CBG).Methods("GET").Name("CreateBadgeGroupPage")
	// groupSubrtr.HandleFunc("/create", controllers.CreateBadgeGroup).Methods("POST").Name("CreateBadgeGroupRequest")

	// badgeSubrtr := rtr.PathPrefix("/badge").Subrouter()
	// badgeSubrtr.HandleFunc("/create", controllers.CreateBadgeT).Methods("GET").Name("CreateBadgePage")
	// badgeSubrtr.HandleFunc("/create", controllers.CreateBadge).Methods("POST").Name("CreateBadgeRequest")
}
