package main

import (
    "github.com/miketheprogrammer/go-thrust/lib/dispatcher"
    "github.com/miketheprogrammer/go-thrust/lib/spawn"
    "github.com/miketheprogrammer/go-thrust"
)

func main() {
    go Run()
    spawn.Run()
    thrustWindow := thrust.NewWindow("http://localhost:8080", nil)
    thrustWindow.Show()
    thrustWindow.Maximize()
    thrustWindow.Focus()
    dispatcher.RunLoop()
}