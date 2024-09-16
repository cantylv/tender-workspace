package queries

import (
	"net/http"
	"strconv"
	e "tender-workspace/internal/utils/myerrors"

	"github.com/gorilla/mux"
)

type UpdateBidData struct {
	BidID    int
	Username string
}

func (q *UpdateBidData) GetParameters(r *http.Request) error {
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
