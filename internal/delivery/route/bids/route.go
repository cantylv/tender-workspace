package bids

import (
	delBids "tender-workspace/internal/delivery/bids"
	repoBids "tender-workspace/internal/repo/bids"
	repoOrgs "tender-workspace/internal/repo/organization"
	repoTender "tender-workspace/internal/repo/tender"
	repoUser "tender-workspace/internal/repo/user"

	// repoUser "tender-workspace/internal/repo/user"
	usecaseBids "tender-workspace/internal/usecase/bids"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func InitHandlers(r *mux.Router, psqlConn *pgx.Conn, logger *zap.Logger) {
	// init repo, usecase, handler
	bRepo := repoBids.NewRepoLayer(psqlConn, logger)
	uRepo := repoUser.NewRepoLayer(psqlConn, logger)
	tRepo := repoTender.NewRepoLayer(psqlConn, logger)
	oRepo := repoOrgs.NewRepoLayer(psqlConn, logger)
	bUsecase := usecaseBids.NewUsecaseLayer(bRepo, uRepo, oRepo, tRepo)
	bDelivery := delBids.NewDeliveryLayer(bUsecase, logger)

	r.HandleFunc("/bids/new", bDelivery.CreateBid)
	r.HandleFunc("/bids/my", bDelivery.GetUserBids)
	r.HandleFunc("/bids/{tenderId}/list", bDelivery.GetTenderListOfBids)
	r.HandleFunc("/bids/{bidId}/status", bDelivery.GetBidStatus).Methods("GET")
	r.HandleFunc("/bids/{bidId}/status", bDelivery.UpdateBidStatus).Methods("PUT")
	r.HandleFunc("/bids/{bidId}/edit", bDelivery.UpdateBid)
	r.HandleFunc("/bids/{bidId}/submit_decision", bDelivery.SubmitDecision)
	// r.HandleFunc("/bids/{bidId}/rollback/{version}")
	// r.HandleFunc("/bids/{tenderId}/reviews")
	// надо еще добавить rollback
}
