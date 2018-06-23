package main

import (
	"github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/anexia-it/go-quickstart"
)

var cmdLockedPut = &cobra.Command{
	Use:   "locked-put <key> <value>",
	Short: "Writes consul KV key while holding the corresponding lock",
	Run: func(cmd *cobra.Command, args []string) {
		logger := quickstart.GetLogger("locked-put")

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

		logger.Info("Obtaining lock and writing key...", zap.String("key", args[0]))
		err = consulAPI.WithLock(args[0], func() (err error) {
			_, err = consulAPI.KV().Put(kvPair, nil)
			if err != nil {
				logger.Error("Could not write consul key", zap.Error(err))
			}
			return
		})

		if err != nil {
			logger.Error("Locked write failed", zap.Error(err))
			return
		}
		logger.Info("Key written", zap.String("key", args[0]), zap.String("value", args[1]))
	},
}

func init() {
	cmdLockedPut.Flags().StringP("address", "a", "192.168.56.101:8500", "Consul address")
	cmdRoot.AddCommand(cmdLockedPut)
}
