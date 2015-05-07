package hipchat

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// RoomService gives access to the room related methods of the API.
type RoomService struct {
	client *Client
}

// Rooms represents a HipChat room list.
type Rooms struct {
	Items      []Room    `json:"items"`
	StartIndex int       `json:"startIndex"`
	MaxResults int       `json:"maxResults"`
	Links      PageLinks `json:"links"`
}

// Room represents a HipChat room.
type Room struct {
	ID                int            `json:"id"`
	Links             RoomLinks      `json:"links"`
	Name              string         `json:"name"`
	XmppJid           string         `json:"xmpp_jid"`
	Statistics        RoomStatistics `json:"statistics"`
	Created           string         `json:"created"`
	IsArchived        bool           `json:"is_archived"`
	Privacy           string         `json:"privacy"`
	IsGuestAccessible bool           `json:"is_guess_accessible"`
	Topic             string         `json:"topic"`
	Participants      []User         `json:"participants"`
	Owner             User           `json:"owner"`
	GuestAccessURL    string         `json:"guest_access_url"`
}

// RoomStatistics represents the HipChat room statistics.
type RoomStatistics struct {
	Links Links `json:"links"`
}

// CreateRoomRequest represents a HipChat room creation request.
type CreateRoomRequest struct {
	Topic       string `json:"topic,omitempty"`
	GuestAccess bool   `json:"guest_access,omitempty"`
	Name        string `json:"name,omitempty"`
	OwnerUserID string `json:"owner_user_id,omitempty"`
	Privacy     string `json:"privacy,omitempty"`
}

// UpdateRoomRequest represents a HipChat room update request.
type UpdateRoomRequest struct {
	Name          string `json:"name"`
	Topic         string `json:"topic"`
	IsGuestAccess bool   `json:"is_guest_access"`
	IsArchived    bool   `json:"is_archived"`
	Privacy       string `json:"privacy"`
	Owner         ID     `json:"owner"`
}

// RoomLinks represents the HipChat room links.
type RoomLinks struct {
	Links
	Webhooks     string `json:"webhooks"`
	Members      string `json:"members"`
	Participants string `json:"participants"`
}

// NotificationRequest represents a HipChat room notification request.
type NotificationRequest struct {
	Color         string `json:"color,omitempty"`
	Message       string `json:"message,omitempty"`
	Notify        bool   `json:"notify,omitempty"`
	MessageFormat string `json:"message_format,omitempty"`
}

// ShareFileRequest represents a HipChat room file share request.
type ShareFileRequest struct {
	Path     string `json:"path"`
	Filename string `json:"filename,omitempty"`
	Message  string `json:"message,omitempty"`
}

// HistoryRequest represents a HipChat room chat history request.
type HistoryRequest struct {
	Date       string `json:"date"`
	Timezone   string `json:"timezone"`
	StartIndex int    `json:"start-index"`
	MaxResults int    `json:"max-results"`
	Reverse    bool   `json:"reverse"`
}

// History represents a HipChat room chat history.
type History struct {
	Items      []Message `json:"items"`
	StartIndex int       `json:"startIndex"`
	MaxResults int       `json:"maxResults"`
	Links      PageLinks `json:"links"`
}

// Message represents a HipChat message.
type Message struct {
	Date          string      `json:"date"`
	From          interface{} `json:"from"` // string | obj <- weak
	Id            string      `json:"id"`
	Mentions      []User      `json:"mentions"`
	Message       string      `json:"message"`
	MessageFormat string      `json:"message_format"`
	Type          string      `json:"type"`
}

// SetTopicRequest represents a hipchat update topic request
type SetTopicRequest struct {
	Topic string `json:"topic"`
}

// InviteRequest represents a hipchat invite to room request
type InviteRequest struct {
	Reason string `json:"reason"`
}

// List returns all the rooms authorized.
//
// HipChat API docs: https://www.hipchat.com/docs/apiv2/method/get_all_rooms
func (r *RoomService) List() (*Rooms, *http.Response, error) {
	req, err := r.client.NewRequest("GET", "room", nil)
	if err != nil {
		return nil, nil, err
	}

	rooms := new(Rooms)
	resp, err := r.client.Do(req, rooms)
	if err != nil {
		return nil, resp, err
	}
	return rooms, resp, nil
}

// Get returns the room specified by the id.
//
// HipChat API docs: https://www.hipchat.com/docs/apiv2/method/get_room
func (r *RoomService) Get(id string) (*Room, *http.Response, error) {
	req, err := r.client.NewRequest("GET", fmt.Sprintf("room/%s", id), nil)
	if err != nil {
		return nil, nil, err
	}

	room := new(Room)
	resp, err := r.client.Do(req, room)
	if err != nil {
		return nil, resp, err
	}
	return room, resp, nil
}

