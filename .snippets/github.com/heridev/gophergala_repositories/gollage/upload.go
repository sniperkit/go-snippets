package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
)

// Stole this from http://sanatgersappa.blogspot.com/2013/03/handling-multiple-file-uploads-in-go.html

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	wallName := vars["id"]
	wall, ok := walls[wallName]
	// You can't add an image to a wall that doesn't exist
	if ok {
		// Wall exists
	} else {
		fmt.Println("This wall doesn't even exist?")
		// Let them know they dun goofed
		http.Redirect(w, r, "/error", 302)
		return
	}
	// Parse the multipart form in the request
	reader, err := r.MultipartReader()
	if err != nil {
		fmt.Println("Failed to parse form: " + err.Error())
		// Let them know they dun goofed
		http.Redirect(w, r, "/error", 302)
		return
	}

	//copy each part to destination.
	link := ""
	for {
		part, err := reader.NextPart()

		if err == io.EOF {
			break
		}

		defer part.Close()

		if part.FormName() == "link" {
			bName, _ := ioutil.ReadAll(part)
			link = string(bName)
		}

		//if part.FileName() is empty, skip this iteration.
		if part.FileName() == "" {
			continue
		}

		buf := new(bytes.Buffer)
		pic, err := Normalize(ImageSize, part, buf)
		if err != nil {
			fmt.Println("Failed to normalize image: " + err.Error())
			// Let them know they dun goofed
			http.Redirect(w, r, "/error", 302)
			return
		}
		img := &Image{
			Pic: pic,
			Url: link,
		}
		wall.AddImage(img)
	}
	wall.ClearPositioning()
	wall.Run()
	go wall.DrawWall()
	http.Redirect(w, r, "/wall/"+wallName, 302)
}
