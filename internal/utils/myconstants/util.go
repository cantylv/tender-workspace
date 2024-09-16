package myconstants

type ContextKey string

const (
	ApplicationJson = "application/json"
	TextPlain       = "text/plain"
)

const (
	RequestID = "request_id"
)

var AvaliableServiceType = map[string]struct{}{
	"construction": {},
	"delivery":     {},
	"manufacture":  {},
}

var AvaliableTenderStatus = map[string]struct{}{
	"created":   {},
	"published": {},
	"closed":    {},
}

var AvaliableBidStatus = map[string]struct{}{
	"created":   {},
	"published": {},
	"canceled":  {},
	"approved":  {},
	"rejected":  {},
}

var AvaliableOrganizationType = map[string]struct{}{
	"IE":  {},
	"LLC": {},
	"JSC": {},
}

var AvaliableBidStatusCreator = map[string]struct{}{
	"created":   {},
	"published": {},
	"canceled":  {},
}

var AvaliableBidStatusApprover = map[string]struct{}{
	"approved": {},
	"rejected": {},
}
