package dto

// INPUT DTO (REQUEST BODY)
type BidInput struct {
	Name            string `json:"name" valid:"name"`
	Description     string `json:"description" valid:"description"`
	Status          string `json:"status" valid:"-"`
	TenderID        int    `json:"tenderId" valid:"-"`
	OrganizationID  int    `json:"organizationId" valid:"-"`
	CreatorUsername string `json:"creatorUsername" valid:"username"`
}

type BidUpdateDataInput struct {
	Name        string `json:"name" valid:"name"`
	Description string `json:"description" valid:"description"`
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

type BidStatus struct {
	Status string `json:"status"`
}
