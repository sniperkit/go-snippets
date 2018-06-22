package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/anexia-it/go-quickstart"
)

var cmdWatch = &cobra.Command{
	Use:   "watch <prefix>",
	Short: "Watches consul KV prefix for changes",
	Run: func(cmd *cobra.Command, args []string) {
		logger := quickstart.GetLogger("watch")

		if len(args) != 1 {
			logger.Error("Required argument missing (prefix)")
			return
		}

		address, _ := cmd.Flags().GetString("address")

		consulAPI, err := quickstart.NewConsulAPI(address)
		if err != nil {
			logger.Error("Could not initialize consul API", zap.Error(err))
			return
		}

		// Handle CTRL+C
		sigChan := make(chan os.Signal, 1)
		defer close(sigChan)
		signal.Notify(sigChan, os.Interrupt)

		cancelContext, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			sig := <-sigChan
			logger.Info("Signal received", zap.Stringer("signal", sig))
			cancel()
		}()

		err = consulAPI.Watch(cancelContext, args[0], func(kvPair *api.KVPair) (err error) {
			logger.Info("KV change detected", zap.String("key", kvPair.Key), zap.ByteString("value", kvPair.Value))
			return
		})

		if err != nil {
			select {
			case <-cancelContext.Done():
			default:
				logger.Error("Watch returned error", zap.Error(err))
			}
		}
	},
}

func init() {
	cmdWatch.Flags().StringP("address", "a", "192.168.56.101:8500", "Consul address")
	cmdRoot.AddCommand(cmdWatch)
}
