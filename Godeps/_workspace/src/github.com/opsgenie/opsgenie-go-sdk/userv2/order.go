package userv2

const (
	AscSortType  Order = "asc"
	DescSortType Order = "desc"
)

// Order is a type of sort.
type Order string

func (o Order) IsValid() bool {
	return o == AscSortType || o == DescSortType
}
