package route

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"gitlab.com/nikko.miu/go_gate/pkg/auth"
	"gitlab.com/nikko.miu/go_gate/pkg/settings"
)

// RequestContext is the context that routing logic maintains
type RequestContext struct {
	*settings.Settings
}

type errorMessage struct {
	Status  int    `json:"-"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

var notFoundError = &errorMessage{
	Status:  http.StatusNotFound,
	Error:   "Not Found",
	Message: "Could not find route",
}

var unauthorizedError = &errorMessage{
	Status:  http.StatusUnauthorized,
	Error:   "Authentication Failed",
	Message: "Could not find or process the authorization header",
}

var badGatewayError = &errorMessage{
	Status:  http.StatusBadGateway,
	Error:   "Could Not Process Request",
	Message: "The server was unavaliable and could not process your request",
}

// New creates new request context object and sets up routing
func New(settings *settings.Settings) *RequestContext {
	return &RequestContext{Settings: settings}
}

// ServiceHandler handles all routing for backing services
func (ctx *RequestContext) ServiceHandler() http.HandlerFunc {
	// Map services for faster access
	services := make(map[string]*url.URL)

	for _, svc := range ctx.Settings.Services {
		services[svc.Name] = svc.BaseURL
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Find matching route
		foundRoute := ctx.findRoute(r.URL.Path)
		if foundRoute == nil {
			errorResponse(w, notFoundError)
			return
		}

		svcURL := services[foundRoute.ServiceName]

		// Validate auth
		_, err := auth.Validate(r.Header.Get("Authorization"), foundRoute.OptionalAuth)
		if err != nil {
			errorResponse(w, unauthorizedError)
			return
		}

		// Get backend the URL
		u, _ := url.Parse(strings.TrimPrefix(r.URL.Path, foundRoute.StripPrefix))

		client := &http.Client{}
		resp, err := client.Do(&http.Request{
			Method: r.Method,
			URL:    svcURL.ResolveReference(u),
			Header: r.Header, // TODO: Allow blacklisting inbound headers
		})
		if err != nil {
			errorResponse(w, badGatewayError)
			return
		}

		defer resp.Body.Close()

		r.Header = resp.Header // TODO: Allow blacklisting outbound headers
		io.Copy(w, resp.Body)
	}
}

func (ctx *RequestContext) findRoute(path string) *settings.RouteSettings {
	var r *settings.RouteSettings

	for _, route := range ctx.Settings.Routes {
		if match, _ := regexp.MatchString(route.Path, path); match {
			r = route
			break
		}
	}

	return r
}

func errorResponse(w http.ResponseWriter, resp *errorMessage) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)

	json.NewEncoder(w).Encode(resp)
}
