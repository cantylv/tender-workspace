package bids

import (
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
	bqp "tender-workspace/internal/entity/dto/queries/bids"
	b "tender-workspace/internal/repo/bids"
)

func newBid(initData *dto.BidInput, user *ent.Employee) *ent.Bid {
	return &ent.Bid{
		Name:           initData.Name,
		Description:    initData.Description,
		Status:         initData.Status,
		Version:        1,
		TenderID:       initData.TenderID,
		OrganizationID: initData.OrganizationID,
		CreatorID:      user.ID,
	}
}

func newUpdateBidProps(updateProps *bqp.UpdateBidData, updateData *dto.BidUpdateDataInput) *b.UpdateBid {
	return &b.UpdateBid{
		BidID:       updateProps.BidID,
		Name:        updateData.Name,
		Description: updateData.Description,
	}
}
