package user

import (
	commonlogger "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/logger"
)

type Logger struct {
	logger commonlogger.Logger
}

// NewLogger is the logger for the user database
func NewLogger(logger commonlogger.Logger) Logger {
	return Logger{logger: logger}
}

// FailedToCreateDocument logs a FailedToCreateDocumentError
func (l Logger) FailedToCreateDocument(err error) {
	l.logger.LogMessageWithDetails("Failed to create document", err.Error())
}

// UserCreated logs a success message when a user is created
func (l Logger) UserCreated(id string) {
	l.logger.LogMessageWithDetails("User created", id)
}
