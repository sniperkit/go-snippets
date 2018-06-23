package app

import (
	"github.com/gophergala/GopherKombat/common/rankings"
	"net/http"
)

func RankingsHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := GetCurrentUser(r)
	data := make(map[string]interface{})
	data["loggedIn"] = ok
	if ok {
		data["user"] = user
		data["today"] = rankings.GetDaily()
		data["month"] = rankings.GetMonthly()
		data["allTime"] = rankings.GetAllTime()
	}
	render(w, "rankings", data)
}
