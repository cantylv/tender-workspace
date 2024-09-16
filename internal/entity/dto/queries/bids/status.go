package queries

import (
	"net/http"
	"strconv"
	"strings"
	mc "tender-workspace/internal/utils/myconstants"
	e "tender-workspace/internal/utils/myerrors"

	"github.com/gorilla/mux"
)

type BidStatus struct {
	BidID    int
	Username string
}

func (q *BidStatus) GetParameters(r *http.Request) error {
	bidIdStr := mux.Vars(r)["bidId"]
	if bidIdStr == "" {
		return e.ErrExistBidID
	}
	bidId, err := strconv.Atoi(bidIdStr)
	if err != nil || bidId < 1 {
		return e.ErrBidID
	}
	q.BidID = bidId

	username := r.URL.Query().Get("username")
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
	bidIdStr := mux.Vars(r)["bidId"]
	if bidIdStr == "" {
		return e.ErrExistBidID
	}
	bidId, err := strconv.Atoi(bidIdStr)
	if err != nil || bidId < 1 {
		return e.ErrBidID
	}
	q.BidID = bidId

	queryParams := r.URL.Query()
	status := queryParams.Get("status")
	if status == "" {
		return e.ErrExistStatus
	}
	
	statusLower := strings.ToLower(status)
	if _, ok := mc.AvaliableBidStatus[statusLower]; !ok {
		return e.ErrQPBidStatusUpdate
	}
	runes := []rune(statusLower)
	q.Status = strings.ToUpper(string(runes[0])) + string(runes[1:])

	username := queryParams.Get("username")
	if username == "" {
		return e.ErrBadPermission
	}
	q.Username = username
	return nil
}
