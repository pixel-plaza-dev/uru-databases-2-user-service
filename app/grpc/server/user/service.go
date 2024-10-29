package user

import (
	commonmessage "github.com/pixel-plaza-dev/uru-databases-2-api-common/message"
	commonbcrypt "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto/bcrypt"
	commoncryptoerror "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto/error"
	commonuser "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb/database/user"
	commonvalidator "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/validator"
	commonvalidatorresponse "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/validator/error/response"
	protobuf "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/compiled-protobuf/user"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb/database/user"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/validator/error_response"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Server struct {
	userDatabase *user.Database
	logger       *user.Logger
	protobuf.UnimplementedUserServer
}

// NewServer creates a new gRPC user server
func NewServer(userDatabase *user.Database, logger *user.Logger) *Server {
	return &Server{userDatabase: userDatabase, logger: logger}
}

// SignUp creates a new user
func (u Server) SignUp(ctx context.Context, request *protobuf.SignUpRequest) (response *protobuf.SignUpResponse, err error) {
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
	commonvalidator.ValidNonEmptyStringFields(&validations, request, &fieldsToValidate)

	// Check if the user already exists
	username := request.GetUsername()
	if len(username) > 0 {
		if _, err := u.userDatabase.FindUserByUsername(username); err == nil {
			userExists = true
			validations["username"] = append(validations["username"], error_response.UsernameTakenError{})
		}
	}

	// Check if the email is valid
	email := request.GetEmail()
	if len(email) > 0 {
		if _, err = commonvalidator.ValidMailAddress(email); err != nil {
			field := "email"
			validations[field] = append(validations[field], commonvalidatorresponse.InvalidMailAddressError{})
		}
	}

	// Check if the password is hashed
	hashedPassword := request.GetHashedPassword()
	if isHashed := commonbcrypt.IsHashed(hashedPassword); !isHashed {
		field := "hashed_password"
		validations[field] = append(validations[field], commoncryptoerror.PasswordNotHashedError{})
	}

	// Check if the birthdate is valid
	birthDateTimestamp := request.GetBirthDate()
	birthDate := birthDateTimestamp.AsTime()
	currentTime := time.Now()
	if birthDateTimestamp == nil || birthDate.After(currentTime) {
		field := "birth_date"
		validations[field] = append(validations[field], error_response.InvalidBirthDateError{BirthDate: birthDate})
	}

	// Check if there are any validation errors
	if len(validations) > 0 {
		err = commonvalidatorresponse.FailedValidationError{FieldsErrors: &validations}

		if userExists {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Create a new user
	userId := primitive.NewObjectID()
	newUser := commonuser.User{
		ID:             userId,
		Username:       username,
		FirstName:      request.GetFirstName(),
		LastName:       request.GetLastName(),
		HashedPassword: hashedPassword,
		BirthDate:      birthDate,
	}

	// Create the user email
	userEmailId := primitive.NewObjectID()
	newUserEmail := commonuser.UserEmail{
		ID:         userEmailId,
		UserID:     userId,
		Email:      email,
		AssignedAt: currentTime,
		IsActive:   true,
	}

	// Create the user phone number
	userPhoneNumberId := primitive.NewObjectID()
	newUserPhoneNumber := commonuser.UserPhoneNumber{
		ID:          userPhoneNumberId,
		UserID:      userId,
		PhoneNumber: request.GetPhoneNumber(),
		AssignedAt:  currentTime,
		IsActive:    true,
	}

	// Insert the user into the user
	if _, err = u.userDatabase.CreateUser(&newUser, &newUserEmail, &newUserPhoneNumber); err != nil {
		// Log the error
		u.logger.FailedToCreateDocument(err)

		return nil, status.Error(codes.Internal, commonmessage.Internal)
	}

	// Log the success
	u.logger.UserCreated(userId.Hex())

	return &protobuf.SignUpResponse{
		Code:    uint32(codes.OK),
		Message: commonmessage.SignUpSuccess,
	}, nil
}

// UpdateProfile updates the user's profile
func (u Server) UpdateProfile(ctx context.Context, request *protobuf.UpdateProfileRequest) (*protobuf.UpdateProfileResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangeUsername changes the user's username
func (u Server) ChangeUsername(ctx context.Context, request *protobuf.ChangeUsernameRequest) (*protobuf.ChangeUsernameResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangePassword changes the user's password
func (u Server) ChangePassword(ctx context.Context, request *protobuf.ChangePasswordRequest) (*protobuf.ChangePasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangeEmail changes the user's email
func (u Server) ChangeEmail(ctx context.Context, request *protobuf.ChangeEmailRequest) (*protobuf.ChangeEmailResponse, error) {
	//TODO implement me
	panic("implement me")
}

// VerifyEmail verifies the user's email
func (u Server) VerifyEmail(ctx context.Context, request *protobuf.VerifyEmailRequest) (*protobuf.VerifyEmailResponse, error) {
	//TODO implement me
	panic("implement me")
}

// GetActiveEmails gets the user's active emails
func (u Server) GetActiveEmails(ctx context.Context, request *protobuf.GetActiveEmailsRequest) (*protobuf.GetActiveEmailsResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangePhoneNumber changes the user's phone number
func (u Server) ChangePhoneNumber(ctx context.Context, request *protobuf.ChangePhoneNumberRequest) (*protobuf.ChangePhoneNumberResponse, error) {
	//TODO implement me
	panic("implement me")
}

// VerifyPhoneNumber verifies the user's phone number
func (u Server) VerifyPhoneNumber(ctx context.Context, request *protobuf.VerifyPhoneNumberRequest) (*protobuf.VerifyPhoneNumberResponse, error) {
	//TODO implement me
	panic("implement me")
}

// GetActivePhoneNumber gets the user's active phone number
func (u Server) GetActivePhoneNumbers(ctx context.Context, request *protobuf.GetActivePhoneNumbersRequest) (*protobuf.GetActivePhoneNumbersResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ForgotPassword sends a password reset link to the user's email
func (u Server) ForgotPassword(ctx context.Context, request *protobuf.ForgotPasswordRequest) (*protobuf.ForgotPasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ResetPassword resets the user's password
func (u Server) ResetPassword(ctx context.Context, request *protobuf.ResetPasswordRequest) (*protobuf.ResetPasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

// DeleteUser deletes the user's account
func (u Server) DeleteUser(ctx context.Context, request *protobuf.DeleteUserRequest) (*protobuf.DeleteUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (Server) mustEmbedUnimplementedServer() {}
