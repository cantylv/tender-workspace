package entity

import "time"

type Bid struct {
	ID             int
	Name           string
	Description    string
	Status         string
	Version        int
	TenderID       int
	CreatorID      int
	AuthorType     string
	OrganizationID int
	CreatedAt      time.Time
}

type BidUpdateData struct {
	Name        string
	Description string
}
