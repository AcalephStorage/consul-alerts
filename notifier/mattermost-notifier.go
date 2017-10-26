package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"io/ioutil"
)

type StringMap map[string]string

type MattermostLoginInfo struct {
	LoginID  string `json:"login_id"`
	Password string `json:"password"`
}

type MattermostAuthInfo struct {
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

type MattermostUserInfo struct {
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

type MattermostTeamInfo struct {
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

type MattermostChannelInfo struct {
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

type MattermostChannelList struct {
	Channels []MattermostChannelInfo
}

type MattermostPostInfo struct {
	PostID        string    `json:"id"`
	CreateAt      int64     `json:"create_at"`
	UpdateAt      int64     `json:"update_at"`
	DeleteAt      int64     `json:"delete_at"`
	UserID        string    `json:"user_id"`
	ChannelID     string    `json:"channel_id"`
	RootID        string    `json:"root_id"`
	ParentID      string    `json:"parent_id"`
	OriginalID    string    `json:"original_id"`
	Message       string    `json:"message"`
	Type          string    `json:"type"`
	Props         StringMap `json:"props"`
	Hashtags      string    `json:"hashtags"`
	Filenames     StringMap `json:"filenames"`
	PendingPostID string    `json:"pending_post_id"`
}

type MattermostNotifier struct {
	ClusterName string  `json:"cluster_name"`
	Url         string  `json:"-"`
	UserName    string  `json:"username"`
	Password    string  `json:"-"`
	Team        string  `json:"team,omitempty"`
	Channel     string  `json:"channel"`
	Detailed    bool    `json:"-"`
	Mode        string  `json:"-"`
	NotifName   string  `json:"-"`
	Enabled     bool    `json:"-"`

	/* Filled in after authentication */
	Initialized bool    `json:"-"`
	Token       string  `json:"-"`
	TeamID      string  `json:"-"`
	UserID      string  `json:"-"`
	ChannelID   string  `json:"-"`
	Text        string  `json:"text"`
}

func (mattermost *MattermostNotifier) GetURL() string {

	proto := "http"
	u := strings.TrimSpace(strings.ToLower(mattermost.Url))
	if u[:5] == "https" && u[5] == ':' {
		proto = "https"
	}

	host := ""
	port := 0
	buf := strings.Split(u, ":")
	if (u[:4] == "http" && u[4] == ':') ||
		(u[:5] == "https" && u[5] == ':') {

		host = strings.Trim(buf[1], "/")
		if len(buf) == 3 {
			port, _ = strconv.Atoi(strings.TrimSpace(buf[2]))
		}

	} else if len(buf) == 2 {
		host = strings.Trim(buf[0], "/")
		port, _ = strconv.Atoi(strings.TrimSpace(buf[1]))

	} else {
		host = strings.TrimSpace(buf[0])
	}

	portstr := ""
	if port > 0 {
		portstr = fmt.Sprintf(":%d", port)
	}

	return fmt.Sprintf("%s://%s%s/api/v3", proto, host, portstr)
}

func (mattermost *MattermostNotifier) Authenticate() bool {

	loginURL := fmt.Sprintf("%s/users/login", mattermost.GetURL())
	loginInfo := MattermostLoginInfo{LoginID: mattermost.UserName,
		Password: mattermost.Password}

	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(loginInfo)

	req, err := http.NewRequest("POST", loginURL, buf)
	if err != nil {
		log.Error("NewRequest: ", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var a MattermostAuthInfo
	err = decoder.Decode(&a)
	if err != nil {
		log.Error("Decode: ", err)
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

func (mattermost *MattermostNotifier) GetAllTeams(teams *[]MattermostTeamInfo) bool {

	teamURL := fmt.Sprintf("%s/teams/all", mattermost.GetURL())
	req, err := http.NewRequest("GET", teamURL, nil)
	if err != nil {
		log.Error("NewRequest: ", err)
		return false
	}

	authorization := fmt.Sprintf("Bearer %s", mattermost.Token)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var buf map[string]*MattermostTeamInfo
	err = decoder.Decode(&buf)
	if err != nil {
		log.Error("Decode: ", err)
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

func (mattermost *MattermostNotifier) GetUser(userID string, userInfo *MattermostUserInfo) bool {

	if userID == "" || userInfo == nil {
		return false
	}

	userURL := fmt.Sprintf("%s/users/%s/get", mattermost.GetURL(), userID)

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		log.Error("NewRequest: ", err)
		return false
	}

	authorization := fmt.Sprintf("Bearer %s", mattermost.Token)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(userInfo)
	if err != nil {
		log.Error("Decode: ", err)
		return false
	}

	return true
}

func (mattermost *MattermostNotifier) GetMe(me *MattermostUserInfo) bool {

	if me == nil {
		return false
	}

	userURL := fmt.Sprintf("%s/users/me", mattermost.GetURL())

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		log.Error("NewRequest: ", err)
		return false
	}

	authorization := fmt.Sprintf("Bearer %s", mattermost.Token)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(me)
	if err != nil {
		log.Error("Decode: ", err)
		return false
	}

	return true
}

func (mattermost *MattermostNotifier) GetTeam(teamID string, teamInfo *MattermostTeamInfo) bool {

	if teamID == "" || teamInfo == nil {
		return false
	}

	teamURL := fmt.Sprintf("%s/teams/%s/me", mattermost.GetURL(), teamID)

	req, err := http.NewRequest("GET", teamURL, nil)
	if err != nil {
		log.Error("NewRequest: ", err)
		return false
	}

	authorization := fmt.Sprintf("Bearer %s", mattermost.Token)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(teamInfo)
	if err != nil {
		log.Error("Decode: ", err)
		return false
	}

	return true
}

func (mattermost *MattermostNotifier) GetChannels(teamID string, channels *[]MattermostChannelInfo) bool {

	if teamID == "" || channels == nil {
		return false
	}

	channelURL := fmt.Sprintf("%s/teams/%s/channels/", mattermost.GetURL(), teamID)
	req, err := http.NewRequest("GET", channelURL, nil)
	if err != nil {
		log.Error("NewRequest: ", err)
		return false
	}

	authorization := fmt.Sprintf("Bearer %s", mattermost.Token)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	fc := &MattermostChannelList{}
	err = decoder.Decode(&fc)
	if err != nil {
		log.Error("Decode: ", err)
		return false
	}
	*channels = fc.Channels

	return true
}

func (mattermost *MattermostNotifier) PostMessage(teamID string, channelID string, postInfo *MattermostPostInfo) bool {

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
		log.Error("NewRequest: ", err)
		return false
	}

	authorization := fmt.Sprintf("Bearer %s", mattermost.Token)
	req.Header.Set("Authorization", authorization)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var p MattermostPostInfo
	err = decoder.Decode(&p)
	if err != nil {
		log.Error("Decode: ", err)
		return false
	}
	*postInfo = p

	return false
}

func (mattermost *MattermostNotifier) Init() bool {
	if mattermost.Initialized == true {
		return true
	}

	if mattermost.Token == "" && !mattermost.Authenticate() {
		log.Println("Mattermost: Unable to authenticate!")
		return false
	}

	if mattermost.TeamID == "" {
		var teams []MattermostTeamInfo

		if !mattermost.GetAllTeams(&teams) {
			log.Println("Mattermost: Unable to get teams!")
			return false
		}

		for i := 0; i < len(teams); i++ {
			if teams[i].Name == mattermost.Team {
				mattermost.TeamID = teams[i].TeamID
				break
			}
		}

		if mattermost.TeamID == "" {
			log.Println("Mattermost: Unable to find team!")
			return false
		}
	}

	if mattermost.UserID == "" {
		var me MattermostUserInfo

		if !mattermost.GetMe(&me) {
			log.Println("Mattermost: Unable to get user!")
			return false
		}

		if me.UserID == "" {
			log.Println("Mattermost: Unable to get user ID!")
			return false
		}

		mattermost.UserID = me.UserID
	}

	if mattermost.ChannelID == "" {
		var channels []MattermostChannelInfo

		if !mattermost.GetChannels(mattermost.TeamID, &channels) {
			log.Println("Mattermost: Unable to get channels!")
			return false
		}

		for i := 0; i < len(channels); i++ {
			if channels[i].Name == mattermost.Channel {
				mattermost.ChannelID = channels[i].ChannelID
				break
			}
		}

		if mattermost.ChannelID == "" {
			log.Println("Mattermost: Unable to find channel!")
			return false
		}
	}

	mattermost.Initialized = true
	return true
}

// NotifierName provides name for notifier selection
func (mattermost *MattermostNotifier) NotifierName() string {
	return "mattermost"
}

func (mattermost *MattermostNotifier) Copy() Notifier {
	notifier := *mattermost
	return &notifier
}

//Notify sends messages to the endpoint notifier
func (mattermost *MattermostNotifier) Notify(messages Messages) bool {
	if mattermost.Mode != "webhook" && !mattermost.Init() {
		return false
	}

	if mattermost.Detailed {
		return mattermost.notifyDetailed(messages)
	}

	return mattermost.notifySimple(messages)
}

func (mattermost *MattermostNotifier) notifySimple(messages Messages) bool {
	overallStatus, pass, warn, fail := messages.Summary()

	text := fmt.Sprintf(header, mattermost.ClusterName, overallStatus, fail, warn, pass)

	for _, message := range messages {
		text += fmt.Sprintf("\n%s:%s:%s is %s.",
			message.Node, message.Service, message.Check, message.Status)
		text += fmt.Sprintf("\n%s\n\n", message.Output)
	}

	mattermost.Text = text

	return mattermost.postToMattermost()
}

func (mattermost *MattermostNotifier) notifyDetailed(messages Messages) bool {

	overallStatus, pass, warn, fail := messages.Summary()

	var emoji string
	switch overallStatus {
	case SYSTEM_HEALTHY:
		emoji = ":white_check_mark:"
	case SYSTEM_UNSTABLE:
		emoji = ":question:"
	case SYSTEM_CRITICAL:
		emoji = ":x:"
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

	mattermost.Text = fmt.Sprintf("%s\n%s\n%s\n\n", title, pretext, detailedBody)

	return mattermost.postToMattermost()

}

func (mattermost *MattermostNotifier) postToMattermost() bool {
	if mattermost.Mode == "webhook" {
                return mattermost.postToMattermostWebHook()
        }

	var postInfo = MattermostPostInfo{
		ChannelID: mattermost.ChannelID,
		Message:   mattermost.Text}

	return mattermost.PostMessage(mattermost.TeamID, mattermost.ChannelID, &postInfo)
}

func (mattermost *MattermostNotifier) postToMattermostWebHook() bool {
	data, err := json.Marshal(mattermost)
        if err != nil {
                log.Println("Unable to marshal mattermost payload:", err)
                return false
        }
        log.Debugf("struct = %+v, json = %s", mattermost, string(data))

        b := bytes.NewBuffer(data)
        if res, err := http.Post(mattermost.Url, "application/json", b); err != nil {
                log.Println("Unable to send data to mattermost:", err)
                return false
        } else {
                defer res.Body.Close()
                statusCode := res.StatusCode
                if statusCode != 200 {
                        body, _ := ioutil.ReadAll(res.Body)
                        log.Println("Unable to notify mattermost:", string(body))
                        return false
                } else {
                        log.Println("Mattermost notification sent.")
                        return true
                }
        }
}
