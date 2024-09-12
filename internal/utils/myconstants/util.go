package myconstants

type ContextKey string

const (
	ApplicationJson = "application/json"
	TextPlain       = "text/plain"
)

const (
	RequestID = "request_id"
)

var AvaliableServiceType = map[string]bool{
	"construction": true,
	"delivery":     true,
	"manufacture":  true,
}

var AvaliableTenderStatus = map[string]bool{
	"created":   true,
	"published": true,
	"closed":    true,
}

var AvaliableBidStatus = map[string]bool{
	"created":   true,
	"published": true,
	"closed":    true,
	"approved":  true,
	"rejected":  true,
}
