package notifier

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"encoding/json"
	"net/http"
	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
)

type StringMap map[string]string

type MatterMostLoginInfo struct {
	LoginID     string `json:"login_id"`
	Password    string `json:"password"`
}

type MatterMostAuthInfo struct {
	UserID             string    `json:"id"`
	CreateAt           int64     `json:"create_at"`
	UpdateAt           int64     `json:"update_at"`
	DeleteAt           int64     `json:"delete_at"`
	UserName           string    `json:"username"` 
	AuthData           string    `json:"auth_data"`
	AuthService        string    `json:"auth_service"`
	Email              string    `json:"email"`
	EmailVerified      bool      `json:"email_verified"`
	NickName           string    `json:"nickname"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Roles              string    `json:"roles"`
	AllowMarketing     bool      `json:"allow_marketing"`
	NotifyProps        StringMap `json:"notify_props"`
	Props              StringMap `json:"props"`
	LastPasswordUpdate int64     `json:"last_password_update"`
	LastPictureUpdate  int64     `json:"last_picture_update"`
}

type MatterMostUserInfo struct {
	UserID             string    `json:"id"`
	CreateAt           int64     `json:"create_at"`
	UpdateAt           int64     `json:"update_at"`
	DeleteAt           int64     `json:"delete_at"`
	Username           string    `json:"username"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	Nickname           string    `json:"nickname"`
	Email              string    `json:"email"`
	EmailVerified      bool      `json:"email_verified"`
	Password           string    `json:"password"`
	AuthData           *string   `json:"auth_data"`
	AuthService        string    `json:"auth_service"`
	Roles              string    `json:"roles"`
	NotifyProps        StringMap `json:"notify_props"`
	Props              StringMap `json:"props,omitempty"`
	LastPasswordUpdate int64     `json:"last_password_update"`
	LastPictureUpdate  int64     `json:"last_picture_update"`
	FailedAttempts     int       `json:"failed_attempts"`
	MfaActive          bool      `json:"mfa_active"`
	MfaSecret          string    `json:"mfa_secret"`
}


type MatterMostTeamInfo struct {
	TeamID          string `json:"id"`
	CreateAt        int64  `json:"create_at"`
	UpdateAt        int64  `json:"update_at"`
	DeleteAt        int64  `json:"delete_at"`
	DisplayName     string `json:"display_name"`
	Name            string `json:"name"`
	Email           string `json:"email"`
	Type            string `json:"type"`
	AllowedDomains  string `json:"allowed_domains"`
	InviteID        string `json:"invite_id"`
	AllowOpenInvite bool   `json:"allow_open_invite"`
}

type MatterMostChannelInfo struct {
	ChannelID     string `json:"id"`
	CreateAt      int64  `json:"create_at"`
	UpdateAt      int64  `json:"update_at"`
	DeleteAt      int64  `json:"delete_at"`
	TeamID        string `json:"team_id"` 
	Type          string `json:"type"`
	DisplayName   string `json:"display_name"`
	Name          string `json:"name"` 
	Header        string `json:"header"`
	Purpose       string `json:"purpose"`
	LastPostAt    int64  `json:"last_post_at"`
	TotalMsgCount int64  `json:"total_msg_count"`
	ExtraUpdateAt int64  `json:"extra_update_at"`
	CreatorID     string `json:"creator_id"`
}

type MatterMostChannelList struct {
	Channels []MatterMostChannelInfo
}

type MatterMostPostInfo struct {
	PostID        string   `json:"id"`
	CreateAt      int64    `json:"create_at"`
	UpdateAt      int64    `json:"update_at"`
	DeleteAt      int64    `json:"delete_at"`
	UserID        string   `json:"user_id"`
	ChannelID     string   `json:"channel_id"`
	RootID        string   `json:"root_id"`
	ParentID      string   `json:"parent_id"`
	OriginalID    string   `json:"original_id"`
	Message       string   `json:"message"`
	Type          string   `json:"type"`
	Props         StringMap`json:"props"`
	Hashtags      string   `json:"hashtags"`
	Filenames     StringMap`json:"filenames"`
	PendingPostID string   `json:"pending_post_id"`
}

type MatterMostNotifier struct {
	ClusterName string
	Url         string
	UserName    string
	Password    string
	TeamName    string
	Channel     string
	Detailed    bool
	NotifName   string

	/* Filled in after authentication */
	Token       string
	TeamID      string
	UserID      string
	ChannelID   string
	Text        string
}

