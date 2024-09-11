package entity

type Bid struct {
	ID          int
	Name        string
	Description string
	Status      string
	Version     int
	TenderID    int
	CreatorID   int
	CreatedAt   string
}

type BidUpdateData struct {
	Name        string
	Description string
}
