package cmd

import (
	"log"
)

var appPort string
var configFile string

var appVersion = "0.0.0"

func init() {
	rootCmd.Flags().StringVar(&configFile, "config", "config/app.yaml", "Configuration file location (default is config/app.yaml)")
	rootCmd.Flags().StringVar(&appPort, "port", "", "Port number to start the service on")
}

// Execute the root command and delegate responsibility to all subcommands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
