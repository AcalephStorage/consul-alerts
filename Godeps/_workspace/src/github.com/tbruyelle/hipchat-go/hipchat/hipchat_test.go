package hipchat

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

// setup sets up a test HTTP server and a hipchat.Client configured to talk
// to that test server.
// Tests should register handlers on mux which provide mock responses for
// the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// github client configured to use test server
	client = NewClient("AuthToken")
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	authToken := "AuthToken"

	c := NewClient(authToken)

	if c.authToken != authToken {
		t.Errorf("NewClient authToken %s, want %s", c.authToken, authToken)
	}
	if c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL %s, want %s", c.BaseURL.String(), defaultBaseURL)
	}
}

func TestNewRequest(t *testing.T) {
	c := NewClient("AuthToken")

	inURL, outURL := "foo", defaultBaseURL+"foo"
	inBody, outBody := &NotificationRequest{Message: "Hello"}, `{"message":"Hello"}`+"\n"
	r, _ := c.NewRequest("GET", inURL, inBody)

	if r.URL.String() != outURL {
		t.Errorf("NewRequest URL %s, want %s", r.URL.String(), outURL)
	}
	body, _ := ioutil.ReadAll(r.Body)
	if string(body) != outBody {
		t.Errorf("NewRequest body %s, want %s", body, outBody)
	}
	authorization := r.Header.Get("Authorization")
	if authorization != "Bearer "+c.authToken {
		t.Errorf("NewRequest authorization header %s, want %s", authorization, "Bearer "+c.authToken)
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("NewRequest Content-Type header %s, want application/json", contentType)
	}
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		Bar int
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprintf(w, `{"Bar":1}`)
	})
	req, _ := client.NewRequest("GET", "/", nil)
	body := new(foo)

	_, err := client.Do(req, body)

	if err != nil {
		t.Fatal(err)
	}
	want := &foo{Bar: 1}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}
