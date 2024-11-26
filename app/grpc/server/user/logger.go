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
func (l Logger) SignedUp(userIdentifier string) {
	l.logger.LogMessageWithDetails("User signed up", userIdentifier)
}

// SignUpFailed logs the user sign up failure
func (l Logger) SignUpFailed(err error) {
	l.logger.LogMessageWithDetails("User sign up failed", err.Error())
}

// PasswordIsCorrect logs the password check success
func (l Logger) PasswordIsCorrect(userIdentifier string) {
	l.logger.LogMessageWithDetails("Password is correct", userIdentifier)
}

// PasswordIsIncorrect logs the password check failure
func (l Logger) PasswordIsIncorrect(userIdentifier string) {
	l.logger.LogMessageWithDetails("Password is incorrect", userIdentifier)
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

// UsernameExistsFailed logs the username check failure
func (l Logger) UsernameExistsFailed(err error) {
	l.logger.LogMessageWithDetails("Username exists check failed", err.Error())
}

// GetUsernameByUserIdFailed logs the username retrieval failure
func (l Logger) GetUsernameByUserIdFailed(err error) {
	l.logger.LogMessageWithDetails("Get username by user ID failed", err.Error())
}

// GetUserIdByUsernameFailed logs the user ID retrieval failure
func (l Logger) GetUserIdByUsernameFailed(err error) {
	l.logger.LogMessageWithDetails("Get user ID by username failed", err.Error())
}

// MissingTokenClaimsSubject logs the missing token claims subject
func (l Logger) MissingTokenClaimsSubject() {
	l.logger.LogMessage("Missing token claims subject")
}

// UpdateProfile logs the user profile update
func (l Logger) UpdateProfile(userIdentifier string) {
	l.logger.LogMessageWithDetails("User profile updated", userIdentifier)
}

// UpdateProfileFailed logs the user profile update failure
func (l Logger) UpdateProfileFailed(err error) {
	l.logger.LogMessageWithDetails("User profile update failed", err.Error())
}

// HashPasswordFailed logs a failed password hash attempt
func (l Logger) HashPasswordFailed(err error) {
	l.logger.LogMessageWithDetails("Failed to hash password", err.Error())
}
