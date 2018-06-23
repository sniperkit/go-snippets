package main

//...
buf := new(bytes.Buffer)
if err := tmpl.Execute(buf, profile); err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
}
templateString := buf.String()
//...