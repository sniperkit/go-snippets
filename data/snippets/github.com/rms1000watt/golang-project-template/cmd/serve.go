package cmd

import (
	"github.com/rms1000watt/golang-project-template/serve"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Start the server",
	Long:    `Start the server`,
	Example: `./golang-project-template serve`,
	Run:     serveFunc,
}

var serveCfg serve.Config

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVar(&serveCfg.Host, "host", "127.0.0.1", "Host to listen on")
	serveCmd.Flags().IntVar(&serveCfg.Port, "port", 7100, "Port to listen on")

	serveCmd.Flags().StringVar(&serveCfg.CertsPath, "certs-path", "certs", "Certs path where the Private/Public key live")
	serveCmd.Flags().StringVar(&serveCfg.CertName, "cert-name", "server.crt", "Public key name")
	serveCmd.Flags().StringVar(&serveCfg.KeyName, "key-name", "server.key", "Private key name")

	setFlagsFromEnv(serveCmd)
}

func serveFunc(cmd *cobra.Command, args []string) {
	configureLogging()

	serve.Serve(serveCfg)
}
