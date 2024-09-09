package tender

import (
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func InitHandlers(r *mux.Router, psqlConn *pgx.Conn) {
	// init repo, usecase, handler
}
