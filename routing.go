package alpinebits

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

	handshakeDataFromRouter func() HandshakeData
}

func (r Request) HandshakeData() HandshakeData {
	if rctx, ok := RouteContextFrom(r.Context); ok {
		return r.handshakeDataFromRouter().Intersect(rctx.HandshakeDataOverride)
	}
	return r.handshakeDataFromRouter()
}

type HandlerFunc func(r Request) (any, error)

type Route struct {
	action               version.Action
	handler              HandlerFunc
	capabilities         []string
	excludeFromHandshake bool
}

type RouteFunc func(*Route)

func WithCapabilities[C ~string](caps ...C) RouteFunc {
	return func(r *Route) {
		for _, c := range caps {
			r.capabilities = append(r.capabilities, string(c))
		}
	}
}

func WithExcludeFromHandshake() RouteFunc {
	return func(r *Route) {
		r.excludeFromHandshake = true
	}
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		preconditionErrorf(w, "expected http method POST, got %s", r.Method)
		return
	}

	clientID := r.Header.Get(HeaderClientID)
	if clientID == "" {
		preconditionErrorf(w, "missing http header: %s", HeaderClientID)
		return
	}

	requestedVersion := r.Header.Get(HeaderClientProtocolVersion)
	if requestedVersion == "" {
		preconditionErrorf(w, "missing http header: %s", HeaderClientProtocolVersion)
		return
	}
	if err := version.ValidateVersionString(requestedVersion); err != nil {
		preconditionError(w, err.Error())
		return
	}

	routes, ok := router.versionRoutes[requestedVersion]
	if !ok {
		supportedVersions := slices.Collect(maps.Keys(router.versionRoutes))
		preconditionErrorf(w,
			"your current version of '%s' does not match one of the servers' supported versions: %s",
			requestedVersion,
			strings.Join(supportedVersions, ", "))
		return
	}

	rctx, hasRouteCtx := RouteContextFrom(r.Context())

	if hasRouteCtx {
		// Check if the requested version is disabled by a handshake override
		if _, ok := rctx.HandshakeDataOverride[requestedVersion]; !ok {
			preconditionErrorf(w,
				"your current version of '%s' was not included in the handshake agreement. Please retry the handshake to ensure compatibility.",
				requestedVersion)
			return
		}
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

	if hasRouteCtx {
		// Check if the action is disabled by a handshake override
		if _, ok := rctx.HandshakeDataOverride[requestedVersion][route.action.HandshakeName()]; !ok {
			preconditionError(w, "unknown or missing action")
			return
		}
	}

	payload := r.Form.Get("request")
	if err := routes.version.ValidateXML(payload); err != nil {
		preconditionErrorf(w,
			"XML validation error for action %s\n\n%s",
			requestedAction,
			err.Error())
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
		handshakeDataFromRouter: func() HandshakeData {
			return NewHandshakeDataFromRouter(*router)
		},
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

func preconditionErrorf(w http.ResponseWriter, msg string, a ...any) {
	preconditionError(w, fmt.Sprintf(msg, a...))
}

func internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	slog.ErrorContext(r.Context(), err.Error())
	http.Error(w, "", http.StatusInternalServerError)
}

type RouteContext struct {
	HandshakeDataOverride HandshakeData
}

type routeContextKey struct{}

func WithRouteContext(ctx context.Context, rctx RouteContext) context.Context {
	return context.WithValue(ctx, routeContextKey{}, rctx)
}

func RouteContextFrom(ctx context.Context) (RouteContext, bool) {
	rctx, ok := ctx.Value(routeContextKey{}).(RouteContext)
	return rctx, ok
}
