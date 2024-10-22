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
