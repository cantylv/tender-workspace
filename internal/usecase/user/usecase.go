package user

import (
	"context"
	"database/sql"
	"errors"
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/entity/dto"
	repoUser "tender-workspace/internal/repo/user"
	e "tender-workspace/internal/utils/myerrors"
)

type Usecase interface {
	GetData(ctx context.Context, username string) (*ent.Employee, error)
	Create(ctx context.Context, initData *dto.UserInput) (*ent.Employee, error)
	GetUserOrganizations(ctx context.Context, username string) ([]int, error)
}

var _ Usecase = (*UsecaseLayer)(nil)

type UsecaseLayer struct {
	repoUser repoUser.Repo
}

func NewUsecaseLayer(repoUser repoUser.Repo) *UsecaseLayer {
	return &UsecaseLayer{
		repoUser: repoUser,
	}
}

func (u *UsecaseLayer) Create(ctx context.Context, initData *dto.UserInput) (*ent.Employee, error) {
	_, err := u.repoUser.GetData(ctx, initData.Username)
	if err == nil {
		return nil, e.ErrUserAlreadyExist
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	data := newUser(initData)
	uData, err := u.repoUser.Create(ctx, data)
	if err != nil {
		return nil, err
	}
	return uData, nil
}

func (u *UsecaseLayer) GetData(ctx context.Context, username string) (*ent.Employee, error) {
	return u.repoUser.GetData(ctx, username)
}

func (u *UsecaseLayer) GetUserOrganizations(ctx context.Context, username string) ([]int, error) {
	user, err := u.repoUser.GetData(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, e.ErrUserExist
		}
		return nil, err
	}
	ids, err := u.repoUser.GetUserOrganizationsIds(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	return ids, nil
}
