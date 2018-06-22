package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var page = `
<script type="text/javascript">
var int=self.setInterval("ajaxFunction()",1000);
function ajaxFunction(){
	var ajaxRequest;
	ajaxRequest = new XMLHttpRequest();
	ajaxRequest.onreadystatechange = function(){
		if(ajaxRequest.readyState == 4){
			var ajaxDisplay = document.getElementById('ajaxDiv');
			ajaxDisplay.innerHTML = ajaxRequest.responseText;
		}
	}
	ajaxRequest.open("GET", "/ajax?api=showclock", true);
	ajaxRequest.send(null);
}
</script>
<pre><a id='ajaxDiv'></a></pre>`

func ajax(w http.ResponseWriter, r *http.Request) {
	//this must add at begin of every session code
	c, err := r.Cookie("session")
	if err != nil || c.Value == "" {
		http.Error(w, "Session expired", 401)
		return
	}

	//return results request by ajax page
	if r.Method == "GET" {
		if r.FormValue("api") == "showclock" {
			fmt.Fprint(w, time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST"))
			return
		}
	}

	//build page content
	b := `<pre>
This module is a simple ajax application who show a real time clock.</pre>` + page

	//finally show the page
	p := Page{
		Title:  "Test page",
		Status: c.Value,
		Body:   template.HTML(b),
	}
	t.ExecuteTemplate(w, "index.html", p)
}

func init() {
	http.HandleFunc("/ajax", ajax)
}
