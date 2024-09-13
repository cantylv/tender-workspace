package entity

import "time"

type Organization struct {
	ID          int
	Name        string
	Description string
	Type        string
	CreatedAt   time.Time
}

type OrganizationResponsible struct {
	OrganizationID int
	UserID         int
}
