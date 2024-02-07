package repository

import "errors"

// ErrNotFound is returned when a requested resource is not found
var ErrNotFound = errors.New("not found")
var ErrDuplicate = errors.New("duplicate found")
var ErrInvalidCredentials = errors.New("invalid credentials")
