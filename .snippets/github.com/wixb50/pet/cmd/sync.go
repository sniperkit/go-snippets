package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"pet/config"
	petSync "pet/sync"
)

// syncCmd represents the sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync snippets",
	Long:  `Sync snippets with alioss`,
	RunE:  sync,
}

func sync(cmd *cobra.Command, args []string) (err error) {
	if config.Conf.AliOSS.AccessID == "" || config.Conf.AliOSS.AccessKey == "" {
		return fmt.Errorf(`access_id or access_key is empty.
Go https://oss.console.aliyun.com/index and create access_id or access_key (only need "AliOSS" scope).
Write access_id or access_key in config file (pet configure).
		`)
	}
	if config.Conf.AliOSS.BucketName == "" || config.Conf.AliOSS.Endpoint == "" {
		return fmt.Errorf(`bucket_name or endpoint is empty.
Go https://oss.console.aliyun.com/index and create bucket_name or endpoint (only need "AliOSS" scope).
Write bucket_name or endpoint in config file (pet configure).
		`)
	}

	if config.Flag.Force {
		return petSync.ForceSync()
	}
	return petSync.AutoSync(config.Conf.General.SnippetFile)
}

func init() {
	RootCmd.AddCommand(syncCmd)
	syncCmd.Flags().BoolVarP(&config.Flag.Force, "force", "f", false,
		`Force sync remote snippets to local`)
}
