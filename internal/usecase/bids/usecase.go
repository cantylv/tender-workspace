package bids

import (
	"context"
	"database/sql"
	"errors"
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
	bqp "tender-workspace/internal/entity/dto/queries/bids"
	"tender-workspace/internal/repo/bids"
	"tender-workspace/internal/repo/user"
	e "tender-workspace/internal/utils/myerrors"
)

type Usecase interface {
	CreateBid(ctx context.Context, initData *dto.BidInput) (*ent.Bid, error)
	GetUserBids(ctx context.Context, params *bqp.ListUserBids) ([]*ent.Bid, error)
	GetBidsForTender(ctx context.Context, tenderID, username string, limit, offset int) ([]*ent.Bid, error)
	GetBidStatus(ctx context.Context, params *bqp.BidStatus) (string, error)
	UpdateBidStatus(ctx context.Context, params *bqp.UpdateBidStatus) (*ent.Bid, error)
	UpdateBid(ctx context.Context, updateData *dto.BidUpdateDataInput, params *bqp.UpdateBidData) (*ent.Bid, error)
	// SubmitBidDecision(ctx context.Context, bidID, decision, username string) (Bid, error)
	// SubmitBidFeedback(ctx context.Context, bidID string, feedback BidFeedback, username string) (Bid, error)
}

type UsecaseLayer struct {
	repoBids bids.Repo
	repoUser user.Repo
}

func NewUsecaseLayer(repoBids bids.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repoBids: repoBids,
	}
}

func (u *UsecaseLayer) CreateBid(ctx context.Context, initData *dto.BidInput) (*ent.Bid, error) {
	// get user id
	userData, err := u.repoUser.GetData(ctx, initData.CreatorUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrUserExist
		}
		return nil, err
	}
	props := newBid(initData, userData)
	return u.repoBids.Create(ctx, props)
}

func (u *UsecaseLayer) GetUserBids(ctx context.Context, params *bqp.ListUserBids) ([]*ent.Bid, error) {
	// get user id
	userData, err := u.repoUser.GetData(ctx, params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrUserExist
		}
		return nil, err
	}
	// get array of organizations ids
	organizationsIds, err := u.repoUser.GetUserOrganizationsIds(ctx, userData.ID)
	if err != nil {
		return nil, e.ErrUserIsNotResponsible
	}
	var userBids []*ent.Bid
	for _, organizationId := range organizationsIds {
		tenders, err := u.repoBids.GetOrganizationBids(ctx, organizationId)
		if err != nil {
			return nil, err
		}
		userBids = append(userBids, tenders...)
	}

	if params.Offset > len(userBids) {
		return nil, e.ErrBigInterval
	}
	if params.Offset+params.Limit > len(userBids) {
		return userBids[params.Offset:], nil
	}
	return userBids[params.Offset : params.Offset+params.Limit], nil
}

func (u *UsecaseLayer) GetBidStatus(ctx context.Context, params *bqp.BidStatus) (string, error) {
	// check that tender exists
	_, err := u.repoBids.GetBid(ctx, params.BidID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", e.ErrNoBids
		}
		return "", err
	}
	// get user id
	_, err = u.repoUser.GetData(ctx, params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", e.ErrUserExist
		}
		return "", err
	}
	status, err := u.repoBids.GetStatus(ctx, params.BidID)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (u *UsecaseLayer) UpdateBidStatus(ctx context.Context, params *bqp.UpdateBidStatus) (*ent.Bid, error) {
	// check that tender exists
	_, err := u.repoBids.GetBid(ctx, params.BidID)
	if err != nil {
		return nil, e.ErrNoBids
	}
	// get user id
	_, err = u.repoUser.GetData(ctx, params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrUserExist
		}
		return nil, err
	}
	// update status
	bid, err := u.repoBids.UpdateStatus(ctx, params.BidID, params.Status)
	if err != nil {
		return nil, e.ErrNoBids
	}
	return bid, err
}

func (u *UsecaseLayer) UpdateBid(ctx context.Context, updateData *dto.BidUpdateDataInput, params *bqp.UpdateBidData) (*ent.Bid, error) {
	// check that bid exists
	_, err := u.repoBids.GetBid(ctx, params.BidID)
	if err != nil {
		return nil, e.ErrNoBids
	}
	// get user id
	_, err = u.repoUser.GetData(ctx, params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrUserExist
		}
		return nil, err
	}
	// update status
	tenderData := newUpdateBidProps(params, updateData)
	bid, err := u.repoBids.Update(ctx, tenderData)
	if err != nil {
		return nil, err
	}
	return bid, err
}
