package main

import (
	"encoding/json"
	"os"

	human "github.com/anexia-it/go-human"
	quickstart "github.com/anexia-it/go-quickstart"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// outputEncoder defines an interface implemented both by go-human and encoding/json's encoder types
type outputEncoder interface {
	Encode(o interface{}) error
}

var cmdDecodeRequest = &cobra.Command{
	Use:   "decode-request <file>",
	Short: "Decodes a request",
	Run: func(cmd *cobra.Command, args []string) {
		logger := quickstart.GetLogger("decode-request")

		// Check if the required argument is present
		if len(args) != 1 {
			logger.Error("Exactly one argument required (file)")
			return
		}

		// Open the file
		f, err := os.Open(args[0])

		// Handle error during file open
		if err != nil {
			logger.Error("Could not open file", zap.Error(err))
			return
		}

		// Ensure that the file handle is closed when leaving this scope
		defer f.Close()

		// Create a JSON decoder
		dec := json.NewDecoder(f)

		// Allocate a new Request struct
		req := &quickstart.Request{}

		// Decode file onto req
		if err := dec.Decode(req); err != nil {
			logger.Error("Could not decode file", zap.Error(err))
			return
		}

		// Retrieve the value of the "json" flag
		wantJSON, _ := cmd.Flags().GetBool("json")

		// Dynamically create our outputEncoder instance
		var enc outputEncoder
		if wantJSON {
			// If JSON was desired, create a JSON decoder...
			jsonEncoder := json.NewEncoder(cmd.OutOrStdout())
			// ... set its indent, so we get pretty-looking output ...
			jsonEncoder.SetIndent("", "  ")
			// ... and set enc to the instance of jsonEncoder
			enc = jsonEncoder
		} else {
			// If JSON was not desired, fall back to our "human" encoder
			enc, _ = human.NewEncoder(cmd.OutOrStdout())
		}

		// Encode the object
		if err := enc.Encode(req); err != nil {
			logger.Error("Encoding failed", zap.Error(err))
		}
	},
}

func init() {
	cmdDecodeRequest.Flags().BoolP("json", "j", false, "Enforces JSON output")
	cmdRoot.AddCommand(cmdDecodeRequest)
}
