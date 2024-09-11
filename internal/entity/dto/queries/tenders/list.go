// Used for get list of tenders
package queries

import (
	"net/http"
	"strconv"
	e "tender-workspace/internal/utils/myerrors"
)

type ListTenders struct {
	Limit       int
	Offset      int
	ServiceType string
}

func (q *ListTenders) GetParameters(r *http.Request) error {
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

	q.ServiceType = "" // explicit
	serviceType := r.Header.Get("service_type")
	if serviceType != "" {
		q.ServiceType = serviceType
	}
	return nil
}
