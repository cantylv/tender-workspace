package organization

import (
	delOrg "tender-workspace/internal/delivery/organization"
	repoOrg "tender-workspace/internal/repo/organization"
	repoUser "tender-workspace/internal/repo/user"
	usecaseOrg "tender-workspace/internal/usecase/organization"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func InitHandlers(r *mux.Router, psqlConn *pgx.Conn, logger *zap.Logger) {
	// init repo, usecase, handler
	orgRepo := repoOrg.NewRepoLayer(psqlConn, logger)
	userRepo := repoUser.NewRepoLayer(psqlConn, logger)
	orgUsecase := usecaseOrg.NewUsecaseLayer(orgRepo, userRepo)
	orgDelivery := delOrg.NewDeliveryLayer(orgUsecase, logger)
	r.HandleFunc("/organizations", orgDelivery.GetListOfOrganizations)
	r.HandleFunc("/organizations/new", orgDelivery.CreateNewOrganization)
	r.HandleFunc("/organizations/{organizationID}/edit", orgDelivery.UpdateOrganization)
	r.HandleFunc("/organizations/{organizationID}/users/{username}/make_responsible", orgDelivery.MakeResponsible)
}
