package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/gophergala/aeris/download"
	"github.com/gophergala/aeris/info"
)

var (
	inputItag int
)

func main() {

	flag.IntVar(&inputItag, "itag", -1, "itag of the stream to download")
	flag.Usage = func() {
		printUsage()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		printUsage()
		os.Exit(0)
	}

	cmd := flag.Arg(0)

	switch cmd {

	case "get":

		if flag.NArg() < 2 {
			printUsageCommand("get")
			os.Exit(1)
		}

		id := flag.Arg(1)

		i := info.NewInfo(id)

		err := i.Fetch()
		if err != nil {
			fmt.Println("there was an error fetching info for the video", id)
			os.Exit(1)
		}

		var downloadStream *info.Stream

		if inputItag != -1 {

			// see if we can satisfy the user
			for _, stream := range i.Streams() {
				if stream.Format.Itag == inputItag {
					downloadStream = stream
					break
				}
			}

			if downloadStream == nil {
				fmt.Printf("%s doesn't have a stream with format with itag %d. Picking default stream (the best one available).\n", id, inputItag)
			}
		}

		// pick the best stream by default
		if downloadStream == nil {
			downloadStream = i.Streams()[0]
		}

		var output io.WriteCloser
		if flag.NArg() > 2 {

			if flag.Arg(2) == "-" {

				output = os.Stdout

			} else {

				fmt.Println("writing download result to:", flag.Arg(2))

				output, err = os.Create(flag.Arg(2))
				if err != nil {
					fmt.Println("failed to create file")
					os.Exit(1)
				}
			}

		} else {

			ext, err := downloadStream.Format.Extension()
			if err != nil {
				// no extension, I don't have all possibilities listed in (*Format).Extension()
				ext = ""
			}

			fmt.Println("writing download result to:", i.Id+ext)

			output, err = os.Create(i.Id + ext)
			if err != nil {
				fmt.Println("failed to create file")
				os.Exit(1)
			}

		}

		download.Download(i, downloadStream, output)

		output.Close()

	case "info":

		if flag.NArg() < 2 {
			printUsageCommand("info")
			os.Exit(1)
		}

		id := flag.Arg(1)

		showVideoInfo(id)

	case "help":

		if flag.NArg() == 2 {
			printUsageCommand(flag.Arg(1))
		} else {
			printUsage()
		}

	default:
		printUsage()

	}

}

func showVideoInfo(id string) {

	i := info.NewInfo(id)

	fmt.Println("fetching video info")

	err := i.Fetch()
	if err != nil {
		fmt.Println("there was an error while fetching info for", id)
	}

	fmt.Println("available itags:")
	fmt.Println("")

	for _, stream := range i.Streams() {
		fmt.Printf("    [%d] %s - %s\n", stream.Format.Itag, stream.Format.Container, stream.Format.Video.Resolution)
	}

}

func printHeader() {
	fmt.Println("")
	fmt.Println("    aeris - (c) Gijsbrecht Hermans 2015")
	fmt.Println("")
}

func printSubHeader() {
	fmt.Println("    Easily download videos in various formats from YouTube")
	fmt.Println("")
}

func printUsage() {

	printHeader()

	printSubHeader()

	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("    aeris [flags] <command> [arguments...]")
	fmt.Println("")
	fmt.Println("A List of possible commands:")
	fmt.Println("")
	fmt.Println("    get         download a video from YouTube")
	fmt.Println("    info        fetch technical info about a video from YouTube")
	fmt.Println("    help        help with aeris")
	fmt.Println("")
	fmt.Println("To view information about the usage of these commands execute")
	fmt.Println("")
	fmt.Println("    aeris help <command>")
}

func printUsageCommand(cmd string) {

	printHeader()

	printSubHeader()

	fmt.Println("Usage:")
	fmt.Println("")

	switch cmd {
	case "get":
		fmt.Println("    aeris get [flags] <video-id> [target]")
		fmt.Println("")
		fmt.Println("    -itag        itag identifying a YouTube stream to download")

	case "info":
		fmt.Println("    aeris info <video-id>")

	case "help":
		fmt.Println("    woah! help-ception!")

	default:
		fmt.Println("command", cmd, "doesn't exist")
	}

}
