package commands

import (
	"github.com/codegangsta/cli"
)

// RegisterDownload adds the "get" command to the command line interface.
func RegisterDownload(app *cli.App) {
	app.Commands = append(app.Commands, cli.Command{
		Name:    "download",
		Aliases: []string{"d"},
		Usage:   "download packages from a repository",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "mirror, m",
				Value:  "http://archive.ubuntu.com/ubuntu",
				Usage:  "use the specified mirror",
				EnvVar: "DEBTOOL_MIRROR",
			},
			cli.StringFlag{
				Name:   "dist, d",
				Value:  "xenial",
				Usage:  "get packages from specified distribution",
				EnvVar: "DEBTOOL_DISTRIBUTION",
			},
		},
	})
}

// DownloadDo downloads a package.
func DownloadDo(ctx *cli.Context) {
	//
}
