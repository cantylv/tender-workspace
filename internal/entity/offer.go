package entity

type Offer struct {
	ID             int `db:"id"`
	TenderID       int `db:"tender_id"`
	PerformerID    int `db:"performer_id"`
	OfferVersionID int `db:"offer_version_id"`
}

type OfferVersion struct {
	ID             int    `db:"id"`
	OrganizationID int    `db:"organization_id"`
	Message        string `db:"message"`
	Price          int    `db:"price"`
	Version        int    `db:"version"`
}
