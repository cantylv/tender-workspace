package bids

// import (
// 	repoBids "tender-workspace/internal/repo/bids"

// 	// repoUser "tender-workspace/internal/repo/user"

// 	"github.com/gorilla/mux"
// 	"github.com/jackc/pgx/v5"
// 	"go.uber.org/zap"
// )

// func InitHandlers(r *mux.Router, psqlConn *pgx.Conn, logger *zap.Logger) {
// 	// init repo, usecase, handler
// 	bRepo := repoBids.NewRepoLayer(psqlConn, logger)
// 	bRepo := repoFeedback.NewRepoLayer(psqlConn, logger)
// 	r.HandleFunc("/bids/{bidId}/feedback", bDelivery.SubmitDecision)
// 	// r.HandleFunc("/bids/{bidId}/rollback/{version}")
// 	// r.HandleFunc("/bids/{tenderId}/reviews")
// 	// надо еще добавить rollback
// }
