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
	dataAddType        string
	dataAddValue       string
	dataAddDescription string
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add data",
		Long:  addDescription,
		Run:   runAdd,
	}

	flags := cmd.Flags()
	flags.StringVarP(&dataAddType, "type", "t", "", "Set type")
	flags.StringVarP(&dataAddValue, "value", "v", "", "Set value")
	flags.StringVarP(&dataAddDescription, "description", "d", "", "Set description")

	return cmd
}

func runAdd(cmd *cobra.Command, args []string) {
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

	if err := s.AddData(args[0], dataAddType, dataAddValue, dataAddDescription); err != nil {
		log.Fatal(err)
	}
}

var addDescription = `
Add data

`
