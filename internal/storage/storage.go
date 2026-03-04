package storage

import "errors"

var (
	ErrUserExists   = errors.New("user alreadyexists")
	ErrUserNotFound = errors.New("ser not found")
	ErrAppNotFound  = errors.New("app not found")
)
