package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type (
	HandshakeClient struct {
		config *HandshakeClientConfig
		client *http.Client
	}
	HandshakeClientConfig struct {
		URL           string
		Username      string
		Password      string
		ClientID      string
		HandshakeData HandshakeData
	}
)

type pingRQ struct {
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_PingRQ"`
	Version  string   `xml:"Version,attr"`
	EchoData string   `xml:",innerxml"`
}

type pingRS struct {
	XMLName  xml.Name `xml:"http://www.opentravel.org/OTA/2003/05 OTA_PingRS"`
	Version  string   `xml:"Version,attr"`
	Success  success  `xml:"Success"`
	Warning  warning  `xml:"Warnings>Warning"`
	EchoData string   `xml:",innerxml"`
}

type success struct{}

type status string

const statusAlpinebitsHandshake status = "ALPINEBITS_HANDSHAKE"

type warning struct {
	Type         int    `xml:"Type,attr"`
	Status       status `xml:"Status,attr"`
	Intersection string `xml:",innerxml"`
}

func (c *HandshakeClientConfig) validate() error {
	if _, err := url.Parse(c.URL); err != nil {
		return err
	}

	if strings.TrimSpace(c.Username) == "" {
		return errors.New("c.Username is empty")
	}

	if strings.TrimSpace(c.Password) == "" {
		return errors.New("c.Password is empty")
	}

	if strings.TrimSpace(c.ClientID) == "" {
		return errors.New("c.ClientID is empty")
	}

	if len(c.HandshakeData) == 0 {
		return errors.New("c.HandshakeData is empty")
	}

	return nil
}

func NewHandshakeClient(config HandshakeClientConfig) (*HandshakeClient, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	return &HandshakeClient{
		config: &config,
		client: &http.Client{},
	}, nil
}

func (c *HandshakeClient) Ping(ctx context.Context) (HandshakeData, *http.Response, error) {
	echoData, err := json.Marshal(c.config.HandshakeData)
	if err != nil {
		return nil, nil, err
	}

	pingRQ := pingRQ{
		Version:  "1.0",
		EchoData: string(echoData),
	}

	req, err := c.newRequest(ctx, "OTA_Ping:Handshaking", pingRQ)
	if err != nil {
		return nil, nil, err
	}

	var resp *http.Response
	var lastErr error
	for version := range c.config.HandshakeData {
		req.Header.Set(HeaderClientProtocolVersion, version)

		var pingRS pingRS
		resp, lastErr = c.do(req, &pingRS)
		if lastErr != nil {
			continue // retry with lower version
		}

		if pingRS.Warning.Status == statusAlpinebitsHandshake {
			var handshakeData HandshakeData
			if lastErr = json.Unmarshal([]byte(pingRS.Warning.Intersection), &handshakeData); lastErr != nil {
				continue // retry with lower version
			}
			return handshakeData, resp, nil
		}

		lastErr = errors.New("no possible version found")
	}

	return nil, resp, lastErr
}

func (c *HandshakeClient) newRequest(ctx context.Context, action string, request any) (*http.Request, error) {
	xml, err := xml.Marshal(request)
	if err != nil {
		return nil, err
	}

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	if err := w.WriteField("action", action); err != nil {
		return nil, err
	}
	if err := w.WriteField("request", string(xml)); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.URL, &body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.config.Username, c.config.Password)

	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set(HeaderClientID, c.config.ClientID)

	return req, nil
}

func (c *HandshakeClient) do(req *http.Request, v any) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, err
	}

	if sc := resp.StatusCode; sc < 200 || sc > 299 {
		return resp, fmt.Errorf("handshake request failed with status code: %d", sc)
	}

	if err = xml.Unmarshal(body, v); err != nil {
		return resp, err
	}

	return resp, nil
}
