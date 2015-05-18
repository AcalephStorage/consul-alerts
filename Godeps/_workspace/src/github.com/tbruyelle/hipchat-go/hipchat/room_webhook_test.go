package hipchat

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestWebhookList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/room/1/webhook", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprintf(w, `
		{
			"items":[
			  {"name":"a", "pattern":"a", "event":"message_received", "url":"h", "id":1, "links":{"self":"s"}},
				{"name":"b", "pattern":"b", "event":"message_received", "url":"h", "id":2, "links":{"self":"s"}}
			],
			"links":{"self":"s", "prev":"a", "next":"b"},
			"startIndex":0,
			"maxResults":10
		}`)
	})

	want := &WebhookList{
		Webhooks: []Webhook{
			Webhook{
				Name:         "a",
				Pattern:      "a",
				Event:        "message_received",
				URL:          "h",
				ID:           1,
				WebhookLinks: WebhookLinks{Links: Links{Self: "s"}},
			},
			Webhook{
				Name:         "b",
				Pattern:      "b",
				Event:        "message_received",
				URL:          "h",
				ID:           2,
				WebhookLinks: WebhookLinks{Links: Links{Self: "s"}},
			},
		},
		StartIndex: 0,
		MaxResults: 10,
		Links:      PageLinks{Links: Links{Self: "s"}, Prev: "a", Next: "b"},
	}

	reqParams := &ListWebhooksRequest{}

	actual, _, err := client.Room.ListWebhooks("1", reqParams)
	if err != nil {
		t.Fatalf("Room.ListWebhooks returns an error %v", err)
	}
	if !reflect.DeepEqual(want, actual) {
		t.Errorf("Room.ListWebhooks returned %+v, want %+v", actual, want)
	}
}

func TestWebhookDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/room/1/webhook/2", func(w http.ResponseWriter, r *http.Request) {
		if m := "DELETE"; m != r.Method {
			t.Errorf("Request method %s, want %s", r.Method, m)
		}
	})

	_, err := client.Room.DeleteWebhook("1", "2")
	if err != nil {
		t.Fatalf("Room.Update returns an error %v", err)
	}
}
