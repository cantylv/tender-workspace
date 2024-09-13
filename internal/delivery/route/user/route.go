package user

import (
	delUser "tender-workspace/internal/delivery/user"
	repoUser "tender-workspace/internal/repo/user"
	usecaseUser "tender-workspace/internal/usecase/user"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func InitHandlers(r *mux.Router, psqlConn *pgx.Conn, logger *zap.Logger) {
	// init repo, usecase, handler
	userRepo := repoUser.NewRepoLayer(psqlConn, logger)
	userUsecase := usecaseUser.NewUsecaseLayer(userRepo)
	userDelivery := delUser.NewDeliveryLayer(userUsecase, logger)
	r.HandleFunc("/users/new", userDelivery.CreateUser)
	r.HandleFunc("/users/{username}", userDelivery.GetUser)
	r.HandleFunc("/users/{username}/organizations", userDelivery.GetUserOrganizationsIds)
}
