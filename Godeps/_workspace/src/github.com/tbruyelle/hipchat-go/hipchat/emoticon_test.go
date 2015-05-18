package hipchat

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestEmoticonList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/emoticon", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method %s, want %s", r.Method, m)
		}
		params := map[string]string{"start-index": "0", "max-results": "100", "type": "all"}
		for k, v := range params {
			if v != r.FormValue(k) {
				t.Errorf("Request query params %s=%s, want %s", k, r.FormValue(k), v)
			}
		}
		fmt.Fprintf(w, `{
			"items": [{"id":1, "url":"u", "shortcut":"s", "links":{"self":"s"}}],
			"startIndex": 1,
			"maxResults": 1,
			"links":{"self":"s", "prev":"p", "next":"n"}
		}`)
	})
	want := &Emoticons{
		Items:      []Emoticon{Emoticon{ID: 1, Url: "u", Shortcut: "s", Links: Links{Self: "s"}}},
		StartIndex: 1,
		MaxResults: 1,
		Links:      PageLinks{Links: Links{Self: "s"}, Prev: "p", Next: "n"},
	}

	emos, _, err := client.Emoticon.List(0, 100, "all")
	if err != nil {
		t.Fatalf("Emoticon.List returned an error %v", err)
	}
	if !reflect.DeepEqual(want, emos) {
		t.Errorf("Emoticon.List returned %+v, want %+v", emos, want)
	}
}
