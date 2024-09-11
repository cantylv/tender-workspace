package queries

import (
	"net/http"
	"strconv"
	e "tender-workspace/internal/utils/myerrors"
)

type TenderStatus struct {
	TenderID int
	Username string
}

func (q *TenderStatus) GetParameters(r *http.Request) error {
	q.TenderID = 0 // explicit
	tenderIdStr := r.Header.Get("tenderId")
	if tenderIdStr == "" {
		return e.ErrExistTenderID
	}
	tenderId, err := strconv.Atoi(tenderIdStr)
	if err != nil {
		return err
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
	return nil
}

type UpdateTenderStatus struct {
	TenderID int
	Username string
	Status   string
}

func (q *UpdateTenderStatus) GetParameters(r *http.Request) error {
	q.TenderID = 0 // explicit
	tenderIdStr := r.Header.Get("tenderId")
	if tenderIdStr == "" {
		return e.ErrExistTenderID
	}
	tenderId, err := strconv.Atoi(tenderIdStr)
	if err != nil {
		return err
	}
	if tenderId < 1 {
		return e.ErrTenderID
	}
	q.TenderID = tenderId

	status := r.Header.Get("status")
	if status == "" {
		return e.ErrTenderStatus
	}
	q.Status = status

	username := r.Header.Get("username")
	if username == "" {
		return e.ErrBadPermission
	}
	q.Username = username
	return nil
}
