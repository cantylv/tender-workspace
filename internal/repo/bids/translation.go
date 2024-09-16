package bids

import (
	"database/sql"
	ent "tender-workspace/internal/entity"
	"time"
)

type bidDB struct {
	ID             int
	Name           string
	Description    string
	Status         string
	Version        int
	TenderID       int
	CreatorID      int
	AuthorType     string
	OrganizationID sql.NullInt32
	CreatedAt      time.Time
}

func newBidDB(bid *ent.Bid) *bidDB {
	orgId := sql.NullInt32{}
	if bid.OrganizationID != 0 {
		orgId.Valid = true
		orgId.Int32 = (int32)(bid.OrganizationID)
	} else {
		orgId.Valid = false
	}
	return &bidDB{
		ID:             bid.ID,
		Name:           bid.Name,
		Description:    bid.Description,
		Status:         bid.Status,
		Version:        bid.Version,
		TenderID:       bid.TenderID,
		CreatorID:      bid.CreatorID,
		AuthorType:     bid.AuthorType,
		OrganizationID: orgId,
		CreatedAt:      bid.CreatedAt,
	}
}

func newBid(bid *bidDB) *ent.Bid {
	orgId := 0
	if bid.OrganizationID.Valid {
		orgId = int(bid.OrganizationID.Int32)
	}
	return &ent.Bid{
		ID:             bid.ID,
		Name:           bid.Name,
		Description:    bid.Description,
		Status:         bid.Status,
		Version:        bid.Version,
		TenderID:       bid.TenderID,
		CreatorID:      bid.CreatorID,
		AuthorType:     bid.AuthorType,
		OrganizationID: orgId,
		CreatedAt:      bid.CreatedAt,
	}
}

func getArrayBidFromDB(rows []*bidDB) []*ent.Bid {
	orgs := make([]*ent.Bid, 0, len(rows))
	for _, row := range rows {
		orgs = append(orgs, newBid(row))
	}
	return orgs
}
