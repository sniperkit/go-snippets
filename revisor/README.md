Revisor is a simple web application for editing the text files behind websites created by applications like Luminos.

The project is based on this tutorial:
https://golang.org/doc/articles/wiki/

At the moment this just runs locally. Here's how to try it:
* go get github.com/gophergala/revisor
* cd to to the code $GOPATH/src/github/gophergala/revisor
* Create directory called "website"

Now you can go run main.go and and Revisor will start on localhost:8080



Todos
* clean the messy recover where the directory list is rebuilt
* implement auth (probably oauth using go-oauth)
* improve directory listing
* handle CRUD for non-text files
