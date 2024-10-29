package error_response

type UsernameTakenError struct{}

// Error returns a formatted error message for UsernameTakenError
func (u UsernameTakenError) Error() (message string) {
	return "Username is already taken"
}
