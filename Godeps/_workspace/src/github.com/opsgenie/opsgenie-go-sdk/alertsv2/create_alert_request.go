package alertsv2

import "net/url"

type CreateAlertRequest struct {
	Message     string `json:"message,omitempty"`
	Alias       string `json:"alias,omitempty"`
	Description string `json:"description,omitempty"`
	Teams       []TeamRecipient `json:"teams,omitempty"`
	VisibleTo   []Recipient `json:"visibleTo,omitempty"`
	Actions     []string `json:"actions,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Details     map[string]string `json:"details,omitempty"`
	Entity      string `json:"entity,omitempty"`
	Source      string `json:"source,omitempty"`
	Priority    Priority `json:"priority,omitempty"`
	User        string `json:"user,omitempty"`
	Note        string `json:"note,omitempty"`
	ApiKey      string `json:"-"`
}

func (r *CreateAlertRequest) GenerateUrl() (string, url.Values, error) {
	return "/v2/alerts", nil, nil
}

func (r *CreateAlertRequest) GetApiKey() string {
	return r.ApiKey
}

func (r *CreateAlertRequest) Init() {
	if r.Teams != nil {
		var convertedTeams []TeamRecipient
		for _, t := range r.Teams {
			recipient := &RecipientDTO{
				Id:   t.getID(),
				Name: t.getName(),
				Type: "team",
			}

			convertedTeams = append(convertedTeams, recipient)
		}
		r.Teams = convertedTeams
	}

	if r.VisibleTo != nil {
		var convertedVisibleTo []Recipient
		for _, r := range r.VisibleTo {
			switch r.(type) {
			case *Team:
				{
					team := r.(*Team)
					recipient := &RecipientDTO{
						Id:   team.ID,
						Name: team.Name,
						Type: "team",
					}
					convertedVisibleTo = append(convertedVisibleTo, recipient)
				}
			case *User:
				{
					user := r.(*User)
					recipient := &RecipientDTO{
						Id:       user.ID,
						Username: user.Username,
						Type:     "user",
					}
					convertedVisibleTo = append(convertedVisibleTo, recipient)
				}
			}
		}
		r.VisibleTo = convertedVisibleTo
	}
}
