package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/anexia-it/go-quickstart"
)

var cmdWorker = &cobra.Command{
	Use:   "worker",
	Short: "Runs a URL check worker",
	Run: func(cmd *cobra.Command, args []string) {
		logger := quickstart.GetLogger("worker")

		address, _ := cmd.Flags().GetString("address")

		consulAPI, err := quickstart.NewConsulAPI(address)
		if err != nil {
			logger.Error("Could not initialize consul API", zap.Error(err))
			return
		}

		// Prepare an HTTP client which does not follow redirects
		httpClient := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
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

		// Wait for the result to be written...
		err = consulAPI.Watch(cancelContext, "requests/", func(kvPair *api.KVPair) (err error) {
			// Obtain the request ID
			requestID := strings.TrimPrefix(kvPair.Key, "requests/")
			logger.Info("Request detected", zap.String("id", requestID))

			// Try obtaining the corresponding request lock
			consulAPI.WithLock(kvPair.Key, func() error {
				logger.Info("Request lock obtained")

				// Delete the request key
				if _, err := consulAPI.KV().Delete(kvPair.Key, nil); err != nil {
					logger.Warn("Failed to delete request key", zap.Error(err))
					// We mask this error
					return nil
				}

				// Decode the request
				req := &quickstart.Request{}

				if err := json.Unmarshal(kvPair.Value, req); err != nil {
					logger.Warn("Failed to unmarshal the request", zap.Error(err))
					return err
				}

				// Prepare the result struct
				result := &quickstart.Result{
					URL: req.URL,
				}

				// Execute the HTTP GET call
				resp, err := httpClient.Get(req.URL)
				if err == nil {
					// No error, record status code
					defer resp.Body.Close()
					result.StatusCode = resp.StatusCode
					logger.Info("Status code retrieved", zap.String("url", req.URL), zap.Int("status_code", resp.StatusCode))
				} else {
					// Error detected
					result.Error = err.Error()
					logger.Warn("Error during check", zap.Error(err))
				}

				// Encode the result
				resultBytes, err := json.Marshal(result)
				if err != nil {
					logger.Error("Failed to marshal result", zap.Error(err))
					return err
				}

				resultKVPair := &api.KVPair{
					Key:   "results/" + requestID,
					Value: resultBytes,
				}

				// Write the result to consul
				if _, err := consulAPI.KV().Put(resultKVPair, nil); err != nil {
					logger.Error("Failed to write result to consul", zap.Error(err))
					return err
				}
				return nil

			})
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
	cmdWorker.Flags().StringP("address", "a", "192.168.56.101:8500", "Consul address")
	cmdRoot.AddCommand(cmdWorker)
}
