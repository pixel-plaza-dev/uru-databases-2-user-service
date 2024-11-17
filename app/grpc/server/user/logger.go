package user

import commonlogger "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/utils/logger"

type Logger struct {
	logger commonlogger.Logger
}

// NewLogger is the logger for the user database
func NewLogger(logger commonlogger.Logger) Logger {
	return Logger{logger: logger}
}

// FailedToHashPassword logs a FailedToHashPasswordError
func (l Logger) FailedToHashPassword(err error) {
	l.logger.LogMessageWithDetails("Failed to hash password", err.Error())
}

// UserSignedUp logs the user signed up
func (l Logger) UserSignedUp(userIdentifier string) {
	l.logger.LogMessageWithDetails("User signed up", userIdentifier)
}

// PasswordCheckSuccess logs the password check success
func (l Logger) PasswordCheckSuccess(userIdentifier string) {
	l.logger.LogMessageWithDetails("Password check success", userIdentifier)
}
