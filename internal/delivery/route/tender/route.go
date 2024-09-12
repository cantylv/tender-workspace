package tender

import (
	"tender-workspace/internal/delivery/tender"
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
	tUsecase := usecaseTender.NewUsecaseLayer(tRepo, uRepo)
	tDelivery := tender.NewDeliveryLayer(tUsecase, logger)
	r.HandleFunc("/tenders", tDelivery.GetListOfTenders)
	r.HandleFunc("/tenders/new", tDelivery.CreateNewTender)
	r.HandleFunc("/tenders/my", tDelivery.GetUserTenders)
	r.HandleFunc("/tender/{tenderId}/status", tDelivery.GetTenderStatus)
	r.HandleFunc("/tender/{tenderId}/status", tDelivery.UpdateTenderStatus)
	r.HandleFunc("/tender/{tenderId}/edit", tDelivery.UpdateTender)
	// надо еще добавить rollback
}
