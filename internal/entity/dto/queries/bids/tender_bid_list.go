package queries

import (
	"net/http"
	"strconv"
	e "tender-workspace/internal/utils/myerrors"

	"github.com/gorilla/mux"
)

type TenderBidList struct {
	TenderID int
	Username string
	Limit    int
	Offset   int
}

func (q *TenderBidList) GetParameters(r *http.Request) error {
	tenderId, err := strconv.Atoi(mux.Vars(r)["tenderId"])
	if err != nil || tenderId < 1 {
		return e.ErrTenderID
	}
	q.TenderID = tenderId

	queryParams := r.URL.Query()

	username := queryParams.Get("username")
	if username == "" {
		return e.ErrBadPermission
	}
	q.Username = username

	q.Limit = 5
	limitStr := queryParams.Get("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 0 {
			return e.ErrQPLimit
		}
		q.Limit = limit
	}

	q.Offset = 0 // explicit
	offsetStr := queryParams.Get("offset")
	if offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			return e.ErrQPOffset
		}
		q.Offset = offset
	}
	return nil
}