func (mattermost *MatterMostNotifier) GetURL() string {
	ssl := false
	proto := "http"

	u := strings.TrimSpace(strings.ToLower(mattermost.Url))
	if u[:5] == "https" && u[5] == ':' {
		ssl = true
		proto = "https"
	}

	host := ""
	port := 80
	buf := strings.Split(u, ":")
	if (u[:4] == "http" && u[4] == ':') ||
		(u[:5] == "https" && u[5] == ':') && len(buf) == 3 {
		host = strings.Trim(buf[1], "/")
		port, _ = strconv.Atoi(strings.TrimSpace(buf[2]))

	} else if len(buf) == 2 {
		host = strings.Trim(buf[0], "/")
		port, _ = strconv.Atoi(strings.TrimSpace(buf[1]))
	} else {
		host = strings.TrimSpace(buf[0])
	}

	return fmt.Sprintf("%s://%s:%d/api/v3", proto, host, port)
}

func (mattermost *MatterMostNotifier) Authenticate() bool {

	loginURL := fmt.Sprintf("%s/users/login", mattermost.GetURL())
	loginInfo := MatterMostLoginInfo{ LoginID: mattermost.UserName,
		Password: mattermost.Password}

    buf := new(bytes.Buffer)
    json.NewEncoder(buf).Encode(loginInfo)

	req, err := http.NewRequest("POST", loginURL, buf)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var a MatterMostAuthInfo
	err = decoder.Decode(&a)
	if err != nil {
		log.Fatal("Decode: ", err)
		return false
	}

	if buf, ok := resp.Header["Token"]; ok {
		if len(buf) > 0 {
			mattermost.Token = buf[0]
			return true
		}
	}

	return false
}

func (mattermost *MatterMostNotifier) GetAllTeams(teams *[]MatterMostTeamInfo) bool {

	teamURL := fmt.Sprintf("%s/teams/all", mattermost.GetURL())
	req, err := http.NewRequest("GET", teamURL, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return false
	}

	authorization := fmt.Sprintf("Bearer %s", mattermost.Token)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var buf map[string]*MatterMostTeamInfo
	err = decoder.Decode(&buf)
	if err != nil {
		log.Fatal("Decode: ", err)
		return false
	}

	if len(buf) > 0 {
		for _, value := range buf {
			*teams = append(*teams, *value)
		}
		return true
	}

	return false
}

func (mattermost *MatterMostNotifier) GetUser(userID string, userInfo *MatterMostUserInfo) bool {

	if userId == "" || userInfo == nil {
		return false
	}

	userURL := fmt.Sprintf("%s/users/%s/get", mattermost.GetURL(), userID)

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return false
	}

	authorization := fmt.Sprintf("Bearer %s", mattermost.Token)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(userInfo)
	if err != nil {
		log.Fatal("Decode: ", err)
		return false
	}

	return true
}

func (mattermost *MatterMostNotifier) GetMe(me *MatterMostUserInfo) bool {

	if me == nil {
		return false
	}

	userURL := fmt.Sprintf("%s/users/me", mattermost.GetURL())

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return false
	}

	authorization := fmt.Sprintf("Bearer %s", mattermost.Token)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(me)
	if err != nil {
		log.Fatal("Decode: ", err)
		return false
	}

	return true
}

func (mattermost *MatterMostNotifier) GetTeam(teamID string, teamInfo *MatterMostTeamInfo) bool {

	if teamID == "" || teamInfo == nil {
		return false
	}

	teamURL := fmt.Sprintf("%s/teams/%s/me", mattermost.GetURL(), teamID)

	req, err := http.NewRequest("GET", teamURL, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return false
	}

	authorization := fmt.Sprintf("Bearer %s", mattermost.Token)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(teamInfo)
	if err != nil {
		log.Fatal("Decode: ", err)
		return false
	}

	return true
}

func (mattermost *MatterMostNotifier) GetChannels(teamID string, channels *[]MatterMostChannelInfo) bool {

	if teamID == "" || channels == nil {
		return false
	}

	channelURL := fmt.Sprintf("%s/teams/%s/channels/", mattermost.GetURL(), teamID)
	req, err := http.NewRequest("GET", channelURL, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return false
	}

	authorization := fmt.Sprintf("Bearer %s", mattermost.Token)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	fc := &MatterMostChannelList{}
	err = decoder.Decode(&fc)
	if err != nil {
		log.Fatal("Decode: ", err)
		return false
	}
	*channels = fc.Channels

	return true
}


