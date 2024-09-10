package ping

import (
	"tender-workspace/internal/delivery/healthcheck"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func InitHTTPHandlers(router *mux.Router, psqlConn *pgx.Conn, logger *zap.Logger) {
	router.HandleFunc("/ping", healthcheck.Ping)
}
