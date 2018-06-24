package data

import (
	"fmt"
	"os"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils"
	"github.com/kassisol/metadata/api/storage"
	"github.com/kassisol/metadata/pkg/adf"
	"github.com/spf13/cobra"
)

var dataListFilter []string

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List whitelisted users",
		Long:    listDescription,
		Run:     runList,
	}

	flags := cmd.Flags()

	flags.StringSliceVarP(&dataListFilter, "filter", "f", []string{}, "Filter output based on conditions provided")

	return cmd
}

func runList(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	if len(args) > 0 {
		cmd.Usage()
		os.Exit(-1)
	}

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		log.Fatal(err)
	}

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	filters := utils.ConvertSliceToMap("=", dataListFilter)

	datas := s.ListData(filters)

	if len(datas) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 20, 1, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tTYPE\tVALUE\tDESCRIPTION")

		for _, data := range datas {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", data.ID, data.Name, data.Type, data.Value, data.Description)
		}

		w.Flush()
	}
}

var listDescription = `
List datas
`
