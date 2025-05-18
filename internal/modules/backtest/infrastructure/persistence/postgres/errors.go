package postgres

import "errors"

// Common database errors
var (
	// ErrDatabaseError is a generic database error
	ErrDatabaseError = errors.New("database error")

	// Repository-specific errors
	ErrPositionNotFound   = errors.New("position not found")
	ErrRunNotFound        = errors.New("run not found")
	ErrRunGroupNotFound   = errors.New("run group not found")
	ErrRunMetricsNotFound = errors.New("run metrics not found")
)
