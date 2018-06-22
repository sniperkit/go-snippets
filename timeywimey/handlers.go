package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	body, err := localFile("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(body))
}

func AuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	profile := r.FormValue("profile")
	passphrase := r.FormValue("p")

	if profile == "" || passphrase == "" {
		// Look for a cookie?
	} else {
		// Set cookies
		expire := time.Now().AddDate(0, 0, 1)
		cookieProfile := &http.Cookie{
			Name:    "profile",
			Value:   profile,
			Expires: expire,
		}
		cookiePassphrase := &http.Cookie{
			Name:    "passphrase",
			Value:   passphrase,
			Expires: expire,
		}
		http.SetCookie(w, cookieProfile)
		http.SetCookie(w, cookiePassphrase)
	}

	if profile == "" || passphrase == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/"+profile, http.StatusSeeOther)
	return
}

func UserIndexHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	projects, err := GetByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := localFile("user_index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := template.Must(template.New("user").Parse(string(body)))
	err = t.Execute(w, projects)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ProjectIndexHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	projectName := mux.Vars(r)["project"]

	var p Project
	err := p.Get(projectName, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := localFile("project_index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(body))
}

func ProjectNewHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	projectName := mux.Vars(r)["project"]

	var p Project
	err := p.Get(projectName, username)
	if err != nil {
		// How should we handle this? Lack of sleep and a screaming child...
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//return
	}

	body, err := localFile("project_new.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := template.Must(template.New("project").Parse(string(body)))
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func ProjectUpdateHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	projectName := mux.Vars(r)["project"]

	p := Project{
		Name:    projectName,
		Owner:   username,
		Members: []string{username},
	}
	p.Get(projectName, username)

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	members := r.Form["member"]
	if p.Owner != username {
		http.Error(w, "Invalid permissions", http.StatusUnauthorized)
		return
	}
	for _, m := range members {
		if m != "" {
			p.Members = append(p.Members, m)
		}
	}
	if len(p.Members) == 0 {
		p.Members = []string{username}
	}

	if err := p.Insert(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%s", username, projectName), http.StatusSeeOther)
}

func MeetingCalendarHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	projectName := mux.Vars(r)["project"]

	var p Project
	err := p.Get(projectName, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	moments, err := p.Meetings.Within(start, end)
	js, err := json.Marshal(moments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func MeetingShowHandler(w http.ResponseWriter, r *http.Request) {
}

func MeetingUpdateHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	projectName := mux.Vars(r)["project"]

	var p Project
	err := p.Get(projectName, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	start := r.Form.Get("startDate")
	startDate, err := time.Parse("2006-01-02T15:04", start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	end := r.Form.Get("endDate")
	endDate, err := time.Parse("2006-01-02T15:04", end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	repeating := r.Form.Get("repeating") != ""
	title := r.Form.Get("title")
	summary := r.Form.Get("summary")

	meeting := Moment{
		ObjectType:   "Event",
		StartDate:    startDate,
		EndDate:      endDate,
		Repeating:    repeating,
		Title:        title,
		Summary:      summary,
		CalendarData: "",
		LastModified: time.Now(),
	}

	p.Meetings.Moments = append(p.Meetings.Moments, meeting)
	if err = p.UpdateMeetings(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%s", username, projectName), http.StatusSeeOther)
}

func MeetingDeleteHandler(w http.ResponseWriter, r *http.Request) {
}

func IssueCalendarHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	projectName := mux.Vars(r)["project"]

	var p Project
	err := p.Get(projectName, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	moments, err := p.Issues.Within(start, end)
	js, err := json.Marshal(moments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func IssueShowHandler(w http.ResponseWriter, r *http.Request) {
}

func IssueUpdateHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	projectName := mux.Vars(r)["project"]

	var p Project
	err := p.Get(projectName, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	start := r.Form.Get("start")
	startDate, err := time.Parse("2006-01-02T15:04", start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	end := r.Form.Get("end")
	endDate, err := time.Parse("2006-01-02T15:04", end)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	title := r.Form.Get("title")
	summary := fmt.Sprintf("Type: %s\nDescription: %s", r.Form.Get("type"), r.Form.Get("summary"))

	meeting := Moment{
		ObjectType:   "Event",
		StartDate:    startDate,
		EndDate:      endDate,
		Repeating:    false,
		Title:        title,
		Summary:      summary,
		CalendarData: "",
		LastModified: time.Now(),
	}

	p.Issues.Moments = append(p.Issues.Moments, meeting)
	if err = p.UpdateIssues(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%s", username, projectName), http.StatusSeeOther)
}

func IssueDeleteHandler(w http.ResponseWriter, r *http.Request) {
}

func TimeCalendarHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	projectName := mux.Vars(r)["project"]

	var p Project
	err := p.Get(projectName, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	moments, err := p.Timesheet.Within(start, end)
	js, err := json.Marshal(moments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func TimeShowHandler(w http.ResponseWriter, r *http.Request) {
}

func TimeUpdateHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	projectName := mux.Vars(r)["project"]

	var p Project
	err := p.Get(projectName, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	date := r.Form.Get("date")
	start := r.Form.Get("start")
	startDate, err := time.Parse("2006-01-02T15:04", date+"T"+start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stop := r.Form.Get("stop")
	endDate, err := time.Parse("2006-01-02T15:04", date+"T"+stop)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	title := fmt.Sprintf("%s punched out", username)
	repeating := false
	summary := r.Form.Get("summary")

	meeting := Moment{
		ObjectType:   "Event",
		StartDate:    startDate,
		EndDate:      endDate,
		Repeating:    repeating,
		Title:        title,
		Summary:      summary,
		CalendarData: "",
		LastModified: time.Now(),
	}

	p.Timesheet.Moments = append(p.Timesheet.Moments, meeting)
	if err = p.UpdateTimesheet(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%s", username, projectName), http.StatusSeeOther)
}

func TimeDeleteHandler(w http.ResponseWriter, r *http.Request) {
}

func InvoiceCalendarHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	projectName := mux.Vars(r)["project"]

	var p Project
	err := p.Get(projectName, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	moments, err := p.Invoices.Within(start, end)
	js, err := json.Marshal(moments)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func InvoiceShowHandler(w http.ResponseWriter, r *http.Request) {
}

func InvoiceUpdateHandler(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	projectName := mux.Vars(r)["project"]

	var p Project
	err := p.Get(projectName, username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	start := r.Form.Get("startDate")
	startDate, err := time.Parse("2006-01-02T15:04", start)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stop := r.Form.Get("stopDate")
	endDate, err := time.Parse("2006-01-02T15:04", stop)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	format := "2006-01-02"
	title := fmt.Sprintf("Invoice for %s to %s", startDate.Format(format), endDate.Format(format))
	lineItems := r.Form["lineItem"]
	summary := ""
	for _, l := range lineItems {
		summary = fmt.Sprintf("%sBillable Item: %s\n", summary, l)
	}
	taxes := r.Form.Get("taxes")
	summary = fmt.Sprintf("%sTaxes: %s\n", summary, taxes)

	meeting := Moment{
		ObjectType:   "Event",
		StartDate:    startDate,
		EndDate:      endDate,
		Repeating:    false,
		Title:        title,
		Summary:      summary,
		CalendarData: "",
		LastModified: time.Now(),
	}

	p.Invoices.Moments = append(p.Invoices.Moments, meeting)
	if err = p.UpdateInvoices(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%s", username, projectName), http.StatusSeeOther)
}

func InvoiceDeleteHandler(w http.ResponseWriter, r *http.Request) {
}

func localFile(filename string) ([]byte, error) {
	page, err := ioutil.ReadFile("html/" + filename)
	if err != nil {
		return nil, err
	}
	return page, nil
}
