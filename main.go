package main

import (
	"github.com/joho/godotenv"
	commonLoader "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/loader"
	commonMongoDb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb"
	"google.golang.org/grpc"
	"net"
	"pixel_plaza/users_service/config"
	"pixel_plaza/users_service/grpc_server"
	protobuf "pixel_plaza/users_service/protobuf/pixel_plaza/users_service"
)

// Load environment variables
func init() {
	if err := godotenv.Load(); err != nil {
		config.EnvironmentLogger.ErrorLoadingEnvironmentVariables(err)
	}
}

func main() {
	// Get the port and listener port
	port, listenerPort := commonLoader.LoadServicePort(config.UsersServicePortKey, config.EnvironmentLogger)

	// Get the MongoDB URI
	mongoDbUri := commonLoader.LoadMongoDBURI(config.MongoDbUriKey, config.EnvironmentLogger)

	// Connect to MongoDB
	mongoDbClient, mongoDbContext, mongoDbCancel, err := commonMongoDb.Connect(mongoDbUri, config.MongoDBLogger, config.MongoDbConnectionTimeout)
	if err != nil {
		config.MongoDBLogger.FailedToConnectToMongoDb(err)
	}
	defer commonMongoDb.Disconnect(mongoDbClient, mongoDbContext, mongoDbCancel, config.MongoDBLogger)

	// Listen on the given port
	listener, err := net.Listen("tcp", listenerPort)
	if err != nil {
		config.ListenerLogger.FailedToListen(err)
	}
	defer func() {
		if err := listener.Close(); err != nil {
			config.ListenerLogger.FailedToClose(err)
		}
	}()

	// Create a new gRPC server
	s := grpc.NewServer()

	// Create a new gRPC UsersServiceServer
	usersServiceServer := grpc_server.NewUsersServiceServer(mongoDbClient)

	// Register the UsersServiceServer with the gRPC server
	protobuf.RegisterUsersServiceServer(s, usersServiceServer)
	config.ListenerLogger.ServerStarted(port)

	// Serve the gRPC server
	if err := s.Serve(listener); err != nil {
		config.ListenerLogger.FailedToServe(err)
	}
	defer s.Stop()
}
