package cmd

import (
	"log"
	"net/http"

	"github.com/spf13/cobra"

	"gitlab.com/nikko.miu/go_gate/pkg/auth"
	"gitlab.com/nikko.miu/go_gate/pkg/settings"
	"gitlab.com/nikko.miu/go_gate/route"
)

const rootLong = `Primarily responsible for routing traffic to
backing services and providing user authentication for routes.`

var rootCmd = &cobra.Command{
	Use:     "go_gate",
	Short:   "Gateway service serving as entrypoint and auth handoff for microservice applications",
	Long:    rootLong,
	Run:     runRoot,
	Version: appVersion,
}

func runRoot(cmd *cobra.Command, args []string) {
	// Load the settings
	appSettings := settings.Load(configFile)

	// Override settings as needed
	if appPort != "" {
		appSettings.Port = appPort
	}

	// Setup Auth
	auth.Setup(appSettings.Auth)

	// Setup route handler
	mux := http.NewServeMux()
	routeContext := route.New(appSettings)

	mux.HandleFunc("/", routeContext.ServiceHandler())

	// Start the server
	log.Printf("Starting server on port %s", appSettings.Port)
	if err := http.ListenAndServe(":"+appSettings.Port, mux); err != nil {
		log.Panic(err)
	}
}
