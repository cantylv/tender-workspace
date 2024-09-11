package myerrors

import "errors"

// HTTP
var (
	ErrQPLimit    = errors.New("parameter 'limit' must be positive")
	ErrQPOffset   = errors.New("parameter 'offset' must be positive")
	ErrQPStatus   = errors.New("parameter 'status' must be in list(Created, Published, Canceled, Approved, Rejected)")
	ErrQPDecision = errors.New("parameter 'decision' must be in list(Approved, Rejected)")

	ErrBadPermission    = errors.New("user doesn't have sufficient rights to obtain the resource")
	ErrMethodNotAllowed = errors.New("method not allowed")

	ErrTenderID      = errors.New("user has specified incorrect tender identifier")
	ErrExistTenderID = errors.New("user must specify tender identifier")
	ErrTenderStatus  = errors.New("user must specify tender status")

	ErrExistDecision = errors.New("user must specify decision")
	ErrExistFeedback = errors.New("user must specify feedback")
)

// DATABASE
var (
	ErrNoRowsAffected = errors.New("no rows affected")
)
