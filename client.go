package alpinebits

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// Client is an AlpineBits HTTP client.
type Client struct {
	baseURL    string
	version    *Version
	username   string
	password   string
	clientID   string
	httpClient *http.Client
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithClientID sets the client ID header.
func WithClientID(id string) ClientOption {
	return func(c *Client) {
		c.clientID = id
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(hc *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = hc
	}
}

// NewClient creates a new AlpineBits client.
func NewClient(baseURL, username, password string, v *Version, opts ...ClientOption) *Client {
	c := &Client{
		baseURL:    baseURL,
		version:    v,
		username:   username,
		password:   password,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Send sends a request and returns the response.
func (c *Client) Send[RQ, RS any](ctx context.Context, a Action[RQ, RS], rq RQ) (*RS, error) {
	// The handshake action has no version and is valid for every client.
	if a.version != nil && a.version != c.version {
		return nil, fmt.Errorf("alpinebits: action %s is of version %s, client uses version %s",
			a.name, a.version.Name(), c.version.Name())
	}

	reqXML, err := xml.Marshal(rq)
	if err != nil {
		return nil, fmt.Errorf("alpinebits: marshal request: %w", err)
	}

	if err := c.version.Validate(reqXML); err != nil {
		return nil, fmt.Errorf("alpinebits: validate request: %w", err)
	}

	body, contentType := io.Reader(bytes.NewReader(reqXML)), "application/xml"
	if !c.version.xmlBody() {
		body, contentType, err = multipartBody(a.Name(), reqXML)
		if err != nil {
			return nil, fmt.Errorf("alpinebits: build request body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, body)
	if err != nil {
		return nil, fmt.Errorf("alpinebits: create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set(HeaderVersion, c.version.Name())
	if c.version.xmlBody() {
		req.Header.Set(HeaderAction, a.Name())
	}
	if c.clientID != "" {
		req.Header.Set(HeaderClientID, c.clientID)
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("alpinebits: send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("alpinebits: read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK || bytes.HasPrefix(respBody, []byte("ERROR:")) {
		msg := strings.TrimPrefix(strings.TrimSpace(string(respBody)), "ERROR:")
		return nil, &ServerError{StatusCode: resp.StatusCode, Message: msg}
	}

	if err := c.version.Validate(respBody); err != nil {
		return nil, fmt.Errorf("alpinebits: validate response: %w", err)
	}

	var rs RS
	if err := xml.Unmarshal(respBody, &rs); err != nil {
		return nil, fmt.Errorf("alpinebits: unmarshal response: %w", err)
	}

	return &rs, nil
}

// ServerError is a server's answer that is not an AlpineBits response: an HTTP
// error status or, with status 200, an "ERROR:" string (AlpineBits §2, §3.4).
// AlpineBits asks clients to retry on status 500.
type ServerError struct {
	StatusCode int
	Message    string
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("alpinebits: server error (status %d): %s", e.StatusCode, e.Message)
}

// Handshake announces caps (keyed by handshake action name) for the client's
// version and returns what the server agreed to for that version.
func (c *Client) Handshake(ctx context.Context, caps ActionCapabilities) (ActionCapabilities, error) {
	version := c.version.Name()
	rq := PingRQ{
		Version:  OTAVersion,
		EchoData: EchoData{Message: encodeCapabilities(map[string]ActionCapabilities{version: caps})},
	}

	rs, err := c.Send(ctx, handshakeAction, rq)
	if err != nil {
		return nil, err
	}

	negotiated, err := parseCapabilities(rs.Warnings.Message)
	if err != nil {
		return nil, fmt.Errorf("alpinebits: parse negotiated capabilities: %w", err)
	}
	actions, ok := negotiated[version]
	if !ok {
		return nil, fmt.Errorf("alpinebits: server does not support version %s", version)
	}
	return actions, nil
}

func multipartBody(action string, reqXML []byte) (*bytes.Buffer, string, error) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	if err := w.WriteField(FormAction, action); err != nil {
		return nil, "", err
	}

	if err := w.WriteField(FormRequest, string(reqXML)); err != nil {
		return nil, "", err
	}

	if err := w.Close(); err != nil {
		return nil, "", err
	}

	return &body, w.FormDataContentType(), nil
}
