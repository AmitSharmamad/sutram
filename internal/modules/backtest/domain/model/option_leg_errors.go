package model

import "errors"

// Common errors for option leg domain
var (
	ErrInvalidOptionDetails = errors.New("option type is required for F&O segment")
	ErrInvalidStrikePrice   = errors.New("strike price is required for F&O segment")
	ErrInvalidEquityDetails = errors.New("option type and strike price should not be set for EQ segment")
	ErrInvalidOrder         = errors.New("order must be greater than 0")
	ErrInvalidTotalLot      = errors.New("total lot must be greater than 0")
)
