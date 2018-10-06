package route

import (
	"encoding/json"
	"io"
	"log"
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
			log.Printf("[ WARNING ] No route for '%s'\n", r.URL.Path)
			errorResponse(w, ctx.ErrorListSettings.NotFound)
			return
		}

		// Get the service for the route
		svcURL := services[foundRoute.ServiceName]
		if svcURL == nil {
			log.Printf("[  ERROR  ] Service '%s' not found\n", foundRoute.ServiceName)
			errorResponse(w, ctx.ErrorListSettings.NotFound)
			return
		}

		// Validate auth
		_, err := auth.Validate(r.Header.Get("Authorization"), !foundRoute.OptionalAuth)
		if err != nil {
			errorResponse(w, ctx.ErrorListSettings.Unauthorized)
			return
		}

		// Build the backend URL path
		u, _ := url.Parse(strings.TrimPrefix(r.URL.Path, foundRoute.StripPrefix))

		// Build the request
		req, _ := http.NewRequest(r.Method, svcURL.ResolveReference(u).String(), r.Body)
		req.Header = r.Header // TODO: Allow blacklisting inbound headers

		// Send the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			errorResponse(w, ctx.ErrorListSettings.ServiceUnavaliable)
			return
		}

		// Close the backend response when done
		defer resp.Body.Close()

		// Respond with backend response
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

func errorResponse(w http.ResponseWriter, resp *settings.ErrorSettings) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)

	json.NewEncoder(w).Encode(resp)
}
