package middlewares

import (
	"fmt"
	"net/http"
	e "tender-workspace/internal/entity"
	f "tender-workspace/internal/utils/functions"
	mc "tender-workspace/internal/utils/myconstants"

	"go.uber.org/zap"
)

func Recover(h http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error(fmt.Sprintf("error while handling request: %v", err))
				props := f.NewResponseProps(w, e.ResponseReason{Reason: "internal server error, please try again"}, http.StatusInternalServerError, mc.ApplicationJson)
				f.Response(props)
				return
			}
		}()
		h.ServeHTTP(w, r)
	})
}
