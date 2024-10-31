package validator

import "errors"

var (
	InvalidBirthDateError = errors.New("invalid birth date")
	UsernameTakenError    = errors.New("username taken")
)
