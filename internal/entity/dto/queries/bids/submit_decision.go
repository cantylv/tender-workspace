package queries

import (
	"net/http"
	"strconv"
	e "tender-workspace/internal/utils/myerrors"
)

var avaliableDecision = map[string]bool{
	"Approved": true,
	"Rejected": true,
}

type SubmitDecision struct {
	BidID    int
	Decision string
	Username string
}

func (q *SubmitDecision) GetParameters(r *http.Request) error {
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

	decision := r.Header.Get("decision")
	if decision == "" {
		return e.ErrExistDecision
	}
	if _, ok := avaliableDecision[decision]; !ok {
		return e.ErrQPDecision
	}
	q.Decision = decision

	username := r.Header.Get("username")
	if username == "" {
		return e.ErrBadPermission
	}
	q.Username = username
	return nil
}
