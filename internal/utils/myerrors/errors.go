package myerrors

import "errors"

// HTTP
var (
	ErrQPLimit           = errors.New("parameter 'limit' must be positive number")
	ErrQPOffset          = errors.New("parameter 'offset' must be positive number")
	ErrQPChangeStatus    = errors.New("parameter 'status' must be in list(Created, Published, Canceled)")
	ErrQPDecision        = errors.New("parameter 'decision' must be in list(Approved, Rejected)")
	ErrQPServiceType     = errors.New("parameter 'service_type' must be in list(Construction, Delivery, Manufacture)")
	ErrQPBidStatus       = errors.New("parameter 'status' must be in list(Created, Published, Canceled)")
	ErrQPBidStatusUpdate = errors.New("parameter 'status' must be in list(Created, Published, Canceled, Approved, Rejected)")
	ErrQPOrgType         = errors.New("parameter 'type' must be in list(IE, LLC, JSC)")

	ErrBadPermission          = errors.New("you doesn't have sufficient rights to obtain the resource")
	ErrBidYourself            = errors.New("you can't offer your own company a service")
	ErrSetDeprecatedStatus    = errors.New("you can't set this status to bid")
	ErrUserAndOrg             = errors.New("you aren't responsible for this organizaton")
	ErrMethodNotAllowed       = errors.New("method not allowed")
	ErrInternal               = errors.New("internal server error, please try again later")
	ErrRequestBody            = errors.New("invalid request body")
	ErrUserExist              = errors.New("you aren't authorized")
	ErrUserAlreadyExist       = errors.New("this username is already reserved")
	ErrUserNotExist           = errors.New("user with this username is not exist")
	ErrOrganizationExist      = errors.New("organization doesn't exist")
	ErrTenderExist            = errors.New("tender doesn't exist")
	ErrUserIsNotResponsible   = errors.New("you aren't responsible for any organizaton")
	ErrUserAlreadyResponsible = errors.New("you are already responsible for this organizaton")
	ErrUserAlreadyHasBid      = errors.New("you are already has bid to this tender")
	ErrOrgAlreadyHasBid       = errors.New("your organization already has bid to this tender")

	ErrBidID           = errors.New("you have specified incorrect parameter 'bidId'")
	ErrTenderID        = errors.New("you have specified incorrect parameter 'tenderId'")
	ErrTenderStatus    = errors.New("you have specified incorrect parameter 'status'")
	ErrNoTenders       = errors.New("there are no tenders specified by your request")
	ErrNoBids          = errors.New("there are no bids specified by your request")
	ErrBadStatusCreate = errors.New("you must specify field 'status' with value 'Created'")
	ErrResponsibilty   = errors.New("you aren't responsible for this organization")
	ErrBigInterval     = errors.New("offset is bigger than size of selected tenders")

	ErrExistServiceType = errors.New("you must specify parameter 'serviceType'")
	ErrExistUsername    = errors.New("you must specify parameter 'username'")
	ErrExistBidID       = errors.New("you must specify parameter 'bidId'")
	ErrExistDecision    = errors.New("you must specify parameter 'decision'")
	ErrExistFeedback    = errors.New("you must specify parameter 'feedback'")
	ErrExistStatus      = errors.New("you must specify parameter 'status'")
	ErrExistType        = errors.New("you must specify parameter 'type'")
	ErrExistTenderID    = errors.New("you must specify parameter 'tenderId'")
)

// DATABASE
var (
	ErrNoRowsAffected = errors.New("no rows affected")
)
