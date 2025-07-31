package errs

import "errors"

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrUsernameTaken   = errors.New("username already taken")
	ErrHashingFailed   = errors.New("failed to hash password")

	ErrInternal = errors.New("internal error")
)