// Notification sends a notification to the room specified by the id.
//
// HipChat API docs: https://www.hipchat.com/docs/apiv2/method/send_room_notification
func (r *RoomService) Notification(id string, notifReq *NotificationRequest) (*http.Response, error) {
	req, err := r.client.NewRequest("POST", fmt.Sprintf("room/%s/notification", id), notifReq)
	if err != nil {
		return nil, err
	}

	return r.client.Do(req, nil)
}

// ShareFile sends a file to the room specified by the id.
//
// HipChat API docs: https://www.hipchat.com/docs/apiv2/method/share_file_with_room
func (r *RoomService) ShareFile(id string, shareFileReq *ShareFileRequest) (*http.Response, error) {
	req, err := r.client.NewFileUploadRequest("POST", fmt.Sprintf("room/%s/share/file", id), shareFileReq)
	if err != nil {
		return nil, err
	}

	return r.client.Do(req, nil)
}

// Create creates a new room.
//
// HipChat API docs: https://www.hipchat.com/docs/apiv2/method/create_room
func (r *RoomService) Create(roomReq *CreateRoomRequest) (*Room, *http.Response, error) {
	req, err := r.client.NewRequest("POST", "room", roomReq)
	if err != nil {
		return nil, nil, err
	}

	room := new(Room)
	resp, err := r.client.Do(req, room)
	if err != nil {
		return nil, resp, err
	}
	return room, resp, nil
}

// Delete deletes an existing room.
//
// HipChat API docs: https://www.hipchat.com/docs/apiv2/method/delete_room
func (r *RoomService) Delete(id string) (*http.Response, error) {
	req, err := r.client.NewRequest("DELETE", fmt.Sprintf("room/%s", id), nil)
	if err != nil {
		return nil, err
	}

	return r.client.Do(req, nil)
}

// Update updates an existing room.
//
// HipChat API docs: https://www.hipchat.com/docs/apiv2/method/update_room
func (r *RoomService) Update(id string, roomReq *UpdateRoomRequest) (*http.Response, error) {
	req, err := r.client.NewRequest("PUT", fmt.Sprintf("room/%s", id), roomReq)
	if err != nil {
		return nil, err
	}

	return r.client.Do(req, nil)
}

// History fetches a room's chat history.
//
// HipChat API docs: https://www.hipchat.com/docs/apiv2/method/view_room_history
func (r *RoomService) History(id string, roomReq *HistoryRequest) (*History, *http.Response, error) {
	u := fmt.Sprintf("room/%s/history", id)
	// Form query parameters
	if roomReq != nil {
		p := url.Values{}
		if roomReq.Date != "" {
			p.Add("date", roomReq.Date)
		}
		if roomReq.Timezone != "" {
			p.Add("timezone", roomReq.Timezone)
		}
		if roomReq.StartIndex != 0 {
			p.Add("start-index", strconv.FormatInt(int64(roomReq.StartIndex), 10))
		}
		if roomReq.MaxResults != 0 {
			p.Add("max-results", strconv.FormatInt(int64(roomReq.MaxResults), 10))
		}
		// There's no way to tell whether caller set a boolean or not. We have to always set
		// it.
		p.Add("reverse", strconv.FormatBool(roomReq.Reverse))
		u += "?" + p.Encode()
	}
	req, err := r.client.NewRequest("GET", u, nil)
	h := new(History)
	resp, err := r.client.Do(req, &h)
	if err != nil {
		return nil, resp, err
	}
	return h, resp, nil
}

// Set Room topic.
//
// HipChat API docs: https://www.hipchat.com/docs/apiv2/method/set_topic
func (r *RoomService) SetTopic(id string, topic string) (*http.Response, error) {
	topicReq := &SetTopicRequest{Topic: topic}

	req, err := r.client.NewRequest("PUT", fmt.Sprintf("room/%s/topic", id), topicReq)
	if err != nil {
		return nil, err
	}

	return r.client.Do(req, nil)
}

// Invite someone to the Room.
//
// HipChat API docs: https://www.hipchat.com/docs/apiv2/method/invite_user
func (r *RoomService) Invite(room string, user string, reason string) (*http.Response, error) {
	reasonReq := &InviteRequest{Reason: reason}

	req, err := r.client.NewRequest("POST", fmt.Sprintf("room/%s/invite/%s", room, user), reasonReq)
	if err != nil {
		return nil, err
	}

	return r.client.Do(req, nil)
}
