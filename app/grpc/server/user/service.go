package user

import (
	commonbcrypt "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto/bcrypt"
	commonuser "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/database/mongodb/model/user"
	commonvalidator "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/utils/validator"
	commonvalidatorerror "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/utils/validator/error"
	pbauth "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/protobuf/compiled/auth"
	protobuf "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/protobuf/compiled/user"
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
	authClient   pbauth.AuthClient
	logger       Logger
	protobuf.UnimplementedUserServer
}

// NewServer creates a new gRPC user server
func NewServer(
	userDatabase *user.Database,
	authClient pbauth.AuthClient,
	logger Logger,
) *Server {
	return &Server{
		userDatabase: userDatabase,
		authClient:   authClient,
		logger:       logger,
	}
}

// SignUp creates a new user
func (s Server) SignUp(
	ctx context.Context,
	request *protobuf.SignUpRequest,
) (response *protobuf.SignUpResponse, err error) {
	// Validation variables
	validations := make(map[string][]error)
	userExists := false

	// Get the request fields
	usernameField := "username"
	emailField := "email"
	birthDateField := "birth_date"

	requestFieldsToValidate := map[string]string{
		"Password":    "password",
		"Email":       emailField,
		"PhoneNumber": "phone_number",
	}

	profileFieldsToValidate := map[string]string{
		"Username":  usernameField,
		"FirstName": "first_name",
		"LastName":  "last_name",
	}

	// Check if the required string fields are empty
	commonvalidator.ValidNonEmptyStringFields(
		&validations,
		request,
		&requestFieldsToValidate,
	)
	commonvalidator.ValidNonEmptyStringFields(
		&validations,
		request.GetProfile(),
		&profileFieldsToValidate,
	)

	// Check if the user already exists
	username := request.GetProfile().GetUsername()
	if len(username) > 0 {
		if _, err := s.userDatabase.FindUserByUsername(
			username,
			nil,
		); err == nil {
			userExists = true
			validations[usernameField] = append(
				validations[usernameField],
				validator.UsernameTakenError,
			)
		}
	}

	// Check if the email is valid
	email := request.GetEmail()
	if len(email) > 0 {
		if _, err = commonvalidator.ValidMailAddress(email); err != nil {
			validations[emailField] = append(
				validations[emailField],
				commonvalidator.InvalidMailAddressError,
			)
		}
	}

	// Check if the birthdate is valid
	birthDateTimestamp := request.GetProfile().GetBirthDate()
	birthDate := birthDateTimestamp.AsTime()
	currentTime := time.Now()
	if birthDateTimestamp == nil || birthDate.After(currentTime) {
		validations[birthDateField] = append(
			validations[birthDateField],
			validator.InvalidBirthDateError,
		)
	}

	// Check if there are any validation errors
	if len(validations) > 0 {
		err = commonvalidatorerror.FailedValidationError{FieldsErrors: &validations}

		if userExists {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Hash the password
	hashedPassword, err := commonbcrypt.HashPassword(request.GetPassword())
	if err != nil {
		s.logger.FailedToHashPassword(err)
		return nil, InternalServerError
	}

	// Create a new user
	userId := primitive.NewObjectID()
	newUser := commonuser.User{
		ID:             userId,
		Username:       username,
		FirstName:      request.GetProfile().GetFirstName(),
		LastName:       request.GetProfile().GetLastName(),
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
	if _, err = s.userDatabase.CreateUser(
		&newUser,
		&newUserEmail,
		&newUserPhoneNumber,
	); err != nil {
		return nil, InternalServerError
	}

	// Log the success
	s.logger.UserSignedUp(userId.Hex())

	return &protobuf.SignUpResponse{
		Code:    uint32(codes.OK),
		Message: SignUpSuccess,
	}, nil
}

func (s Server) IsPasswordCorrect(
	ctx context.Context,
	request *protobuf.IsPasswordCorrectRequest,
) (*protobuf.IsPasswordCorrectResponse, error) {
	// Validation variables
	validations := make(map[string][]error)

	// Get the request fields
	fieldsToValidate := map[string]string{
		"Username": "username",
		"Password": "password",
	}

	// Check if the required string fields are empty
	commonvalidator.ValidNonEmptyStringFields(
		&validations,
		request,
		&fieldsToValidate,
	)

	// Check if there are any validation errors
	if len(validations) > 0 {
		err := commonvalidatorerror.FailedValidationError{FieldsErrors: &validations}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Validate the password and get the user ID
	userIdentifier, err := s.userDatabase.IsPasswordCorrect(
		request.GetUsername(),
		request.GetPassword(),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Log the password check success
	s.logger.PasswordCheckSuccess(userIdentifier)

	return &protobuf.IsPasswordCorrectResponse{
		Code:    uint32(codes.OK),
		Message: IsPasswordCorrectSuccess,
		UserId:  userIdentifier,
	}, nil
}

// UpdateProfile updates the user's profile
func (s Server) UpdateProfile(
	ctx context.Context,
	request *protobuf.UpdateProfileRequest,
) (*protobuf.UpdateProfileResponse, error) {
	return nil, InDevelopmentError
}

// GetProfile gets the user's profile
func (s Server) GetProfile(
	ctx context.Context,
	request *protobuf.GetProfileRequest,
) (*protobuf.GetProfileResponse, error) {
	return nil, InDevelopmentError
}

// GetFullProfile gets the user's full profile
func (s Server) GetFullProfile(
	ctx context.Context,
	request *protobuf.GetFullProfileRequest,
) (*protobuf.GetFullProfileResponse, error) {
	return nil, InDevelopmentError
}

// GetUserIdByUsername gets the user's ID by username
func (s Server) GetUserIdByUsername(
	ctx context.Context,
	request *protobuf.GetUserIdByUsernameRequest,
) (*protobuf.GetUserIdByUsernameResponse, error) {
	return nil, InDevelopmentError
}

// UsernameExists checks if the username exists
func (s Server) UsernameExists(
	ctx context.Context,
	request *protobuf.UsernameExistsRequest,
) (*protobuf.UsernameExistsResponse, error) {
	return nil, InDevelopmentError
}

// GetUsernameByUserId gets the user's username by ID
func (s Server) GetUsernameByUserId(
	ctx context.Context,
	request *protobuf.GetUsernameByUserIdRequest,
) (*protobuf.GetUsernameByUserIdResponse, error) {
	return nil, InDevelopmentError
}

// ChangeUsername changes the user's username
func (s Server) ChangeUsername(
	ctx context.Context,
	request *protobuf.ChangeUsernameRequest,
) (*protobuf.ChangeUsernameResponse, error) {
	return nil, InDevelopmentError
}

// ChangePassword changes the user's password
func (s Server) ChangePassword(
	ctx context.Context,
	request *protobuf.ChangePasswordRequest,
) (*protobuf.ChangePasswordResponse, error) {
	return nil, InDevelopmentError
}

// AddEmail adds an email to the user's account
func (s Server) AddEmail(
	ctx context.Context,
	request *protobuf.AddEmailRequest,
) (*protobuf.AddEmailResponse, error) {
	return nil, InDevelopmentError
}

// DeleteEmail deletes an email from the user's account
func (s Server) DeleteEmail(
	ctx context.Context,
	request *protobuf.DeleteEmailRequest,
) (*protobuf.DeleteEmailResponse, error) {
	return nil, InDevelopmentError
}

// SendVerificationEmail sends a verification email to the user
func (s Server) SendVerificationEmail(
	ctx context.Context,
	request *protobuf.SendVerificationEmailRequest,
) (*protobuf.SendVerificationEmailResponse, error) {
	return nil, InDevelopmentError
}

// GetPrimaryEmail gets the user's primary email
func (s Server) GetPrimaryEmail(
	ctx context.Context,
	request *protobuf.GetPrimaryEmailRequest,
) (*protobuf.GetPrimaryEmailResponse, error) {
	return nil, InDevelopmentError
}

// ChangePrimaryEmail changes the user's primary email
func (s Server) ChangePrimaryEmail(
	ctx context.Context,
	request *protobuf.ChangePrimaryEmailRequest,
) (*protobuf.ChangePrimaryEmailResponse, error) {
	return nil, InDevelopmentError
}

// VerifyEmail verifies the user's email
func (s Server) VerifyEmail(
	ctx context.Context,
	request *protobuf.VerifyEmailRequest,
) (*protobuf.VerifyEmailResponse, error) {
	return nil, InDevelopmentError
}

// GetActiveEmails gets the user's active emails
func (s Server) GetActiveEmails(
	ctx context.Context,
	request *protobuf.GetActiveEmailsRequest,
) (*protobuf.GetActiveEmailsResponse, error) {
	return nil, InDevelopmentError
}

// ChangePhoneNumber changes the user's phone number
func (s Server) ChangePhoneNumber(
	ctx context.Context,
	request *protobuf.ChangePhoneNumberRequest,
) (*protobuf.ChangePhoneNumberResponse, error) {
	return nil, InDevelopmentError
}

// GetPhoneNumber gets the user's phone number
func (s Server) GetPhoneNumber(
	ctx context.Context,
	request *protobuf.GetPhoneNumberRequest,
) (*protobuf.GetPhoneNumberResponse, error) {
	return nil, InDevelopmentError
}

// VerifyPhoneNumber verifies the user's phone number
func (s Server) VerifyPhoneNumber(
	ctx context.Context,
	request *protobuf.VerifyPhoneNumberRequest,
) (*protobuf.VerifyPhoneNumberResponse, error) {
	return nil, InDevelopmentError
}

// SendVerificationSMS sends a verification SMS to the user
func (s Server) SendVerificationSMS(
	ctx context.Context,
	request *protobuf.SendVerificationSMSRequest,
) (*protobuf.SendVerificationSMSResponse, error) {
	return nil, InDevelopmentError
}

// ForgotPassword sends a password reset link to the user's email
func (s Server) ForgotPassword(
	ctx context.Context,
	request *protobuf.ForgotPasswordRequest,
) (*protobuf.ForgotPasswordResponse, error) {
	return nil, InDevelopmentError
}

// ResetPassword resets the user's password
func (s Server) ResetPassword(
	ctx context.Context,
	request *protobuf.ResetPasswordRequest,
) (*protobuf.ResetPasswordResponse, error) {
	return nil, InDevelopmentError
}

// DeleteUser deletes the user's account
func (s Server) DeleteUser(
	ctx context.Context,
	request *protobuf.DeleteUserRequest,
) (*protobuf.DeleteUserResponse, error) {
	return nil, InDevelopmentError
}
