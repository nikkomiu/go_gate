package main

import (
	"log"
	"net/http"

	"gitlab.com/nikko.miu/go_gate/pkg/settings"
)

type headBlockPlugin struct{}

var inboundBlocked []interface{}
var outboundBlocked []interface{}

// Setup the plugin
func (*headBlockPlugin) Setup(settings interface{}) {
	blockSettings := settings.(map[interface{}]interface{})

	inbound, ok := blockSettings["inbound"].([]interface{})
	if !ok {
		log.Fatal("Could not process inbound blocked headers list")
	}

	outbound, ok := blockSettings["outbound"].([]interface{})
	if !ok {
		log.Fatal("Could not process outbound blocked headers list")
	}

	inboundBlocked = inbound
	outboundBlocked = outbound
}

// PreRequest handler for the plugin
func (*headBlockPlugin) PreRequest(w http.ResponseWriter, r *http.Request, route *settings.RouteSettings) error {
	for _, header := range inboundBlocked {
		strHeader := header.(string)

		r.Header.Del(strHeader)
	}

	return nil
}

// PostRequest handler for the plugin
func (*headBlockPlugin) PostRequest(w http.ResponseWriter, r *http.Request, route *settings.RouteSettings) error {
	for _, header := range outboundBlocked {
		strHeader := header.(string)

		w.Header().Del(strHeader)
	}

	return nil
}

func main() {
	log.Fatal("[ ERROR ] This module is not meant to be loaded directly!")
}

// Plugin sample
var Plugin headBlockPlugin
