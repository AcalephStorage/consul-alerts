// handling Webhook data

package hipchat

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Response Types

type WebhookLinks struct {
	Links
}

type Webhook struct {
	WebhookLinks WebhookLinks `json:"links"`
	Name         string       `json:"name"`
	Event        string       `json:"event"`
	Pattern      string       `json:"pattern"`
	URL          string       `json:"url"`
	ID           int          `json:"id,omitempty"`
}

type WebhookList struct {
	Webhooks   []Webhook `json:"items"`
	StartIndex int       `json:"startIndex"`
	MaxResults int       `json:"maxResults"`
	Links      PageLinks `json:"links"`
}

// Request Types

type ListWebhooksRequest struct {
	MaxResults int `json:"max-results"`
	StartIndex int `json:"start-index"`
}

type CreateWebhookRequest struct {
	Name    string `json:"name"`
	Event   string `json:"event"`
	Pattern string `json:"pattern"`
	URL     string `json:"url"`
}

// ListWebhooks returns all the webhooks for a given room.
//
// HipChat API docs: https://www.hipchat.com/docs/apiv2/method/get_all_webhooks
func (r *RoomService) ListWebhooks(id interface{}, roomReq *ListWebhooksRequest) (*WebhookList, *http.Response, error) {
	u := fmt.Sprintf("room/%v/webhook", id)
	if roomReq != nil {
		p := url.Values{}
		if roomReq.MaxResults != 0 {
			p.Add("max-results", strconv.FormatInt(int64(roomReq.MaxResults), 10))
		}
		if roomReq.StartIndex != 0 {
			p.Add("start-index", strconv.FormatInt(int64(roomReq.StartIndex), 10))
		}
		if len(p) > 0 {
			u += "?" + p.Encode()
		}
	}
	req, err := r.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	whList := new(WebhookList)

	resp, err := r.client.Do(req, whList)
	if err != nil {
		return nil, resp, err
	}
	return whList, resp, nil
}

// DeleteWebhook removes the given webhook.
//
// HipChat API docs: https://www.hipchat.com/docs/apiv2/method/delete_webhook
func (r *RoomService) DeleteWebhook(id interface{}, webhookId interface{}) (*http.Response, error) {
	req, err := r.client.NewRequest("DELETE", fmt.Sprintf("room/%v/webhook/%v", id, webhookId), nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req, nil)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

// CreateWebhook creates a new webhook.
//
// HipChat API docs: https://www.hipchat.com/docs/apiv2/method/create_webhook
func (r *RoomService) CreateWebhook(id interface{}, roomReq *CreateWebhookRequest) (*Webhook, *http.Response, error) {
	req, err := r.client.NewRequest("POST", fmt.Sprintf("room/%v/webhook", id), roomReq)
	if err != nil {
		return nil, nil, err
	}

	wh := new(Webhook)

	resp, err := r.client.Do(req, wh)
	if err != nil {
		return nil, resp, err
	}

	return wh, resp, nil
}
