package tender

import (
	"tender-workspace/internal/delivery/tender"
	repoOrg "tender-workspace/internal/repo/organization"
	repoTender "tender-workspace/internal/repo/tender"
	repoUser "tender-workspace/internal/repo/user"
	usecaseTender "tender-workspace/internal/usecase/tender"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func InitHandlers(r *mux.Router, psqlConn *pgx.Conn, logger *zap.Logger) {
	// init repo, usecase, handler
	tRepo := repoTender.NewRepoLayer(psqlConn, logger)
	uRepo := repoUser.NewRepoLayer(psqlConn, logger)
	oRepo := repoOrg.NewRepoLayer(psqlConn, logger)
	tUsecase := usecaseTender.NewUsecaseLayer(tRepo, uRepo, oRepo)
	tDelivery := tender.NewDeliveryLayer(tUsecase, logger)
	r.HandleFunc("/tenders", tDelivery.GetListOfTenders)
	r.HandleFunc("/tenders/new", tDelivery.CreateNewTender)
	r.HandleFunc("/tenders/my", tDelivery.GetUserTenders)
	r.HandleFunc("/tenders/{tenderId}/status", tDelivery.GetTenderStatus).Methods("GET")
	r.HandleFunc("/tenders/{tenderId}/status", tDelivery.UpdateTenderStatus).Methods("PUT")
	r.HandleFunc("/tenders/{tenderId}/edit", tDelivery.UpdateTender)
	// надо еще добавить rollback
}
