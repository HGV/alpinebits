package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"slices"
	"strings"

	"github.com/HGV/alpinebits/version"
)

const (
	HeaderServerAcceptEncoding  = "X-AlpineBits-Server-Accept-Encoding"
	HeaderClientID              = "X-AlpineBits-ClientID"
	HeaderClientProtocolVersion = "X-AlpineBits-ClientProtocolVersion"
)

type Router struct {
	http.Handler

	versionRoutes map[string]Routes
}

type Routes struct {
	version      version.Version[version.Action]
	actionRoutes map[string]Route
}

func NewRouter() *Router {
	return &Router{
		versionRoutes: make(map[string]Routes),
	}
}

func (r *Router) Version(version version.Version[version.Action], fn func(s *Subrouter)) *Router {
	subrouter := newSubrouter()
	fn(subrouter)
	r.versionRoutes[version.String()] = Routes{
		version:      version,
		actionRoutes: subrouter.actionRoutes,
	}
	return r
}

type Subrouter struct {
	actionRoutes map[string]Route
}

func newSubrouter() *Subrouter {
	return &Subrouter{
		actionRoutes: make(map[string]Route),
	}
}

func (s *Subrouter) Action(action version.Action, handlerFn HandlerFunc, opts ...RouteFunc) {
	route := Route{
		handler: handlerFn,
		action:  action,
	}
	for _, opt := range opts {
		opt(&route)
	}
	s.actionRoutes[action.String()] = route
}

type Request struct {
	Context      context.Context
	ClientID     string
	Data         any
	Capabilities []string
}

type HandlerFunc func(r Request) (any, error)

type Route struct {
	action       version.Action
	handler      HandlerFunc
	capabilities []string
}

type RouteFunc func(*Route)

func WithCapabilities[C ~string](caps ...C) RouteFunc {
	return func(r *Route) {
		for _, c := range caps {
			r.capabilities = append(r.capabilities, string(c))
		}
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		preconditionError(w, fmt.Sprintf("expected http method POST, got %s", r.Method))
		return
	}

	clientID := r.Header.Get(HeaderClientID)
	if clientID == "" {
		preconditionError(w, fmt.Sprintf("missing http header: %s", HeaderClientID))
		return
	}

	requestedVersion := r.Header.Get(HeaderClientProtocolVersion)
	if requestedVersion == "" {
		preconditionError(w, fmt.Sprintf("missing http header: %s", HeaderClientProtocolVersion))
		return
	}
	if err := version.ValidateVersionString(requestedVersion); err != nil {
		preconditionError(w, err.Error())
		return
	}

	routes, ok := router.versionRoutes[requestedVersion]
	if !ok {
		supportedVersions := slices.Collect(maps.Keys(router.versionRoutes))
		preconditionError(w,
			fmt.Sprintf("your current alpinebits version of '%s' does not match one of the servers' supported versions: %s",
				requestedVersion,
				strings.Join(supportedVersions, ", ")),
		)
	}

	if err := r.ParseMultipartForm(1 << 20); err != nil {
		preconditionError(w, err.Error())
		return
	}

	requestedAction := r.Form.Get("action")
	route, ok := routes.actionRoutes[requestedAction]
	if !ok {
		preconditionError(w, "unknown or missing action")
		return
	}

	payload := r.Form.Get("request")
	if err := routes.version.ValidateXML(payload); err != nil {
		preconditionError(w,
			fmt.Sprintf("XML validation error for action %s\n\n%s",
				requestedAction,
				err.Error()))
		return
	}

	data, err := route.action.Unmarshal([]byte(payload))
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	req := Request{
		Context:      r.Context(),
		ClientID:     clientID,
		Data:         data,
		Capabilities: route.capabilities,
	}
	resp, err := route.handler(req)
	if err != nil {
		internalServerError(w, r, err)
		return
	}

	b, err := xml.Marshal(resp)
	if err != nil {
		internalServerError(w, r, err)
		return
	}
	if err = routes.version.ValidateXML(string(b)); err != nil {
		internalServerError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write(b)
}

func preconditionError(w http.ResponseWriter, msg string) {
	http.Error(w, fmt.Sprintf("ERROR: %s", msg), http.StatusBadRequest)
}

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	slog.ErrorContext(r.Context(), err.Error())
	http.Error(w, "", http.StatusInternalServerError)
}
