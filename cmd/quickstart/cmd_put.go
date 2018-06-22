package main

import (
	"github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/anexia-it/go-quickstart"
)

var cmdPut = &cobra.Command{
	Use:   "put <key> <value>",
	Short: "Writes consul KV key",
	Run: func(cmd *cobra.Command, args []string) {
		logger := quickstart.GetLogger("put")

		if len(args) != 2 {
			logger.Error("Required arguments missing (key value)")
			return
		}

		address, _ := cmd.Flags().GetString("address")

		consulAPI, err := quickstart.NewConsulAPI(address)
		if err != nil {
			logger.Error("Could not initialize consul API", zap.Error(err))
			return
		}

		kvPair := &api.KVPair{
			Key:   args[0],
			Value: []byte(args[1]),
		}

		_, err = consulAPI.KV().Put(kvPair, nil)
		if err != nil {
			logger.Error("Could not retrieve consul KV key", zap.Error(err))
			return
		}
		logger.Info("Key written", zap.String("key", args[0]), zap.String("value", args[1]))
	},
}

func init() {
	cmdPut.Flags().StringP("address", "a", "192.168.56.101:8500", "Consul address")
	cmdRoot.AddCommand(cmdPut)
}
