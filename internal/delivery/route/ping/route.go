package ping

import (
	"tender-workspace/internal/delivery/healthcheck"

	"github.com/gorilla/mux"
)

func InitHandlers(router *mux.Router) {
	router.HandleFunc("/ping", healthcheck.Ping)
}
