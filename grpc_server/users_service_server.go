package grpc_server

import (
	commonBcrypt "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto/bcrypt"
	commonBcryptError "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/custom_error/crypto/bcrypt"
	commonValidatorErrorResponse "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/custom_error_response/validator"
	commonValidator "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	validatorErrorResponse "pixel_plaza/users_service/custom_error_response/validator"
	"pixel_plaza/users_service/logger"
	"pixel_plaza/users_service/mongodb"
	protobuf "pixel_plaza/users_service/protobuf/pixel_plaza/users_service"
	"time"
)

const (
	// Internal is the message for internal server error
	Internal = "Internal server error"

	// SignUpFailedMessage is the message for failed sign up
	SignUpFailedMessage = "Failed to sign up"

	// SignUpSuccessMessage is the message for successful sign up
	SignUpSuccessMessage = "Successfully signed up"
)

type UsersServiceServer struct {
	userDatabase *mongodb.UserDatabase
	logger       *logger.UsersServiceLogger
	protobuf.UnimplementedUsersServiceServer
}

// NewUsersServiceServer creates a new users service server
func NewUsersServiceServer(userDatabase *mongodb.UserDatabase, logger *logger.UsersServiceLogger) *UsersServiceServer {
	return &UsersServiceServer{userDatabase: userDatabase, logger: logger}
}

// SignUp creates a new user
func (u UsersServiceServer) SignUp(ctx context.Context, request *protobuf.SignUpRequest) (response *protobuf.SignUpResponse, err error) {
	// Validation variables
	validations := make(map[string][]error)
	userExists := false

	// Get the request fields
	fieldsToValidate := map[string]string{
		"Username":       "username",
		"FirstName":      "first_name",
		"LastName":       "last_name",
		"HashedPassword": "hashed_password",
		"Email":          "email",
		"PhoneNumber":    "phone_number",
	}

	// Check if the required string fields are empty
	commonValidator.ValidNonEmptyStringFields(&validations, request, &fieldsToValidate)

	// Check if the user already exists
	username := request.GetUsername()
	if len(username) > 0 {
		if _, err := u.userDatabase.FindUserByUsername(username); err == nil {
			userExists = true
			validations["username"] = append(validations["username"], validatorErrorResponse.UsernameTakenError{})
		}
	}

	// Check if the email is valid
	email := request.GetEmail()
	if len(email) > 0 {
		if _, err = commonValidator.ValidMailAddress(email); err != nil {
			field := "email"
			validations[field] = append(validations[field], commonValidatorErrorResponse.InvalidMailAddressError{})
		}
	}

	// Check if the password is hashed
	hashedPassword := request.GetHashedPassword()
	if isHashed := commonBcrypt.IsHashed(hashedPassword); !isHashed {
		field := "hashed_password"
		validations[field] = append(validations[field], commonBcryptError.PasswordNotHashedError{})
	}

	// Check if the birthdate is valid
	birthDateTimestamp := request.GetBirthDate()
	birthDate := birthDateTimestamp.AsTime()
	currentTime := time.Now()
	if birthDateTimestamp == nil || birthDate.After(currentTime) {
		field := "birth_date"
		validations[field] = append(validations[field], validatorErrorResponse.InvalidBirthDateError{BirthDate: birthDate})
	}

	// Check if there are any validation errors
	if len(validations) > 0 {
		err = commonValidatorErrorResponse.FailedValidationError{FieldsErrors: &validations}

		if userExists {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Create a new user
	userId := primitive.NewObjectID()
	user := mongodb.User{
		ID:             userId,
		Username:       username,
		FirstName:      request.GetFirstName(),
		LastName:       request.GetLastName(),
		HashedPassword: hashedPassword,
		BirthDate:      birthDate,
	}

	// Create the user email
	userEmailId := primitive.NewObjectID()
	userEmail := mongodb.UserEmail{
		ID:         userEmailId,
		UserID:     userId,
		Email:      email,
		AssignedAt: currentTime,
		IsActive:   true,
	}

	// Create the user phone number
	userPhoneNumberId := primitive.NewObjectID()
	userPhoneNumber := mongodb.UserPhoneNumber{
		ID:          userPhoneNumberId,
		UserID:      userId,
		PhoneNumber: request.GetPhoneNumber(),
		AssignedAt:  currentTime,
		IsActive:    true,
	}

	// Insert the user into the database
	if _, err = u.userDatabase.CreateUser(&user, &userEmail, &userPhoneNumber); err != nil {
		// Log the error
		u.logger.FailedToCreateDocument(err)

		return nil, status.Error(codes.Internal, SignUpFailedMessage)
	}

	// Log the success
	u.logger.UserCreated(userId.Hex())

	return &protobuf.SignUpResponse{
		Code:    uint32(codes.OK),
		Message: SignUpSuccessMessage,
	}, nil
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
