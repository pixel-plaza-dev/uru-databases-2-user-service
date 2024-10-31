package main

import (
	"flag"
	"github.com/joho/godotenv"
	commonenv "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/env"
	commonflag "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/flag"
	commonlistener "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/listener"
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
		panic(commonenv.FailedToLoadEnvironmentVariablesError)
	}
}

func main() {
	// Declare flags and parse them
	commonflag.SetModeFlag()
	flag.Parse()
	logger.FlagLogger.ModeFlagSet(commonflag.Mode)

	// Get the listener port
	servicePort, err := commonlistener.LoadServicePort(listener.UserServicePortKey)
	if err != nil {
		panic(err)
	}
	logger.EnvironmentLogger.EnvironmentVariableLoaded(listener.UserServicePortKey)

	// Get the MongoDB URI
	mongoDbUri, err := commonmongodb.LoadMongoDBURI(mongodb.UriKey)
	if err != nil {
		panic(err)
	}
	logger.EnvironmentLogger.EnvironmentVariableLoaded(mongodb.UriKey)

	// Get the required MongoDB database name
	mongoDbName, err := commonmongodb.LoadMongoDBName(mongodb.DbNameKey)
	if err != nil {

		panic(err)
	}
	logger.EnvironmentLogger.EnvironmentVariableLoaded(mongodb.DbNameKey)

	// Get the MongoDB configuration
	mongoDbConfig := &commonmongodb.Config{Uri: mongoDbUri, Timeout: mongodb.ConnectionCtxTimeout}

	// Connect to MongoDB
	mongodbConnection := commonmongodb.NewDefaultConnectionHandler(mongoDbConfig)
	if _, err = mongodbConnection.Connect(); err != nil {
		panic(err)
	}
	defer func() {
		// Disconnect from MongoDB
		mongodbConnection.Disconnect()
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
		panic(commonlistener.FailedToListenError)
	}
	defer func() {
		if err := portListener.Close(); err != nil {
			panic(commonlistener.FailedToCloseError)
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
	if err = s.Serve(portListener); err != nil {
		panic(commonlistener.FailedToServeError)
	}
	logger.ListenerLogger.ServerStarted(servicePort.Port)
	defer s.Stop()
}
