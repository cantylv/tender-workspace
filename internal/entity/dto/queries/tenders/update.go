package queries

import (
	"net/http"
	"strconv"
	e "tender-workspace/internal/utils/myerrors"

	"github.com/gorilla/mux"
)

type TenderUpdate struct {
	TenderID int
	Username string
}

func (q *TenderUpdate) GetParameters(r *http.Request) error {
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
