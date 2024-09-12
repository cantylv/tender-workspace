package route

import (
	"net/http"
	"tender-workspace/internal/delivery/route/ping"
	"tender-workspace/internal/delivery/route/tender"
	"tender-workspace/internal/middlewares"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func InitHTTPHandlers(router *mux.Router, psqlConn *pgx.Conn, logger *zap.Logger) http.Handler {
	router = router.PathPrefix("/api").Subrouter()
	ping.InitHandlers(router)
	tender.InitHandlers(router, psqlConn, logger)

	return middlewares.Init(router, logger)
}
