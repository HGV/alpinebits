package alpinebits

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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
func Send[RQ, RS any](c *Client, ctx context.Context, a Action[RQ, RS], rq RQ) (*RS, error) {
	reqXML, err := xml.Marshal(rq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	if err := c.version.Validate(string(reqXML)); err != nil {
		return nil, fmt.Errorf("validate request: %w", err)
	}

	body, contentType, err := c.buildMultipartBody(a.Name(), reqXML)
	if err != nil {
		return nil, fmt.Errorf("build request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set(HeaderVersion, c.version.Name())
	if c.clientID != "" {
		req.Header.Set(HeaderClientID, c.clientID)
	}
	req.SetBasicAuth(c.username, c.password)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server error %d: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if err := c.version.Validate(string(respBody)); err != nil {
		return nil, fmt.Errorf("validate response: %w", err)
	}

	var rs RS
	if err := xml.Unmarshal(respBody, &rs); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return &rs, nil
}

// Handshake performs capability negotiation with the server.
func (c *Client) Handshake(ctx context.Context, caps ActionCapabilities) (ActionCapabilities, error) {
	rq := PingRQ{
		Version:  OTAVersion,
		EchoData: EchoData{Message: encodeCapabilities(caps)},
	}

	rs, err := Send(c, ctx, handshakeAction, rq)
	if err != nil {
		return nil, err
	}

	// Parse negotiated capabilities from response
	negotiated, err := parseCapabilities(rs.Warnings.Message)
	if err != nil {
		return nil, fmt.Errorf("parse negotiated capabilities: %w", err)
	}

	return negotiated, nil
}

func (c *Client) buildMultipartBody(action string, reqXML []byte) (*bytes.Buffer, string, error) {
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
