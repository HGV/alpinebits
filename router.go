package alpinebits

import (
	"context"
	"encoding/xml"
	"net/http"
)

const (
	HeaderVersion  = "X-AlpineBits-ClientProtocolVersion"
	HeaderClientID = "X-AlpineBits-ClientID"

	FormAction  = "action"
	FormRequest = "request"
)

// Router is an HTTP handler that routes AlpineBits requests by version and action.
type Router struct {
	versions    map[string]*versionRouter
	maxFormSize int64
	clientStore ClientStore
}

// RouterOption configures a Router.
type RouterOption func(*Router)

// WithMaxFormSize sets the maximum size for multipart form data.
func WithMaxFormSize(size int64) RouterOption {
	return func(r *Router) {
		r.maxFormSize = size
	}
}

// WithClientStore sets the client store for capability negotiation.
func WithClientStore(store ClientStore) RouterOption {
	return func(r *Router) {
		r.clientStore = store
	}
}

type versionRouter struct {
	version *Version
	actions map[string]*route
}

// serverCapabilities returns all registered action capabilities for this version.
func (vr *versionRouter) serverCapabilities() ActionCapabilities {
	caps := make(ActionCapabilities)
	for name, route := range vr.actions {
		caps[name] = route.capabilities
	}
	return caps
}

type route struct {
	action       action
	handler      handlerFunc
	capabilities []Capability
}

// NewRouter creates a new AlpineBits router.
func NewRouter(opts ...RouterOption) *Router {
	r := &Router{
		versions:    make(map[string]*versionRouter),
		maxFormSize: 32 << 20, // 32MB default
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Request is the typed request passed to handlers.
type Request[RQ any] struct {
	Context      context.Context
	ClientID     string
	Data         RQ
	Capabilities []Capability
}

// HandlerFunc is a typed handler function.
type HandlerFunc[RQ, RS any] func(r Request[RQ]) (RS, error)

// request is the internal untyped request for handler storage.
type request struct {
	ctx          context.Context
	clientID     string
	data         any
	capabilities []Capability
}

// handlerFunc is the internal untyped handler for storage.
type handlerFunc func(r *request) (any, error)

// RouteOption configures an action route within a version.
type RouteOption interface {
	apply(*versionRouter)
}

type routeOption struct {
	action       action
	handler      handlerFunc
	capabilities []Capability
}

func (o *routeOption) apply(vr *versionRouter) {
	// // Prevent manual registration of handshake action
	// if o.action.Name() == handshakeAction.Name() {
	// 	panic("alpinebits: cannot manually register handshake action - it is auto-registered by the router")
	// }
	vr.actions[o.action.Name()] = &route{
		action:       o.action,
		handler:      o.handler,
		capabilities: o.capabilities,
	}
}

// HandleOption configures an action handler.
type HandleOption func(*routeOption)

// WithCapabilities sets the capabilities supported by this route.
func WithCapabilities(caps ...Capability) HandleOption {
	return func(o *routeOption) {
		o.capabilities = caps
	}
}

// Handle creates a RouteOption that registers a typed handler for an action.
func Handle[RQ, RS any](a Action[RQ, RS], fn HandlerFunc[RQ, RS], opts ...HandleOption) RouteOption {
	handler := func(r *request) (any, error) {
		rq := r.data.(*RQ)
		return fn(Request[RQ]{
			Context:      r.ctx,
			ClientID:     r.clientID,
			Data:         *rq,
			Capabilities: r.capabilities,
		})
	}
	ro := &routeOption{
		action:  a,
		handler: handler,
	}
	for _, opt := range opts {
		opt(ro)
	}
	return ro
}

// ServerCapabilities returns the registered capabilities for a version.
func (r *Router) ServerCapabilities(version string) ActionCapabilities {
	vr, ok := r.versions[version]
	if !ok {
		return nil
	}
	return vr.serverCapabilities()
}

// Version registers a version with its action handlers.
// The handshake action (OTA_Ping) is automatically registered.
func (r *Router) Version(v *Version, opts ...RouteOption) *Router {
	vr := &versionRouter{
		version: v,
		actions: make(map[string]*route),
	}

	// Auto-register handshake action with closure handler
	vr.actions[handshakeAction.Name()] = &route{
		action: handshakeAction,
		handler: func(req *request) (any, error) {
			rq := req.data.(*PingRQ)

			clientCaps := rq.ClientCapabilities()
			serverCaps := vr.serverCapabilities()
			negotiated := Intersect(serverCaps, clientCaps)

			if r.clientStore != nil {
				r.clientStore.Set(req.clientID, v.Name(), negotiated)
			}

			return rq.BuildResponse(negotiated), nil
		},
	}

	for _, opt := range opts {
		opt.apply(vr)
	}
	r.versions[v.Name()] = vr
	return r
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	versionName := req.Header.Get(HeaderVersion)
	if versionName == "" {
		http.Error(w, "missing "+HeaderVersion+" header", http.StatusBadRequest)
		return
	}

	vr, ok := r.versions[versionName]
	if !ok {
		http.Error(w, "unsupported version: "+versionName, http.StatusBadRequest)
		return
	}

	if err := req.ParseMultipartForm(r.maxFormSize); err != nil {
		http.Error(w, "failed to parse multipart form", http.StatusBadRequest)
		return
	}

	actionName := req.FormValue(FormAction)
	if actionName == "" {
		http.Error(w, "missing "+FormAction+" form field", http.StatusBadRequest)
		return
	}

	xmlPayload := req.FormValue(FormRequest)
	if xmlPayload == "" {
		http.Error(w, "missing "+FormRequest+" form field", http.StatusBadRequest)
		return
	}

	route, ok := vr.actions[actionName]
	if !ok {
		http.Error(w, "unsupported action: "+actionName, http.StatusBadRequest)
		return
	}

	if err := vr.version.Validate(xmlPayload); err != nil {
		http.Error(w, "XML validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	data, err := route.action.Unmarshal([]byte(xmlPayload))
	if err != nil {
		http.Error(w, "failed to unmarshal request: "+err.Error(), http.StatusBadRequest)
		return
	}

	clientID := req.Header.Get(HeaderClientID)

	// Determine capabilities for this request
	var caps []Capability
	if r.clientStore != nil {
		// Use negotiated capabilities if handshake was performed
		negotiated := r.clientStore.Get(clientID, versionName)
		if actionCaps, ok := negotiated[actionName]; ok {
			caps = actionCaps
		} else if len(negotiated) > 0 {
			// Client negotiated but not this action
			http.Error(w, "action not negotiated: "+actionName, http.StatusForbidden)
			return
		} else {
			// No negotiation yet - use server capabilities (first request or no handshake required)
			caps = route.capabilities
		}
	} else {
		// No client store - use server capabilities directly
		caps = route.capabilities
	}

	handlerReq := &request{
		ctx:          req.Context(),
		clientID:     clientID,
		data:         data,
		capabilities: caps,
	}

	resp, err := route.handler(handlerReq)
	if err != nil {
		http.Error(w, "handler error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respXML, err := xml.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to marshal response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := vr.version.Validate(string(respXML)); err != nil {
		http.Error(w, "response XML validation failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write(respXML)
}
