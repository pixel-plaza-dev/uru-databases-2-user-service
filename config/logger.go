package config

import (
	commonLogger "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/logger"
	"pixel_plaza/users_service/logger"
)

var (
	// ListenerLogger is the logger for the listener
	ListenerLogger = commonLogger.NewListenerLogger(ListenerLoggerName)

	// EnvironmentLogger is the logger for the environment
	EnvironmentLogger = commonLogger.NewEnvironmentLogger(EnvironmentLoggerName)

	// MongoDBLogger is the logger for the MongoDB client
	MongoDBLogger = commonLogger.NewMongoDbLogger(MongoDbLoggerName)

	// UsersServiceLogger is the logger for the Users Service server
	UsersServiceLogger = logger.NewUsersServiceLogger(UsersServiceLoggerName)
)
