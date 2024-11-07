package user

import "errors"

var (
	PasswordDoesNotMatchError = errors.New("password does not match")
)
