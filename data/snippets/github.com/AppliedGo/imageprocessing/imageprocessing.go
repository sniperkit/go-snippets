/*
<!--
Copyright (c) 2016 Christoph Berger. Some rights reserved.
Use of this text is governed by a Creative Commons Attribution Non-Commercial
Share-Alike License that can be found in the LICENSE.txt file.

The source code contained in this file may import third-party source code
whose licenses are provided in the respective license files.
-->

<!--
NOTE: The comments in this file are NOT godoc compliant. This is not an oversight.

Comments and code in this file are used for describing and explaining a particular topic to the reader. While this file is a syntactically valid Go source file, its main purpose is to get converted into a blog article. The comments were created for learning and not for code documentation.
-->

+++
title = "Picturesque!"
description = "Image processing with Go libraries"
author = "Christoph Berger"
email = "chris@appliedgo.net"
date = "2016-12-22"
publishdate = "2016-12-22"
draft = "false"
domains = ["Image Processing"]
tags = ["Image", "Picture", "processing", "enhancement"]
categories = ["Tools And Libraries"]
+++

Let's face it: Pictures taken with a smartphone usually aren't quite like Ansel Adams masterpieces. But with a little post-processing, some of them might still reveal their true beauty. A couple of Go libraries can help.

<!--more-->

Almost all of the posts on AppliedGo.net so far are about applying Go to various problem domains and building a basic implementation from scratch. Today's post is a bit different. I picked a goal - processing an image - and searched for Go libraries to help me with that job. These are the libraries I will be using here:

* [`artyom/smartcrop`](https://github.com/artyom/smartcrop)
* [`anthonynsimon/bild`](https://github.com/anthonynsimon/bild)
* [`fogleman/primitive`](https://github.com/fogleman/primitive)

- - -

**Update:** `artyom/smartcrop` is replacing `muesli/smartcrop` that was used for the previous version of this article. `artyom/smartcrop` is a fork of `muesli/smartcrop` with no external dependencies and a simpler API.

- - -

And last not least, the `image` package from the Go standard library.

So let's start coding!

## The inevitable Imports And Globals section
*/

// Imports
package main

import (
	// basic image handling
	"image"
	// The `jpeg` package decodes and encodes JPG images.
	"image/jpeg"

	// The third-party libraries used here.
	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/transform"
	"github.com/artyom/smartcrop"
	"github.com/fogleman/primitive/primitive"
	"github.com/pkg/errors"

	//...and the rest.
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"runtime"
	"time"
)

/*

## Loading and saving

First, we want to load an image. Here is our test image. If I am not wrong, it shows a [red kite](https://en.wikipedia.org/wiki/Red_kite).

![Red Kite](original.jpg)

The `image` library provides a `Decode` function that can read JPG, GIF, and PNG data, provided that the appropriate sub-package has been loaded (see the import section).

And while we are at it, let's also define a function for saving an image.

*/

//openImage imports an image from a given path.
func openImage(path string) (image.Image, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot open "+path)
	}

	// Decode from JPG into image.Image format.
	img, err := jpeg.Decode(imgFile)
	if err != nil {
		return nil, errors.Wrap(err, "Decoding the image failed.")
	}

	return img, nil
}

// saveImage saves the image to `pname/fname.jpg`.
func saveImage(img image.Image, pname, fname string) error {
	fpath := path.Join(pname, fname)

	f, err := os.Create(fpath)
	if err != nil {
		return errors.Wrap(err, "Cannot create file: "+fpath)
	}
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 85})
	if err != nil {
		return errors.Wrap(err, "Failed to encode the image as JPEG")
	}
	return nil
}

