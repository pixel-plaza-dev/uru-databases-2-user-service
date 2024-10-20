package main

import (
	"github.com/joho/godotenv"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"users_service/config"
	"users_service/logger"
	protobuf "users_service/protobuf/pixel_plaza/users_service"
)

// Load environment variables
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

type usersServiceServer struct {
	protobuf.UnimplementedUsersServiceServer
}

// SignUp creates a new user
func (u usersServiceServer) SignUp(ctx context.Context, request *protobuf.SignUpRequest) (*protobuf.SignUpResponse, error) {
	//TODO implement me
	panic("implement me")
}

// UpdateProfile updates the user's profile
func (u usersServiceServer) UpdateProfile(ctx context.Context, request *protobuf.UpdateProfileRequest) (*protobuf.UpdateProfileResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangeUsername changes the user's username
func (u usersServiceServer) ChangeUsername(ctx context.Context, request *protobuf.ChangeUsernameRequest) (*protobuf.ChangeUsernameResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangePassword changes the user's password
func (u usersServiceServer) ChangePassword(ctx context.Context, request *protobuf.ChangePasswordRequest) (*protobuf.ChangePasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangeEmail changes the user's email
func (u usersServiceServer) ChangeEmail(ctx context.Context, request *protobuf.ChangeEmailRequest) (*protobuf.ChangeEmailResponse, error) {
	//TODO implement me
	panic("implement me")
}

// VerifyEmail verifies the user's email
func (u usersServiceServer) VerifyEmail(ctx context.Context, request *protobuf.VerifyEmailRequest) (*protobuf.VerifyEmailResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangePhoneNumber changes the user's phone number
func (u usersServiceServer) ChangePhoneNumber(ctx context.Context, request *protobuf.ChangePhoneNumberRequest) (*protobuf.ChangePhoneNumberResponse, error) {
	//TODO implement me
	panic("implement me")
}

// VerifyPhoneNumber verifies the user's phone number
func (u usersServiceServer) VerifyPhoneNumber(ctx context.Context, request *protobuf.VerifyPhoneNumberRequest) (*protobuf.VerifyPhoneNumberResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ForgotPassword sends a password reset link to the user's email
func (u usersServiceServer) ForgotPassword(ctx context.Context, request *protobuf.ForgotPasswordRequest) (*protobuf.ForgotPasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ResetPassword resets the user's password
func (u usersServiceServer) ResetPassword(ctx context.Context, request *protobuf.ResetPasswordRequest) (*protobuf.ResetPasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

// DeleteUser deletes the user's account
func (u usersServiceServer) DeleteUser(ctx context.Context, request *protobuf.DeleteUserRequest) (*protobuf.DeleteUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (usersServiceServer) mustEmbedUnimplementedUsersServiceServer() {}

func main() {
	// Listen on the given port
	port, listenerPort := config.LoadUsersServicePort()
	listener, err := net.Listen("tcp", listenerPort)
	if err != nil {
		logger.ListenerLogger.FailedToListen(err)
	}

	// Create a new gRPC server
	s := grpc.NewServer()

	// Register the server with the gRPC server
	protobuf.RegisterUsersServiceServer(s, &usersServiceServer{})

	// Register reflection service on gRPC server.
	logger.ListenerLogger.ServerStarted(port)

	// Serve the gRPC server
	if err := s.Serve(listener); err != nil {
		logger.ListenerLogger.FailedToServe(err)
	}
}
