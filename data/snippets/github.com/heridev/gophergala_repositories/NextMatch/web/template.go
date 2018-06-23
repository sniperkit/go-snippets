package main

import (
	"html/template"
	"net/http"
)

var (
	tLoader *template.Template
	funcMap = make(template.FuncMap)
)

type args map[string]interface{}

func loadTmpl() {
	// cache template parsing
	tLoader = template.Must(template.New(``).Funcs(funcMap).ParseGlob(`resources/templates/*`))
}

func execT(w http.ResponseWriter, name string, data interface{}) {
	if err := tLoader.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func head(w http.ResponseWriter, data interface{}) {
	execT(w, `header`, data)
}

func footer(w http.ResponseWriter, data interface{}) {
	execT(w, `footer`, data)
}

func content(w http.ResponseWriter, data interface{}) {
	execT(w, `content`, data)
}

func defaultT(w http.ResponseWriter, data interface{}) {
	head(w, data)
	content(w, data)
	footer(w, data)
}

func customT(w http.ResponseWriter, name string, data interface{}) {
	head(w, data)
	execT(w, name, data)
	footer(w, data)
}

func addTfunc(name string, fn interface{}) {
	funcMap[name] = fn
}

var sampleData = `
{"event":[{"event_id":"20150124-new-york-knicks-at-charlotte-hornets","event_status":"completed","away_team":{"abbreviation":"NY","active":true,"city":"New York","conference":"East","division":"Atlantic","first_name":"New York","full_name":"New York Knicks","last_name":"Knicks","site_name":"Madison Square Garden","state":"New York","team_id":"new-york-knicks"},"home_team":{"abbreviation":"CHA","active":true,"city":"Charlotte","conference":"East","division":"Southeast","first_name":"Charlotte","full_name":"Charlotte Hornets","last_name":"Hornets","site_name":"Time Warner Cable Arena","state":"North Carolina","team_id":"charlotte-hornets"},"season_type":"regular","site":{"capacity":19026,"city":"Charlotte","name":"Time Warner Cable Arena","state":"North Carolina","surface":"Hardwood"},"sport":"NBA","start_date_time":"2015-01-24T19:00:00-05:00"},{"event_id":"20150124-detroit-pistons-at-milwaukee-bucks","event_status":"completed","away_team":{"abbreviation":"DET","active":true,"city":"Auburn Hills","conference":"East","division":"Central","first_name":"Detroit","full_name":"Detroit Pistons","last_name":"Pistons","site_name":"The Palace of Auburn Hills","state":"Michigan","team_id":"detroit-pistons"},"home_team":{"abbreviation":"MIL","active":true,"city":"Milwaukee","conference":"East","division":"Central","first_name":"Milwaukee","full_name":"Milwaukee Bucks","last_name":"Bucks","site_name":"BMO Harris Bradley Center","state":"Wisconsin","team_id":"milwaukee-bucks"},"season_type":"regular","site":{"capacity":18717,"city":"Milwaukee","name":"BMO Harris Bradley Center","state":"Wisconsin","surface":""},"sport":"NBA","start_date_time":"2015-01-24T19:30:00-05:00"},{"event_id":"20150124-philadelphia-76ers-at-memphis-grizzlies","event_status":"completed","away_team":{"abbreviation":"PHI","active":true,"city":"Philadelphia","conference":"East","division":"Atlantic","first_name":"Philadelphia","full_name":"Philadelphia 76ers","last_name":"76ers","site_name":"Wachovia Center","state":"Pennsylvania","team_id":"philadelphia-76ers"},"home_team":{"abbreviation":"MEM","active":true,"city":"Memphis","conference":"West","division":"Southwest","first_name":"Memphis","full_name":"Memphis Grizzlies","last_name":"Grizzlies","site_name":"FedExForum","state":"Tennessee","team_id":"memphis-grizzlies"},"season_type":"regular","site":{"capacity":18165,"city":"Memphis","name":"FedExForum","state":"Tennessee","surface":"Hardwood"},"sport":"NBA","start_date_time":"2015-01-24T20:00:00-05:00"},{"event_id":"20150124-brooklyn-nets-at-utah-jazz","event_status":"scheduled","away_team":{"abbreviation":"BKN","active":true,"city":"Brooklyn","conference":"East","division":"Atlantic","first_name":"Brooklyn","full_name":"Brooklyn Nets","last_name":"Nets","site_name":"Barclays Center","state":"New York","team_id":"brooklyn-nets"},"home_team":{"abbreviation":"UTA","active":true,"city":"Salt Lake City","conference":"West","division":"Northwest","first_name":"Utah","full_name":"Utah Jazz","last_name":"Jazz","site_name":"EnergySolutions Arena","state":"Utah","team_id":"utah-jazz"},"season_type":"regular","site":{"capacity":20000,"city":"Salt Lake City","name":"EnergySolutions Arena","state":"Utah","surface":"Hardwood"},"sport":"NBA","start_date_time":"2015-01-24T21:00:00-05:00"},{"event_id":"20150124-washington-wizards-at-portland-trail-blazers","event_status":"scheduled","away_team":{"abbreviation":"WAS","active":true,"city":"Washington","conference":"East","division":"Southeast","first_name":"Washington","full_name":"Washington Wizards","last_name":"Wizards","site_name":"Verizon Center","state":"District of Columbia","team_id":"washington-wizards"},"home_team":{"abbreviation":"POR","active":true,"city":"Portland","conference":"West","division":"Northwest","first_name":"Portland","full_name":"Portland Trail Blazers","last_name":"Trail Blazers","site_name":"Moda Center","state":"Oregon","team_id":"portland-trail-blazers"},"season_type":"regular","site":{"capacity":20636,"city":"Portland","name":"Moda Center","state":"Oregon","surface":""},"sport":"NBA","start_date_time":"2015-01-24T22:00:00-05:00"}],"events_date":"2015-01-24T00:00:00-05:00"}
`
