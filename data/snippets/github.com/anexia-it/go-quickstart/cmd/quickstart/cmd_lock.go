package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/anexia-it/go-quickstart"
)

var cmdLock = &cobra.Command{
	Use:   "lock <key>",
	Short: "Acquires and holds lock for given key until canceled",
	Run: func(cmd *cobra.Command, args []string) {
		logger := quickstart.GetLogger("watch")

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

		// Obtain the lock...
		err = consulAPI.WithLock(args[0], func() error {
			logger.Info("Lock obtained")

			// ... and hold it until this command is canceled
			<-cancelContext.Done()
			return nil
		})

		if err != nil {
			logger.Error("Failed to obtain lock", zap.Error(err))
			return
		}
		logger.Info("Lock released")
	},
}

func init() {
	cmdLock.Flags().StringP("address", "a", "192.168.56.101:8500", "Consul address")
	cmdRoot.AddCommand(cmdLock)
}