/*

## Smartcrop

The test image has quite some space around the Red Kite with nothing interesting to see. So let's crop the image. But not manually; instead, let `smartcrop` do the job. `smartcrop` attempts to detect the most interesting part of an image.

Note that Smartcrop can use face recognition for finding the right crop. Obviously, we cannot use this feature here on the bird picture, so we switch it off.

Smartcrop does not crop the image itself, it only returns the suggested cropping rectangle. We can use the standard `Image` library for cropping the image. The `Image` type has no crop method, but the color types have a method called `SubImage`, like, for example, the RGBA type:

    func (p *RGBA) SubImage(r Rectangle) Image

How can we access this method? We could type-assert the `Image` to the appropriate color type (like RGBA, CMYK, etc.) but `Image`'s properties give us no clue which color type our JPEG image has been decoded to.

The [solution](https://stackoverflow.com/questions/16072910/trouble-getting-a-subimage-of-an-image-in-go) is to type-assert `Image` to an interface that consists of just the SubImage method. Then we can call `SubImage` without knowing the exact color type.

*/

// The SubImager interface exposes the SubImage method to facilitate the type conversion from `Image` to the appropriate color type.
type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

// `crop` auto-crops the image in-place.
func crop(img image.Image, width, height int) (image.Image, error) {

	rect, err := smartcrop.Crop(img, width, height)
	if err != nil {
		return nil, errors.Wrap(err, "Smartcrop failed")
	}

	// Now let's crop the image to the suggested area.
	// First, we need to apply the aforementioned type assertion.
	si, ok := (img).(SubImager)
	if !ok {
		return nil, errors.New("crop(): img does not support SubImage()")
	}
	// Then we pass the cropping borders to `SubImage()`. Note that the returned sub-image shares pixels with the original image, so make a copy if you want to manipulate only the sub-image.
	subImg := si.SubImage(rect)

	return subImg, nil
}

/*

The result is not too bad! The algorithm found the interesting part of the image, although I would have put the bird a tad bit more towards the center. But hey, that's an automated algorithm that is certainly not specialized for identifying birds, so the result is perfectly ok.

![Smartcrop result](cropped.jpg)


## bild

Next, let's try fixing the somewhat over-exposured foreground and grass. For this, I'll use [`anthonynsimon/bild`](https://github.com/anthonynsimon/bild), a comprehensive image manipulation library. (If you wonder about the name, "Bild" is the German word for "picture" or "image".)

`bild` uses `image.Image` as image format, so we can reuse the `img` variable without having to save and re-open the image.

`bild` is organized as sub-packages that group related operations. For example, the adjust package provides adjustments, the blend package provides image blending operations, and so on.

Let's try a few things for fun (each time starting from the unmodified image).

As the colors seem a bit pale, let's try increasing the saturation.
*/

// Apply 50% saturation
func saturate(img image.Image) image.Image {
	return adjust.Saturation(img, 0.5)
}

/*
Before:

![Cropped](cropped.jpg)

After:

![Saturated](saturated.jpg)

Already looks better, after just one simple adjustment!


Next test: What happens if we multiply the image with itself?
*/

// Multiply the image with itself
func multiply(img image.Image) image.Image {
	return blend.Multiply(img, img)
}

/*
That's interesting: Dark colors are darker, and so are the lighter colors, but not that much as the darker ones, and they also seem more intense.

![Multiplied](multiplied.jpg)

Try more effects for yourself! Especially, try to combine two or more effects to get new results.

As a last test with `bild`, let's sharpen the saturated image.

*/

// Sharpen the image using unsharp masking.
func sharpen(img image.Image) image.Image {
	return effect.UnsharpMask(img, 0.6, 1.2)
}

/*

For better comparison, I zoomed in and put the before and after images side-by-side.

![Sharpened](sharpenedBeforeAfter.jpg)

## primitive

The next package is `fogleman\primitive`. Don't be fooled by the name; this package is anything but primitive. The name has a meaning though: This package "reproduces" an image by applying geometric primitives like rectangles, ellipses, etc. to it.

This package comes as a binary package; however, it is well structured and includes sub-packages, so after peeking into `main.go` we can integrate the algorithm in our code.
*/

