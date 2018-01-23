package userv2

const (
	UsernameSortField   Sort = "username"
	FullNameSortField   Sort = "fullName"
	InsertedAtSortField Sort = "insertedAt"
)

// Sort is a name of user field, which is used for sort users.
type Sort string

func (s Sort) IsValid() bool {
	return s == UsernameSortField || s == FullNameSortField || s == InsertedAtSortField
}
