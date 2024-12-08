package user

import "errors"

var (
	NilDatabaseError        = errors.New("user database cannot be nil")
	EmailAlreadyExistsError = errors.New("user email already exists")
)
