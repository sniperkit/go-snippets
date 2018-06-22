package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:   "golang-project-template",
	Short: "Golang Project Template done right",
	Long: `Golang Project Template done right
	
	Starting point Golang project with Gometaliter, Govendor,
	TDD, Precommit git hooks, security best practices, etc.`,
}

var logLevel string

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Set log level (debug, info, warn, error, fatal)")

	setPFlagsFromEnv(rootCmd)
}

// Execute gets called from the main project directory for convenience
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func setPFlagsFromEnv(cmd *cobra.Command) {
	// Courtesy of https://github.com/coreos/pkg/blob/master/flagutil/env.go
	cmd.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		key := strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
		if val := os.Getenv(key); val != "" {
			if err := cmd.PersistentFlags().Set(f.Name, val); err != nil {
				fmt.Println("Failed setting flag from env:", err)
			}
		}
	})
}

func setFlagsFromEnv(cmd *cobra.Command) {
	// Courtesy of https://github.com/coreos/pkg/blob/master/flagutil/env.go
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		key := strings.ToUpper(strings.Replace(f.Name, "-", "_", -1))
		if val := os.Getenv(key); val != "" {
			if err := cmd.Flags().Set(f.Name, val); err != nil {
				fmt.Println("Failed setting flag from env:", err)
			}
		}
	})
}

func configureLogging() {
	if level, err := log.ParseLevel(logLevel); err != nil {
		log.Error("log-level argument malformed: ", logLevel, ": ", err)
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}
}
