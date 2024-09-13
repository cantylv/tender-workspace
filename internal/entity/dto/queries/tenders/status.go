package queries

import (
	"net/http"
	"strconv"
	"strings"
	mc "tender-workspace/internal/utils/myconstants"
	e "tender-workspace/internal/utils/myerrors"

	"github.com/gorilla/mux"
)

type TenderStatus struct {
	TenderID int
	Username string
}

func (q *TenderStatus) GetParameters(r *http.Request) error {
	vars := mux.Vars(r)
	q.TenderID = 0 // explicit
	tenderIdStr := vars["tenderId"]
	if tenderIdStr == "" {
		return e.ErrExistTenderID
	}
	tenderId, err := strconv.Atoi(tenderIdStr)
	if err != nil || tenderId < 1 {
		return e.ErrTenderID
	}
	q.TenderID = tenderId

	username := r.URL.Query().Get("username")
	if username == "" {
		return e.ErrBadPermission
	}
	q.Username = username
	return nil
}

type UpdateTenderStatus struct {
	TenderID int
	Username string
	Status   string
}

func (q *UpdateTenderStatus) GetParameters(r *http.Request) error {
	vars := mux.Vars(r)
	q.TenderID = 0 // explicit
	tenderIdStr := vars["tenderId"]
	if tenderIdStr == "" {
		return e.ErrExistTenderID
	}
	tenderId, err := strconv.Atoi(tenderIdStr)
	if err != nil || tenderId < 1 {
		return e.ErrTenderID
	}
	q.TenderID = tenderId

	queryParams := r.URL.Query()

	q.Status = "" // explicit
	status := queryParams.Get("status")
	if status == "" {
		return e.ErrTenderStatus
	}
	statusLower := strings.ToLower(status)
	if _, ok := mc.AvaliableTenderStatus[statusLower]; !ok {
		return e.ErrQPChangeStatus
	}
	runes := []rune(statusLower)
	q.Status = strings.ToUpper(string(runes[0])) + string(statusLower[1:len(runes)])

	username := queryParams.Get("username")
	if username == "" {
		return e.ErrBadPermission
	}
	q.Username = username
	return nil
}
