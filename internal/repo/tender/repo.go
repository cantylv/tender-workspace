package tender

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	ent "tender-workspace/internal/entity"
	tqp "tender-workspace/internal/entity/dto/queries/tenders"
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
	Get(ctx context.Context, tenderId int) (*ent.Tender, error)
	Create(ctx context.Context, initData *ent.Tender) (*ent.Tender, error)
	GetUserTenders(ctx context.Context, params *UserTendersProps) ([]*ent.Tender, error)
	ChangeStatus(ctx context.Context, tenderId int, status string) (*ent.Tender, error)
	Update(ctx context.Context, newTenderData *ent.UpdateTenderData, params *UpdateTenderProps) (*ent.Tender, error)
	GetTenderStatus(ctx context.Context, tenderId int) (string, error)
	GetTender(ctx context.Context, tenderId int) (*ent.Tender, error)
	GetOrganizationTenders(ctx context.Context, organizationId int) ([]*ent.Tender, error)
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
    ) RETURNING id, name, description, type, status, version, organization_id, creator_id, created_at`
	sqlRowUpdateTenderStatus = `UPDATE tender SET status=$1 WHERE id=$2 RETURNING id, name, description, type, status, version, organization_id, creator_id, created_at`
)

func (r *RepoLayer) Get(ctx context.Context, tenderId int) (*ent.Tender, error) {
	row := r.Client.QueryRow(ctx, `SELECT id, name, description, type, status, version, organization_id, creator_id, created_at FROM tender WHERE id=$1`, tenderId)
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
			requestId := ctx.Value(mc.ContextKey(mc.RequestID)).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(mc.RequestID, requestId))
			continue
		}
		tenders = append(tenders, &t)
	}
	return tenders, nil
}

func getGetAllSqlQuery(params *tqp.ListTenders) (string, []any) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder().
		Select("id, name, description, type, status, version, organization_id, creator_id, created_at").
		From("tender")
	if params.ServiceType != "" {
		sb = sb.Where(sb.Equal("type", params.ServiceType))
	}
	sb = sb.Where(sb.NotEqual("status", "Created"))
	sb = sb.Where(sb.NotEqual("status", "Closed"))
	sb = sb.Offset(params.Offset).Limit(params.Limit).OrderBy("name").Asc()
	return sb.Build()
}

func (r *RepoLayer) Create(ctx context.Context, initData *ent.Tender) (*ent.Tender, error) {
	timeNow := time.Now()
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
	row := r.Client.QueryRow(ctx, sqlRowUpdateTenderStatus, status, tenderId)
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
	query, args := getUpdateSqlQuery(newTenderData, params)
	fmt.Println()
	fmt.Println()
	fmt.Println(query)
	tx, err := r.Client.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	// Откат транзакции в случае ошибки
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	tag, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if tag.RowsAffected() == 0 {
		return nil, err
	}
	row := tx.QueryRow(ctx, `SELECT id, name, description, type, status, version, organization_id, creator_id, created_at FROM tender WHERE id=$1`, params.TenderID)
	var t ent.Tender
	err = row.Scan(
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
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return &t, nil
}

func getUpdateSqlQuery(newTenderData *ent.UpdateTenderData, params *UpdateTenderProps) (string, []any) {
	sb := sqlbuilder.PostgreSQL.NewUpdateBuilder().Update("tender")
	// Собираем все изменения
	var updates []string

	if newTenderData.Type != "" {
		updates = append(updates, sb.Assign("type", newTenderData.Type))
	}
	if newTenderData.Name != "" {
		updates = append(updates, sb.Assign("name", newTenderData.Name))
	}
	if newTenderData.Description != "" {
		updates = append(updates, sb.Assign("description", newTenderData.Description))
	}
	updates = append(updates, sb.Assign("updated_at", time.Now()))
	sb.Set(updates...)
	sb.Where(sb.Equal("id", params.TenderID))
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
			requestId := ctx.Value(mc.ContextKey(mc.RequestID)).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(mc.RequestID, requestId))
			continue
		}
		tenders = append(tenders, &t)
	}
	return tenders, nil
}

func getGetUserTendersSqlQuery(params *UserTendersProps) (string, []any) {
	sb := sqlbuilder.NewSelectBuilder().Select("id, name, description, type, status, version, organization_id, creator_id, created_at").From("tender")
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

func (r *RepoLayer) GetOrganizationTenders(ctx context.Context, organizationId int) ([]*ent.Tender, error) {
	rows, err := r.Client.Query(ctx, `SELECT id, name, description, type, status, version, organization_id, creator_id, created_at FROM tender WHERE organization_id=$1`, organizationId)
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
			requestId := ctx.Value(mc.ContextKey(mc.RequestID)).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(mc.RequestID, requestId))
			continue
		}
		tenders = append(tenders, &t)
	}
	return tenders, nil
}

func (r *RepoLayer) GetTender(ctx context.Context, tenderId int) (*ent.Tender, error) {
	row := r.Client.QueryRow(ctx, `SELECT id, name, description, type, status, version, organization_id, creator_id, created_at FROM tender WHERE id=$1`, tenderId)
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
