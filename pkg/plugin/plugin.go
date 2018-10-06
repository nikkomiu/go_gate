package plugin

import (
	"log"
	"net/http"
	"plugin"

	"gitlab.com/nikko.miu/go_gate/pkg/settings"
)

var setupPlugins []SetupFunc
var preRequestPlugins []PreRequestFunc
var postRequestPlugins []PostRequestFunc

// Load the plugin
func Load(plugins []*settings.PluginSettings) {
	for _, configuredPlugin := range plugins {
		// Read the plugin file
		p, err := plugin.Open(configuredPlugin.Path)
		if err != nil {
			log.Printf("[ WARNING ] Could not load plugin '%s' %s", configuredPlugin.Path, err)
			continue
		}

		// Find the plugin export in the plugin library
		symPlugin, err := p.Lookup("Plugin")
		if err != nil {
			log.Printf("[   INFO  ] Could not find exported setup for plugin %s\n", configuredPlugin.Path)
			continue
		}

		// Load the Setup
		plugSetup, ok := symPlugin.(Setup)
		if ok {
			plugSetup.Setup(configuredPlugin.Settings)
		}

		// Load the PreRequest
		plugPre, ok := symPlugin.(PreRequest)
		if ok {
			preRequestPlugins = append(preRequestPlugins, plugPre.PreRequest)
		}

		// Load the PostRequest
		plugPost, ok := symPlugin.(PostRequest)
		if ok {
			postRequestPlugins = append(postRequestPlugins, plugPost.PostRequest)
		}
	}
}

// HandlePreReq handles all pre request plugins that are configured
func HandlePreReq(plugins []*settings.PluginSettings, w http.ResponseWriter, r *http.Request, route *settings.RouteSettings) error {
	for _, configuredPlugin := range preRequestPlugins {
		err := configuredPlugin(w, r, route)
		if err != nil {
			return err
		}
	}

	return nil
}

// HandlePostReq handles all post request plugins that are configured
func HandlePostReq(plugins []*settings.PluginSettings, w http.ResponseWriter, r *http.Request, route *settings.RouteSettings) error {
	for _, configuredPlugin := range postRequestPlugins {
		err := configuredPlugin(w, r, route)
		if err != nil {
			return err
		}
	}

	return nil
}
