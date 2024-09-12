package user

import (
	"context"
	"fmt"
	ent "tender-workspace/internal/entity"
	mc "tender-workspace/internal/utils/myconstants"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Repo interface {
	GetData(ctx context.Context, username string) (*ent.Employee, error)
	GetUserOrganizationsIds(ctx context.Context, userId int) ([]int, error)
	IsResponsible(ctx context.Context, resp *ent.OrganizationResponsible) error
}

type RepoLayer struct {
	Client *pgx.Conn
	Logger *zap.Logger
}

var _ Repo = (*RepoLayer)(nil)

func NewRepoLayer(client *pgx.Conn, logger *zap.Logger) *RepoLayer {
	return &RepoLayer{
		Client: client,
		Logger: logger,
	}
}

var (
	sqlRowCheckResponsibility = `SELECT 1 FROM organization_responsible WHERE organization_id=$1 AND user_id=$2`
)

func (r *RepoLayer) GetData(ctx context.Context, username string) (*ent.Employee, error) {
	row := r.Client.QueryRow(ctx, "SELECT * FROM employee WHERE email=$1", username)
	var user ent.Employee
	err := row.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *RepoLayer) GetUserOrganizationsIds(ctx context.Context, userId int) ([]int, error) {
	rows, err := r.Client.Query(ctx, `SELECT organization_id FROM organization_responsible WHERE user_id=$1`, userId)
	if err != nil {
		return nil, err
	}

	var ids []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			requestId := ctx.Value(mc.RequestID).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(mc.RequestID, requestId))
			continue
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *RepoLayer) IsResponsible(ctx context.Context, resp *ent.OrganizationResponsible) error {
	row := r.Client.QueryRow(ctx, sqlRowCheckResponsibility, resp.OrganizationID, resp.UserID)
	var isExist int
	err := row.Scan(&isExist)
	if err != nil {
		return err
	}
	return nil
}
