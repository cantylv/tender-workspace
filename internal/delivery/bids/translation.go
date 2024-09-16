package bids

import bqp "tender-workspace/internal/entity/dto/queries/bids"

func newUpdateBidStatusParams(params *bqp.SubmitDecision) *bqp.UpdateBidStatus {
	return &bqp.UpdateBidStatus{
		BidID:    params.BidID,
		Status:   params.Decision,
		Username: params.Username,
	}
}
