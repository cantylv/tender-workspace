package middlewares

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func Init(r *mux.Router, logger *zap.Logger) (h http.Handler) {
	h = Access(r, logger)
	h = Recover(h, logger)
	return h
}
