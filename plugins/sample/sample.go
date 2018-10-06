package main

import (
	"log"
	"net/http"

	"gitlab.com/nikko.miu/go_gate/pkg/settings"
)

type samplePlugin struct{}

// Setup the plugin
func (*samplePlugin) Setup(settings interface{}) {
	log.Println("Setting up sample plugin...")
}

// PreRequest handler for the plugin
func (*samplePlugin) PreRequest(w http.ResponseWriter, r *http.Request, route *settings.RouteSettings) error {
	log.Printf("Running sample pre request for %s\n", r.URL.Path)

	return nil
}

// PostRequest handler for the plugin
func (*samplePlugin) PostRequest(w http.ResponseWriter, r *http.Request, route *settings.RouteSettings) error {
	log.Printf("Running sample post request for %s\n", r.URL.Path)

	return nil
}

func main() {
	log.Fatal("[ ERROR ] This module is not meant to be loaded directly!")
}

// Plugin sample
var Plugin samplePlugin
