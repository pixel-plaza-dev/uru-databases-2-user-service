package user

import commonlogger "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/utils/logger"

type Logger struct {
	logger commonlogger.Logger
}

// NewLogger is the logger for the user database
func NewLogger(logger commonlogger.Logger) Logger {
	return Logger{logger: logger}
}

// SignedUp logs the user sign up
func (l Logger) SignedUp(userId string) {
	l.logger.LogMessageWithDetails("User signed up", userId)
}

// SignUpFailed logs the user sign up failure
func (l Logger) SignUpFailed(err error) {
	l.logger.LogMessageWithDetails("User sign up failed", err.Error())
}

// PasswordIsCorrect logs the password check success
func (l Logger) PasswordIsCorrect(userId string) {
	l.logger.LogMessageWithDetails("Password is correct", userId)
}

// PasswordIsIncorrect logs the password check failure
func (l Logger) PasswordIsIncorrect(userId string) {
	l.logger.LogMessageWithDetails("Password is incorrect", userId)
}

// PasswordIsCorrectFailed logs the password check failure
func (l Logger) PasswordIsCorrectFailed(err error) {
	l.logger.LogMessageWithDetails("Password check failed", err.Error())
}

// UserFoundByUsername logs the user retrieval success
func (l Logger) UserFoundByUsername(username string) {
	l.logger.LogMessageWithDetails("User found by username", username)
}

// UserNotFoundByUsername logs the user retrieval failure
func (l Logger) UserNotFoundByUsername(username string) {
	l.logger.LogMessageWithDetails("User not found by username", username)
}

// UserFoundByUserId logs the user retrieval success
func (l Logger) UserFoundByUserId(userId string) {
	l.logger.LogMessageWithDetails("User found by user ID", userId)
}

// UserNotFoundByUserId logs the user retrieval failure
func (l Logger) UserNotFoundByUserId(userId string) {
	l.logger.LogMessageWithDetails("User not found by user ID", userId)
}

// UsernameExists logs the username check success
func (l Logger) UsernameExists(username string) {
	l.logger.LogMessageWithDetails("Username exists", username)
}

// UsernameExistsFailed logs the username check failure
func (l Logger) UsernameExistsFailed(err error) {
	l.logger.LogMessageWithDetails("Username exists check failed", err.Error())
}

// GetUsernameByUserIdFailed logs the username retrieval failure
func (l Logger) GetUsernameByUserIdFailed(err error) {
	l.logger.LogMessageWithDetails("Failed to fetch username by user ID", err.Error())
}

// GetUserIdByUsernameFailed logs the user ID retrieval failure
func (l Logger) GetUserIdByUsernameFailed(err error) {
	l.logger.LogMessageWithDetails("Failed to fetch user ID by username", err.Error())
}

// MissingTokenClaimsSubject logs the missing token claims subject
func (l Logger) MissingTokenClaimsSubject() {
	l.logger.LogMessage("Missing token claims subject")
}

// UpdateUser logs the user update
func (l Logger) UpdateUser(userId string) {
	l.logger.LogMessageWithDetails("User updated", userId)
}

// UpdateUserFailed logs the user update failure
func (l Logger) UpdateUserFailed(err error) {
	l.logger.LogMessageWithDetails("User update failed", err.Error())
}

// GetPhoneNumber logs the user phone number retrieval
func (l Logger) GetPhoneNumber(userId string) {
	l.logger.LogMessageWithDetails("Fetched user phone number", userId)
}

// GetPhoneNumberFailed logs the user phone number retrieval failure
func (l Logger) GetPhoneNumberFailed(err error) {
	l.logger.LogMessageWithDetails("Failed to fetch user phone number", err.Error())
}

// GetProfile logs the user profile update
func (l Logger) GetProfile(userId string) {
	l.logger.LogMessageWithDetails("Fetched user profile", userId)
}

// GetProfileFailed logs the user profile update failure
func (l Logger) GetProfileFailed(err error) {
	l.logger.LogMessageWithDetails("Failed to fetch user profile", err.Error())
}

// UpdateUsername logs the user username update
func (l Logger) UpdateUsername(userId string) {
	l.logger.LogMessageWithDetails("User username updated", userId)
}

// UpdateUsernameFailed logs the user username update failure
func (l Logger) UpdateUsernameFailed(err error) {
	l.logger.LogMessageWithDetails("User username update failed", err.Error())
}

// UpdatePassword logs the user password update
func (l Logger) UpdatePassword(userId string) {
	l.logger.LogMessageWithDetails("User password updated", userId)
}

// UpdatePasswordFailed logs the user password update failure
func (l Logger) UpdatePasswordFailed(err error) {
	l.logger.LogMessageWithDetails("User password update failed", err.Error())
}

// HashPasswordFailed logs a failed password hash attempt
func (l Logger) HashPasswordFailed(err error) {
	l.logger.LogMessageWithDetails("Failed to hash password", err.Error())
}

// UpdatePhoneNumber logs the user phone number update
func (l Logger) UpdatePhoneNumber(userId string) {
	l.logger.LogMessageWithDetails("User phone number updated", userId)
}

// UpdatePhoneNumberFailed logs the user phone number update failure
func (l Logger) UpdatePhoneNumberFailed(err error) {
	l.logger.LogMessageWithDetails("User phone number update failed", err.Error())
}
