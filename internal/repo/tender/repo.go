package tender

import (
	"context"
	"fmt"
	ent "tender-workspace/internal/entity"
	"tender-workspace/internal/utils/myconstants"
	"tender-workspace/internal/utils/myerrors"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Repo interface {
	GetAll(ctx context.Context, params *ent.TenderQueryParameters) ([]ent.Tender, error)
	Create(ctx context.Context, initData *ent.Tender) error
	CanInteract(ctx context.Context, resp *ent.OrganizationResponsible) (bool, error)
	ChangeStatus(ctx context.Context, status string, tenderId string) error
	Update(ctx context.Context, newTenderData *ent.UpdateTenderData) error
}

type RepoLayer struct {
	Client *pgx.Conn
	Logger *zap.Logger
}

func NewRepoLayer(client *pgx.Conn, logger *zap.Logger) *RepoLayer {
	return &RepoLayer{
		Client: client,
		Logger: logger,
	}
}

var (
	sqlRowCreateTender = `INSERT INTO tender (
        organization_id, 
        creator_username, 
        service_type, 
        name, 
        description, 
        status, 
        version
    ) 
    VALUES (
        $1, $2, $3, $4, $5, $6, $7
    )`
	sqlRowCheckResponsibility = `SELECT 1 FROM organization_responsible WHERE organization_id=$1 AND user_id=$2`
	sqlRowUpdateStatus        = `UPDATE tender SET status=$1 WHERE id=$2`
)

func (r *RepoLayer) GetAll(ctx context.Context, params *ent.TenderQueryParameters) ([]ent.Tender, error) {
	query, args := getGetAllSqlQuery(params)
	rows, err := r.Client.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenders []ent.Tender
	for rows.Next() {
		var t ent.Tender
		err = rows.Scan(
			&t.OrganizationID,
			&t.CreatorUsername,
			&t.ServiceType,
			&t.Name,
			&t.Description,
			&t.Status,
			&t.Version,
		)
		if err != nil {
			requestId := ctx.Value(myconstants.RequestID).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(myconstants.RequestID, requestId))
			continue
		}
		tenders = append(tenders, t)
	}
	return tenders, nil
}

func getGetAllSqlQuery(params *ent.TenderQueryParameters) (string, []any) {
	sb := sqlbuilder.NewSelectBuilder().Select("*").From("tender")
	if params.ServiceType != "" {
		sb = sb.Where(sb.Equal("type", params.ServiceType))
	}
	sb = sb.Offset(int(params.Offset)).Limit(int(params.Limit)).OrderBy("name").Asc()
	return sb.Build()
}

func (r *RepoLayer) Create(ctx context.Context, initData *ent.Tender) error {
	tag, err := r.Client.Exec(ctx, sqlRowCreateTender,
		initData.OrganizationID,
		initData.CreatorUsername,
		initData.ServiceType,
		initData.Name,
		initData.Description,
		initData.Status,
		1,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return myerrors.ErrNoRowsAffected
	}
	return nil
}

func (r *RepoLayer) CanInteract(ctx context.Context, resp *ent.OrganizationResponsible) (bool, error) {
	row := r.Client.QueryRow(ctx, sqlRowCheckResponsibility, resp.OrganizationID, resp.UserID)
	var isExist int
	err := row.Scan(&isExist)
	if err != nil {
		return false, err
	}
	return isExist == 1, nil
}

func (r *RepoLayer) ChangeStatus(ctx context.Context, status string, tenderId string) error {
	tag, err := r.Client.Exec(ctx, sqlRowUpdateStatus, status, tenderId)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return myerrors.ErrNoRowsAffected
	}
	return nil
}

func (r *RepoLayer) Update(ctx context.Context, newTenderData *ent.UpdateTenderData) error {
	query, args := getUpdateSqlQuery(newTenderData)
	tag, err := r.Client.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return myerrors.ErrNoRowsAffected
	}
	return nil
}

func getUpdateSqlQuery(newTenderData *ent.UpdateTenderData) (string, []any) {
	sb := sqlbuilder.NewUpdateBuilder().Update("tender")
	if newTenderData.ServiceType != "" {
		sb = sb.Set("service_type", newTenderData.ServiceType)
	}
	if newTenderData.Name != "" {
		sb = sb.Set("name", newTenderData.Name)
	}
	if newTenderData.Description != "" {
		sb = sb.Set("description", newTenderData.Description)
	}
	return sb.Build()
}
