package validator

type InvalidBirthDateError struct {
	BirthDate interface{}
}

// Error returns the error message
func (e InvalidBirthDateError) Error() string {
	return "Invalid birth date"
}
