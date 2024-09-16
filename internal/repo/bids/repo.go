package bids

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	ent "tender-workspace/internal/entity"
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
		tender_id, 
		creator_id,
		author_type,
		organization_id,
		created_at,
		updated_at)
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $9
	) RETURNING id, name, description, status, version, tender_id, creator_id, author_type, organization_id, created_at`
)

type Repo interface {
	Create(ctx context.Context, initData *ent.Bid) (*ent.Bid, error)
	GetBid(ctx context.Context, bidID int) (*ent.Bid, error)
	GetStatus(ctx context.Context, bidID int) (string, error)
	UpdateStatus(ctx context.Context, bidID int, status string, newBidVersion int) (*ent.Bid, error)
	Update(ctx context.Context, newData *UpdateBid, newBidVersion int) (*ent.Bid, error)
	GetTenderBids(ctx context.Context, tenderID int) ([]*ent.Bid, error)
	GetOrganizationBids(ctx context.Context, organizationID int) ([]*ent.Bid, error)
	GetUserBids(ctx context.Context, creatorID int) ([]*ent.Bid, error)
	UserHasBid(ctx context.Context, userID, tenderID int) (bool, error)
	OrganizationHasBid(ctx context.Context, orgID, tenderID int) (bool, error)
}

var _ Repo = (*RepoLayer)(nil)

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

func (r *RepoLayer) Create(ctx context.Context, initData *ent.Bid) (*ent.Bid, error) {
	timeNow := time.Now()
	initDataDB := newBidDB(initData)
	row := r.Client.QueryRow(ctx, sqlRowCreateBid,
		initDataDB.Name,
		initDataDB.Description,
		initDataDB.Status,
		initDataDB.Version,
		initDataDB.TenderID,
		initDataDB.CreatorID,
		initDataDB.AuthorType,
		initDataDB.OrganizationID,
		timeNow,
	)
	var bidDB bidDB
	err := row.Scan(
		&bidDB.ID,
		&bidDB.Name,
		&bidDB.Description,
		&bidDB.Status,
		&bidDB.Version,
		&bidDB.TenderID,
		&bidDB.CreatorID,
		&bidDB.AuthorType,
		&bidDB.OrganizationID,
		&bidDB.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return newBid(&bidDB), nil
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

func (r *RepoLayer) UpdateStatus(ctx context.Context, bidId int, status string, newBidVersion int) (*ent.Bid, error) {
	row := r.Client.QueryRow(ctx, `UPDATE bids SET status=$1, version=$2 WHERE id=$3 RETURNING id, name, description, status, version, tender_id, creator_id, author_type, organization_id, created_at`, status, newBidVersion, bidId)
	var bidDB bidDB
	err := row.Scan(
		&bidDB.ID,
		&bidDB.Name,
		&bidDB.Description,
		&bidDB.Status,
		&bidDB.Version,
		&bidDB.TenderID,
		&bidDB.CreatorID,
		&bidDB.AuthorType,
		&bidDB.OrganizationID,
		&bidDB.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return newBid(&bidDB), nil
}

func (r *RepoLayer) Update(ctx context.Context, newData *UpdateBid, newBidVersion int) (*ent.Bid, error) {
	query, args := updateSqlQuery(newData, newBidVersion)
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
	row := tx.QueryRow(ctx, `SELECT id, name, description, status, version, tender_id, creator_id, author_type, organization_id, created_at FROM bids WHERE id=$1`, newData.BidID)
	var bidDB bidDB
	err = row.Scan(
		&bidDB.ID,
		&bidDB.Name,
		&bidDB.Description,
		&bidDB.Status,
		&bidDB.Version,
		&bidDB.TenderID,
		&bidDB.CreatorID,
		&bidDB.AuthorType,
		&bidDB.OrganizationID,
		&bidDB.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return newBid(&bidDB), nil
}

func updateSqlQuery(newData *UpdateBid, newBidVersion int) (string, []any) {
	sb := sqlbuilder.PostgreSQL.NewUpdateBuilder().Update("bids")
	// Собираем все изменения
	var updates []string

	if newData.Name != "" {
		updates = append(updates, sb.Assign("name", newData.Name))
	}
	if newData.Description != "" {
		updates = append(updates, sb.Assign("description", newData.Description))
	}
	updates = append(updates, sb.Assign("updated_at", time.Now()), sb.Assign("version", newBidVersion))
	sb.Set(updates...)
	sb.Where(sb.Equal("id", newData.BidID))
	return sb.Build()
}

func (r *RepoLayer) GetTenderBids(ctx context.Context, tenderID int) ([]*ent.Bid, error) {
	rows, err := r.Client.Query(ctx, `SELECT id, name, description, status, version, tender_id, creator_id, author_type, organization_id, created_at FROM bids WHERE tender_id=$1 AND status='Published'`, tenderID)
	if err != nil {
		return nil, err
	}
	var bidsDB []*bidDB
	for rows.Next() {
		var bidDB bidDB
		err := rows.Scan(
			&bidDB.ID,
			&bidDB.Name,
			&bidDB.Description,
			&bidDB.Status,
			&bidDB.Version,
			&bidDB.TenderID,
			&bidDB.CreatorID,
			&bidDB.AuthorType,
			&bidDB.OrganizationID,
			&bidDB.CreatedAt,
		)
		if err != nil {
			requestId := ctx.Value(mc.ContextKey(mc.RequestID)).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(mc.RequestID, requestId))
			continue
		}
		bidsDB = append(bidsDB, &bidDB)
	}
	return getArrayBidFromDB(bidsDB), nil
}

func (r *RepoLayer) GetOrganizationBids(ctx context.Context, organizationId int) ([]*ent.Bid, error) {
	rows, err := r.Client.Query(ctx, `SELECT id, name, description, status, version, tender_id, creator_id, author_type, organization_id, created_at FROM bids WHERE organization_id=$1`, organizationId)
	if err != nil {
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
			&bid.AuthorType,
			&bid.OrganizationID,
			&bid.CreatedAt,
		)
		if err != nil {
			requestId := ctx.Value(mc.ContextKey(mc.RequestID)).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(mc.RequestID, requestId))
			continue
		}
		bids = append(bids, &bid)
	}
	return bids, nil
}

func (r *RepoLayer) GetUserBids(ctx context.Context, creatorID int) ([]*ent.Bid, error) {
	rows, err := r.Client.Query(ctx, `SELECT id, name, description, status, version, tender_id, creator_id, author_type, organization_id, created_at FROM bids WHERE creator_id=$1 AND organization_id IS NULL`, creatorID)
	if err != nil {
		return nil, err
	}

	var bids []*bidDB
	for rows.Next() {
		var bid bidDB
		err := rows.Scan(
			&bid.ID,
			&bid.Name,
			&bid.Description,
			&bid.Status,
			&bid.Version,
			&bid.TenderID,
			&bid.CreatorID,
			&bid.AuthorType,
			&bid.OrganizationID,
			&bid.CreatedAt,
		)
		if err != nil {
			requestId := ctx.Value(mc.ContextKey(mc.RequestID)).(string)
			r.Logger.Error(fmt.Sprintf("error while scanning sql result: %v", err), zap.String(mc.RequestID, requestId))
			continue
		}
		bids = append(bids, &bid)
	}
	return getArrayBidFromDB(bids), nil
}

func (r *RepoLayer) GetBid(ctx context.Context, bidId int) (*ent.Bid, error) {
	row := r.Client.QueryRow(ctx, `SELECT id, name, description, status, version, tender_id, creator_id, author_type, organization_id, created_at FROM bids WHERE id=$1`, bidId)
	var bid bidDB
	err := row.Scan(
		&bid.ID,
		&bid.Name,
		&bid.Description,
		&bid.Status,
		&bid.Version,
		&bid.TenderID,
		&bid.CreatorID,
		&bid.AuthorType,
		&bid.OrganizationID,
		&bid.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return newBid(&bid), nil
}

func (r *RepoLayer) UserHasBid(ctx context.Context, userID, tenderID int) (bool, error) {
	row := r.Client.QueryRow(ctx, `SELECT 1 FROM bids WHERE creator_id = $1 AND tender_id = $2`, userID, tenderID)
	var hasBid int
	err := row.Scan(&hasBid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *RepoLayer) OrganizationHasBid(ctx context.Context, orgID, tenderID int) (bool, error) {
	row := r.Client.QueryRow(ctx, `SELECT 1 FROM bids WHERE organization_id = $1 AND tender_id = $2`, orgID, tenderID)
	var hasBid int
	err := row.Scan(&hasBid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
