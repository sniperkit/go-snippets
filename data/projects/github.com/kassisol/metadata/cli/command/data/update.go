package data

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils"
	"github.com/kassisol/metadata/api/storage"
	"github.com/kassisol/metadata/pkg/adf"
	"github.com/spf13/cobra"
)

var (
	dataUpdateValue       string
	dataUpdateDescription string
)

func newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update [name]",
		Aliases: []string{"edit"},
		Short:   "Update data value and description",
		Long:    updateDescription,
		Run:     runUpdate,
	}

	flags := cmd.Flags()
	flags.StringVarP(&dataUpdateValue, "value", "v", "", "Update value")
	flags.StringVarP(&dataUpdateDescription, "description", "d", "", "Update description")

	return cmd
}

func runUpdate(cmd *cobra.Command, args []string) {
	defer utils.RecoverFunc()

	if len(args) < 1 || len(args) > 1 {
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

	if err := s.UpdateData(args[0], dataUpdateValue, dataUpdateDescription); err != nil {
		log.Fatal(err)
	}
}

var updateDescription = `
Update data

`
