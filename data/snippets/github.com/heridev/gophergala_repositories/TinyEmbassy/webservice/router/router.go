/*
* @Author: souravray
* @Date:   2015-01-24 11:30:37
* @Last Modified by:   souravray
* @Last Modified time: 2015-01-25 03:10:49
 */

package router

import (
	"github.com/gophergala/tinyembassy/webservice/controllers"
	"github.com/gorilla/mux"
	//"net/http"
)

func Routes(rtr *mux.Router) {

	// web pages routes
	// webSubrtr := rtr.PathPrefix("/web/{camp}").Subrouter()
	// webSubrtr.HandleFunc("/login", controllers.Landing).Methods("GET").Name("LogIn")
	// webSubrtr.HandleFunc("/login", controllers.SignIn).Methods("POST").Name("LogIn")
	// webSubrtr.HandleFunc("/badge", controllers.GetBadge).Methods("GET").Name("LogIn")

	//track routes
	apiSubrtr := rtr.PathPrefix("/track").Subrouter()
	apiSubrtr.HandleFunc("/img/{camp}/{refr}/{badge}", controllers.RedirectImage).Methods("GET").Name("RedirectImage")
	apiSubrtr.HandleFunc("/click/{camp}/{refr}/{badge}", controllers.RedirectTargetURL).Methods("GET").Name("RedirectTargetURL")
	apiSubrtr.HandleFunc("/action/{marker}", controllers.DataPoint).Methods("GET").Name("DataPoint")
}
