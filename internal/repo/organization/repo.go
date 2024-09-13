package organization

import (
	"context"
	ent "tender-workspace/internal/entity"
	e "tender-workspace/internal/utils/myerrors"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Repo interface {
	IsUserResponsible(ctx context.Context, resp *ent.OrganizationResponsible) error
	MakeUserResponsible(ctx context.Context, resp *ent.OrganizationResponsible) error
}

type RepoLayer struct {
	Client *pgx.Conn
}

var _ Repo = (*RepoLayer)(nil)

func NewRepoLayer(client *pgx.Conn, logger *zap.Logger) *RepoLayer {
	return &RepoLayer{
		Client: client,
	}
}

var (
	sqlRowCheckResponsibility = `SELECT 1 FROM organization_responsible WHERE organization_id=$1 AND user_id=$2`
	sqlRowMakeResponsibility  = `INSERT INTO organization_responsible(organization_id, user_id) VALUES($1, $2)`
)

func (r *RepoLayer) IsUserResponsible(ctx context.Context, resp *ent.OrganizationResponsible) error {
	row := r.Client.QueryRow(ctx, sqlRowCheckResponsibility, resp.OrganizationID, resp.UserID)
	var isExist int
	err := row.Scan(&isExist)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoLayer) MakeUserResponsible(ctx context.Context, resp *ent.OrganizationResponsible) error {
	tag, err := r.Client.Exec(ctx, sqlRowMakeResponsibility, resp.OrganizationID, resp.UserID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return e.ErrNoRowsAffected
	}

	return nil
}
