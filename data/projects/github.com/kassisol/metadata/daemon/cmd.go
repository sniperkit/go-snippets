package daemon

import (
	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/metadata/pkg/adf"
	"github.com/spf13/cobra"
)

var (
	serverBindAddress string
	serverBindPort    int
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metadatad",
		Short: "Starts a metadata server",
		Long:  daemonDescription,
		Run:   runDaemon,
	}

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		log.Fatal(err)
	}

	flags := cmd.Flags()
	flags.StringVarP(&serverBindAddress, "bind-address", "a", "0.0.0.0", "Bind Address")
	flags.IntVarP(&serverBindPort, "bind-port", "p", 80, "Bind Port")

	return cmd
}

var daemonDescription = `
Starts a metadata server

`
