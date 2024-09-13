package dto

// INPUT DTO (REQUEST BODY) -
type TenderInput struct {
	Name            string `json:"name" valid:"name"`
	Description     string `json:"description" valid:"description"`
	Type            string `json:"serviceType" valid:"-"`
	Status          string `json:"status" valid:"-"`
	OrganizationID  int    `json:"organizationId" valid:"-"`
	CreatorUsername string `json:"creatorUsername" valid:"username"`
}

type TenderUpdateDataInput struct {
	Name        string `json:"name" valid:"name"`
	Description string `json:"description" valid:"description"`
	ServiceType string `json:"serviceType" valid:"-"`
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

type TenderStatus struct {
	Status string `json:"status"`
}
