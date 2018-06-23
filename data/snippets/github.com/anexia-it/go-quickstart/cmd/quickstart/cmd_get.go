package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/anexia-it/go-quickstart"
)

var cmdGet = &cobra.Command{
	Use:   "get <key>",
	Short: "Retrieves consul KV key",
	Run: func(cmd *cobra.Command, args []string) {
		logger := quickstart.GetLogger("get")

		if len(args) != 1 {
			logger.Error("Required argument missing (key)")
			return
		}

		address, _ := cmd.Flags().GetString("address")

		consulAPI, err := quickstart.NewConsulAPI(address)
		if err != nil {
			logger.Error("Could not initialize consul API", zap.Error(err))
			return
		}

		kvPair, _, err := consulAPI.KV().Get(args[0], nil)
		if err != nil {
			logger.Error("Could not retrieve consul KV key", zap.Error(err))
			return
		} else if kvPair == nil {
			logger.Warn("Key not present in consul", zap.String("key", args[0]))
			return
		}
		cmd.Printf("%s=%s\n", args[0], kvPair.Value)
	},
}

func init() {
	cmdGet.Flags().StringP("address", "a", "192.168.56.101:8500", "Consul address")
	cmdRoot.AddCommand(cmdGet)
}
