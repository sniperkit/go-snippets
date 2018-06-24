package controllers

import (
	"github.com/rsmohamad/comp4321/database"
	"html/template"
	"net/http"
)

var historyTemplate *template.Template

func clearHandler(w http.ResponseWriter, r *http.Request) {
	userId := database.GetCookieInstance().GetCookieId(r)
	database.GetCookieInstance().ClearSearchHistory(userId)
	historyHandler(w, r)
}

func historyHandler(w http.ResponseWriter, r *http.Request) {
	userId := database.GetCookieInstance().GetCookieId(r)
	history := database.GetCookieInstance().GetSearchHistory(userId)
	historyTemplate.Execute(w, history)
}

func LoadHistory() {
	historyTemplate, _ = template.ParseFiles("views/historyView.html")
	http.HandleFunc("/history", historyHandler)
	http.HandleFunc("/history/clear", clearHandler)
}
