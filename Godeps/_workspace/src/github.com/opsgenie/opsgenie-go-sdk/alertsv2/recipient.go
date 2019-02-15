package alertsv2

type Recipient interface {
	SetID(id string)
}

type TeamRecipient interface {
	SetID(id string)
	SetName(name string)
	getID() string
	getName() string
}

type Team struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (t *Team) SetID(id string) {
	t.ID = id
}

func (t *Team) SetName(name string) {
	t.Name = name
}

func (t *Team) getID() string {
	return t.ID
}

func (t *Team) getName() string {
	return t.Name
}

type User struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

func (u *User) SetID(id string) {
	u.ID = id
}

func (u *User) SetUsername(username string) {
	u.Username = username
}

type Escalation struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type RecipientDTO struct {
	Id       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
	Type     string `json:"type"`
}

func (r *RecipientDTO) SetID(id string) {
	r.Id = id
}

func (r *RecipientDTO) SetName(id string) {
	r.Id = id
}

func (r *RecipientDTO) getName() string {
	return r.Name
}

func (r *RecipientDTO) getID() string {
	return r.Id
}
