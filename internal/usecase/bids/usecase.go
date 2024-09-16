package bids

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
	bqp "tender-workspace/internal/entity/dto/queries/bids"
	"tender-workspace/internal/repo/bids"
	"tender-workspace/internal/repo/organization"
	"tender-workspace/internal/repo/tender"
	"tender-workspace/internal/repo/user"
	mc "tender-workspace/internal/utils/myconstants"
	e "tender-workspace/internal/utils/myerrors"
)

type Usecase interface {
	CreateBid(ctx context.Context, initData *dto.BidInput) (*ent.Bid, error)
	GetUserBids(ctx context.Context, params *bqp.ListUserBids) ([]*ent.Bid, error)
	GetTenderBids(ctx context.Context, params *bqp.TenderBidList) ([]*ent.Bid, error)

	GetBidStatus(ctx context.Context, params *bqp.BidStatus) (string, error)
	UpdateBidStatus(ctx context.Context, params *bqp.UpdateBidStatus) (*ent.Bid, error)
	UpdateBid(ctx context.Context, updateData *dto.BidUpdateDataInput, params *bqp.UpdateBidData) (*ent.Bid, error)
	// SubmitBidDecision(ctx context.Context, bidID, decision, username string) (Bid, error)
	// SubmitBidFeedback(ctx context.Context, bidID string, feedback BidFeedback, username string) (Bid, error)
}

var _ Usecase = (*UsecaseLayer)(nil)

type UsecaseLayer struct {
	repoBids         bids.Repo
	repoUser         user.Repo
	repoOrganization organization.Repo
	repoTender       tender.Repo
}

func NewUsecaseLayer(repoBids bids.Repo, repoUser user.Repo, repoOrganization organization.Repo, repoTender tender.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repoBids:         repoBids,
		repoUser:         repoUser,
		repoOrganization: repoOrganization,
		repoTender:       repoTender,
	}
}

func (u *UsecaseLayer) CreateBid(ctx context.Context, initData *dto.BidInput) (*ent.Bid, error) {
	// check validation of req body fields
	bidStatus := strings.ToLower(initData.Status)
	if bidStatus != "created" {
		return nil, e.ErrBadStatusCreate
	}
	runes := []rune(bidStatus)
	initData.Status = strings.ToUpper(string(runes[0])) + string(runes[1:])
	// get user id
	userData, err := u.repoUser.GetData(ctx, initData.CreatorUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrUserExist
		}
		return nil, err
	}
	// check tender existing
	t, err := u.repoTender.GetTender(ctx, initData.TenderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrTenderExist
		}
		return nil, err
	}
	isUser := true
	if initData.OrganizationID != 0 {
		// check org existing
		_, err = u.repoOrganization.Get(ctx, initData.OrganizationID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, e.ErrOrganizationExist
			}
			return nil, err
		}
		isResponsible, err := u.repoOrganization.IsUserResponsible(ctx, userData.ID, initData.OrganizationID)
		if err != nil {
			return nil, err
		}
		if !isResponsible {
			return nil, e.ErrUserAndOrg
		}
		isUser = false
	}
	if isUser {
		has, err := u.repoBids.UserHasBid(ctx, userData.ID, initData.TenderID)
		if err != nil {
			return nil, err
		}
		if has {
			return nil, e.ErrUserAlreadyHasBid
		}
	} else {
		if t.OrganizationID == initData.OrganizationID {
			return nil, e.ErrBidYourself
		}
		has, err := u.repoBids.OrganizationHasBid(ctx, initData.OrganizationID, initData.TenderID)
		if err != nil {
			return nil, err
		}
		if has {
			return nil, e.ErrOrgAlreadyHasBid
		}
	}

	props := newBid(initData, userData)
	if isUser {
		props.AuthorType = "User"
	} else {
		props.AuthorType = "Responsible"
	}
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
	var userBids []*ent.Bid
	// get array of organizations ids
	organizationsIds, err := u.repoUser.GetUserOrganizationsIds(ctx, userData.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if len(organizationsIds) != 0 {
		for _, organizationId := range organizationsIds {
			bids, err := u.repoBids.GetOrganizationBids(ctx, organizationId)
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					continue
				}
				return nil, err
			}
			if len(bids) != 0 {
				userBids = append(userBids, bids...)
			}
		}
	}
	bs, err := u.repoBids.GetUserBids(ctx, userData.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if len(bs) != 0 {
		userBids = append(userBids, bs...)
	}
	if params.Offset >= len(userBids) {
		return nil, nil
	}
	return userBids[params.Offset:min(params.Offset+params.Limit, len(userBids))], nil
}

