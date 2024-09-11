package queries

import (
	"net/http"
	"strconv"
	e "tender-workspace/internal/utils/myerrors"
)

var avaliableStatus = map[string]bool{
	"Created":   true,
	"Published": true,
	"Canceled":  true,
	"Approved":  true,
	"Rejected":  true,
}

type BidStatus struct {
	BidID    int
	Username string
}

func (q *BidStatus) GetParameters(r *http.Request) error {
	q.BidID = 0 // explicit
	bidIdStr := r.Header.Get("bidId")
	if bidIdStr == "" {
		return e.ErrExistTenderID
	}
	bidId, err := strconv.Atoi(bidIdStr)
	if err != nil {
		return err
	}
	if bidId < 1 {
		return e.ErrTenderID
	}
	q.BidID = bidId

	username := r.Header.Get("username")
	if username == "" {
		return e.ErrBadPermission
	}
	q.Username = username
	return nil
}

type UpdateBidStatus struct {
	BidID    int
	Status   string
	Username string
}

func (q *UpdateBidStatus) GetParameters(r *http.Request) error {
	q.BidID = 0 // explicit
	bidIdStr := r.Header.Get("bidId")
	if bidIdStr == "" {
		return e.ErrExistTenderID
	}
	bidId, err := strconv.Atoi(bidIdStr)
	if err != nil {
		return err
	}
	if bidId < 1 {
		return e.ErrTenderID
	}
	q.BidID = bidId

	status := r.Header.Get("status")
	if status == "" {
		return e.ErrBadPermission
	}
	if _, ok := avaliableStatus[status]; !ok {
		return e.ErrQPStatus
	}
	q.Status = status

	username := r.Header.Get("username")
	if username == "" {
		return e.ErrBadPermission
	}
	q.Username = username
	return nil
}
