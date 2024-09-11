package entity

type Tender struct {
	ID             int
	Name           string
	Description    string
	Type           string
	Status         string
	Version        int
	OrganizationID int
	CreatorID      int
	CreatedAt      string
}

type UpdateTenderData struct {
	Name        string
	Description string
	Type        string
}
