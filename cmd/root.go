package cmd

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/spf13/cobra"

	gatePlugin "gitlab.com/nikko.miu/go_gate/pkg/plugin"
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
	if appPort != defaultPort {
		appSettings.Port = appPort
	}

	// Setup Plugins
	gatePlugin.Load(appSettings.Plugins)

	// Setup route handler
	mux := http.NewServeMux()
	routeContext := route.New(appSettings)

	mux.HandleFunc("/", routeContext.ServiceHandler())

	loggedMux := handlers.CombinedLoggingHandler(os.Stdout, mux)

	// Start the server
	log.Printf("Starting server on port %s", appSettings.Port)
	if err := http.ListenAndServe(":"+appSettings.Port, loggedMux); err != nil {
		log.Fatal(err)
	}
}
