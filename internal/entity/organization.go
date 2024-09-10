package entity

type Organization struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Type        string `db:"type"`
}

type OrganizationResponsible struct {
	ID             int `db:"id"`
	OrganizationID int `db:"organization_id"`
	UserID         int `db:"user_id"`
}
