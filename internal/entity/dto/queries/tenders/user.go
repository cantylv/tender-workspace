package queries

import (
	"net/http"
	"strconv"
	e "tender-workspace/internal/utils/myerrors"
)

// Used for get list of user tenders
type ListUserTenders struct {
	Limit    int
	Offset   int
	Username string
}

func (q *ListUserTenders) GetParameters(r *http.Request) error {
	q.Limit = 5
	queryParams := r.URL.Query()
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

	username := queryParams.Get("username")
	if username == "" {
		return e.ErrBadPermission
	}
	q.Username = username
	return nil
}
