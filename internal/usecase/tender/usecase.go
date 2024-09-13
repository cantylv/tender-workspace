package tender

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
	tqp "tender-workspace/internal/entity/dto/queries/tenders"
	"tender-workspace/internal/repo/organization"
	"tender-workspace/internal/repo/tender"
	"tender-workspace/internal/repo/user"
	mc "tender-workspace/internal/utils/myconstants"
	"tender-workspace/internal/utils/myerrors"
	e "tender-workspace/internal/utils/myerrors"
)

type CreateTenderData struct {
	Name            string
	Description     string
	Type            string
	Status          string
	OrganizationID  int
	CreatorUsername string
}

type Usecase interface {
	GetTenders(ctx context.Context, params *tqp.ListTenders) ([]*ent.Tender, error)
	CreateTender(ctx context.Context, initData *dto.TenderInput) (*ent.Tender, error)
	GetUserTenders(ctx context.Context, params *tqp.ListUserTenders) ([]*ent.Tender, error)
	GetTenderStatus(ctx context.Context, params *tqp.TenderStatus) (string, error)
	UpdateTenderStatus(ctx context.Context, params *tqp.UpdateTenderStatus) (*ent.Tender, error)
	UpdateTender(ctx context.Context, updateData *dto.TenderUpdateDataInput, params *tqp.TenderUpdate) (*ent.Tender, error)

	// RollbackTender откатывает параметры тендера к указанной версии
	// RollbackTender(ctx context.Context, tenderID int, version int, username string) (Tender, error)
}

var _ Usecase = (*UsecaseLayer)(nil)

type UsecaseLayer struct {
	repoTenders      tender.Repo
	repoUser         user.Repo
	repoOrganization organization.Repo
}

func NewUsecaseLayer(repoTenders tender.Repo, repoUser user.Repo, repoOrganization organization.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repoTenders:      repoTenders,
		repoUser:         repoUser,
		repoOrganization: repoOrganization,
	}
}

func (u *UsecaseLayer) GetTenders(ctx context.Context, params *tqp.ListTenders) ([]*ent.Tender, error) {
	tenders, err := u.repoTenders.GetAll(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return tenders, nil
}

func (u *UsecaseLayer) CreateTender(ctx context.Context, initData *dto.TenderInput) (*ent.Tender, error) {
	// check validation of req body fields
	tenderStatus := strings.ToLower(initData.Status)
	if tenderStatus != "created" {
		return nil, e.ErrBadStatusCreate
	}
	runes := []rune(tenderStatus)
	initData.Status = strings.ToUpper(string(runes[0])) + string(runes[1:len(tenderStatus)])

	serviceType := strings.ToLower(initData.Type)
	if _, ok := mc.AvaliableServiceType[serviceType]; !ok {
		return nil, e.ErrQPServiceType
	}
	runes = []rune(serviceType)
	initData.Type = strings.ToUpper(string(runes[0])) + string(runes[1:len(serviceType)])
	// get user id
	userData, err := u.repoUser.GetData(ctx, initData.CreatorUsername)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrUserExist
		}
		return nil, err
	}
	// check org existing
	_, err = u.repoOrganization.Get(ctx, initData.OrganizationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrOrganizationExist
		}
		return nil, err
	}
	isResponsible, err := u.repoOrganization.IsUserResponsible(ctx, initData.OrganizationID, userData.ID)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, e.ErrResponsibilty
	}
	tenderProps := newTender(userData, initData)
	t, err := u.repoTenders.Create(ctx, tenderProps)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (u *UsecaseLayer) GetUserTenders(ctx context.Context, params *tqp.ListUserTenders) ([]*ent.Tender, error) {
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
	var userTenders []*ent.Tender
	for _, organizationId := range organizationsIds {
		tenders, err := u.repoTenders.GetOrganizationTenders(ctx, organizationId)
		if err != nil {
			return nil, err
		}
		userTenders = append(userTenders, tenders...)
	}

	if params.Offset > len(userTenders) {
		return nil, e.ErrBigInterval
	}
	if params.Offset+params.Limit > len(userTenders) {
		return userTenders[params.Offset:], nil
	}
	return userTenders[params.Offset : params.Offset+params.Limit], nil
}

func (u *UsecaseLayer) GetTenderStatus(ctx context.Context, params *tqp.TenderStatus) (string, error) {
	// check that tender exists
	tender, err := u.repoTenders.GetTender(ctx, params.TenderID)
	if err != nil {
		return "", e.ErrNoTenders
	}
	// get user id
	userData, err := u.repoUser.GetData(ctx, params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", e.ErrUserExist
		}
		return "", err
	}
	// check if user is responsible for the organization
	isResponsible, err := u.repoOrganization.IsUserResponsible(ctx, userData.ID, tender.OrganizationID)
	if err != nil {
		return "", err
	}
	if !isResponsible {
		return "", e.ErrResponsibilty
	}
	return tender.Status, nil
}

func (u *UsecaseLayer) UpdateTenderStatus(ctx context.Context, params *tqp.UpdateTenderStatus) (*ent.Tender, error) {
	// check that tender exists
	tender, err := u.repoTenders.GetTender(ctx, params.TenderID)
	if err != nil {
		return nil, e.ErrNoTenders
	}
	// get user id
	userData, err := u.repoUser.GetData(ctx, params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrUserExist
		}
		return nil, err
	}
	// check if user is responsible for the organization
	isResponsible, err := u.repoOrganization.IsUserResponsible(ctx, userData.ID, tender.OrganizationID)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, e.ErrResponsibilty
	}
	// update status
	tender, err = u.repoTenders.ChangeStatus(ctx, params.TenderID, params.Status)
	if err != nil {
		return nil, err
	}
	return tender, err
}

func (u *UsecaseLayer) UpdateTender(ctx context.Context, updateData *dto.TenderUpdateDataInput, params *tqp.TenderUpdate) (*ent.Tender, error) {
	serviceType := strings.ToLower(updateData.ServiceType)
	if _, ok := mc.AvaliableServiceType[serviceType]; !ok {
		return nil, myerrors.ErrQPServiceType
	}
	runes := []rune(serviceType)
	updateData.ServiceType = strings.ToUpper(string(runes[0])) + string(serviceType[1:len(runes)])
	// check that tender exists
	tender, err := u.repoTenders.GetTender(ctx, params.TenderID)
	if err != nil {
		return nil, e.ErrNoTenders
	}
	// get user id
	userData, err := u.repoUser.GetData(ctx, params.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrUserExist
		}
		return nil, err
	}
	// check if user is responsible for the organization
	isResponsible, err := u.repoOrganization.IsUserResponsible(ctx, userData.ID, tender.OrganizationID)
	if err != nil {
		return nil, err
	}
	if !isResponsible {
		return nil, e.ErrResponsibilty
	}
	// update status
	tenderData := newUpdateTenderData(updateData)
	tenderProps := newUpdateTenderProps(params, userData)
	tender, err = u.repoTenders.Update(ctx, tenderData, tenderProps)
	if err != nil {
		return nil, err
	}
	return tender, err
}
