package alpinebits

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"slices"
)

const (
	HeaderVersion  = "X-AlpineBits-ClientProtocolVersion"
	HeaderClientID = "X-AlpineBits-ClientID"
	// HeaderAction carries the action from AlpineBits 2026-10 on.
	HeaderAction = "X-AlpineBits-Action"

	// Until AlpineBits 2024-10 the action and the request are form fields.
	FormAction  = "action"
	FormRequest = "request"
)

// Router is an HTTP handler that routes AlpineBits requests by version and
// action. It does not authenticate: wrap it in a middleware that checks the
// credentials and answers failures with Unauthorized.
type Router struct {
	versions    map[string]*versionRouter
	maxFormSize int64
	clientStore ClientStore
}

// RouterOption configures a Router.
type RouterOption func(*Router)

// WithMaxFormSize limits the request body; larger requests are answered with
// 413. The default is 32 MiB.
func WithMaxFormSize(size int64) RouterOption {
	return func(r *Router) {
		r.maxFormSize = size
	}
}

// WithClientStore makes the router remember handshakes: a client that did one
// may only use what was negotiated. Without a store, every client may use
// everything the server offers.
func WithClientStore(store ClientStore) RouterOption {
	return func(r *Router) {
		r.clientStore = store
	}
}

type versionRouter struct {
	version *Version
	actions map[string]*route
}

