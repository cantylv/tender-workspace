package healthcheck

import (
	"net/http"
	ent "tender-workspace/internal/entity"
	f "tender-workspace/internal/utils/functions"
	mc "tender-workspace/internal/utils/myconstants"
	e "tender-workspace/internal/utils/myerrors"
)

// Ping этот эндпоинт используется для проверки готовности сервера обрабатывать запросы.
// Чекер программа будет ждать первый успешный ответ и затем начнет выполнение тестовых сценариев.
func Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		propsError := f.NewResponseProps(w, ent.ResponseReason{Reason: e.ErrMethodNotAllowed.Error()}, http.StatusMethodNotAllowed, mc.ApplicationJson)
		f.Response(propsError)
		return
	}
	propsResponse := f.NewResponseProps(w, "ok", http.StatusOK, mc.TextPlain)
	f.Response(propsResponse)
}
