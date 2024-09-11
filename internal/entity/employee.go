package entity

import "time"

type Employee struct {
	ID        int
	Username  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
