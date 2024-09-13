package user

import (
	"context"
	"fmt"
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/utils/functions"
	mc "tender-workspace/internal/utils/myconstants"
	"time"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Repo interface {
	Create(ctx context.Context, initData *ent.Employee) (*ent.Employee, error)
	GetData(ctx context.Context, username string) (*ent.Employee, error)
	GetUserOrganizationsIds(ctx context.Context, userId int) ([]int, error)
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
	sqlRowCreateUser          = `INSERT INTO employee (
											username, 
											first_name, 
											last_name,
											created_at,
											updated_at
										) VALUES ($1, $2, $3, $4, $4) RETURN *`
)

func (r *RepoLayer) Create(ctx context.Context, initData *ent.Employee) (*ent.Employee, error) {
	timeNow := functions.FormatTime(time.Now())
	row := r.Client.QueryRow(ctx, sqlRowCreateUser, initData.Username, initData.FirstName, initData.LastName, timeNow)
	var u ent.Employee
	err := row.Scan(&u.ID, &u.Username, &u.FirstName, &u.LastName, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *RepoLayer) GetData(ctx context.Context, username string) (*ent.Employee, error) {
	row := r.Client.QueryRow(ctx, "SELECT * FROM employee WHERE username=$1", username)
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
