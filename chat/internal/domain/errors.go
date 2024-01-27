package domain

import "errors"

var (
	ErrUserNotInChat = errors.New("user not in chat")

	ErrEntityIsNotExists = errors.New("no such entity")
)
