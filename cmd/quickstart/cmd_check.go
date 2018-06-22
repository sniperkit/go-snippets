package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"

	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/anexia-it/go-quickstart"
)

var cmdCheck = &cobra.Command{
	Use:   "check <url>",
	Short: "Instructs workers to check the given URL",
	Run: func(cmd *cobra.Command, args []string) {
		logger := quickstart.GetLogger("check")

		if len(args) != 1 {
			logger.Error("Required argument missing (url)")
			return
		}

		address, _ := cmd.Flags().GetString("address")

		consulAPI, err := quickstart.NewConsulAPI(address)
		if err != nil {
			logger.Error("Could not initialize consul API", zap.Error(err))
			return
		}

		// Build a new request ID and prepare the consul keys
		requestID := uuid.NewV4().String()
		requestKey := "requests/" + requestID
		resultKey := "results/" + requestID

		// Encode the request
		req := &quickstart.Request{
			URL: args[0],
		}

		reqBytes, err := json.Marshal(req)
		if err != nil {
			logger.Error("Failed to encode request", zap.Error(err))
			return
		}

		// Send the request by writing the KV key
		kvPair := &api.KVPair{
			Key:   requestKey,
			Value: reqBytes,
		}

		_, err = consulAPI.KV().Put(kvPair, nil)
		if err != nil {
			logger.Error("Could not write request to consul", zap.Error(err))
			return
		}
		logger.Info("Request written", zap.String("url", req.URL))

		// Start watching for the result
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

		// Wait for the result to be written...
		err = consulAPI.Watch(cancelContext, resultKey, func(kvPair *api.KVPair) (err error) {
			// Decode the result
			res := &quickstart.Result{}

			if err = json.Unmarshal(kvPair.Value, res); err != nil {
				logger.Error("Failed to unmarshal result", zap.Error(err))
				return
			}

			// Cancel the watcher
			defer cancel()

			// check if the worker reported an error
			if res.Error != "" {
				logger.Error("Worker reported error", zap.String("error", res.Error))
				return
			}

			// Result received: print it
			logger.Info("Result received", zap.String("url", res.URL), zap.Int("status_code", res.StatusCode))

			// Delete the result again
			if _, err = consulAPI.KV().Delete(resultKey, nil); err != nil {
				logger.Warn("Failed to remove result key", zap.Error(err))
			}
			return
		})

		if err != nil {
			select {
			case <-cancelContext.Done():
			default:
				logger.Error("Watch returned error", zap.Error(err))
			}
		}

		return
	},
}

func init() {
	cmdCheck.Flags().StringP("address", "a", "192.168.56.101:8500", "Consul address")
	cmdRoot.AddCommand(cmdCheck)
}
