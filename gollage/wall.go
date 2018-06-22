package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/png"
	"math/rand"
	"net/http"
	"strings"

	_ "code.google.com/p/vp8-go/webp"
	_ "image/jpeg"
)

const GridWidth = 960
const GridHeight = 540

// Make this computed
const RowMaxHeight = 100
const HorizontalMargin = 5
const VerticalMargin = 5

type Wall struct {
	Images  []*Image
	Url     string
	Name    string
	Heights []int
}

type Image struct {
	Pic        image.Image
	Url        string
	XOffset    int
	YOffset    int
	DispWidth  int
	DispHeight int
}

type Metadata struct {
	Url        string
	XOffset    int
	YOffset    int
	DispWidth  int
	DispHeight int
}

func (w *Wall) AddImage(img *Image) {
	w.Images = append(w.Images, img)
}

func newWallHandler(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	name = strings.Replace(name, " ", "", -1)

	if len(name) == 0 {
		fmt.Println("No name entered")
		// Let them know they dun goofed
		http.Redirect(w, r, "/error", 302)
		return
	}
	// If it exists, they can't have it
	if _, ok := walls[name]; ok {
		// Sorry brah, this wall's taken
		http.Redirect(w, r, "/error", 302)
		return
	} else {
		//
		err := NewWallBucket(name)
		if err != nil {
			fmt.Println("Error making bucket", err)
			// Let them know we couldn't persist it
			http.Redirect(w, r, "/error", 302)
			return
		} else {
			//Don't make the wall until we're sure we can persist it
			walls[name] = &Wall{
				Images: []*Image{},
				Name:   name,
			}
			http.Redirect(w, r, "/wall/"+name, 302)
			return
		}
	}
}

func wallHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Make sure the wall exists
	wall, ok := walls[vars["id"]]
	if ok {
		data := struct {
			Wall     Wall
			Channel  string
			Width    int
			Height   int
			Host     string
			LinkJSON string
		}{
			*wall,
			wall.Name,
			GridWidth,
			GridHeight,
			r.Host,
			wall.ImageLocJSON(),
		}
		err := templates.ExecuteTemplate(w, "wall.html", data)
		if err != nil {
			fmt.Println("Error executing template:", err)
		}
	} else {
		fmt.Println("Tried to view non-existent wall")
		http.Redirect(w, r, "/error", 302)
	}
}

// Stole this javascript from http://blog.vjeux.com/wp-content/uploads/2012/05/google-layout.html

func GetHeight(images []*Image, width int) int {
	width -= len(images) * HorizontalMargin
	height := 0
	for _, img := range images {
		height += img.Pic.Bounds().Dx() / img.Pic.Bounds().Dy()
	}
	return width / height
}

func (w *Wall) CalcYOffset(rowNum int) int {
	yOffset := 0
	if rowNum > 0 {
		for _, rowHeight := range w.Heights[:rowNum] {
			yOffset += rowHeight
		}
	}
	return yOffset
}

func (w *Wall) SetRow(images []*Image, height, rowNum int) {
	w.Heights = append(w.Heights, height+VerticalMargin)
	xOffset, yOffset := 0, w.CalcYOffset(rowNum)
	for _, image := range images {
		bounds := image.Pic.Bounds()
		image.DispWidth = height * bounds.Dx() / bounds.Dy()
		image.DispHeight = height

		image.XOffset = xOffset
		image.YOffset = yOffset

		xOffset += image.DispWidth + HorizontalMargin
	}
}

func (w *Wall) Run() {
	maxHeight := RowMaxHeight
	var slice []*Image
	var height int
	n := 0
	images := w.Images
OuterLoop:
	for len(images) > 0 {
		for i := 1; i < len(images)+1; i++ {
			slice = images[:i]
			height = GetHeight(slice, GridWidth)
			if height < maxHeight {
				w.SetRow(slice, height, n)
				images = images[i:]
				n++
				continue OuterLoop
			}
		}
		w.SetRow(images, Min(maxHeight, height), n)
		break
	}
}

func Min(inputs ...int) int {
	smallest := inputs[0]
	for _, num := range inputs {
		if num < smallest {
			smallest = num
		}
	}
	return smallest
}

func (w *Wall) DrawWall() {

	b := image.Rect(0, 0, GridWidth, GridHeight)
	m := image.NewRGBA(b)
	var original = make(chan Image, 100)
	var resized = make(chan Image, 100)

	// We make worker threads for resizing images whenever we draw a new wall
	for i := 0; i < 5; i++ {
		go ResizeWorker(original, resized)
	}

	for _, img := range w.Images {
		original <- *img
	}
	close(original)

	for _, _ = range w.Images {
		img := <-resized
		loc := image.Rect(img.XOffset, img.YOffset, img.XOffset+img.DispWidth, img.YOffset+img.DispHeight)
		draw.Draw(m, loc, img.Pic, image.ZP, draw.Src)
	}

	out := new(bytes.Buffer)
	encoder := png.Encoder{png.BestCompression}
	encoder.Encode(out, m)

	// Out the full sized one
	AddWallImage(w.Name, "full", out)

	// Downsize it, out the thumbnail
	thumb := resize.Resize(GridWidth/4, 0, m, resize.NearestNeighbor)
	out = new(bytes.Buffer)
	encoder.Encode(out, thumb)
	AddWallImage(w.Name, "thumb", out)
}

func (w *Wall) ClearPositioning() {
	w.Heights = []int{}
	for _, img := range w.Images {
		img.XOffset = 0
		img.YOffset = 0
		img.DispWidth = 0
		img.DispHeight = 0
	}
}

func ResizeWorker(originals <-chan Image, resized chan<- Image) {
	for img := range originals {
		newImage := Image{
			XOffset:    img.XOffset,
			YOffset:    img.YOffset,
			DispWidth:  img.DispWidth,
			DispHeight: img.DispHeight,
		}
		newImage.Pic = resize.Resize(uint(img.DispWidth), uint(img.DispHeight), img.Pic, resize.NearestNeighbor)
		resized <- newImage
	}
}

func (w *Wall) ImageLocJSON() string {
	data := make([]Metadata, len(w.Images))
	for i, img := range w.Images {
		data[i].Url = img.Url
		data[i].XOffset = img.XOffset
		data[i].YOffset = img.YOffset
		data[i].DispWidth = img.DispWidth
		data[i].DispHeight = img.DispHeight
	}
	// Should probably check that error, or nah, it is a hackathon
	js, _ := json.Marshal(data)
	return string(js)
}

func RandomWalls(num int) []*Wall {
	retWalls := make([]*Wall, num)
	for i := 0; i < num; i++ {
		retWalls[i] = RandomWall()
	}
	return retWalls
}

func RandomWall() *Wall {
	// produce a pseudo-random number between 0 and len(a)-1
	i := int(float32(len(walls)) * rand.Float32())
	for _, v := range walls {
		if i == 0 {
			return v
		} else {
			i--
		}
	}
	return nil
}
