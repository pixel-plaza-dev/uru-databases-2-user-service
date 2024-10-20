package main

import (
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
	"users_service/config"
	"users_service/grpc_server"
	"users_service/logger"
	protobuf "users_service/protobuf/pixel_plaza/users_service"
)

// Load environment variables
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Listen on the given port
	port, listenerPort := config.LoadUsersServicePort()
	listener, err := net.Listen("tcp", listenerPort)
	if err != nil {
		logger.ListenerLogger.FailedToListen(err)
	}

	// Create a new gRPC grpc_server
	s := grpc.NewServer()

	// Register the grpc_server with the gRPC grpc_server
	protobuf.RegisterUsersServiceServer(s, grpc_server.UsersServiceServer)

	// Register reflection service on gRPC grpc_server.
	logger.ListenerLogger.ServerStarted(port)

	// Serve the gRPC grpc_server
	if err := s.Serve(listener); err != nil {
		logger.ListenerLogger.FailedToServe(err)
	}
}
