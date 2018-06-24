package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/datatogether/dt/core"
	"github.com/ipfs/go-datastore"
	"github.com/spf13/cobra"
)

var (
	archiveCmdUrlsFile string
)

// archiveCmd represents the export command
var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "archive one ore more urls",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && archiveCmdUrlsFile == "" {
			fmt.Println("please specify one or more urls to archive, or a file with the --file flag")
			return
		}

		urls := args
		if archiveCmdUrlsFile != "" {
			urls = []string{}
			f, err := os.Open(archiveCmdUrlsFile)
			ExitIfErr(err)

			s := bufio.NewScanner(f)
			for s.Scan() {
				urls = append(urls, s.Text())
			}
		}

		store, err := GetFilestore(false)
		ExitIfErr(err)

		ar := core.ArchiveRequests{Store: store}
		p := &core.ArchiveUrlsParams{
			Urls: urls,
		}
		path := datastore.NewKey("")

		spinner.Start()
		err = ar.ArchiveUrls(p, &path)
		spinner.Stop()
		ExitIfErr(err)

		PrintSuccess(path.String())
	},
}

func init() {
	RootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().StringVarP(&archiveCmdUrlsFile, "file", "f", "", "file of urls, one per line")
}
