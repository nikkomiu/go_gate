package plugin

import (
	"net/http"

	"gitlab.com/nikko.miu/go_gate/pkg/settings"
)

// SetupFunc is the setup function to be called
type SetupFunc func(interface{})

// PreRequestFunc is the pre request function to be called
type PreRequestFunc func(http.ResponseWriter, *http.Request, *settings.RouteSettings) error

// PostRequestFunc is the post request function to be called
type PostRequestFunc func(http.ResponseWriter, *http.Request, *settings.RouteSettings) error

// Setup wraps the Setup func
type Setup interface {
	Setup(interface{})
}

// PreRequest wraps the PreRequest func
type PreRequest interface {
	PreRequest(http.ResponseWriter, *http.Request, *settings.RouteSettings) error
}

// PostRequest wraps the PostRequest func
type PostRequest interface {
	PostRequest(http.ResponseWriter, *http.Request, *settings.RouteSettings) error
}
