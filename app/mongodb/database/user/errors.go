package user

import "errors"

var (
	EmailAlreadyExistsError = errors.New("user email already exists")
)
