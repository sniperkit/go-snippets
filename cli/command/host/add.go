package host

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils"
	"github.com/kassisol/metadata/api/storage"
	"github.com/kassisol/metadata/pkg/adf"
	"github.com/spf13/cobra"
)

var (
	hostAddEnable    bool
	hostAddFQDN      string
	hostAddProfile   string
	hostAddInterface []string
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add host",
		Long:  addDescription,
		Run:   runAdd,
	}

	flags := cmd.Flags()
	flags.BoolVarP(&hostAddEnable, "enable", "e", true, "Enable host entry")
	flags.StringVarP(&hostAddFQDN, "fqdn", "f", "", "Set FQDN")
	flags.StringVarP(&hostAddProfile, "profile", "p", "", "Set profile")
	flags.StringSliceVarP(&hostAddInterface, "interface", "i", []string{}, "Set interfaces")

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

	if err := s.AddHost(hostAddEnable, args[0], hostAddFQDN, hostAddProfile, hostAddInterface); err != nil {
		log.Fatal(err)
	}
}

var addDescription = `
Add host

`
