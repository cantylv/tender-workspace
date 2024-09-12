package queries

import (
	"net/http"
	"strconv"
	mc "tender-workspace/internal/utils/myconstants"
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
		return e.ErrTenderID
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
		return e.ErrExistBidID
	}
	bidId, err := strconv.Atoi(bidIdStr)
	if err != nil {
		return e.ErrBidID
	}
	if bidId < 1 {
		return e.ErrBidID
	}
	q.BidID = bidId

	q.Status = "" // explicit
	status := r.Header.Get("status")
	if status == "" {
		return e.ErrExistStatus
	}
	if _, ok := mc.AvaliableBidStatus[status]; !ok {
		return e.ErrQPBidStatus
	}
	q.Status = status

	username := r.Header.Get("username")
	if username == "" {
		return e.ErrBadPermission
	}
	q.Username = username
	return nil
}
