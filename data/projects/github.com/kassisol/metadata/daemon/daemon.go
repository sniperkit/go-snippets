package daemon

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/juliengk/go-utils/validation"
	"github.com/kassisol/metadata/api/server"
	"github.com/kassisol/metadata/api/storage"
	"github.com/kassisol/metadata/pkg/adf"
	"github.com/spf13/cobra"
)

/*
func serverInitConfig(appDir string) error {
	s, err := storage.NewDriver("sqlite", appDir)
	if err != nil {
		return err
	}
	defer s.End()

	return nil
}
*/

func runDaemon(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		cmd.Usage()
		os.Exit(-1)
	}

	cfg := adf.NewDaemon()
	if err := cfg.Init(); err != nil {
		log.Fatal(err)
	}

	/*
	if err := serverInitConfig(cfg.App.Dir.Root); err != nil {
		log.Fatal(err)
	}
	*/

	s, err := storage.NewDriver("sqlite", cfg.App.Dir.Root)
	if err != nil {
		log.Fatal(err)
	}
	defer s.End()

	// Input validations
	// IV - API Bind address
	if err := validation.IsValidIP(serverBindAddress); err != nil {
		log.Fatal(err)
	}

	// IV - API Port
	if err := validation.IsValidPort(serverBindPort); err != nil {
		log.Fatal(err)
	}

	addr := fmt.Sprintf("%s:%d", serverBindAddress, serverBindPort)

	server.API(addr)
}
