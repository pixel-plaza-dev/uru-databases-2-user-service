package logger

import (
	commonenv "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/env"
	commonflag "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/flag"
	commonlistener "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/listener"
	commonlogger "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/logger"
	commonmongodb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb"
	userserver "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/grpc/server/user"
	userdatabase "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb/database/user"
)

var (
	// FlagLogger is the logger for the flag
	FlagLogger = commonflag.NewLogger(commonlogger.NewDefaultLogger("Flag"))

	// ListenerLogger is the logger for the listener
	ListenerLogger = commonlistener.NewLogger(commonlogger.NewDefaultLogger("Net Listener"))

	// EnvironmentLogger is the logger for the environment
	EnvironmentLogger = commonenv.NewLogger(commonlogger.NewDefaultLogger("Environment"))

	// MongoDbLogger is the logger for the MongoDB client
	MongoDbLogger = commonmongodb.NewLogger(commonlogger.NewDefaultLogger("MongoDB"))

	// UserServerLogger is the logger for the user server
	UserServerLogger = userserver.NewLogger(commonlogger.NewDefaultLogger("User Server"))

	// UserDatabaseLogger is the logger for the user database
	UserDatabaseLogger = userdatabase.NewLogger(commonlogger.NewDefaultLogger("User Database"))
)
