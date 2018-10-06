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

// ErrorSettings contain information about individual errors
type ErrorSettings struct {
	Status int    `yaml:"status" json:"-"`
	Short  string `yaml:"short" json:"error"`
	Long   string `yaml:"long" json:"message"`
}

// ErrorListSettings contains the list of avaliable errors from the config
type ErrorListSettings struct {
	NotFound           *ErrorSettings `yaml:"notFound"`
	Unauthorized       *ErrorSettings `yaml:"unauthorized"`
	ServiceUnavaliable *ErrorSettings `yaml:"serviceUnavaliable"`
}

// Settings are the root configuration settings for the application
type Settings struct {
	// Non config file values (must be loaded directly)
	ConfigFile string

	Port string `yaml:"port"`

	ErrorListSettings *ErrorListSettings `yaml:"errors"`

	Auth     *AuthSettings      `yaml:"auth"`
	Routes   []*RouteSettings   `yaml:"routes"`
	Services []*ServiceSettings `yaml:"services"`
}

func getDefaultSettings() *Settings {
	return &Settings{
		Port: "3000",
		ErrorListSettings: &ErrorListSettings{
			NotFound:           &ErrorSettings{Status: 404, Short: "Not Found", Long: "Could not find route"},
			ServiceUnavaliable: &ErrorSettings{Status: 502, Short: "Could Not Process Request", Long: "The server was unable to process your request"},
			Unauthorized:       &ErrorSettings{Status: 401, Short: "Authentication Failed", Long: "Could not find or process the authentication"},
		},
	}
}

// Load will create a settings object and load the config from the settings file
func Load(configFile string) *Settings {
	settings := getDefaultSettings()
	settings.ConfigFile = configFile

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