// serverCapabilities is keyed by handshake name; actions sharing one (OTA_Read
// and OTA_NotifReport) get the union of their capabilities.
func (vr *versionRouter) serverCapabilities() ActionCapabilities {
	caps := make(ActionCapabilities)
	for _, rt := range vr.actions {
		name := rt.action.HandshakeName()
		// OTA_Ping is implicit from 2024-10 on (AlpineBits 2024-10 §3.2).
		if name == handshakeAction.handshakeName && vr.version.Name() >= "2024-10" {
			continue
		}
		merged := caps[name]
		for _, c := range rt.capabilities {
			if !slices.Contains(merged, c) {
				merged = append(merged, c)
			}
		}
		caps[name] = merged
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
		maxFormSize: 32 << 20,
	}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Request is the typed request passed to handlers.
type Request[RQ any] struct {
	// User is the user of the basic access authentication, which AlpineBits
	// requires (§2).
	User string
	// ClientID is the optional X-AlpineBits-ClientID header.
	ClientID     string
	Data         RQ
	Capabilities []Capability
}

// HandlerFunc is a typed handler function.
type HandlerFunc[RQ, RS any] func(ctx context.Context, r Request[RQ]) (RS, error)

// request is the internal untyped request for handler storage.
type request struct {
	user         string
	clientID     string
	data         any
	capabilities []Capability
}

// handlerFunc is the internal untyped handler for storage.
type handlerFunc func(ctx context.Context, r *request) (any, error)

// HandleOption configures an action route.
type HandleOption func(*route)

// WithCapabilities sets the capabilities supported by this route.
func WithCapabilities(caps ...Capability) HandleOption {
	return func(rt *route) {
		rt.capabilities = caps
	}
}

// ServerCapabilities returns the registered capabilities for a version.
func (r *Router) ServerCapabilities(version string) ActionCapabilities {
	vr, ok := r.versions[version]
	if !ok {
		return nil
	}
	return vr.serverCapabilities()
}

// Handle registers fn for a, under the version of a. fn only sees requests
// that passed the action's rules, and rejects requests by returning
// Reject(errs...).
func (r *Router) Handle[RQ, RS any](a Action[RQ, RS], fn HandlerFunc[RQ, RS], opts ...HandleOption) *Router {
	rt := &route{
		action: a,
		handler: func(ctx context.Context, r *request) (any, error) {
			rq := *r.data.(*RQ)
			if a.validate != nil {
				if errs := a.validate(rq, r.capabilities); len(errs) > 0 {
					return a.reject(errs), nil
				}
			}
			rs, err := fn(ctx, Request[RQ]{
				User:         r.user,
				ClientID:     r.clientID,
				Data:         rq,
				Capabilities: r.capabilities,
			})
			if rej, ok := errors.AsType[*Rejection](err); ok {
				return a.reject(rej.Errors), nil
			}
			return rs, err
		},
	}
	for _, opt := range opts {
		opt(rt)
	}
	r.version(a.version).actions[a.name] = rt
	return r
}

// version returns the routes of v, registering v with its handshake action
// (OTA_Ping) on first use.
func (r *Router) version(v *Version) *versionRouter {
	if vr, ok := r.versions[v.Name()]; ok {
		if vr.version != v {
			panic("alpinebits: two versions named " + v.Name())
		}
		return vr
	}
	vr := &versionRouter{
		version: v,
		actions: make(map[string]*route),
	}
	vr.actions[handshakeAction.Name()] = &route{
		action: handshakeAction,
		handler: func(ctx context.Context, req *request) (any, error) {
			rq := req.data.(*PingRQ)
			negotiated, err := r.handshake(ctx, req.user, rq.ClientCapabilities())
			if err != nil {
				return nil, err
			}
			return rq.BuildResponse(negotiated), nil
		},
	}
	r.versions[v.Name()] = vr
	return vr
}

// latest returns the newest version whose requests have an XML body, or the
// newest of those whose requests are multipart/form-data; nil if there is none.
func (r *Router) latest(xmlBody bool) *versionRouter {
	var latest *versionRouter
	for _, vr := range r.versions {
		if vr.version.xmlBody() == xmlBody && (latest == nil || vr.version.Name() > latest.version.Name()) {
			latest = vr
		}
	}
	return latest
}

// handshake negotiates every version both sides support, whichever version
// header the OTA_Ping came with (AlpineBits §3.2), and saves the result.
func (r *Router) handshake(ctx context.Context, user string, clientCaps map[string]ActionCapabilities) (map[string]ActionCapabilities, error) {
	negotiated := make(map[string]ActionCapabilities)
	for name, vr := range r.versions {
		if announced, ok := clientCaps[name]; ok {
			negotiated[name] = Intersect(vr.serverCapabilities(), announced)
		}
	}
	if r.clientStore != nil {
		if err := r.clientStore.Save(ctx, user, negotiated); err != nil {
			return nil, err
		}
	}
	return negotiated, nil
}

// ServeHTTP implements http.Handler. A request sent as application/xml is for
// 2026-10 or later, one sent as multipart/form-data for an earlier version,
// so one URL serves both (AlpineBits 2026-10 §2.1).
//
// It answers as AlpineBits 2024-10 §2.2 says, for every version: requests that
// fail before their data is checked get 400 (405 for a wrong method, 413 for a
// body that is too large) and "ERROR:" followed by the reason, server failures
// get 500. A request that breaks the action's rules or is rejected by its
// handler gets the action's error response with 200.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method must be POST")
		return
	}
	if ce := req.Header.Get("Content-Encoding"); ce != "" && ce != "identity" {
		writeError(w, http.StatusBadRequest, "unsupported GZIP compression")
		return
	}
	mediaType, _, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
	xmlBody := mediaType == "application/xml"
	latest := r.latest(xmlBody)
	if latest == nil || !xmlBody && mediaType != "multipart/form-data" {
		writeError(w, http.StatusBadRequest, "unsupported content type")
		return
	}

	req.Body = http.MaxBytesReader(w, req.Body, r.maxFormSize)
	var actionName string
	var payload []byte
	var err error
	// For multipart, maxMemory equal to the body limit keeps the form out of
	// temporary files.
	if xmlBody {
		actionName = req.Header.Get(HeaderAction)
		payload, err = io.ReadAll(req.Body)
	} else if err = req.ParseMultipartForm(r.maxFormSize); err == nil {
		actionName = req.PostFormValue(FormAction)
		payload = []byte(req.PostFormValue(FormRequest))
	}
	if err != nil {
		if _, ok := errors.AsType[*http.MaxBytesError](err); ok {
			writeError(w, http.StatusRequestEntityTooLarge, "request too large")
			return
		}
		writeError(w, http.StatusBadRequest, "unknown or missing action")
		return
	}

	versionName := req.Header.Get(HeaderVersion)
	vr, ok := r.versions[versionName]
	switch {
	case ok && vr.version.xmlBody() != xmlBody:
		writeError(w, http.StatusBadRequest, "content type does not match the client protocol version")
		return
	case ok:
	case actionName == handshakeAction.Name():
		// The handshake negotiates all common versions whatever the header
		// says (AlpineBits 2022-10 §3.2).
		vr = latest
	case versionName == "":
		writeError(w, http.StatusBadRequest, "no valid client protocol version provided")
		return
	default:
		writeError(w, http.StatusBadRequest, "your current alpinebits version does not match one of the servers supported versions")
		return
	}
	route, ok := vr.actions[actionName]
	if !ok {
		writeError(w, http.StatusBadRequest, "unknown or missing action")
		return
	}
	if len(payload) == 0 {
		writeError(w, http.StatusBadRequest, "missing request")
		return
	}

	err = vr.version.Validate(payload)
	var data any
	if err == nil {
		data, err = route.action.unmarshal(payload)
	}
	if err != nil {
		writeError(w, http.StatusBadRequest, "XML validation error: "+err.Error())
		return
	}

	user, _, _ := req.BasicAuth()
	// Without a handshake, a client may use everything the server offers.
	caps := route.capabilities
	if r.clientStore != nil && actionName != handshakeAction.Name() {
		negotiated, ok, err := r.clientStore.Load(req.Context(), user)
		if err != nil {
			fail(w, req, actionName, err)
			return
		}
		if ok {
			actionCaps, ok := negotiated[vr.version.Name()][route.action.HandshakeName()]
			if !ok {
				writeError(w, http.StatusBadRequest, "action not negotiated")
				return
			}
			caps = actionCaps
		}
	}

	resp, err := route.handler(req.Context(), &request{
		user:         user,
		clientID:     req.Header.Get(HeaderClientID),
		data:         data,
		capabilities: caps,
	})
	if err != nil {
		fail(w, req, actionName, err)
		return
	}
	write(w, req, vr, actionName, resp)
}

func write(w http.ResponseWriter, req *http.Request, vr *versionRouter, actionName string, resp any) {
	respXML, err := xml.Marshal(resp)
	if err == nil {
		err = vr.version.Validate(respXML)
	}
	if err != nil {
		fail(w, req, actionName, fmt.Errorf("invalid response: %w", err))
		return
	}
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.Write(respXML)
}

// fail answers 500 without exposing err to the client, and logs it.
func fail(w http.ResponseWriter, req *http.Request, actionName string, err error) {
	slog.ErrorContext(req.Context(), "alpinebits: request failed", "action", actionName, "err", err)
	writeError(w, http.StatusInternalServerError, "internal server error")
}

// Unauthorized answers a request whose credentials or X-AlpineBits-ClientID
// the server does not accept, as AlpineBits §2 requires. AlpineBits 2024-10
// §2.2.1 suggests the reason "invalid or missing username/password".
func Unauthorized(w http.ResponseWriter, reason string) {
	writeError(w, http.StatusUnauthorized, reason)
}

func writeError(w http.ResponseWriter, status int, reason string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	io.WriteString(w, "ERROR:"+reason)
}
