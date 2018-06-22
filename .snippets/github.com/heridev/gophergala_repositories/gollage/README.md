# gollage

I've never used Pinterest, but I think this might be similar. The idea was to
have a system where people could create "walls" and put images and links on
them, where each wall has a theme. Gollage accomplishes this to some degree.

## To Use
1. Make a wall on [Gollage](http://www.gollage.io)
2. Give out the link to your wall to people. It'll be [www.gollage.io/wall/YourWallName](www.gollage.io/wall/Default)
3. Upload images and set links on them. You can upload multiple files at once.
4. Refresh the page (I background the dynamic image generation/resizing process)
5. Enjoy your Gollage!

## What Works
* Wall Making
* Linked Images
* Thumbnails
* Dynamic Image Tiling
* Multiple Image Upload

## What Doesn't Work
* Some issues with the tiling algorithm
* WebSockets are implemented, but don't let you know when the wall you're looking at is updated
* It isn't as pretty as I'd like
* Tiling algorithm doesn't scale nicely yet
* Originally I wanted to be able to zoom in on huge collages, but you can't currently do that

## Taking advantage of Go
Go is used for pretty much everything, with only a few lines of Javascript
used. I use Go for:

* Dynamic image composition
* Image Resizing
* HTTP serving/AWS interaction
* Templating pages

I use Goroutines to make worker threads for processing image composition. All
image manipulation happens in Go, no external tools (like ImageMagick) are
used. Standard packages and Gorilla packages are used extensively.

# Stream of Consciousness Updates
These are the things I was writing as I worked on this at various points throughout this weekend.

I don't really know what this is yet, I'm going to start with a web server and
some WebSockets and the basis of an idea. Then I'm going to get some scorpion
bowls, and hope the Ballmer Peak does the rest.

## Follow up

Totally missed the Ballmer Peak, was conviced code should have been in English
and not Go, got frustrated and passed out on couch. 

## Update

Ate lots of greasy breakfast foods, hangover successfully quelled. I think I've
written code things, but I have no guarantee that anything does literally
anything at all.

## Day 7

Ideas are starting to flow, pages are starting to form, morale is at an
all-time high. Bathroom break imminent, caffeine likely to blame.

## Unholy things are happening

You know it's bad when you've resorted to turing a Javascript algorithm into
Go, gross. Also I have no idea what I'm doing, how do I image?

## Back on Track

After hours of flopping around like an electricuted fish out of water in an
earthquake, thrashing like a system with 1 MB of RAM trying to load Unity, I've
finally oriented myself and put together something that kind of works. Current
functionality is ability to upload images, functional image placement
algorithm, functional thumbnails and other little things. Next steps include
magical things like WebSocket updates of the page, smarter collage scaling,
functional URLs, and div clickbox mapping with ImageJSON method. But first, a
meeting with some clients. I'm definitely going to be cutting this too close
for comfort, but I wouldn't really be comfortable with anything else.