func (u *UsecaseLayer) GetTenderBids(ctx context.Context, params *bqp.TenderBidList) ([]*ent.Bid, error) {
	// get user id
	userData, err := u.repoUser.GetData(ctx, params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrUserExist
		}
		return nil, err
	}
	// check tender existing
	t, err := u.repoTender.GetTender(ctx, params.TenderID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrTenderExist
		}
		return nil, err
	}
	// check user responsibility
	isResponsible, err := u.repoOrganization.IsUserResponsible(ctx, userData.ID, t.OrganizationID)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, e.ErrResponsibilty
	}
	// get tender bids
	return u.repoBids.GetTenderBids(ctx, params.TenderID)
}

func (u *UsecaseLayer) GetBidStatus(ctx context.Context, params *bqp.BidStatus) (string, error) {
	// get user id
	userData, err := u.repoUser.GetData(ctx, params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", e.ErrUserExist
		}
		return "", err
	}
	// check that bid exists
	bid, err := u.repoBids.GetBid(ctx, params.BidID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", e.ErrNoBids
		}
		return "", err
	}
	// check responsibility
	isResponsible, err := u.repoOrganization.IsUserResponsible(ctx, userData.ID, bid.OrganizationID)
	if err != nil {
		return "", err
	}
	if isResponsible {
		return bid.Status, nil
	}
	// if authorType is 'User'
	if bid.CreatorID == userData.ID {
		return bid.Status, nil
	}
	// for tender side
	t, err := u.repoTender.GetTender(ctx, bid.TenderID)
	if err != nil {
		return "", err
	}
	isResponsible, err = u.repoOrganization.IsUserResponsible(ctx, userData.ID, t.OrganizationID)
	if err != nil {
		return "", err
	}
	if isResponsible {
		if bid.Status == "Published" {
			return bid.Status, nil
		}
		return "", e.ErrBadPermission
	}
	return "", e.ErrBadPermission
}

func (u *UsecaseLayer) UpdateBidStatus(ctx context.Context, params *bqp.UpdateBidStatus) (*ent.Bid, error) {
	// get user id
	userData, err := u.repoUser.GetData(ctx, params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrUserExist
		}
		return nil, err
	}
	// check that bid exists
	bid, err := u.repoBids.GetBid(ctx, params.BidID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrNoBids
		}
		return nil, err
	}
	hasCreatorPrivilege := false
	// if authorType is 'User'
	if bid.CreatorID == userData.ID {
		hasCreatorPrivilege = true
	}
	// check responsibility
	isResponsible, err := u.repoOrganization.IsUserResponsible(ctx, userData.ID, bid.OrganizationID)
	if err != nil {
		return nil, err
	}
	if isResponsible || hasCreatorPrivilege {
		statusLower := strings.ToLower(params.Status)
		if _, ok := mc.AvaliableBidStatusCreator[statusLower]; !ok {
			return nil, e.ErrSetDeprecatedStatus
		}
		runes := []rune(statusLower)
		params.Status = strings.ToUpper(string(runes[0])) + string(runes[1:])
		return u.repoBids.UpdateStatus(ctx, params.BidID, params.Status, bid.Version+1)
	}
	// for tender side
	t, err := u.repoTender.GetTender(ctx, bid.TenderID)
	if err != nil {
		return nil, err
	}
	isResponsible, err = u.repoOrganization.IsUserResponsible(ctx, userData.ID, t.OrganizationID)
	if err != nil {
		return nil, err
	}
	if isResponsible {
		statusLower := strings.ToLower(params.Status)
		if _, ok := mc.AvaliableBidStatusApprover[statusLower]; !ok {
			return nil, e.ErrSetDeprecatedStatus
		}
		runes := []rune(statusLower)
		params.Status = strings.ToUpper(string(runes[0])) + string(runes[1:])
		return u.repoBids.UpdateStatus(ctx, params.BidID, params.Status, bid.Version+1)
	}
	return nil, e.ErrBadPermission
}

func (u *UsecaseLayer) UpdateBid(ctx context.Context, updateData *dto.BidUpdateDataInput, params *bqp.UpdateBidData) (*ent.Bid, error) {
	// check that bid exists
	bid, err := u.repoBids.GetBid(ctx, params.BidID)
	if err != nil {
		return nil, e.ErrNoBids
	}
	// get user id
	user, err := u.repoUser.GetData(ctx, params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrUserExist
		}
		return nil, err
	}
	// check user responsibility
	isResponsible, err := u.repoOrganization.IsUserResponsible(ctx, user.ID, bid.OrganizationID)
	if err != nil {
		return nil, err
	}
	if isResponsible || bid.CreatorID == user.ID {
		// update status
		bidData := newUpdateBidProps(params, updateData)
		bid, err := u.repoBids.Update(ctx, bidData, bid.Version+1)
		if err != nil {
			return nil, err
		}
		return bid, err
	}
	return nil, e.ErrResponsibilty
}
