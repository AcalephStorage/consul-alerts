package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/AcalephStorage/consul-alerts/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

type TelegramUser struct {
	ID				int64  `json:"id"`
	FirstName		string `json:"first_name"`
	LastName		string `json:"last_name"`
	UserName		string `json:"username"`
	LanguageCode	string `json:"language_code"`
}

type TelegramUserResponse struct {
	OK			bool         `json:"ok"`
	Description	string       `json:"description"`
	User		TelegramUser `json:"result"`
}

type TelegramChat struct {
	ID				int64  `json:"id"`
	Type			string `json:"type"`
	Title			string `json:"title"`
	UserName		string `json:"username"`
	FirstName		string `json:"first_name"`
	LastName		string `json:"last_name"`
	AllAdmins		bool   `json:"last_name"`
}

type TelegramMessageEntity struct {
	Type	string       `json:"type"`
	Offset	int64        `json:"offset"`
	Length	int64        `json:"length"`
	URL		string       `json:"url"`
	User	TelegramUser `json:"user"`
}

type TelegramAudio struct {
	FileID		string `json:"file_id"`
	Duration	int64  `json:"duration"`
	Performer	string `json:"performer"`
	Title		string `json:"title"`
	MimeType	string `json:"mime_type"`
	FileSize	int64  `json:"file_size"`
}

type TelegramPhotoSize struct {
	FileID		string `json:"file_id"`
	Width		int64  `json:"width"`
	Height		int64  `json:"height"`
	FileSize	int64  `json:"file_size"`
}

type TelegramDocument struct {
	FileID		string            `json:"file_id"`
	Thumb		TelegramPhotoSize `json:"thumb"`
	FileName	string            `json:"file_name"`
	MimeType	string            `json:"mime_type"`
	FileSize	int64             `json:"file_size"`
}

type TelegramAnimation struct {
	FileID		string            `json:"file_id"`
	Thumb		TelegramPhotoSize `json:"thumb"`
	FileName	string            `json:"file_name"`
	MimeType	string            `json:"mime_type"`
	FileSize	int64             `json:"file_size"`
}

type TelegramGame struct {
	Title			string                  `json:"title"`
	Description		string                  `json:"description"`
	Photo			[]TelegramPhotoSize     `json:"photo"`
	Text			string                  `json:"text"`
	TextEntities	[]TelegramMessageEntity `json:"text_entities"`
	Animation		TelegramAnimation       `json:"animation"`
}

type TelegramSticker struct {
	FileID		string            `json:"file_id"`
	Width		int64             `json:"width"`
	Height		int64             `json:"height"`
	Thumb		TelegramPhotoSize `json:"thumb"`
	Emoji		string            `json:"emoji"`
	FileSize	int64             `json:"file_size"`
}

type TelegramVideo struct {
	FileID		string            `json:"file_id"`
	Width		int64             `json:"width"`
	Height		int64             `json:"height"`
	Duration	int64             `json:"duration"`
	Thumb		TelegramPhotoSize `json:"thumb"`
	MimeType	string            `json:"mime_type"`
	FileSize	int64             `json:"file_size"`
}

type TelegramVoice struct {
	FileID		string `json:"file_id"`
	Duration	int64  `json:"duration"`
	MimeType	string `json:"mime_type"`
	FileSize	int64  `json:"file_size"`
}

type TelegramVideoNote struct {
	FileID		string            `json:"file_id"`
	Length		int64             `json:"length"`
	Duration	int64             `json:"duration"`
	Thumb		TelegramPhotoSize `json:"thumb"`
	FileSize	int64             `json:"file_size"`
}

type TelegramContact struct {
	PhoneNumber	string `json:"phone_number"`
	FirstName	string `json:"first_name"`
	LastName	string `json:"last_name"`
	UserID		int64  `json:"user_id"`
}

type TelegramLocation struct {
	Longitude	float64 `json:"longitude"`
	Latitude	float64 `json:"latitude"`
}

type TelegramVenue struct {
	Location		TelegramLocation `json:"location"`
	Title			string           `json:"title"`
	Address			string           `json:"address"`
	FoursquareID	string           `json:"foursquare_id"`
}

