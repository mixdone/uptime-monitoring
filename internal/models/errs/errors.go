package errs

import "errors"

var (
	ErrTokenExpired     = errors.New("token is expired")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrTokenWrongFormat = errors.New("wrong token format")

	ErrSessionNotFound = errors.New("session not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrUsernameTaken   = errors.New("username already taken")
	ErrHashingFailed   = errors.New("failed to hash password")

	ErrInternal = errors.New("internal error")

	ErrNotFound = errors.New("resource not found ")
)
