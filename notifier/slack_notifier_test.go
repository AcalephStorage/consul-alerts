package notifier

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestSlackJsonUnmarshalling(t *testing.T) {
	cluster_name := "my_cluster"
	channel := "test_channel"
	username := "test_username"
	icon_url := "test_icon_url"
	icon_emoji := "test_icon_emoji"
	text := "test_text"

	expectedNotifier := SlackNotifier{
		ClusterName: cluster_name,
		Channel:     channel,
		Username:    username,
		IconUrl:     icon_url,
		IconEmoji:   icon_emoji,
		Text:        text,
	}
	var unmarshalledNotifier SlackNotifier

	data := []byte(fmt.Sprintf(`{
    "cluster_name": "%s",
    "channel": "%s",
    "username": "%s",
    "icon_url": "%s",
    "icon_emoji": "%s",
    "text": "%s"
  }`, cluster_name, channel, username, icon_url, icon_emoji, text))

	fmt.Printf("%s\n", data)
	if err := json.Unmarshal(data, &unmarshalledNotifier); err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(expectedNotifier, unmarshalledNotifier) {
		t.Fatalf("Expected slackNotifier to be %s, got %s\n", expectedNotifier, unmarshalledNotifier)
	}
}
