package hipchat

import (
	"fmt"
	"net/http"
)

// EmoticonService gives access to the emoticon related part of the API.
type EmoticonService struct {
	client *Client
}

// Emoticons represents a list of hipchat emoticons.
type Emoticons struct {
	Items      []Emoticon `json:"items"`
	StartIndex int        `json:"startIndex"`
	MaxResults int        `json:"maxResults"`
	Links      PageLinks  `json:"links"`
}

// Emoticon represents a hipchat emoticon.
type Emoticon struct {
	ID       int    `json:"id"`
	Url      string `json:"url"`
	Links    Links  `json:"links"`
	Shortcut string `json:"shortcut"`
}

// List returns the list of all the emoticons
//
// HipChat api docs : https://www.hipchat.com/docs/apiv2/method/get_all_emoticons
func (e *EmoticonService) List(start, max int, type_ string) (*Emoticons, *http.Response, error) {
	req, err := e.client.NewRequest("GET",
		fmt.Sprintf("emoticon?start-index=%d&max-results=%d&type=%s", start, max, type_), nil)
	if err != nil {
		return nil, nil, err
	}

	emoticons := new(Emoticons)
	resp, err := e.client.Do(req, emoticons)
	if err != nil {
		return nil, resp, err
	}
	return emoticons, resp, nil
}
