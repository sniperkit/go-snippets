package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"math/rand"
	"os"
	"time"
)

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())

	out, err := os.Create("./output.gif")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// generate some QR code look a like image

	imgRect := image.Rect(0, 0, 100, 100)
	img := image.NewGray(imgRect)
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.ZP, draw.Src)
	for y := 0; y < 100; y += 10 {
		for x := 0; x < 100; x += 10 {
			fill := &image.Uniform{color.Black}
			if rand.Intn(10)%2 == 0 {
				fill = &image.Uniform{color.White}
			}
			draw.Draw(img, image.Rect(x, y, x+10, y+10), fill, image.ZP, draw.Src)
		}
	}

	var opt gif.Options

	opt.NumColors = 256 // you can add more parameters if you want

	// ok, write out the data into the new GIF file

	err = gif.Encode(out, img, &opt) // put num of colors to 256   <------- here
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Generated image to output.gif \n")
}
