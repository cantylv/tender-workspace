package dto

// INPUT DTO (REQUEST BODY)
type TenderInput struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	Type            string `json:"serviceType"`
	Status          string `json:"status"`
	OrganizationID  string `json:"organizationId"`
	CreatorUsername string `json:"creatorUsername"`
}

type TenderUpdateDataInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"serviceType"`
}

// OUTPUT DTO (RESPONSE BODY)
type TenderOutput struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Type        string `json:"serviceType"`
	Version     int    `json:"version"`
	CreatedAt   string `json:"createdAt"`
}
