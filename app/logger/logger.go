package logger

import (
	commonenv "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/env"
	commonflag "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/flag"
	commonlistener "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/listener"
	commonmongodb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb/database/user"
)

const (
	// FlagLoggerName is the name of the flag logger
	FlagLoggerName = "Flag"

	// ListenerLoggerName is the name of the listener logger
	ListenerLoggerName = "Net Listener"

	// EnvironmentLoggerName is the name of the environment logger
	EnvironmentLoggerName = "Environment"

	// MongoDbLoggerName is the name of the MongoDB logger
	MongoDbLoggerName = "MongoDB"

	// UserDatabaseLoggerName is the name of the user database logger
	UserDatabaseLoggerName = "User Database"
)

var (
	// FlagLogger is the logger for the flag
	FlagLogger = commonflag.NewLogger(FlagLoggerName)

	// ListenerLogger is the logger for the listener
	ListenerLogger = commonlistener.NewLogger(ListenerLoggerName)

	// EnvironmentLogger is the logger for the environment
	EnvironmentLogger = commonenv.NewLogger(EnvironmentLoggerName)

	// MongoDbLogger is the logger for the MongoDB client
	MongoDbLogger = commonmongodb.NewLogger(MongoDbLoggerName)

	// UserDatabaseLogger is the logger for the user database
	UserDatabaseLogger = user.NewLogger(UserDatabaseLoggerName)
)