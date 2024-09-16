package queries

import (
	"net/http"
	"strconv"
	"strings"
	mc "tender-workspace/internal/utils/myconstants"
	e "tender-workspace/internal/utils/myerrors"
)

type OrganizationList struct {
	Type   string
	Limit  int
	Offset int
}

func (q *OrganizationList) GetParameters(r *http.Request) error {
	queryParams := r.URL.Query()
	q.Limit = 5
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

	q.Type = "" // explicit
	orgType := queryParams.Get("type")
	if orgType != "" {
		typeUpper := strings.ToUpper(orgType)
		if _, ok := mc.AvaliableOrganizationType[typeUpper]; !ok {
			return e.ErrQPOrgType
		}
		q.Type = typeUpper
	}
	return nil
}
