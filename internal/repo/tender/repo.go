package tender

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	ent "tender-workspace/internal/entity"
	queries "tender-workspace/internal/entity/dto/queries/tenders"
	tqp "tender-workspace/internal/entity/dto/queries/tenders"
	f "tender-workspace/internal/utils/functions"
	mc "tender-workspace/internal/utils/myconstants"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// QUERY
type UserTendersProps struct {
	Limit  int
	Offset int
	UserID int
}

type UpdateTenderProps struct {
	TenderID int
	UserID   int
}

// необходимо сделать откат версии
// /tenders/{tenderId}/rollback/{version}:
type Repo interface {
	GetAll(ctx context.Context, params *tqp.ListTenders) ([]*ent.Tender, error)
	Create(ctx context.Context, initData *ent.Tender) (*ent.Tender, error)
	GetUserTenders(ctx context.Context, params *UserTendersProps) ([]*ent.Tender, error)
	ChangeStatus(ctx context.Context, tenderId int, status string) (*ent.Tender, error)
	Update(ctx context.Context, newTenderData *ent.UpdateTenderData, params *UpdateTenderProps) (*ent.Tender, error)
	GetTenderStatus(ctx context.Context, tenderId int) (string, error)
	GetTender(ctx context.Context, tenderId int) (*ent.Tender, error)
	GetOrganizationTenders(ctx context.Context, organizationId int) ([]*ent.Tender, error)
	UpdateTenderStatus(ctx context.Context, tenderId int, status string) (*ent.Tender, error)
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
	sqlRowCreateTender = `INSERT INTO tender (
	    name, 
        description, 
		type, 
		status, 
		version,
        organization_id, 
        creator_id, 
		created_at,
		updated_at
    ) 
    VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $8
    ) RETURNING *`
)

func (r *RepoLayer) GetAll(ctx context.Context, params *tqp.ListTenders) ([]*ent.Tender, error) {
	query, args := getGetAllSqlQuery(params)
	rows, err := r.Client.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []*ent.Tender
	for rows.Next() {
		var t ent.Tender
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.Type,
			&t.Status,
			&t.Version,
			&t.OrganizationID,
			&t.CreatorID,
			&t.CreatedAt,
		)
		if err != nil {
			requestId := ctx.Value(mc.RequestID).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(mc.RequestID, requestId))
			continue
		}
		tenders = append(tenders, &t)
	}
	return tenders, nil
}

func getGetAllSqlQuery(params *queries.ListTenders) (string, []any) {
	sb := sqlbuilder.NewSelectBuilder().Select("*").From("tender")
	if params.ServiceType != "" {
		sb = sb.Where(sb.Equal("type", params.ServiceType))
	}
	sb = sb.Offset(int(params.Offset)).Limit(int(params.Limit)).OrderBy("name").Asc()
	return sb.Build()
}

func (r *RepoLayer) Create(ctx context.Context, initData *ent.Tender) (*ent.Tender, error) {
	timeNow := f.FormatTime(time.Now())
	row := r.Client.QueryRow(ctx, sqlRowCreateTender,
		initData.Name,
		initData.Description,
		initData.Type,
		initData.Status,
		1,
		initData.OrganizationID,
		initData.CreatorID,
		timeNow,
	)
	var t ent.Tender
	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&t.Type,
		&t.Status,
		&t.Version,
		&t.OrganizationID,
		&t.CreatorID,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *RepoLayer) ChangeStatus(ctx context.Context, tenderId int, status string) (*ent.Tender, error) {
	row := r.Client.QueryRow(ctx, `UPDATE tender SET status=$1 WHERE id=$2 RETURN *`, status, tenderId)
	var t ent.Tender
	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&t.Type,
		&t.Status,
		&t.Version,
		&t.OrganizationID,
		&t.CreatorID,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *RepoLayer) Update(ctx context.Context, newTenderData *ent.UpdateTenderData, params *UpdateTenderProps) (*ent.Tender, error) {
	query, args := getUpdateSqlQuery(newTenderData)
	row := r.Client.QueryRow(ctx, query+"RETURN *", args...)
	var t ent.Tender
	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&t.Type,
		&t.Status,
		&t.Version,
		&t.OrganizationID,
		&t.CreatorID,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func getUpdateSqlQuery(newTenderData *ent.UpdateTenderData) (string, []any) {
	sb := sqlbuilder.NewUpdateBuilder().Update("tender")
	if newTenderData.Type != "" {
		sb = sb.Set("service_type", newTenderData.Type)
	}
	if newTenderData.Name != "" {
		sb = sb.Set("name", newTenderData.Name)
	}
	if newTenderData.Description != "" {
		sb = sb.Set("description", newTenderData.Description)
	}
	sb = sb.Set("updated_at", f.FormatTime(time.Now()))
	return sb.Build()
}

func (r *RepoLayer) GetUserTenders(ctx context.Context, params *UserTendersProps) ([]*ent.Tender, error) {
	query, args := getGetUserTendersSqlQuery(params)
	rows, err := r.Client.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tenders := make([]*ent.Tender, 0)
	for rows.Next() {
		var t ent.Tender
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.Type,
			&t.Status,
			&t.Version,
			&t.OrganizationID,
			&t.CreatedAt,
		)
		if err != nil {
			requestId := ctx.Value(mc.RequestID).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(mc.RequestID, requestId))
			continue
		}
		tenders = append(tenders, &t)
	}
	return tenders, nil
}

func getGetUserTendersSqlQuery(params *UserTendersProps) (string, []any) {
	sb := sqlbuilder.NewSelectBuilder().Select("*").From("tender")
	sb = sb.Where(sb.Equal("creator_id", params.UserID))
	sb = sb.Offset(int(params.Offset)).Limit(int(params.Limit)).OrderBy("name").Asc()
	return sb.Build()
}

func (r *RepoLayer) GetTenderStatus(ctx context.Context, tenderId int) (string, error) {
	row := r.Client.QueryRow(ctx, `SELECT status FROM tender WHERE id=$1`, tenderId)
	var status string
	err := row.Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (r *RepoLayer) UpdateTenderStatus(ctx context.Context, tenderId int, status string) (*ent.Tender, error) {
	row := r.Client.QueryRow(ctx, `UPDATE tender SET status=$1 WHERE id=$2 SELECT *`, status, tenderId)
	var t ent.Tender
	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&t.Type,
		&t.Status,
		&t.Version,
		&t.OrganizationID,
		&t.CreatorID,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *RepoLayer) GetOrganizationTenders(ctx context.Context, organizationId int) ([]*ent.Tender, error) {
	rows, err := r.Client.Query(ctx, `SELECT * FROM tender WHERE organization_id=$1`, organizationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var tenders []*ent.Tender
	for rows.Next() {
		var t ent.Tender
		err := rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.Type,
			&t.Status,
			&t.Version,
			&t.OrganizationID,
			&t.CreatorID,
			&t.CreatedAt,
		)
		if err != nil {
			requestId := ctx.Value(mc.RequestID).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(mc.RequestID, requestId))
			continue
		}
		tenders = append(tenders, &t)
	}
	return tenders, nil
}

func (r *RepoLayer) GetTender(ctx context.Context, tenderId int) (*ent.Tender, error) {
	row := r.Client.QueryRow(ctx, `SELECT * FROM tender WHERE id=$1`, tenderId)
	var t ent.Tender
	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.Description,
		&t.Type,
		&t.Status,
		&t.Version,
		&t.OrganizationID,
		&t.CreatorID,
		&t.CreatedAt,
	)
	if err != nil {
		return nil, nil
	}
	return &t, nil
}
