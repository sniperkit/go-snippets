package main

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/anexia-it/go-quickstart"
)

var cmdDelete = &cobra.Command{
	Use:   "delete <key>",
	Short: "Deletes consul KV key",
	Run: func(cmd *cobra.Command, args []string) {
		logger := quickstart.GetLogger("delete")

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

		_, err = consulAPI.KV().Delete(args[0], nil)
		if err != nil {
			logger.Error("Could not delete consul KV key", zap.Error(err))
			return
		}
		logger.Info("Key deleted", zap.String("key", args[0]))
	},
}

func init() {
	cmdDelete.Flags().StringP("address", "a", "192.168.56.101:8500", "Consul address")
	cmdRoot.AddCommand(cmdDelete)
}
