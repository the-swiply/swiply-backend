package domain

import "errors"

var (
	ErrResendIsNotAllowed = errors.New("too few time after previous code sending")
	ErrEntityIsNotExists  = errors.New("no such entity")

	ErrCodeIsIncorrect = errors.New("code is not correct")
	ErrValidateToken   = errors.New("token is not valid")
	ErrTooMuchAttempts = errors.New("too much attempts with incorrect code")
)
