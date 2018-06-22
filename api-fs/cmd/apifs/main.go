package main

import (
	"flag"
	"log"
	"os"

	"github.com/gophergala/api-fs/filesystem"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

var mountpoint string

func init() {
	flag.StringVar(&mountpoint, "mountpoint", "", "mount point for apifs")
}

func main() {
	flag.Parse()

	if mountpoint == "" {
		flag.Usage()
		os.Exit(2)
	}

	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("apifs"),
		fuse.Subtype("apifs"),
		fuse.LocalVolume(),
		fuse.VolumeName("API FS"),
	)

	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	if err = fs.Serve(c, filesystem.NewFS()); err != nil {
		log.Fatal(err)
	}

	<-c.Ready
	if err = c.MountError; err != nil {
		log.Fatal(err)
	}
}
