package organization

import (
	"context"
	"database/sql"
	"errors"
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
	oqp "tender-workspace/internal/entity/dto/queries/organizations"
	repoOrg "tender-workspace/internal/repo/organization"
	repoUser "tender-workspace/internal/repo/user"
	mc "tender-workspace/internal/utils/myconstants"
	e "tender-workspace/internal/utils/myerrors"
)

type Usecase interface {
	GetAll(ctx context.Context, params *oqp.OrganizationList) ([]*ent.Organization, error)
	Create(ctx context.Context, initData *dto.OrganizationInput) (*ent.Organization, error)
	Update(ctx context.Context, updateData *dto.OrganizationInput, organizationId int) (*ent.Organization, error)
	MakeResponsible(ctx context.Context, username string, organizationID int) error
}

var _ Usecase = (*UsecaseLayer)(nil)

type UsecaseLayer struct {
	repoOrg  repoOrg.Repo
	repoUser repoUser.Repo
}

func NewUsecaseLayer(repoOrg repoOrg.Repo, repoUser repoUser.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repoOrg:  repoOrg,
		repoUser: repoUser,
	}
}

func (u *UsecaseLayer) GetAll(ctx context.Context, params *oqp.OrganizationList) ([]*ent.Organization, error) {
	return u.repoOrg.GetAll(ctx, params)
}

func (u *UsecaseLayer) Create(ctx context.Context, initData *dto.OrganizationInput) (*ent.Organization, error) {
	// check type existing
	if _, ok := mc.AvaliableOrganizationType[initData.Type]; !ok {
		return nil, e.ErrQPOrgType
	}
	dataCreate := newOrganization(initData)
	org, err := u.repoOrg.Create(ctx, dataCreate)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (u *UsecaseLayer) Update(ctx context.Context, updateData *dto.OrganizationInput, organizationId int) (*ent.Organization, error) {
	// check type existing
	if _, ok := mc.AvaliableOrganizationType[updateData.Type]; !ok {
		return nil, e.ErrQPOrgType
	}
	dataCreate := newOrganization(updateData)
	dataCreate.ID = organizationId
	org, err := u.repoOrg.Update(ctx, dataCreate)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (u *UsecaseLayer) MakeResponsible(ctx context.Context, username string, organizationID int) error {
	user, err := u.repoUser.GetData(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrUserExist
		}
		return err
	}
	_, err = u.repoOrg.Get(ctx, organizationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return e.ErrOrganizationExist
		}
		return err
	}
	isResponsible, err := u.repoOrg.IsUserResponsible(ctx, user.ID, organizationID)
	if err != nil {
		return err
	}
	if isResponsible {
		return e.ErrUserAlreadyResponsible
	}
	err = u.repoOrg.MakeUserResponsible(ctx, user.ID, organizationID)
	if err != nil {
		return err
	}
	return nil
}
