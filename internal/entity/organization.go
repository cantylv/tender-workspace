package entity

type Organization struct {
	ID          int
	Name        string
	Description string
	Type        string
}

type OrganizationResponsible struct {
	OrganizationID int
	UserID         int
}
