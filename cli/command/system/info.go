package system

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/kassisol/metadata/api/storage"
	"github.com/kassisol/metadata/pkg/adf"
	"github.com/kassisol/metadata/version"
	"github.com/spf13/cobra"
)

func NewInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Display information about Metadata",
		Long:  infoDescription,
		Run:   runInfo,
	}

	return cmd
}

func runInfo(cmd *cobra.Command, args []string) {
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

	fmt.Println("Datas:", s.CountData())
	fmt.Println("Profiles:", s.CountProfile())
	fmt.Println("IPs:", s.CountIP())
	fmt.Println("Interfaces:", s.CountInterface())
	fmt.Println("Hosts:", s.CountHost())

	fmt.Println("Server Version:", version.Version)
	fmt.Println("Storage Driver: sqlite")
	fmt.Println("Metadata Root Dir:", cfg.App.Dir.Root)
}

var infoDescription = `
Display information about Metadata

`
