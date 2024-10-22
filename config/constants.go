package config

import "time"

// Constants for the configuration of the application
const (
	// UsersServicePortKey is the key of the default port for the application
	UsersServicePortKey = "USERS_SERVICE_PORT"

	// MongoDbUriKey is the key of the MongoDB host
	MongoDbUriKey = "MONGO_DB_HOST"

	// MongoDbNameKey is the key of the MongoDB database name
	MongoDbNameKey = "MONGO_DB_NAME"

	// ListenerLoggerName is the name of the listener logger
	ListenerLoggerName = "Net Listener"

	// EnvironmentLoggerName is the name of the environment logger
	EnvironmentLoggerName = "Environment"

	// MongoDbLoggerName is the name of the MongoDB logger
	MongoDbLoggerName = "MongoDB"

	// UsersServiceLoggerName is the name of the UsersService logger
	UsersServiceLoggerName = "UsersService"

	// MongoDbConnectionTimeout is the timeout for the MongoDB connection
	MongoDbConnectionTimeout = 60 * time.Second
)