func (mattermost *MatterMostNotifier) PostMessage(teamID string, channelID string, postInfo *MatterMostPostInfo) bool {

	if teamID == "" || channelID == "" || postInfo == nil {
		return false
	}

	postURL := fmt.Sprintf("%s/teams/%s/channels/%s/posts/create",
		mattermost.GetURL(), teamID, channelID)

	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(*postInfo)

	req, err := http.NewRequest("POST", postURL, buf)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return false
	}

	authorization := fmt.Sprintf("Bearer %s", mattermost.Token)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var p MatterMostPostInfo
	err = decoder.Decode(&p)
	if err != nil {
		log.Fatal("Decode: ", err)
		return false
	}
	*postInfo = p

	return false
}

func (mattermost *MatterMostNotifier) Init() bool {
	if mattermost.Initialized == true {
		return true
	}

	if mattermost.Token == "" && !mattermost.Authenticate() {
		return false
	}

	if mattermost.TeamID == "" {
		var teams []MatterMostTeamInfo

		if !mattermost.GetAllTeams(&teams) {
			return false
		}

		for i := 0;i < len(teams);i++ {
			if teams[i].Name == mattermost.TeamName {
				mattermost.TeamID = teams[i].TeamID
				break
			}
		}

		if mattermost.TeamID == "" {
			return false
		}
	}

	if mattermost.UserID == "" {
		var me MatterMostUserInfo

		if !mattermost.GetMe(&me) {
			return false
		}

		if me.UserID == "" {
			return false
		}

		mattermost.UserID = me.UserID
	}

	if mattermost.ChannelID == "" {
		var channels []MatterMostChannelInfo

		if !mattermost.GetChannels(mattermost.TeamID, &channels) {
			return false
		}

		for i := 0;i < len(channels);i++ {
			if channels[i].Name == mattermost.Channel {
				mattermost.ChannelID = channels[i].ChannelID
				break
			}
		}

		if mattermost.ChannelID == "" {
			return false
		}
	}

	mattermost.Initialized = true
	return true
}


// NotifierName provides name for notifier selection
func (mattermost *MatterMostNotifier) NotifierName() string {
	return mattermost.NotifName
}

//Notify sends messages to the endpoint notifier
func (mattermost *MatterMostNotifier) Notify(messages Messages) bool {
	if !mattermost.Init() {
		return false
	}

	if mattermost.Detailed {
		return mattermost.notifyDetailed(messages)
	}

	return mattermost.notifySimple(messages)
}

func (mattermost *MatterMostNotifier) notifySimple(messages Messages) bool {
	overallStatus, pass, warn, fail := messages.Summary()

	text := fmt.Sprintf(header, mattermost.ClusterName, overallStatus, fail, warn, pass)

	for _, message := range messages {
		text += fmt.Sprintf("\n%s:%s:%s is %s.",
			message.Node, message.Service, message.Check, message.Status)
		text += fmt.Sprintf("\n%s", message.Output)
	}

	mattermost.Text = text

	return mattermost.postToMatterMost()
}

func (mattermost *MatterMostNotifier) notifyDetailed(messages Messages) bool {

	overallStatus, pass, warn, fail := messages.Summary()

	var emoji, color string
	switch overallStatus {
	case SYSTEM_HEALTHY:
		emoji = ":white_check_mark:"
		color = "good"
	case SYSTEM_UNSTABLE:
		emoji = ":question:"
		color = "warning"
	case SYSTEM_CRITICAL:
		emoji = ":x:"
		color = "danger"
	default:
		emoji = ":question:"
	}
	title := "Consul monitoring report"
	pretext := fmt.Sprintf("%s %s is *%s*", emoji, mattermost.ClusterName, overallStatus)

	detailedBody := ""
	detailedBody += fmt.Sprintf("*Changes:* Fail = %d, Warn = %d, Pass = %d",
		fail, warn, pass)
	detailedBody += fmt.Sprintf("\n")

	for _, message := range messages {
		detailedBody += fmt.Sprintf("\n*[%s:%s]* %s is *%s.*",
			message.Node, message.Service, message.Check, message.Status)
		detailedBody += fmt.Sprintf("\n`%s`", strings.TrimSpace(message.Output))
	}

	mattermost.Text := fmt.Sprintf("%s\n%s\n%s\n", title, pretext, detailedBody)

	return mattermost.postToMatterMost()

}

func (mattermost *MatterMostNotifier) postToMatterMost() bool {
	var postInfo = MatterMostPostInfo{
		ChannelID: mattermost.ChannelID,
		Message: mattermost.Text }

	return mattermost.PostMessage(mattermost.TeamID, mattermost.ChannelID, &postInfo)
}
