package userv2

const (
	// OwnerRoleId is the text value of standard role "owner"
	OwnerRoleId = "Owner"
	// AdminRole is the text value of standard role "admin"
	AdminRoleId = "Admin"
	// UserRoleId is the text value of standard role "user"
	UserRoleId = "User"
	// StakeHolderRoleId is the text value of standard role "user"
	StakeHolderRoleId = "Stakeholder"
	// ObserverRoleId is the text value of standard role "user"
	ObserverRoleId = "Observer"
)

// UserRole contains data of role.
type UserRole struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
