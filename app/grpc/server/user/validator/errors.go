package validator

import "errors"

var (
	UsernameTakenError        = errors.New("username taken")
	NewPasswordSameAsOldError = errors.New("new password same as old")
)
