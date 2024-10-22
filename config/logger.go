package config

import (
	"github.com/pixel-plaza-dev/uru-databases-2-go-service-common/logger"
)

var (
	// ListenerLogger is the logger for the listener
	ListenerLogger = logger.NewListenerLogger(ListenerLoggerName)

	// EnvironmentLogger is the logger for the environment
	EnvironmentLogger = logger.NewEnvironmentLogger(EnvironmentLoggerName)

	// MongoDBLogger is the logger for the MongoDB client
	MongoDBLogger = logger.NewMongoDbLogger(MongoDbLoggerName)
)
