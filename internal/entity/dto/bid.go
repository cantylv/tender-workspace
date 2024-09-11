package dto

// INPUT DTO (REQUEST BODY)
type BidInput struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Status          string `json:"status"`
	TenderID        int    `json:"tenderId"`
	OrganizationID  int    `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}

type BidUpdateDataInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// OUTPUT DTO (RESPONSE BODY)
type BidOutput struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	AuthorType string `json:"authorType"`
	AuthorID   int    `json:"authorId"`
	Version    int    `json:"version"`
	CreatedAt  string `json:"createdAt"`
}
