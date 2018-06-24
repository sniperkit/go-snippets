package nic

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils"
	"github.com/kassisol/metadata/api/storage"
	"github.com/kassisol/metadata/pkg/adf"
	"github.com/spf13/cobra"
)

var nicUpdateType string
var nicUpdateValue string

func newUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [mac]",
		Short: "Update interface IP or floating IP",
		Long:  updateDescription,
		Run:   runUpdate,
	}

	flags := cmd.Flags()
	flags.StringVarP(&nicUpdateType, "type", "t", "ip", "Type")
	flags.StringVarP(&nicUpdateValue, "value", "v", "", "Value")

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

	validTypes := []string{
		"ip",
		"floatingip",
	}

	if !utils.StringInSlice(nicUpdateType, validTypes, false) {
		log.Fatalf("Type '%s' is not valid", nicUpdateType)
	}

	if err := s.UpdateInterface(args[0], nicUpdateType, nicUpdateValue); err != nil {
		log.Fatal(err)
	}
}

var updateDescription = `
Update interface IP or floating IP

`
