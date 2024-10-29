package main

import (
	"flag"
	"github.com/joho/godotenv"
	commonenverror "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/env/error"
	commonflag "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/flag"
	commonlistener "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/listener"
	commonlistenererror "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/listener/error"
	commonmongodb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb"
	protobuf "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/compiled-protobuf/user"
	userserver "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/grpc/server/user"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/listener"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/logger"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb"
	userdatabase "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb/database/user"
	"google.golang.org/grpc"
	"net"
)

// Load environment variables
func init() {
	if err := godotenv.Load(); err != nil {
		panic(commonenverror.FailedToLoadEnvironmentVariablesError{Err: err})
	}
}

func main() {
	// Declare flags and parse them
	commonflag.SetModeFlag()
	flag.Parse()
	logger.FlagLogger.ModeFlagSet(commonflag.Mode)

	// Get the listener port
	servicePort, err := commonlistener.LoadServicePort(listener.UsersServicePortKey)
	if err != nil {
		panic(err)
	}
	logger.EnvironmentLogger.EnvironmentVariableLoaded(listener.UsersServicePortKey)

	// Get the MongoDB URI
	mongoDbUri, err := commonmongodb.LoadMongoDBURI(mongodb.MongoDbUriKey)
	if err != nil {
		panic(err)
	}
	logger.EnvironmentLogger.EnvironmentVariableLoaded(mongodb.MongoDbUriKey)

	// Get the required MongoDB database name
	mongoDbName, err := commonmongodb.LoadMongoDBName(mongodb.MongoDbNameKey)
	if err != nil {

		panic(err)
	}
	logger.EnvironmentLogger.EnvironmentVariableLoaded(mongodb.MongoDbNameKey)

	// Get the MongoDB configuration
	mongoDbConfig := &commonmongodb.Config{Uri: mongoDbUri, Timeout: mongodb.ConnectionCtxTimeout}

	// Connect to MongoDB
	mongodbConnection, err := commonmongodb.Connect(mongoDbConfig)
	if err != nil {
		panic(err)
	}
	defer func() {
		// Disconnect from MongoDB
		commonmongodb.Disconnect(mongodbConnection)
		logger.MongoDbLogger.DisconnectedFromMongoDB()
	}()
	logger.MongoDbLogger.ConnectedToMongoDB()

	// Create user database handler
	userDatabase, err := userdatabase.NewDatabase(mongodbConnection, mongoDbName)
	if err != nil {
		panic(err)
	}

	// Listen on the given port
	portListener, err := net.Listen("tcp", servicePort.FormattedPort)
	if err != nil {
		panic(commonlistenererror.FailedToListenError{Err: err})
	}
	defer func() {
		if err := portListener.Close(); err != nil {
			panic(commonlistenererror.FailedToCloseError{Err: err})
		}
	}()

	// Create a new gRPC server
	s := grpc.NewServer()

	// Create a new gRPC UsersServiceServer
	usersServiceServer := userserver.NewServer(userDatabase, logger.UserDatabaseLogger)

	// Register the user server with the gRPC server
	protobuf.RegisterUserServer(s, usersServiceServer)
	logger.ListenerLogger.ServerStarted(servicePort.Port)

	// Serve the gRPC server
	if err := s.Serve(portListener); err != nil {
		panic(commonlistenererror.FailedToServeError{Err: err})
	}
	logger.ListenerLogger.ServerStarted(servicePort.Port)
	defer s.Stop()
}
