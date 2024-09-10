package entity

type Tender struct {
	ID              int    `db:"id"`
	OrganizationID  int    `db:"organization_id"`
	CreatorUsername string `db:"creator_username"`
	ServiceType     string `db:"service_type"`
	Name            string `db:"name"`
	Description     string `db:"description"`
	Status          string `db:"status"`
	Version         int    `db:"version"`
}

type UpdateTenderData struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	ServiceType string `db:"service_type"`
}
