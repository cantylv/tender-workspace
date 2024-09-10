package entity

import (
	"net/http"
	"strconv"
	e "tender-workspace/internal/utils/myerrors"
)

type TenderQueryParameters struct {
	Limit       int32
	Offset      int32
	ServiceType string
}

func (q *TenderQueryParameters) GetParameters(r *http.Request) error {
	limitStr := r.Header.Get("limit")
	q.Limit = 5
	if limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return err
		}
		if limit < 0 {
			return e.ErrInvalidQueryParameterLimit
		}
		q.Limit = int32(limit)
	}
	offsetStr := r.Header.Get("offset")
	if offsetStr != "" {
		offset, err := strconv.Atoi(limitStr)
		if err != nil {
			return err
		}
		if offset < 0 {
			return e.ErrInvalidQueryParameterOffset
		}
		q.Offset = int32(offset)
	}
	serviceType := r.Header.Get("service_type")
	if serviceType != "" {
		q.ServiceType = serviceType
	}
	return nil
}
