package main

import (
	"log"

	gosxnotifier "github.com/deckarep/gosx-notifier"
)

func main() {

	// message (required)
	note := gosxnotifier.NewNotification("Follow me on Github ...")

	// title
	note.Title = "Maik Ellerbrock ⛅  ✨"

	// link
	note.Link = "http://github.com/ellerbrock" //or BundleID like: com.apple.Terminal

	// group indentifier (show one notification per group)
	note.Group = "com.frapsoft.notifier.identifier"

	// content icon
	note.ContentImage = "frapsoft.png"

	/* play Sound

	   Basso
	   Blow
	   Bottle
	   Frog
	   Funk
	   Glass
	   Hero
	   Morse
	   Ping
	   Pop
	   Purr
	   Sosumi
	   Tink
	*/

	note.Sound = gosxnotifier.Glass

	// subtitle
	// note.Subtitle = "Lets hack on stuff ..."

	// use system icons
	// note.Sender = "com.apple.Safari"

	// app icon
	// note.AppIcon = "frapsoft.png"

	// push notification
	err := note.Push()

	// error
	if err != nil {
		log.Println("Uh oh!")
	}
}
