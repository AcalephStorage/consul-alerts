package alerts

import "strconv"

type CreateAlertResponse struct	{
	Message	string	`json:"message"`
	AlertId	string	`json:"alertId"`
	Status	string	`json:"status"`
	Code	int		`json:"code"`
}

type CloseAlertResponse struct {
	Status 	string	`json:"status"`
	Code 	int 	`json:"code"`
}

type DeleteAlertResponse struct {
	Status 	string 	`json:"status"`
	Code 	int 	`json:"code"`
}

type GetAlertResponse struct {
	Tags 	[]string				`json:"tags"`
	Count 	int 					`json:"count"`
	Status 	string					`json:"status"`
	Teams 	[]string				`json:"teams"`
	Recipients	[]string			`json:"recipients"`
	TinyId 	string					`json:"tinyId"`
	Alias	string					`json:"alias"`
	Entity	string					`json:"entity"`
	Id 		string					`json:"id"`
	UpdatedAt	int					`json:"updatedAt"`
	Message 	string				`json:"message"`
	Details		map[string]string	`json:"details"`
	Source 	string					`json:"source"`
	Description	string				`json:"description"`
	CreatedAt int					`json:"createdAt"`
	IsSeen bool						`json:"isSeen"`
	Acknowledged bool				`json:"acknowledged"`
	Owner	string					`json:"owner"`
	Actions []string				`json:"actions"`
	SystemData map[string]string	`json:"systemData"`
}

func (res *GetAlertResponse) GetIntegrationType() string {
	if val, ok := res.SystemData["integrationType"]; ok {
		return val
	}
	return ""
}
func (res *GetAlertResponse) GetIntegrationId() string {
	if val, ok := res.SystemData["integrationId"]; ok {
		return val
	}
	return ""
}
func (res *GetAlertResponse) GetIntegrationName() string {
	if val, ok := res.SystemData["integrationName"]; ok {
		return val
	}
	return ""
}
func (res *GetAlertResponse) GetAckTime() int {
	if val, ok := res.SystemData["ackTime"]; ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return -1
		} else {
			return i
		}
	}
	return -1
}
func (res *GetAlertResponse) GetAcknowledgedBy() string {
	if val, ok := res.SystemData["acknowledgedBy"]; ok {
		return val
	}
	return ""
}
func (res *GetAlertResponse) GetCloseTime() int {
	if val, ok := res.SystemData["closeTime"]; ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return -1
		} else {
			return i
		}
	}
	return -1
}
func (res *GetAlertResponse) GetClosedBy() string {
	if val, ok := res.SystemData["closedBy"]; ok {
		return val
	}
	return ""
}

type ListAlertsResponse struct {
	Alerts []struct {
		Id 				string `json:"id"`
		Alias 			string `json:"alias"`
		Message 		string `json:"message"`
		Status 			string `json:"status"`
		IsSeen 			bool `json:"isSeen"`
		Acknowledged 	bool `json:"acknowledged"`
		CreatedAt 		int `json:"createdAt"`
		UpdatedAt 		int `json:"updatedAt"`
		TinyId 			string `json:"tinyId"`
	} `json:"alerts"`
}

type ListAlertNotesResponse struct {
	Took int 			`json:"took"`
	LastKey	string 		`json:"lastKey"`
	Notes []struct {
		Note string 	`json:"note"`
		Owner string 	`json:"owner"`
		CreatedAt int	`json:"createdAt"`
	} `json:"notes"`
}

type ListAlertLogsResponse struct {
	LastKey string `json:"lastKey"`
	Logs []struct {
		Log string `json:"log"`
		LogType string `json:"logType"`
		Owner string `json:"owner"`
		CreatedAt int `json:"createdAt"`
	}`json:"logs"`
}

type ListAlertRecipientsResponse struct {
	Users []struct {
		Username string `json:"username"`
		State string `json:"state"`
		Method string `json:"method"`
		StateChangedAt int `json:"stateChangedAt"`
	}`json:"users"`

	Groups []struct {
		Group map[string]struct {
			Username string `json:"username"`
			State string `json:"state"`
			Method string `json:"method"`
			StateChangedAt int `json:"stateChangedAt"`
		}
	}`json:"groups"`
}

type AcknowledgeAlertResponse struct {
	Status 	string 	`json:"status"`
	Code 	int 	`json:"code"`
}

type RenotifyAlertResponse struct {
	Status 	string 	`json:"status"`
	Code 	int 	`json:"code"`	
}

type TakeOwnershipAlertResponse struct {
	Status 	string 	`json:"status"`
	Code 	int 	`json:"code"`	
}

type AssignOwnerAlertResponse struct {
	Status 	string 	`json:"status"`
	Code 	int 	`json:"code"`	
}

type AddTeamAlertResponse struct {
	Status 	string 	`json:"status"`
	Code 	int 	`json:"code"`	
}

type AddRecipientAlertResponse struct {
	Status 	string 	`json:"status"`
	Code 	int 	`json:"code"`		
}

type AddNoteAlertResponse struct {
	Status 	string 	`json:"status"`
	Code 	int 	`json:"code"`		
}

type ExecuteActionAlertResponse struct {
	Result 	string 	`json:"result"`
	Code 	int 	`json:"code"`		
}

type AttachFileAlertResponse struct {
	Status 	string 	`json:"status"`
	Code 	int 	`json:"code"`		
}



