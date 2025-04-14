package v_2020_10

import (
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/HGV/alpinebits/v_2020_10/common"
	"github.com/HGV/alpinebits/v_2020_10/freerooms"
	"github.com/HGV/alpinebits/v_2020_10/guestrequests"
	"github.com/HGV/alpinebits/v_2020_10/inventory"
	"github.com/HGV/alpinebits/v_2020_10/rateplans"
	"github.com/HGV/alpinebits/version"
	"github.com/HGV/x"
)

type (
	Client struct {
		config *ClientConfig
		client *http.Client
	}
	ClientConfig struct {
		URL               string
		Username          string
		Password          string
		ClientID          string
		Version           version.Version[version.Action]
		NegotiatedVersion map[string][]string
		HttpClient        *http.Client
	}
	ClientResponse[RS any] struct {
		*http.Response

		Data          *RS
		SendInventory bool
		SendFreeRooms bool
		SendRatePlans bool
	}
)

func NewClient(config ClientConfig) (*Client, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}

	return &Client{
		config: &config,
		client: x.If(config.HttpClient != nil, config.HttpClient, &http.Client{}),
	}, nil
}

func (c *ClientConfig) validate() error {
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

	if c.Version == nil {
		return errors.New("c.Version is empty")
	}

	if err := version.ValidateVersionString(c.Version.String()); err != nil {
		return err
	}

	if len(c.NegotiatedVersion) == 0 {
		return errors.New("c.NegotiatedVersion is empty")
	}

	return nil
}

func (c *Client) PushHotelInvCountNotif(ctx context.Context, r freerooms.HotelInvCountNotifRQ) (*ClientResponse[freerooms.HotelInvCountNotifRS], error) {
	return sendRequest[freerooms.HotelInvCountNotifRS](ctx, c, ActionHotelInvCountNotif, r)
}

func (c *Client) PullGuestRequests(ctx context.Context, r guestrequests.ReadRQ) (*ClientResponse[guestrequests.ResRetrieveRS], error) {
	return sendRequest[guestrequests.ResRetrieveRS](ctx, c, ActionReadGuestRequests, r)
}

func (c *Client) PushAcknowledgement(ctx context.Context, r guestrequests.NotifReportRQ) (*ClientResponse[guestrequests.NotifReportRS], error) {
	return sendRequest[guestrequests.NotifReportRS](ctx, c, ActionNotifReportGuestRequests, r)
}

func (c *Client) PushHotelDescriptiveContentNotif(ctx context.Context, r inventory.HotelDescriptiveContentNotifRQ) (*ClientResponse[inventory.HotelDescriptiveContentNotifRS], error) {
	return sendRequest[inventory.HotelDescriptiveContentNotifRS](ctx, c, ActionHotelDescriptiveContentNotifInventory, r)
}

func (c *Client) PushRatePlans(ctx context.Context, r rateplans.HotelRatePlanNotifRQ) (*ClientResponse[rateplans.HotelRatePlanNotifRS], error) {
	return sendRequest[rateplans.HotelRatePlanNotifRS](ctx, c, ActionHotelRatePlanNotifRatePlans, r)
}

func sendRequest[RS any, RQ any](ctx context.Context, c *Client, action Action, rq RQ) (*ClientResponse[RS], error) {
	req, err := c.newRequest(ctx, action, rq)
	if err != nil {
		return nil, err
	}

	var rs RS
	resp, err := c.do(req, &rs)
	if err != nil {
		return nil, err
	}

	return newClientResponse(resp, &rs)
}

func (c *Client) newRequest(ctx context.Context, action Action, request any) (*http.Request, error) {
	if _, ok := c.config.NegotiatedVersion[action.HandshakeName()]; !ok {
		return nil, errors.New("unsupported action")
	}

	xml, err := xml.Marshal(request)
	if err != nil {
		return nil, err
	}

	if err = c.config.Version.ValidateXML(string(xml)); err != nil {
		return nil, err
	}

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	if err := w.WriteField("action", action.String()); err != nil {
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
	req.Header.Set("X-AlpineBits-ClientID", c.config.ClientID)
	req.Header.Set("X-AlpineBits-ClientProtocolVersion", c.config.Version.String())

	return req, nil
}

func (c *Client) do(req *http.Request, v any) (*http.Response, error) {
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
		return nil, fmt.Errorf("request failed with status code: %d", sc)
	}

	if err = c.config.Version.ValidateXML(string(body)); err != nil {
		return nil, fmt.Errorf("xml validation failed: %w", err)
	}

	if err = xml.Unmarshal(body, v); err != nil {
		return resp, err
	}

	return resp, nil
}

func newClientResponse[T any](r *http.Response, v *T) (*ClientResponse[T], error) {
	response := &ClientResponse[T]{Response: r, Data: v}
	response.populateCompleteSetRequests(v)
	return response, nil
}

func (r *ClientResponse[T]) populateCompleteSetRequests(v any) {
	if rs, ok := v.(common.Response); ok {
		var statuses []common.Status

		if rs.Errors != nil {
			for _, e := range *rs.Errors {
				statuses = append(statuses, e.Status)
			}
		}

		if rs.Warnings != nil {
			for _, w := range *rs.Warnings {
				statuses = append(statuses, w.Status)
			}
		}

		for _, status := range statuses {
			switch status {
			case common.StatusSendInventory:
				r.SendInventory = true
			case common.StatusSendFreeRooms:
				r.SendFreeRooms = true
			case common.StatusSendRatePlans:
				r.SendRatePlans = true
			}
		}
	}
}
