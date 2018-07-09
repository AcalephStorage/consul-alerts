package userv2

const (
	// ContactExpandableField is the query parameter, which is needed to load fully contact data of user.
	ContactExpandableField = "contact"
)

// Expand is a type of expand user data, which is used for loading fully data of user.
type Expand string
