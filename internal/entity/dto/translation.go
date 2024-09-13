package dto

import (
	ent "tender-workspace/internal/entity"
	f "tender-workspace/internal/utils/functions"
)

func NewArrayTenderOutput(tenders []*ent.Tender) []*TenderOutput {
	res := make([]*TenderOutput, 0, len(tenders))
	for _, tender := range tenders {
		tenderOutput := NewTenderOutput(tender)
		res = append(res, tenderOutput)
	}
	return res
}

func NewTenderOutput(tenders *ent.Tender) *TenderOutput {
	return &TenderOutput{
		ID:          tenders.ID,
		Name:        tenders.Name,
		Description: tenders.Description,
		Status:      tenders.Status,
		Type:        tenders.Type,
		Version:     tenders.Version,
		CreatedAt:   f.FormatTime(tenders.CreatedAt),
	}
}

func NewArrayBidOutput(bids []*ent.Bid) []*BidOutput {
	res := make([]*BidOutput, 0, len(bids))
	for _, bid := range bids {
		bidOutput := NewBidOutput(bid)
		res = append(res, bidOutput)
	}
	return res
}

func NewBidOutput(bid *ent.Bid) *BidOutput {
	return &BidOutput{
		ID:         bid.ID,
		Name:       bid.Name,
		Status:     bid.Status,
		AuthorType: bid.AuthorType,
		AuthorID:   bid.CreatorID,
		Version:    bid.Version,
		CreatedAt:  f.FormatTime(bid.CreatedAt),
	}
}
