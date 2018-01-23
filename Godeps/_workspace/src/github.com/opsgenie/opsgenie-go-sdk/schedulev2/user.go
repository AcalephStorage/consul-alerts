package schedulev2

type User struct {
	ID			string		`json:"id"`
	Username	string		`json:"username"`
	Type		UserType	`json:"type"`
}

const (
	UserUserType	UserType = "user"
	NoneUserType    UserType = "none"
)

type UserType string