//Making art.
func primitivePicture(img image.Image) image.Image {

	// Resize the image to 256x256 to save processing time.
	// `transform` is a `bild` package.

	img = transform.Resize(img, 256, 256, transform.Linear)

	// Seed random number generator.
	rand.Seed(time.Now().UTC().UnixNano())

	// Set the background color.
	bg := primitive.MakeColor(primitive.AverageImageColor(img))

	// NewModel(image, background color, output size, # of workers)
	model := primitive.NewModel(img, bg, 1024, runtime.NumCPU())

	for i := 0; i < 100; i++ {
		// 5 = rotated rectangles,
		// 128 = default alpha,
		// 0 = default repeat
		fmt.Print(".")
		model.Step(primitive.ShapeType(5), 128, 0)
	}

	return model.Context.Image()
}

/*

Here is the result:

![Primitive](primitive.jpg)

Niiice!

If the first result is not satisfying, simply run the code again. The results will be different each time.

Also try other mode values (replace the "5" in `primitive.ShapeType(5)` in the call to `model.Step()` above).

Valid values are (from `primitive`'s help text):

    0=combo 1=triangle 2=rect 3=ellipse 4=circle
	5=rotatedrect 6=beziers 7=rotatedellipse 8=polygon

Now the test image is not really suited for dramatic effects, so feel free to visit [`primitive`'s GitHub repository](https://github.com/fogleman/primitive) to see a couple of awesome Primitive Pictures!

Last not least, the `main` function connects all the code snippets.
*/

// main
func main() {
	img, err := openImage("original.jpg")
	if err != nil {
		log.Fatal(err)
	}

	// If you don't want to install opencv, just comment out the crop() and saveImage() calls and the related error checks.
	//
	// Crop attempts to find the best crop of img based on the given width and height values.
	img, err = crop(img, 1000, 1000)
	if err != nil {
		log.Fatal(err)
	}
	err = saveImage(img, ".", "cropped.jpg")
	if err != nil {
		log.Fatal(err)
	}

	// Let's continue with a manually cropped image.
	img, err = openImage("cropped.jpg")
	if err != nil {
		log.Fatal(err)
	}

	// Fun with `bild`.
	sat := saturate(img)
	err = saveImage(sat, ".", "saturated.jpg")
	if err != nil {
		log.Fatal(err)
	}

	mult := multiply(img)
	err = saveImage(mult, ".", "multiplied.jpg")
	if err != nil {
		log.Fatal(err)
	}

	shrp := sharpen(sat)
	err = saveImage(shrp, ".", "sharpened.jpg")
	if err != nil {
		log.Fatal(err)
	}

	// Create "primitive" art.
	pri := primitivePicture(sat)
	err = saveImage(pri, ".", "primitive.jpg")
	if err != nil {
		log.Fatal(err)
	}

}

/*

Get the full code from [GitHub](https://github.com/appliedgo/imageprocessing):

```
go get -d github.com/appliedgo/imageprocessing
cd $GOPATH/src/github.com/appliedgo/imageprocessing
go run imageprocessing.go
```

## Odds and Ends

I planned to include a section on halftoning; however, I quickly found that although code is available, it is not a library - at least not yet. So I invite you to head over to [Halftoning with Go - Part 1](https://maxhalford.github.io/blog/halftoning-1/), which is a very intersting read about halftoning and dithering, starting with average dithering and ending with the Floyd-Steinberg algorithm.

Find more libaries and tools at -

* [Awesome Go](http://awesome-go.com/#images)
* [LibHunt](https://go.libhunt.com/categories/498-images)

This is the last post for this year. Enjoy the holidays! See you again in January.

Until then, happy coding!


- - -

Changelog

2016-12-23: Replaced `muesli/smartcrop` by `artyom/smartcrop`.

2016-12-24: Added missing steps for getting the source code
*/
