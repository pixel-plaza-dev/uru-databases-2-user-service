package grpc_server

import (
	commonMongoDb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb"
	"golang.org/x/net/context"
	protobuf "pixel_plaza/users_service/protobuf/pixel_plaza/users_service"
)

type UsersServiceServer struct {
	mongodbConnection *commonMongoDb.Connection
	protobuf.UnimplementedUsersServiceServer
}

// NewUsersServiceServer creates a new users service server
func NewUsersServiceServer(mongoDbConnection *commonMongoDb.Connection) *UsersServiceServer {
	// Create a new logger and MongoDB service
	return &UsersServiceServer{mongodbConnection: mongoDbConnection}
}

// SignUp creates a new user
func (u UsersServiceServer) SignUp(ctx context.Context, request *protobuf.SignUpRequest) (*protobuf.SignUpResponse, error) {

	panic("implement me")
}

// UpdateProfile updates the user's profile
func (u UsersServiceServer) UpdateProfile(ctx context.Context, request *protobuf.UpdateProfileRequest) (*protobuf.UpdateProfileResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangeUsername changes the user's username
func (u UsersServiceServer) ChangeUsername(ctx context.Context, request *protobuf.ChangeUsernameRequest) (*protobuf.ChangeUsernameResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangePassword changes the user's password
func (u UsersServiceServer) ChangePassword(ctx context.Context, request *protobuf.ChangePasswordRequest) (*protobuf.ChangePasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangeEmail changes the user's email
func (u UsersServiceServer) ChangeEmail(ctx context.Context, request *protobuf.ChangeEmailRequest) (*protobuf.ChangeEmailResponse, error) {
	//TODO implement me
	panic("implement me")
}

// VerifyEmail verifies the user's email
func (u UsersServiceServer) VerifyEmail(ctx context.Context, request *protobuf.VerifyEmailRequest) (*protobuf.VerifyEmailResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangePhoneNumber changes the user's phone number
func (u UsersServiceServer) ChangePhoneNumber(ctx context.Context, request *protobuf.ChangePhoneNumberRequest) (*protobuf.ChangePhoneNumberResponse, error) {
	//TODO implement me
	panic("implement me")
}

// VerifyPhoneNumber verifies the user's phone number
func (u UsersServiceServer) VerifyPhoneNumber(ctx context.Context, request *protobuf.VerifyPhoneNumberRequest) (*protobuf.VerifyPhoneNumberResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ForgotPassword sends a password reset link to the user's email
func (u UsersServiceServer) ForgotPassword(ctx context.Context, request *protobuf.ForgotPasswordRequest) (*protobuf.ForgotPasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ResetPassword resets the user's password
func (u UsersServiceServer) ResetPassword(ctx context.Context, request *protobuf.ResetPasswordRequest) (*protobuf.ResetPasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

// DeleteUser deletes the user's account
func (u UsersServiceServer) DeleteUser(ctx context.Context, request *protobuf.DeleteUserRequest) (*protobuf.DeleteUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (UsersServiceServer) mustEmbedUnimplementedUsersServiceServer() {}
