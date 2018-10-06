package settings

import (
	"io/ioutil"
	"log"
	"net/url"

	yaml "gopkg.in/yaml.v2"
)

// AuthSettings for dealing with Auth0 tokens
type AuthSettings struct {
	Domain  string `yaml:"domain"`
	JWKSURL string `yaml:"jwksUrl"`
}

// ServiceSettings contain all of the rules for services
type ServiceSettings struct {
	Name          string `yaml:"name"`
	BaseURLString string `yaml:"url"`

	BaseURL *url.URL
}

// RouteSettings contain route settings inside of a service
type RouteSettings struct {
	Path         string `yaml:"path"`
	ServiceName  string `yaml:"service"`
	StripPrefix  string `yaml:"stripPrefix"`
	OptionalAuth bool   `yaml:"optionalAuth"`
}

// Settings are the root configuration settings for the application
type Settings struct {
	// Non config file values (must be loaded directly)
	ConfigFile string

	Port string `yaml:"port"`

	Auth     *AuthSettings      `yaml:"auth"`
	Routes   []*RouteSettings   `yaml:"routes"`
	Services []*ServiceSettings `yaml:"services"`
}

// Load will create a settings object and load the config from the settings file
func Load(configFile string) *Settings {
	settings := &Settings{
		ConfigFile: configFile,
		Port:       "3000",
	}

	yamlFile, err := ioutil.ReadFile(settings.ConfigFile)
	if err != nil {
		log.Fatalf("Could not read app config: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, settings)
	if err != nil {
		log.Fatalf("Could not parse app config YAML: %v", err)
	}

	settings.loadServiceURLs()

	return settings
}

func (settings *Settings) loadServiceURLs() {
	for _, svc := range settings.Services {
		url, err := url.Parse(svc.BaseURLString)
		if err != nil {
			log.Fatal(err)
		}

		svc.BaseURL = url
	}
}
