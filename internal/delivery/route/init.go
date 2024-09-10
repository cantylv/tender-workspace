package route

import (
	"net/http"
	"tender-workspace/internal/delivery/route/tender"
	"tender-workspace/internal/middlewares"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func InitHTTPHandlers(router *mux.Router, psqlConn *pgx.Conn, logger *zap.Logger) http.Handler {
	router = router.PathPrefix("/api").Subrouter()
	tender.InitHandlers(router, psqlConn)

	return middlewares.Init(router, logger)
}
