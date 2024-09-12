package bids

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	ent "tender-workspace/internal/entity"
	f "tender-workspace/internal/utils/functions"
	mc "tender-workspace/internal/utils/myconstants"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type UpdateBid struct {
	BidID       int
	Name        string
	Description string
}

var (
	sqlRowCreateBid = `INSERT INTO bids (
		name, 
		description, 
		status, 
		version, 
		tendor_id, 
		creator_id,
		organization_id,
		created_at,
		updated_at)
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $7
	) RETURN *`
)

type Repo interface {
	Create(ctx context.Context, initData *ent.Bid) (*ent.Bid, error)
	GetBid(ctx context.Context, bidId int) (*ent.Bid, error)
	GetStatus(ctx context.Context, bidId int) (string, error)
	UpdateStatus(ctx context.Context, bidId int, status string) (*ent.Bid, error)
	Update(ctx context.Context, newData *UpdateBid) (*ent.Bid, error)
	GetOrganizationBids(ctx context.Context, organizationId int) ([]*ent.Bid, error)
}

type RepoLayer struct {
	Client *pgx.Conn
	Logger *zap.Logger
}

var _ Repo = (*RepoLayer)(nil)

func NewRepoLayer(client *pgx.Conn) *RepoLayer {
	return &RepoLayer{
		Client: client,
	}
}

func (r *RepoLayer) Create(ctx context.Context, initData *ent.Bid) (*ent.Bid, error) {
	timeNow := f.FormatTime(time.Now())
	row := r.Client.QueryRow(ctx, sqlRowCreateBid,
		initData.Name,
		initData.Description,
		initData.Status,
		initData.Version,
		initData.TenderID,
		initData.CreatorID,
		initData.OrganizationID,
		timeNow,
	)
	var bid ent.Bid
	err := row.Scan(
		&bid.ID,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.Version,
		&bid.TenderID,
		&bid.CreatorID,
		&bid.OrganizationID,
		&bid.Status,
	)
	if err != nil {
		return nil, err
	}
	return &bid, nil
}

func (r *RepoLayer) GetStatus(ctx context.Context, bidId int) (string, error) {
	row := r.Client.QueryRow(ctx, `SELECT status from bids WHERE id=$1`, bidId)
	var status string
	err := row.Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (r *RepoLayer) UpdateStatus(ctx context.Context, bidId int, status string) (*ent.Bid, error) {
	row := r.Client.QueryRow(ctx, `UPDATE bids SET status=$1 WHERE id=$2 RETURN *`, status, bidId)
	var b ent.Bid
	err := row.Scan(
		&b.ID,
		&b.Name,
		&b.Status,
		&b.Version,
		&b.TenderID,
		&b.CreatorID,
		&b.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *RepoLayer) Update(ctx context.Context, newData *UpdateBid) (*ent.Bid, error) {
	query, args := getUpdateSqlQuery(newData)
	row := r.Client.QueryRow(ctx, query+"RETURN *", args...)
	var b ent.Bid
	err := row.Scan(
		&b.ID,
		&b.Name,
		&b.Status,
		&b.Version,
		&b.TenderID,
		&b.CreatorID,
		&b.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &b, err
}

func getUpdateSqlQuery(newTenderData *UpdateBid) (string, []any) {
	sb := sqlbuilder.NewUpdateBuilder().Update("bids").Where("bidId", fmt.Sprintf("%d", newTenderData.BidID))
	if newTenderData.Name != "" {
		sb = sb.Set("name", newTenderData.Name)
	}
	if newTenderData.Description != "" {
		sb = sb.Set("description", newTenderData.Description)
	}
	sb = sb.Set("updated_at", f.FormatTime(time.Now()))
	return sb.Build()
}

func (r *RepoLayer) GetOrganizationBids(ctx context.Context, organizationId int) ([]*ent.Bid, error) {
	rows, err := r.Client.Query(ctx, `SELECT * FROM bids WHERE organization_id=$1`, organizationId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var bids []*ent.Bid
	for rows.Next() {
		var bid ent.Bid
		err := rows.Scan(
			&bid.ID,
			&bid.Name,
			&bid.Description,
			&bid.Status,
			&bid.Version,
			&bid.TenderID,
			&bid.CreatorID,
			&bid.OrganizationID,
			&bid.Status,
		)
		if err != nil {
			requestId := ctx.Value(mc.RequestID).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(mc.RequestID, requestId))
			continue
		}
		bids = append(bids, &bid)
	}
	return bids, nil
}

func (r *RepoLayer) GetBid(ctx context.Context, bidId int) (*ent.Bid, error) {
	row := r.Client.QueryRow(ctx, `SELECT * FROM bids WHERE id=$1`, bidId)
	var bid ent.Bid
	err := row.Scan(
		&bid.ID,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.Version,
		&bid.TenderID,
		&bid.CreatorID,
		&bid.OrganizationID,
		&bid.Status,
	)
	if err != nil {
		return nil, err
	}
	return &bid, nil
}
