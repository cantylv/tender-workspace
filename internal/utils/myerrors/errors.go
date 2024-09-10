package myerrors

import "errors"

// QUERY PARAMETERS
var (
	ErrInvalidQueryParameterLimit  = errors.New("parameter limit must be positive")
	ErrInvalidQueryParameterOffset = errors.New("parameter offset must be positive")
	ErrMethodNotAllowed            = errors.New("method not allowed")
)

// DATABASE
var (
	ErrNoRowsAffected = errors.New("no rows affected")
)
