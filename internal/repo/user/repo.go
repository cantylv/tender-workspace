package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type Repo interface {
	GetData(ctx context.Context, userID int) ([]entity.User, error)
	Create(ctx context.Context, initData *ent.Tender) error
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