type TelegramInvoice struct {
	Title			string `json:"title"`
	Description		string `json:"description"`
	StartParameter	string `json:"start_parameter"`
	Currency		string `json:"currency"`
	TotalAmount		int64  `json:"total_amount"`
}

type TelegramShippingAddress struct {
	CountryCode	string `json:"country_code"`
	State		string `json:"state"`
	City		string `json:"city"`
	StreetLine1	string `json:"street_line1"`
	StreetLine2	string `json:"street_line2"`
	PostCode	string `json:"post_code"`
}

type TelegramOrderInfo struct {
	Name			string                  `json:"name"`
	PhoneNumber		string                  `json:"phone_number"`
	Email			string                  `json:"email"`
	ShippingAddress	TelegramShippingAddress `json:"shipping_address"`
}

type TelegramSuccessfulPayment struct {
	Currency				string            `json:"currency"`
	TotalAmount				int64             `json:"total_amount"`
	InvoicePayload			string            `json:"invoice_payload"`
	ShippingOptionID		string            `json:"shipping_option_id"`
	OrderInfo				TelegramOrderInfo `json:"order_info"`
	TelegramPaymentChargeID	string            `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeID	string            `json:"provider_payment_charge_id"`
}

type TelegramMessage struct {
	MessageID				int64                     `json:"message_id"`
	From					TelegramUser              `json:"from"`
	Date					int64                     `json:"date"`
	Chat					TelegramChat              `json:"chat"`
	ForwardFrom				TelegramUser              `json:"forward_from"`
	ForwardFromChat			TelegramChat              `json:"forward_from_chat"`
	ForwardFromMessageID	int64                     `json:"forward_from_message_id"`
	ForwardDate				int64                     `json:"forward_date"`
	ReplyToMessage			*TelegramMessage          `json:"reply_to_message"`
	EditDate				int64                     `json:"edit_date"`
	Text					string                    `json:"text"`
	Entities				[]TelegramMessageEntity   `json:"entities"`
	Audio					TelegramAudio             `json:"audio"`
	Document				TelegramDocument          `json:"document"`
	Game					TelegramGame              `json:"game"`
	Photo					[]TelegramPhotoSize       `json:"photo"`
	Sticker					TelegramSticker           `json:"sticker"`
	Video					TelegramVideo             `json:"video"`
	Voice					TelegramVoice             `json:"voice"`
	VideoNote				TelegramVideoNote         `json:"video_note"`
	NewChatMembers			[]TelegramUser            `json:"new_chat_members"`
	Caption					string                    `json:"caption"`
	Contact					TelegramContact           `json:"contact"`
	Location				TelegramLocation          `json:"location"`
	Venue					TelegramVenue             `json:"venue"`
	NewChatMember			TelegramUser              `json:"new_chat_member"`
	LeftChatMember			TelegramUser              `json:"left_chat_member"`
	NewChatTitle			string                    `json:"new_chat_title"`
	NewChatPhoto			[]TelegramPhotoSize       `json:"new_chat_photo"`
	DeleteChatPhoto			bool                      `json:"delete_chat_photo"`
	GroupChatCreated		bool                      `json:"group_chat_created"`
	SupergroupChatCreated	bool                      `json:"supergroup_chat_created"`
	ChannelChatCreated		bool                      `json:"channel_chat_created"`
	MigreateToChatID		int64                     `json:"migrate_to_chat_id"`
	MigrateFromChatID		int64                     `json:"migrate_from_chat_id"`
	PinnedMessage			*TelegramMessage          `json:"pinned_message"`
	Invoice					TelegramInvoice           `json:"invoice"`
	SuccessfulPayment		TelegramSuccessfulPayment `json:"successful_payment"`
}

type TelegramMessageResponse struct {
	OK			bool            `json:"ok"`
	Description	string          `json:"description"`
	Message		TelegramMessage `json:"result"`
}

type TelegramCallbackGame struct {
	/* placeholder */
}

type TelegramKeyboardButton struct {
	Text			string `json:"text"`
	RequestContact	bool   `json:"request_contact"`
	RequestLocation	bool   `json:"request_location"`
}

type TelegramReplyKeyboardRemove struct {
	RemoveKeyboard	bool `json:"remove_keyboard"`
	Selective		bool `json:"selective"`
}

type TelegramInlineKeyboardMarkup struct {
	InlineKeyboard	[][]TelegramInlineKeyboardButton `json:"inline_keyboard"`
}

type TelegramInlineKeyboardButton struct {
	Text							string               `json:"text"`
	URL								string               `json:"callback_data"`
	CallbackData					string               `json:"callback_data"`
	SwitchInlineQuery				string               `json:"switch_inline_query"`
	SwitchInlineQueryCurrentChat	string               `json:"switch_inline_query_current_chat"`
	CallbackGame					TelegramCallbackGame `json:"callback_game"`
	Pay								bool                 `json:"pay"`
}

type TelegramReplyKeyboardMarkup struct {
	Keyboard		[]TelegramKeyboardButton `json:"keyboard"`
	ResizeKeyboard	bool                     `json:"resize_keybaord"`
	OneTimeKeyboard	bool                     `json:"one_time_keyboard"`
	Selective		bool                     `json:"selective"`
}

type TelegramForceReply struct {
	ForceReply	bool `json:"force_reply"`
	Selective	bool `json:"selective"`
}

type TelegramSendMessage struct {
	ChatID					string `json:"chat_id"`
	Text					string `json:"text"`
	ParseMode				string `json:"parse_mode"`
	DisableWebPagePreview	bool   `json:"disable_web_page_preview"`
	DisableNotification		bool   `json:"disable_notification"`
	ReplyToMessageID		int64  `json:"reply_to_message_id"`

	/* XXX */
	ReplyMarkup				string `json:"-,"`
}

type TelegramForwardMessage struct {
	ChatID				string `json:"chat_id"`
	FromChatID			string `json:"from_chat_id"`
	DisableNotification	bool   `json:"disable_notification"`
	MessageID			int64  `json:"message_id"`
}

type TelegramSendPhoto struct {
	ChatID				string `json:"chat_id"`
	Photo				string `json:"photo"`
	Caption				string `json:"caption"`
	DisableNotification	bool   `json:"disable_notification"`
	ReplyToMessageID	int64  `json:"reply_to_message_id"`

	/* XXX */
	ReplyMarkup				string `json:"-,"`
}

type TelegramSendAudio struct {
	ChatID				string `json:"chat_id"`
	Audio				string `json:"audio"`
	Caption				string `json:"caption"`
	Duration			int64  `json:"duration"`
	Performer			string `json:"performer"`
	Title				string `json:"title"`
	DisableNotification	bool   `json:"disable_notification"`
	ReplyToMessageID	int64  `json:"reply_to_message_id"`

	/* XXX */
	ReplyMarkup				string `json:"-,"`
}

type TelegramSendDocument struct {
	ChatID				string `json:"chat_id"`
	Document			string `json:"document"`
	Caption				string `json:"caption"`
	Title				string `json:"title"`
	DisableNotification	bool   `json:"disable_notification"`
	ReplyToMessageID	int64  `json:"reply_to_message_id"`

	/* XXX */
	ReplyMarkup				string `json:"-,"`
}

type TelegramSendSticker struct {
	ChatID				string `json:"chat_id"`
	Sticker				string `json:"sticker"`
	Caption				string `json:"caption"`
	Title				string `json:"title"`
	DisableNotification	bool   `json:"disable_notification"`
	ReplyToMessageID	int64  `json:"reply_to_message_id"`

	/* XXX */
	ReplyMarkup				string `json:"-,"`
}

type TelegramSendVideo struct {
	ChatID				string `json:"chat_id"`
	Video				string `json:"video"`
	Duration			int64  `json:"duration"`
	Width				int64  `json:"width"`
	Height				int64  `json:"height"`
	Caption				string `json:"caption"`
	DisableNotification	bool   `json:"disable_notification"`
	ReplyToMessageID	int64  `json:"reply_to_message_id"`

	/* XXX */
	ReplyMarkup				string `json:"-,"`
}

type TelegramSendVoice struct {
	ChatID				string `json:"chat_id"`
	Voice				string `json:"voice"`
	Caption				string `json:"caption"`
	Duration			int64  `json:"duration"`
	Title				string `json:"title"`
	DisableNotification	bool   `json:"disable_notification"`
	ReplyToMessageID	int64  `json:"reply_to_message_id"`

	/* XXX */
	ReplyMarkup				string `json:"-,"`
}

type TelegramSendVideoNote struct {
	ChatID				string `json:"chat_id"`
	VideoNote			string `json:"video_note"`
	Duration			int64  `json:"duration"`
	Length				int64  `json:"length"`
	Height				int64  `json:"height"`
	Caption				string `json:"caption"`
	DisableNotification	bool   `json:"disable_notification"`
	ReplyToMessageID	int64  `json:"reply_to_message_id"`

	/* XXX */
	ReplyMarkup				string `json:"-,"`
}

type TelegramSendLocation struct {
	ChatID				string  `json:"chat_id"`
	Latitude			float64 `json:"latitude"`
	Longitude			float64 `json:"longitude"`
	DisableNotification	bool    `json:"disable_notification"`
	ReplyToMessageID	int64   `json:"reply_to_message_id"`

	/* XXX */
	ReplyMarkup				string `json:"-,"`
}

type TelegramSendVenue struct {
	ChatID				string  `json:"chat_id"`
	Latitude			float64 `json:"latitude"`
	Longitude			float64 `json:"longitude"`
	Title				string  `json:"title"`
	Address				string  `json:"address"`
	FoursquareID		string  `json:"foursquare_id"`
	DisableNotification	bool    `json:"disable_notification"`
	ReplyToMessageID	int64   `json:"reply_to_message_id"`

	/* XXX */
	ReplyMarkup				string `json:"-,"`
}

type TelegramSendContact struct {
	ChatID				string `json:"chat_id"`
	PhoneNumber			string `json:"phone_number"`
	FirstName			string `json:"first_name"`
	LastName			string `json:"last_name"`
	FoursquareID		string `json:"foursquare_id"`
	DisableNotification	bool   `json:"disable_notification"`
	ReplyToMessageID	int64  `json:"reply_to_message_id"`

	/* XXX */
	ReplyMarkup				string `json:"-,"`
}

type TelegramSendChatAction struct {
	ChatID				string `json:"chat_id"`
	Action				string `json:"action"`
}

type TelegramInlineQuery struct {
	ID			string           `json:"id"`
	From		TelegramUser     `json:"from"`
	Location	TelegramLocation `json:"location"`
	Query		string           `json:"query"`
	Offset		string           `json:"offset"`
}

type TelegramChosenInlineResult struct {
	ResultID		string           `json:"result_id"`
	From			TelegramUser     `json:"from"`
	Location		TelegramLocation `json:"location"`
	InlineMessageID	string           `json:"inline_message_id"`
	Query			string           `json:"query"`
}

type TelegramCallbackQuery struct {
	ID				string          `json:"id"`
	From			TelegramUser    `json:"from"`
	Message			TelegramMessage `json:"message"`
	InlineMessageID	string          `json:"inline_message_id"`
	ChatInstance	string          `json:"chat_instance"`
	Data			string          `json:"data"`
	GameShortName	string          `json:"game_short_name"`
}

type TelegramShippingQuery struct {
	ID				string                  `json:"id"`
	From			TelegramUser            `json:"from"`
	InvoicePayload	string                  `json:"invoice_payload"`
	ShippingAddress	TelegramShippingAddress `json:"shipping_address"`
}

type TelegramPreCheckoutQuery struct {
	ID					string            `json:"id"`
	From				TelegramUser      `json:"from"`
	Currency			string            `json:"currency"`
	TotalAmount			int64             `json:"total_amount"`
	InvoicePayload		string            `json:"invoice_payload"`
	ShippingOptionID	string            `json:"shipping_option_id"`
	OrderInfo			TelegramOrderInfo `json:"order_info"`
}

type TelegramUpdate struct {
	UpdateID			int64                      `json:"update_id"`
	Message				TelegramMessage            `json:"message"`
	EditedMessage		TelegramMessage            `json:"edited_message"`
	ChannelPost			TelegramMessage            `json:"channel_post"`
	EditedChannelPost	TelegramMessage            `json:"edited_channel_post"`
	InlineQuery			TelegramInlineQuery	       `json:"inline_query"`
	ChosenInlineResult	TelegramChosenInlineResult `json:"chosen_inline_result"`
	CallbackQuery		TelegramCallbackQuery      `json:"callback_query"`
	ShippingQuery		TelegramShippingQuery      `json:"shipping_query"`
	PreCheckoutQuery	TelegramPreCheckoutQuery   `json:"pre_checkout_query"`
}

type TelegramUpdateResponse struct {
	OK			bool             `json:"ok"`
	Description	string           `json:"description"`
	Updates		[]TelegramUpdate  `json:"result"`
}

type TelegramGetUpdates struct {
	Offset			int64    `json:"offset"`
	Limit			int64	 `json:"limit"`
	Timeout			int64    `json:"timeout"`
	AllowedUpdates	[]string `json:"allowed_updates"`
}

type TelegramNotifier struct {
	ClusterName string
	Url         string
	Token       string
	ChatID      string
	NotifName   string
	Enabled     bool
	Initialized bool
}

func (telegram *TelegramNotifier) GetURL() string {

	proto := "http"
	u := strings.TrimSpace(strings.ToLower(telegram.Url))
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

	return fmt.Sprintf("%s://%s%s/bot%s", proto, host, portstr, telegram.Token)
}

func (telegram *TelegramNotifier) GetMe(tuptr *TelegramUser) bool {

	checkURL := fmt.Sprintf("%s/getMe", telegram.GetURL())

	req, err := http.NewRequest("GET", checkURL, nil)
	if err != nil {
		log.Println("NewRequest: ", err)
		return false
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var r TelegramUserResponse
	err = decoder.Decode(&r)
	if err != nil {
		log.Println("Decode: ", err)
		return false
	}

	ret := false
	if r.OK {
		ret = true
	}

	if tuptr != nil {
		*tuptr = r.User
	}

	return ret
}

func (telegram *TelegramNotifier) SendMessage(who string, text string, tmptr *TelegramMessage) bool {
	if who == "" || text == "" {
		return false
	}

	sendMessageURL := fmt.Sprintf("%s/sendMessage", telegram.GetURL())

	message := TelegramSendMessage{ChatID: who, Text: text}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(message)

	req, err := http.NewRequest("POST", sendMessageURL, buf)
	if err != nil {
		log.Println("NewRequest: ", err)
		return false
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Do: ", err)
		return false
	}

	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var r TelegramMessageResponse
	err = decoder.Decode(&r)
	if err != nil {
		log.Println("Decode: ", err)
		return false
	}

	ret := false
	if r.OK {
		ret = true
	}

	if tmptr != nil {
		*tmptr = r.Message
	}

	return ret
}

func (telegram *TelegramNotifier) Init() bool {

	if telegram.Initialized == true {
		return true
	}

	if !telegram.GetMe(nil) {
		log.Println("Telegram: unable to verify bot!")
		return false
	}

	telegram.Initialized = true
	return true
}

func (telegram *TelegramNotifier) NotifierName() string {
	return "telegram"
}

func (telegram *TelegramNotifier) Copy() Notifier {
	notifier := *telegram
	return &notifier
}

func (telegram *TelegramNotifier) Notify(messages Messages) bool {
	if !telegram.Init() {
		return false
	}

	return telegram.NotifyDetailed(messages)
}

func (telegram *TelegramNotifier) NotifyDetailed(messages Messages) bool {

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
	pretext := fmt.Sprintf("%s %s is *%s*", emoji, telegram.ClusterName, overallStatus)

	detailedBody := ""
	detailedBody += fmt.Sprintf("*Changes:* Fail = %d, Warn = %d, Pass = %d",
		fail, warn, pass)
	detailedBody += fmt.Sprintf("\n")

	for _, message := range messages {
		detailedBody += fmt.Sprintf("\n*[%s:%s]* %s is *%s.*",
			message.Node, message.Service, message.Check, message.Status)
		detailedBody += fmt.Sprintf("\n`%s`", strings.TrimSpace(message.Output))
	}

	text := fmt.Sprintf("%s\n%s\n%s\n\n", title, pretext, detailedBody)

	return telegram.SendMessage(telegram.ChatID, text, nil)
}
