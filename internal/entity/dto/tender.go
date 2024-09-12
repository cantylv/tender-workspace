package dto

// INPUT DTO (REQUEST BODY) -
type TenderInput struct {
	Name            string `json:"name" valid:"name"`
	Description     string `json:"description" valid:"description"`
	Type            string `json:"serviceType" valid:"serviceType"`
	Status          string `json:"status" valid:"status"`
	OrganizationID  int    `json:"organizationId" valid:"organizationId"`
	CreatorUsername string `json:"creatorUsername" valid:"username"`
}

type TenderUpdateDataInput struct {
	Name        string `json:"name" valid:"tenderName"`
	Description string `json:"description" valid:"tenderDescription"`
	ServiceType string `json:"serviceType" valid:"tenderStatus"`
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
