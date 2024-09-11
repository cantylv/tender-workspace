package bids

import (
	"context"
	"fmt"
	ent "tender-workspace/internal/entity"
	f "tender-workspace/internal/utils/functions"
	mc "tender-workspace/internal/utils/myconstants"
	e "tender-workspace/internal/utils/myerrors"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type TenderBidsRepoQP struct {
	TenderID string
	UserID   int
	Limit    int32
	Offset   int32
}

var (
	sqlRowCreateBid = `INSERT INTO bids (
		name, 
		description, 
		status, 
		version, 
		tendor_id, 
		creator_id,
		created_at,
		updated_at)
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $7
	)`
)

type Repo interface {
	Create(ctx context.Context, initData *ent.Bid) error
	GetUserBids(ctx context.Context, params *ent.UserRepoQP) ([]ent.Bid, error)
	GetTenderBids(ctx context.Context, params *ent.TenderBidsRepoQP) ([]ent.Bid, error)
	GetStatus(ctx context.Context, bidId int) (string, error)
	UpdateStatus(ctx context.Context, bidId int, status string) error
	Update(ctx context.Context, newData *ent.UpdateBidData, bidId int) error
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

func (r *RepoLayer) Create(ctx context.Context, initData *ent.Bid) error {
	timeNow := f.FormatTime(time.Now())
	tag, err := r.Client.Exec(ctx, sqlRowCreateBid,
		initData.Name,
		initData.Description,
		initData.Status,
		initData.Version,
		initData.TenderID,
		initData.CreatorID,
		timeNow,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return e.ErrNoRowsAffected
	}
	return nil
}

func (r *RepoLayer) GetUserBids(ctx context.Context, params *ent.UserRepoQP) ([]ent.Bid, error) {
	query, args := getGetUserBidsSqlQuery(params)
	rows, err := r.Client.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []ent.Bid
	for rows.Next() {
		var b ent.Bid
		err = rows.Scan(
			&b.ID,
			&b.Name,
			&b.Description,
			&b.Status,
			&b.Version,
			&b.TenderID,
			&b.CreatorID,
			&b.CreatedAt,
		)
		if err != nil {
			requestId := ctx.Value(mc.RequestID).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(mc.RequestID, requestId))
			continue
		}
		bids = append(bids, b)
	}
	return bids, nil
}

func getGetUserBidsSqlQuery(params *ent.UserRepoQP) (string, []any) {
	sb := sqlbuilder.NewSelectBuilder().Select("*").From("bids")
	sb = sb.Where(sb.Equal("creator_id", params.UserID))
	sb = sb.Offset(int(params.Offset)).Limit(int(params.Limit)).OrderBy("name").Asc()
	return sb.Build()
}

func (r *RepoLayer) GetTenderBids(ctx context.Context, params *ent.TenderBidsRepoQP) ([]ent.Bid, error) {
	query, args := getGetAllSqlQuery(params)
	rows, err := r.Client.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bids []ent.Bid
	for rows.Next() {
		var b ent.Bid
		err = rows.Scan(
			&b.ID,
			&b.Name,
			&b.Status,
			&b.Version,
			&b.TenderID,
			&b.CreatorID,
			&b.CreatedAt,
		)
		if err != nil {
			requestId := ctx.Value(mc.RequestID).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(mc.RequestID, requestId))
			continue
		}
		bids = append(bids, b)
	}
	return bids, nil
}

func getGetAllSqlQuery(params *ent.TenderBidsRepoQP) (string, []any) {
	sb := sqlbuilder.NewSelectBuilder().Select("*").From("bids")
	sb = sb.Offset(int(params.Offset)).Limit(int(params.Limit)).OrderBy("name").Asc()
	return sb.Build()
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

func (r *RepoLayer) UpdateStatus(ctx context.Context, bidId int, status string) error {
	tag, err := r.Client.Exec(ctx, `UPDATE bids SET status=$1 WHERE id=$2`, status, bidId)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return e.ErrNoRowsAffected
	}
	return nil
}

func (r *RepoLayer) Update(ctx context.Context, newData *ent.UpdateBidData, bidId int) error {
	query, args := getUpdateSqlQuery(newData, bidId)
	tag, err := r.Client.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return e.ErrNoRowsAffected
	}
	return nil
}

func getUpdateSqlQuery(newTenderData *ent.UpdateBidData, bidId int) (string, []any) {
	sb := sqlbuilder.NewUpdateBuilder().Update("bids").Where("bidId", fmt.Sprintf("%d", bidId))
	if newTenderData.Name != "" {
		sb = sb.Set("name", newTenderData.Name)
	}
	if newTenderData.Description != "" {
		sb = sb.Set("description", newTenderData.Description)
	}
	sb = sb.Set("updated_at", f.FormatTime(time.Now()))
	return sb.Build()
}
