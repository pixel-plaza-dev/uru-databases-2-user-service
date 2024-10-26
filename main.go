package main

import (
	"github.com/joho/godotenv"
	customEnvironmentError "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/custom_error/environment"
	customListenerError "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/custom_error/listener"
	commonLoader "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/loader"
	commonMongoDb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb"
	"google.golang.org/grpc"
	"net"
	"pixel_plaza/users_service/config"
	"pixel_plaza/users_service/grpc_server"
	"pixel_plaza/users_service/mongodb"
	protobuf "pixel_plaza/users_service/protobuf/pixel_plaza/users_service"
)

// Load environment variables
func init() {
	if err := godotenv.Load(); err != nil {
		panic(customEnvironmentError.FailedToLoadEnvironmentVariablesError{Err: err})
	}
}

func main() {
	// Get the port and listener port
	servicePort, err := commonLoader.LoadServicePort(config.UsersServicePortKey)
	if err != nil {
		panic(err)
	}
	config.EnvironmentLogger.EnvironmentVariableLoaded(config.UsersServicePortKey)

	// Get the MongoDB URI
	mongoDbUri, err := commonLoader.LoadMongoDBURI(config.MongoDbUriKey)
	if err != nil {
		panic(err)
	}
	config.EnvironmentLogger.EnvironmentVariableLoaded(config.MongoDbUriKey)

	// Get the required MongoDB database name
	mongoDbName, err := commonLoader.LoadMongoDBName(config.MongoDbNameKey)
	if err != nil {

		panic(err)
	}
	config.EnvironmentLogger.EnvironmentVariableLoaded(config.MongoDbNameKey)

	// Get the MongoDB configuration
	mongoDbConfig := &commonMongoDb.Config{Uri: mongoDbUri, Timeout: mongodb.ConnectionCtxTimeout}

	// Connect to MongoDB
	mongodbConnection, err := commonMongoDb.Connect(mongoDbConfig)
	if err != nil {
		panic(err)
	}
	defer func() {
		// Disconnect from MongoDB
		commonMongoDb.Disconnect(mongodbConnection)
		config.MongoDbLogger.DisconnectedFromMongoDB()
	}()
	config.MongoDbLogger.ConnectedToMongoDB()

	// Create users service database handler
	usersServiceDatabase, err := mongodb.NewUserDatabase(mongodbConnection, mongoDbName)
	if err != nil {
		panic(err)
	}

	// Listen on the given port
	listener, err := net.Listen("tcp", servicePort.FormattedPort)
	if err != nil {
		panic(customListenerError.FailedToListenError{Err: err})
	}
	defer func() {
		if err := listener.Close(); err != nil {
			panic(customListenerError.FailedToCloseError{Err: err})
		}
	}()

	// Create a new gRPC server
	s := grpc.NewServer()

	// Create a new gRPC UsersServiceServer
	usersServiceServer := grpc_server.NewUsersServiceServer(usersServiceDatabase, config.UsersServiceLogger)

	// Register the UsersServiceServer with the gRPC server
	protobuf.RegisterUsersServiceServer(s, usersServiceServer)
	config.ListenerLogger.ServerStarted(servicePort.Port)

	// Serve the gRPC server
	if err := s.Serve(listener); err != nil {
		panic(customListenerError.FailedToServeError{Err: err})
	}
	config.ListenerLogger.ServerStarted(servicePort.Port)
	defer s.Stop()
}
