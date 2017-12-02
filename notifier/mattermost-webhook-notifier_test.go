package notifier

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestMattermostWebhookJsonUnmarshalling(t *testing.T) {
	cluster_name := "my_cluster"
	channel := "test_channel"
	username := "test_username"
	icon_url := "test_icon_url"
	text := "test_text"

	expectedNotifier := MattermostWebhookNotifier{
		ClusterName: cluster_name,
		Channel:     channel,
		Username:    username,
		IconUrl:     icon_url,
		Text:        text,
	}
	var unmarshalledNotifier MattermostWebhookNotifier

	data := []byte(fmt.Sprintf(`{
    "cluster_name": "%s",
    "channel": "%s",
    "username": "%s",
    "icon_url": "%s",
    "text": "%s"
  }`, cluster_name, channel, username, icon_url, text))

	fmt.Printf("%s\n", data)
	if err := json.Unmarshal(data, &unmarshalledNotifier); err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(expectedNotifier, unmarshalledNotifier) {
		t.Fatalf("Expected mattermostWebhookNotifier to be %s, got %s\n", expectedNotifier, unmarshalledNotifier)
	}
}

func TestMattermostWebhookPost(t *testing.T) {
	cluster_name := "my_cluster"
	channel := "test_channel"
	username := "test_username"
	icon_url := "test_icon_url"
	text := "test_text"
	enabled := true

	expectedValues := url.Values{}
	expectedValues.Set("payload", fmt.Sprintf(`{"cluster_name":"%s","channel":"%s","username":"%s","icon_url":"%s","text":"%s","enabled":%t}`, cluster_name, channel, username, icon_url, text, enabled))
	expectedPayload := []byte(expectedValues.Encode())

	var actualPayload []byte

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		actualPayload, err = ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}))
	defer ts.Close()

	n := MattermostWebhookNotifier{
		Url:         ts.URL,
		ClusterName: cluster_name,
		Channel:     channel,
		Username:    username,
		IconUrl:     icon_url,
		Text:        text,
		Enabled:     enabled,
	}
	n.postToMattermostWebhook()
	if string(expectedPayload) != string(actualPayload) {
		t.Errorf("Expected request body to be %s, got %s\n", expectedPayload, actualPayload)
	}
}
