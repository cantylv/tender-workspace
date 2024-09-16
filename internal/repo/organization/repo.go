package organization

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	ent "tender-workspace/internal/entity"
	oqp "tender-workspace/internal/entity/dto/queries/organizations"
	mc "tender-workspace/internal/utils/myconstants"
	e "tender-workspace/internal/utils/myerrors"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Repo interface {
	Create(ctx context.Context, initData *ent.Organization) (*ent.Organization, error)
	Get(ctx context.Context, organizationId int) (*ent.Organization, error)
	GetAll(ctx context.Context, params *oqp.OrganizationList) ([]*ent.Organization, error)
	Update(ctx context.Context, updateData *ent.Organization) (*ent.Organization, error)
	IsUserResponsible(ctx context.Context, userId, organizationId int) (bool, error)
	MakeUserResponsible(ctx context.Context, userId, organizationId int) error
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
	sqlRowCreateOrganization = `INSERT INTO organization (
											name, 
											description,
											type, 
											created_at, 
											updated_at
											) VALUES($1, $2, $3, $4, $4) RETURNING id, name, description, type, created_at`
	sqlRowCheckResponsibility = `SELECT 1 FROM organization_responsible WHERE organization_id=$1 AND user_id=$2`
	sqlRowMakeResponsibility  = `INSERT INTO organization_responsible(organization_id, user_id) VALUES($1, $2)`
	sqlRowUpdateOrganization  = `UPDATE organization
	SET name = $1,
		description = $2,
		type = $3,
		updated_at = $4
	WHERE id = $5
	RETURNING id, name, description, type, created_at`
)

func (r *RepoLayer) GetAll(ctx context.Context, params *oqp.OrganizationList) ([]*ent.Organization, error) {
	query, args := getAllSqlQuery(params)
	rows, err := r.Client.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orgsDB []*organizationDB
	for rows.Next() {
		var o organizationDB
		err := rows.Scan(&o.ID, &o.Name, &o.Description, &o.Type, &o.CreatedAt)
		if err != nil {
			requestId := ctx.Value(mc.ContextKey(mc.RequestID)).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(mc.RequestID, requestId))
			continue
		}
		orgsDB = append(orgsDB, &o)
	}
	orgs := getArrayOrganizationFromDB(orgsDB)
	return orgs, nil
}

func getAllSqlQuery(params *oqp.OrganizationList) (string, []any) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder().Select("id, name, description, type, created_at").From("organization")
	if params.Type != "" {
		sb = sb.Where(sb.Equal("type", params.Type))
	}
	sb = sb.Offset(params.Offset).Limit(params.Limit).OrderBy("name").Asc()
	return sb.Build()
}

func (r *RepoLayer) Create(ctx context.Context, initData *ent.Organization) (*ent.Organization, error) {
	row := r.Client.QueryRow(ctx, sqlRowCreateOrganization, initData.Name, initData.Description, initData.Type, time.Now())
	var oDB organizationDB
	err := row.Scan(&oDB.ID, &oDB.Name, &oDB.Description, &oDB.Type, &oDB.CreatedAt)
	if err != nil {
		return nil, err
	}
	return getOrganizationFromDB(&oDB), nil
}

func (r *RepoLayer) Update(ctx context.Context, updateData *ent.Organization) (*ent.Organization, error) {
	row := r.Client.QueryRow(ctx, sqlRowUpdateOrganization,
		updateData.Name,
		updateData.Description,
		updateData.Type,
		time.Now(),
		updateData.ID)
	var oDB organizationDB
	err := row.Scan(&oDB.ID, &oDB.Name, &oDB.Description, &oDB.Type, &oDB.CreatedAt)
	if err != nil {
		return nil, err
	}
	return getOrganizationFromDB(&oDB), nil
}

func (r *RepoLayer) Get(ctx context.Context, organizationId int) (*ent.Organization, error) {
	row := r.Client.QueryRow(ctx, `SELECT id, name, description, type, created_at FROM organization WHERE id=$1`, organizationId)
	var oDB organizationDB
	err := row.Scan(&oDB.ID, &oDB.Name, &oDB.Description, &oDB.Type, &oDB.CreatedAt)
	if err != nil {
		return nil, err
	}
	return getOrganizationFromDB(&oDB), nil
}

func (r *RepoLayer) IsUserResponsible(ctx context.Context, userID, organizationID int) (bool, error) {
	row := r.Client.QueryRow(ctx, sqlRowCheckResponsibility, organizationID, userID)
	var isExist int
	err := row.Scan(&isExist)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *RepoLayer) MakeUserResponsible(ctx context.Context, userID, organizationID int) error {
	tag, err := r.Client.Exec(ctx, sqlRowMakeResponsibility, organizationID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return e.ErrNoRowsAffected
	}
	return nil
}
