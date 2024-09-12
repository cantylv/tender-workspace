package queries

import (
	"net/http"
	"strconv"
	"strings"
	mc "tender-workspace/internal/utils/myconstants"
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
	if serviceType == "" {
		return nil
	}
	serviceType = strings.ToLower(serviceType)
	if _, ok := mc.AvaliableServiceType[serviceType]; !ok {
		return e.ErrQPServiceType
	}
	runes := []rune(serviceType)
	q.ServiceType = strings.ToUpper(string(runes[0])) + string(serviceType[1:len(runes)])
	return nil
}
