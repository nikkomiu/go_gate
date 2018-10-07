package cmd

import (
	"log"
)

const defaultPort = "3000"
const appVersion = "0.0.0"

var appPort string
var configFile string

func init() {
	rootCmd.Flags().StringVar(&configFile, "config", "config/app.yaml", "Configuration file location")
	rootCmd.Flags().StringVar(&appPort, "port", defaultPort, "Port number to start the service on")
}

// Execute the root command and delegate responsibility to all subcommands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
