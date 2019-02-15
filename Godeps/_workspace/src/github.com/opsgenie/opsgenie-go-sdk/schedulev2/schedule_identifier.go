package schedulev2

import (
	"errors"
	"net/url"
)

type ScheduleIdentifier struct {
	ID   string
	Name string
}

func (request *ScheduleIdentifier) GenerateUrl() (string, url.Values, error) {
	baseUrl := "/v2/schedules/"
	params := url.Values{}

	if request.Name != "" {
		baseUrl += request.Name
		params.Add("scheduleIdentifierType", "name")
	} else if request.ID != "" {
		baseUrl += request.ID
		params.Add("scheduleIdentifierType", "id")
	} else {
		return "", nil, errors.New("Identifier should be ID or Name")
	}

	return baseUrl, params, nil
}
