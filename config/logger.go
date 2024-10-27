package config

import (
	commonLogger "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/logger"
	"github.com/pixel-plaza-dev/uru-databases-2-users-service/logger"
)

const (
	// ListenerLoggerName is the name of the listener logger
	ListenerLoggerName = "Net Listener"

	// EnvironmentLoggerName is the name of the environment logger
	EnvironmentLoggerName = "Environment"

	// MongoDbLoggerName is the name of the MongoDB logger
	MongoDbLoggerName = "MongoDB"

	// UsersServiceLoggerName is the name of the UsersService logger
	UsersServiceLoggerName = "UsersService"
)

var (
	// ListenerLogger is the logger for the listener
	ListenerLogger = commonLogger.NewListenerLogger(ListenerLoggerName)

	// EnvironmentLogger is the logger for the environment
	EnvironmentLogger = commonLogger.NewEnvironmentLogger(EnvironmentLoggerName)

	// MongoDbLogger is the logger for the MongoDB client
	MongoDbLogger = commonLogger.NewMongoDbLogger(MongoDbLoggerName)

	// UsersServiceLogger is the logger for the Users Service server
	UsersServiceLogger = logger.NewUsersServiceLogger(UsersServiceLoggerName)
)
