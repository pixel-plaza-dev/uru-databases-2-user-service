package user

import (
	commonmessage "github.com/pixel-plaza-dev/uru-databases-2-api-common/message"
	commoncrypto "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto"
	commonbcrypt "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto/bcrypt"
	commonuser "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb/database/user"
	commonvalidator "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/validator"
	commonvalidatorerror "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/validator/error"
	protobuf "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/compiled-protobuf/user"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb/database/user"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// Server is the gRPC user server
type Server struct {
	userDatabase *user.Database
	logger       user.Logger
	protobuf.UnimplementedUserServer
}

// NewServer creates a new gRPC user server
func NewServer(userDatabase *user.Database, logger user.Logger) *Server {
	return &Server{userDatabase: userDatabase, logger: logger}
}

// SignUp creates a new user
func (s Server) SignUp(ctx context.Context, request *protobuf.SignUpRequest) (response *protobuf.SignUpResponse, err error) {
	// Validation variables
	validations := make(map[string][]error)
	userExists := false

	// Get the request fields
	usernameField := "username"
	emailField := "email"
	hashedPasswordField := "hashed_password"
	birthDateField := "birth_date"

	fieldsToValidate := map[string]string{
		"Username":       usernameField,
		"FirstName":      "first_name",
		"LastName":       "last_name",
		"HashedPassword": hashedPasswordField,
		"Email":          emailField,
		"PhoneNumber":    "phone_number",
	}

	// Check if the required string fields are empty
	commonvalidator.ValidNonEmptyStringFields(&validations, request, &fieldsToValidate)

	// Check if the user already exists
	username := request.GetUsername()
	if len(username) > 0 {
		if _, err := s.userDatabase.FindUserByUsername(username); err == nil {
			userExists = true
			validations[usernameField] = append(validations[usernameField], validator.UsernameTakenError)
		}
	}

	// Check if the email is valid
	email := request.GetEmail()
	if len(email) > 0 {
		if _, err = commonvalidator.ValidMailAddress(email); err != nil {
			validations[emailField] = append(validations[emailField], commonvalidator.InvalidMailAddressError)
		}
	}

	// Check if the password is hashed
	hashedPassword := request.GetHashedPassword()
	if isHashed := commonbcrypt.IsHashed(hashedPassword); !isHashed {
		validations[hashedPasswordField] = append(validations[hashedPasswordField], commoncrypto.PasswordNotHashedError)
	}

	// Check if the birthdate is valid
	birthDateTimestamp := request.GetBirthDate()
	birthDate := birthDateTimestamp.AsTime()
	currentTime := time.Now()
	if birthDateTimestamp == nil || birthDate.After(currentTime) {
		validations[birthDateField] = append(validations[birthDateField], validator.InvalidBirthDateError)
	}

	// Check if there are any validation errors
	if len(validations) > 0 {
		err = commonvalidatorerror.FailedValidationError{FieldsErrors: &validations}

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
	if _, err = s.userDatabase.CreateUser(&newUser, &newUserEmail, &newUserPhoneNumber); err != nil {
		// Log the error
		s.logger.FailedToCreateDocument(err)

		return nil, InternalError
	}

	// Log the success
	s.logger.UserCreated(userId.Hex())

	return &protobuf.SignUpResponse{
		Code:    uint32(codes.OK),
		Message: commonmessage.SignUpSuccess,
	}, nil
}

// UpdateProfile updates the user's profile
func (s Server) UpdateProfile(ctx context.Context, request *protobuf.UpdateProfileRequest) (*protobuf.UpdateProfileResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangeUsername changes the user's username
func (s Server) ChangeUsername(ctx context.Context, request *protobuf.ChangeUsernameRequest) (*protobuf.ChangeUsernameResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangePassword changes the user's password
func (s Server) ChangePassword(ctx context.Context, request *protobuf.ChangePasswordRequest) (*protobuf.ChangePasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangeEmail changes the user's email
func (s Server) ChangeEmail(ctx context.Context, request *protobuf.ChangeEmailRequest) (*protobuf.ChangeEmailResponse, error) {
	//TODO implement me
	panic("implement me")
}

// VerifyEmail verifies the user's email
func (s Server) VerifyEmail(ctx context.Context, request *protobuf.VerifyEmailRequest) (*protobuf.VerifyEmailResponse, error) {
	//TODO implement me
	panic("implement me")
}

// GetActiveEmails gets the user's active emails
func (s Server) GetActiveEmails(ctx context.Context, request *protobuf.GetActiveEmailsRequest) (*protobuf.GetActiveEmailsResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ChangePhoneNumber changes the user's phone number
func (s Server) ChangePhoneNumber(ctx context.Context, request *protobuf.ChangePhoneNumberRequest) (*protobuf.ChangePhoneNumberResponse, error) {
	//TODO implement me
	panic("implement me")
}

// VerifyPhoneNumber verifies the user's phone number
func (s Server) VerifyPhoneNumber(ctx context.Context, request *protobuf.VerifyPhoneNumberRequest) (*protobuf.VerifyPhoneNumberResponse, error) {
	//TODO implement me
	panic("implement me")
}

// GetActivePhoneNumbers gets the user's active phone numbers
func (s Server) GetActivePhoneNumbers(ctx context.Context, request *protobuf.GetActivePhoneNumbersRequest) (*protobuf.GetActivePhoneNumbersResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ForgotPassword sends a password reset link to the user's email
func (s Server) ForgotPassword(ctx context.Context, request *protobuf.ForgotPasswordRequest) (*protobuf.ForgotPasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

// ResetPassword resets the user's password
func (s Server) ResetPassword(ctx context.Context, request *protobuf.ResetPasswordRequest) (*protobuf.ResetPasswordResponse, error) {
	//TODO implement me
	panic("implement me")
}

// DeleteUser deletes the user's account
func (s Server) DeleteUser(ctx context.Context, request *protobuf.DeleteUserRequest) (*protobuf.DeleteUserResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (Server) mustEmbedUnimplementedServer() {}
