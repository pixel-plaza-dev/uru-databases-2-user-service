package logger

import (
	commonLogger "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/logger"
)

type UsersServiceLogger struct {
	logger *commonLogger.Logger
}

// NewUsersServiceLogger is the logger for the Users Service server
func NewUsersServiceLogger(name string) *UsersServiceLogger {
	return &UsersServiceLogger{logger: commonLogger.NewLogger(name)}
}

// FailedToCreateDocument logs a FailedToCreateDocumentError
func (l UsersServiceLogger) FailedToCreateDocument(err error) {
	l.logger.LogMessageWithDetails("Failed to create document", err.Error())
}

// UserCreated logs a success message when a user is created
func (l UsersServiceLogger) UserCreated(id string) {
	l.logger.LogMessageWithDetails("User created", id)
}
