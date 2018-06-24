package ip

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils"
	"github.com/kassisol/metadata/api/storage"
	"github.com/kassisol/metadata/pkg/adf"
	"github.com/spf13/cobra"
)

var (
	ipAddNetmask string
	ipAddGateway string
)

func newAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add [ip]",
		Short: "Add IP address",
		Long:  addDescription,
		Run:   runAdd,
	}

	flags := cmd.Flags()
	flags.StringVarP(&ipAddNetmask, "netmask", "n", "", "Netmask")
	flags.StringVarP(&ipAddGateway, "gateway", "g", "", "Gateway")

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

	if err := s.AddIP(args[0], ipAddNetmask, ipAddGateway); err != nil {
		log.Fatal(err)
	}
}

var addDescription = `
Add IP address

`
