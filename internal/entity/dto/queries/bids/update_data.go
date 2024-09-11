package queries

import (
	"net/http"
	"strconv"
	e "tender-workspace/internal/utils/myerrors"
)

type UpdateBidData struct {
	BidID    int
	Username string
}

func (q *UpdateBidData) GetParameters(r *http.Request) error {
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
