package dto

// INPUT DATA FLOW
type OrganizationInput struct {
	Name        string `json:"name" valid:"name"`
	Description string `json:"description" valid:"-"`
	Type        string `json:"type" valid:"-"`
}

// OUTPUT DATA FLOW
type OrganizationOutput struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	CreatedAt   string `json:"createdAt"`
}
