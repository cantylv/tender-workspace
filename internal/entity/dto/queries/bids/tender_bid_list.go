package queries

import (
	"net/http"
	"strconv"
	e "tender-workspace/internal/utils/myerrors"
)

type TenderBidList struct {
	TenderID int
	Username string
	Limit    int
	Offset   int
}

func (q *TenderBidList) GetParameters(r *http.Request) error {
	q.TenderID = 0 // explicit
	tenderIdStr := r.Header.Get("tenderId")
	if tenderIdStr == "" {
		return e.ErrExistTenderID
	}
	tenderId, err := strconv.Atoi(tenderIdStr)
	if err != nil {
		return e.ErrTenderID
	}
	if tenderId < 1 {
		return e.ErrTenderID
	}
	q.TenderID = tenderId

	username := r.Header.Get("username")
	if username == "" {
		return e.ErrBadPermission
	}
	q.Username = username

	q.Limit = 5
	limitStr := r.Header.Get("limit")
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return err
		}
		if limit < 0 {
			return e.ErrQPLimit
		}
		q.Limit = limit
	}

	q.Offset = 0 // explicit
	offsetStr := r.Header.Get("offset")
	if offsetStr != "" {
		offset, err := strconv.Atoi(limitStr)
		if err != nil {
			return err
		}
		if offset < 0 {
			return e.ErrQPOffset
		}
		q.Offset = offset
	}
	return nil
}
