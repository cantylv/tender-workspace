package queries

import (
	"net/http"
	"strconv"
	"strings"
	mc "tender-workspace/internal/utils/myconstants"
	e "tender-workspace/internal/utils/myerrors"

	"github.com/gorilla/mux"
)

type SubmitDecision struct {
	BidID    int
	Decision string
	Username string
}

func (q *SubmitDecision) GetParameters(r *http.Request) error {
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
	decision := queryParams.Get("decision")
	decisionLower := strings.ToLower(decision)
	if _, ok := mc.AvaliableBidStatusApprover[decision]; !ok {
		return e.ErrQPDecision
	}
	runes := []rune(decisionLower)
	q.Decision = strings.ToUpper(string(runes[0])) + string(runes[1:])

	username := queryParams.Get("username")
	if username == "" {
		return e.ErrBadPermission
	}
	q.Username = username
	return nil
}
