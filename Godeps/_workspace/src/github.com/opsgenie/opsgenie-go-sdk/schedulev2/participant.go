package schedulev2

import "errors"

type Participant interface {}

type participant struct {
	ID 			string			`json:"id"`
	Username 	string			`json:"username"`
	Name	 	string			`json:"name"`
	Type 		ParticipantType	`json:"type"`
}

/*
** If participants' type is escalation or team, you can use name or id fields for referring.
** Otherwise (type is user), we use username or id for referencing
*/
func NewParticipant(participantType ParticipantType, ID string, name string, username string) (Participant, error){
	if participantType == "" {
		return nil, errors.New("Participant Type must not be empty.")
	}

	if participantType == UserParticipant {

		if ID != "" {
			return participant{ID:ID, Type:participantType}, nil
		} else if username != "" {
			return participant{Username:username, Type:participantType}, nil
		} else {
			return nil, errors.New("Username or ID must not be empty for UserParticipant")
		}

	} else if participantType == TeamParticipant || participantType == EscalationParticipant {

		if ID != "" {
			return participant{ID:ID, Type:participantType}, nil
		} else if name != "" {
			return participant{Name:name, Type:participantType}, nil
		} else {
			return nil, errors.New("Name or ID must not be empty for TeamParticipant or EscalationParticipant")
		}

	} else if participantType  == NoneParticipant {

		return participant{Type:participantType}, nil

	} else {

		return nil, errors.New("ParticipantType must not be empty")
	}
}

const (
	UserParticipant 		ParticipantType = "user"
	TeamParticipant 		ParticipantType = "team"
	EscalationParticipant 	ParticipantType	= "escalation"
	NoneParticipant 		ParticipantType = "none"
)
type ParticipantType string
