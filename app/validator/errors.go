package validator

import "errors"

var InvalidBirthDateError = errors.New("invalid birth date")
var UsernameTakenError = errors.New("username taken")
