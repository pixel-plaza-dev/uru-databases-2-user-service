package user

import (
	"errors"
	commonbcrypt "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto/bcrypt"
	commonuser "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/database/mongodb/model/user"
	commongrpcclientctx "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/grpc/client/context"
	commongrpcserverctx "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/grpc/server/context"
	pbauth "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/protobuf/compiled/auth"
	protobuf "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/protobuf/compiled/user"
	userservervalidator "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/grpc/server/user/validator"
	mongodbuser "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb/database/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// Server is the gRPC user server
type Server struct {
	userDatabase *mongodbuser.Database
	authClient   pbauth.AuthClient
	logger       Logger
	validator    *userservervalidator.Validator
	protobuf.UnimplementedUserServer
}

// NewServer creates a new gRPC user server
func NewServer(
	userDatabase *mongodbuser.Database,
	authClient pbauth.AuthClient,
	logger Logger,
	validator *userservervalidator.Validator,
) *Server {
	return &Server{
		userDatabase: userDatabase,
		authClient:   authClient,
		logger:       logger,
		validator:    validator,
	}
}

// SignUp creates a new user
func (s Server) SignUp(
	ctx context.Context,
	request *protobuf.SignUpRequest,
) (response *protobuf.SignUpResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateSignUpRequest(request); err != nil {
		s.logger.SignUpFailed(err)
		return nil, err
	}

	// Hash the password
	hashedPassword, err := commonbcrypt.HashPassword(request.GetPassword())
	if err != nil {
		s.logger.HashPasswordFailed(err)
		return nil, InternalServerError
	}

	// Create a new user
	userId := primitive.NewObjectID()
	newUser := commonuser.User{
		ID:             userId,
		Username:       request.GetUsername(),
		FirstName:      request.GetFirstName(),
		LastName:       request.GetLastName(),
		HashedPassword: hashedPassword,
	}

	// Add the birthdate if it exists
	if request.GetBirthdate() != nil {
		newUser.Birthdate = request.GetBirthdate().AsTime()
	}

	// Create the user email
	currentTime := time.Now()
	userEmailId := primitive.NewObjectID()
	newUserEmail := commonuser.UserEmail{
		ID:         userEmailId,
		UserID:     userId,
		Email:      request.GetEmail(),
		AssignedAt: currentTime,
	}

	// Create the user phone number
	userPhoneNumberId := primitive.NewObjectID()
	newUserPhoneNumber := commonuser.UserPhoneNumber{
		ID:          userPhoneNumberId,
		UserID:      userId,
		PhoneNumber: request.GetPhoneNumber(),
		AssignedAt:  currentTime,
	}

	// Insert the user into the user
	if err = s.userDatabase.CreateUser(
		&newUser,
		&newUserEmail,
		&newUserPhoneNumber,
	); err != nil {
		s.logger.SignUpFailed(err)
		return nil, InternalServerError
	}

	// User signed up successfully
	s.logger.SignedUp(userId.Hex())

	return &protobuf.SignUpResponse{
		Message: SignUpSuccess,
	}, nil
}

func (s Server) IsPasswordCorrect(
	ctx context.Context,
	request *protobuf.IsPasswordCorrectRequest,
) (response *protobuf.IsPasswordCorrectResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateIsPasswordCorrectRequest(request); err != nil {
		s.logger.PasswordIsCorrectFailed(err)
		return nil, err
	}

	// Get the user ID and hashed password by username
	user, err := s.userDatabase.GetUserHashedPassword(
		context.Background(), request.GetUsername(),
	)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		s.logger.PasswordIsCorrectFailed(err)
		return nil, InternalServerError
	}

	// Check if the password matches
	matches := commonbcrypt.CheckPasswordHash(user.HashedPassword, request.GetPassword())

	// Get the user ID
	userId := user.ID.Hex()

	// Check if the password doesn't match or the user doesn't exist
	if err != nil || !matches {
		// User checked password unsuccessfully
		s.logger.PasswordIsIncorrect(userId)

		return nil, status.Error(codes.OK, IsPasswordCorrectFailure)
	}

	// User checked password successfully
	s.logger.PasswordIsCorrect(userId)

	return &protobuf.IsPasswordCorrectResponse{
		Message: IsPasswordCorrectSuccess,
		UserId:  userId,
	}, nil
}

// UsernameExists checks if the username exists
func (s Server) UsernameExists(
	ctx context.Context,
	request *protobuf.UsernameExistsRequest,
) (response *protobuf.UsernameExistsResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateUsernameExistsRequest(request); err != nil {
		s.logger.UsernameExistsFailed(err)
		return nil, err
	}

	// Check if the username exists
	exists, err := s.userDatabase.UsernameExists(context.Background(), request.GetUsername())
	if err != nil {
		s.logger.UsernameExistsFailed(err)
		return nil, InternalServerError
	}

	// Check if the username doesn't exist
	if !exists {
		// Username does not exist
		s.logger.UserNotFoundByUsername(request.GetUsername())

		return nil, status.Error(codes.NotFound, FoundByUsernameFailure)
	}

	// User found by username
	s.logger.UserFoundByUsername(request.GetUsername())

	return &protobuf.UsernameExistsResponse{
		Message: FoundByUsernameSuccess,
		Exists:  true,
	}, nil
}

// GetUserIdByUsername gets the user's ID by username
func (s Server) GetUserIdByUsername(
	ctx context.Context,
	request *protobuf.GetUserIdByUsernameRequest,
) (response *protobuf.GetUserIdByUsernameResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateGetUserIdByUsernameRequest(request); err != nil {
		s.logger.GetUserIdByUsernameFailed(err)
		return nil, err
	}

	// Get the user ID by username
	userId, err := s.userDatabase.GetUserIdByUsername(context.Background(), request.GetUsername())
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.GetUserIdByUsernameFailed(err)
		return nil, InternalServerError
	}

	// Check if the username doesn't exist
	if userId == "" {
		// Username does not exist
		s.logger.UserNotFoundByUsername(request.GetUsername())

		return nil, status.Error(codes.NotFound, FoundByUsernameFailure)
	}

	// User found by username
	s.logger.UserFoundByUsername(request.GetUsername())

	return &protobuf.GetUserIdByUsernameResponse{
		Message: FoundByUsernameSuccess,
		UserId:  userId,
	}, nil
}

// GetUsernameByUserId gets the user's username by ID
func (s Server) GetUsernameByUserId(
	ctx context.Context,
	request *protobuf.GetUsernameByUserIdRequest,
) (response *protobuf.GetUsernameByUserIdResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateGetUsernameByUserIdRequest(request); err != nil {
		s.logger.GetUsernameByUserIdFailed(err)
		return nil, err
	}

	// Get the username by user ID
	username, err := s.userDatabase.GetUsernameByUserId(context.Background(), request.GetUserId())
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.GetUsernameByUserIdFailed(err)
		return nil, InternalServerError
	}

	// Check if the user ID doesn't exist
	if username == "" {
		// User ID does not exist
		s.logger.UserNotFoundByUserId(request.GetUserId())

		return nil, status.Error(codes.NotFound, FoundByUserIdFailure)
	}

	// User found by user ID
	s.logger.UserFoundByUsername(request.GetUserId())

	return &protobuf.GetUsernameByUserIdResponse{
		Message:  FoundByUserIdSuccess,
		Username: username,
	}, nil
}

// GetProfile gets the user's profile
func (s Server) GetProfile(
	ctx context.Context,
	request *protobuf.GetProfileRequest,
) (response *protobuf.GetProfileResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateGetProfileRequest(request); err != nil {
		s.logger.GetProfileFailed(err)
		return nil, err
	}

	// Get the profile by user ID
	profile, err := s.userDatabase.GetProfile(context.Background(), request.GetUserId())
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.GetProfileFailed(err)
		return nil, InternalServerError
	}

	// Check if the user ID doesn't exist
	if err != nil {
		// User ID does not exist
		s.logger.UserNotFoundByUserId(request.GetUserId())

		return nil, status.Error(codes.NotFound, FoundByUserIdFailure)
	}

	// User profile found by user ID
	s.logger.GetProfile(request.GetUserId())

	return &protobuf.GetProfileResponse{
		Message:   GetProfileSuccess,
		Username:  profile.Username,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
	}, nil
}

// UpdateUser updates the user
func (s Server) UpdateUser(
	ctx context.Context,
	request *protobuf.UpdateUserRequest,
) (response *protobuf.UpdateUserResponse, err error) {
	// Get the user ID from the access token
	subject, err := commongrpcserverctx.GetCtxTokenClaimsSubject(ctx)
	if err != nil {
		s.logger.MissingTokenClaimsSubject()
		return nil, InternalServerError
	}

	// Create the update fields BSON
	var update bson.M

	// Iterate over the request string fields
	for key, value := range map[string]interface{}{
		"first_name": request.GetFirstName(),
		"last_name":  request.GetLastName(),
	} {
		// Skip empty values
		if value != "" {
			update[key] = value
		}
	}

	// Iterate over the request time fields
	if request.GetBirthdate() != nil {
		update["birthdate"] = request.GetBirthdate().AsTime()
	}

	// Update the user
	_, err = s.userDatabase.UpdateUser(context.Background(), subject, &update)
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.UpdateUserFailed(err)
		return nil, InternalServerError
	}

	// User found by user ID
	s.logger.UpdateUser(subject)

	return &protobuf.UpdateUserResponse{
		Message: UpdatedSuccess,
	}, nil
}

// ChangeUsername changes the user's username
func (s Server) ChangeUsername(
	ctx context.Context,
	request *protobuf.ChangeUsernameRequest,
) (response *protobuf.ChangeUsernameResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateChangeUsernameRequest(request); err != nil {
		s.logger.GetProfileFailed(err)
		return nil, err
	}

	// Get the user ID from the access token
	subject, err := commongrpcserverctx.GetCtxTokenClaimsSubject(ctx)
	if err != nil {
		s.logger.MissingTokenClaimsSubject()
		return nil, InternalServerError
	}

	// Update the user's username
	err = s.userDatabase.UpdateUserUsername(subject, request.GetUsername())
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		s.logger.UpdateUsernameFailed(err)
		return nil, InternalServerError
	}

	// Check if the username already exists
	if err != nil {
		// Username exists
		s.logger.UsernameExists(request.GetUsername())

		return nil, status.Error(codes.AlreadyExists, UsernameExistsSuccess)
	}

	// Updated the user's username
	s.logger.UpdateUsername(subject)

	return &protobuf.ChangeUsernameResponse{
		Message: ChangeUsernameSuccess,
	}, nil
}

// ChangePassword changes the user's password
func (s Server) ChangePassword(
	ctx context.Context,
	request *protobuf.ChangePasswordRequest,
) (response *protobuf.ChangePasswordResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateChangePasswordRequest(request); err != nil {
		s.logger.UpdatePasswordFailed(err)
		return nil, err
	}

	// Get the user ID from the access token
	subject, err := commongrpcserverctx.GetCtxTokenClaimsSubject(ctx)
	if err != nil {
		s.logger.MissingTokenClaimsSubject()
		return nil, InternalServerError
	}

	// Check if the old password is correct
	userHashedPassword, err := s.userDatabase.GetUserHashedPassword(context.Background(), subject)
	if err != nil {
		s.logger.PasswordIsCorrectFailed(err)
		return nil, InternalServerError
	}

	// Check if the password matches
	matches := commonbcrypt.CheckPasswordHash(userHashedPassword.HashedPassword, request.GetOldPassword())
	if !matches {
		s.logger.PasswordIsIncorrect(subject)
		return nil, status.Error(codes.InvalidArgument, PasswordIsIncorrect)
	}

	// Get the user's hashed password
	hashedNewPassword, err := commonbcrypt.HashPassword(request.GetNewPassword())
	if err != nil {
		s.logger.HashPasswordFailed(err)
		return nil, InternalServerError
	}

	// Get outgoing gRPC context
	grpcCtx, err := commongrpcclientctx.GetOutgoingCtx(ctx)
	if err != nil {
		return nil, InternalServerError
	}

	// Update the user's password
	err = s.userDatabase.UpdateUserPassword(grpcCtx, subject, hashedNewPassword)
	if err != nil {
		s.logger.UpdatePasswordFailed(err)
		return nil, InternalServerError
	}

	// Updated the user's password
	s.logger.UpdatePassword(subject)

	return &protobuf.ChangePasswordResponse{
		Message: ChangePasswordSuccess,
	}, nil
}

// GetPhoneNumber gets the user's phone number
func (s Server) GetPhoneNumber(
	ctx context.Context,
	request *protobuf.GetPhoneNumberRequest,
) (*protobuf.GetPhoneNumberResponse, error) {
	// Get the user ID from the access token
	subject, err := commongrpcserverctx.GetCtxTokenClaimsSubject(ctx)
	if err != nil {
		s.logger.MissingTokenClaimsSubject()
		return nil, InternalServerError
	}

	// Get the current phone number by user ID
	phoneNumber, err := s.userDatabase.GetUserPhoneNumber(context.Background(), subject)
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.GetPhoneNumberFailed(err)
		return nil, InternalServerError
	}

	// User found by user ID
	s.logger.GetPhoneNumber(subject)

	return &protobuf.GetPhoneNumberResponse{
		Message:     GetPhoneNumbersSuccess,
		PhoneNumber: phoneNumber,
	}, nil
}

// ChangePhoneNumber changes the user's phone number
func (s Server) ChangePhoneNumber(
	ctx context.Context,
	request *protobuf.ChangePhoneNumberRequest,
) (response *protobuf.ChangePhoneNumberResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateChangePhoneNumberRequest(request); err != nil {
		s.logger.UpdatePhoneNumberFailed(err)
		return nil, err
	}

	// Get the user ID from the access token
	subject, err := commongrpcserverctx.GetCtxTokenClaimsSubject(ctx)
	if err != nil {
		s.logger.MissingTokenClaimsSubject()
		return nil, InternalServerError
	}

	// Update the user's phone number
	err = s.userDatabase.UpdateUserPhoneNumber(subject, request.GetPhoneNumber())
	if err != nil {
		s.logger.UpdatePhoneNumberFailed(err)
		return nil, InternalServerError
	}

	// Updated the user's phone number
	s.logger.UpdatePhoneNumber(subject)

	return &protobuf.ChangePhoneNumberResponse{
		Message: ChangePhoneNumberSuccess,
	}, nil
}

// GetMyProfile gets the user's profile
func (s Server) GetMyProfile(
	ctx context.Context,
	request *protobuf.GetMyProfileRequest,
) (response *protobuf.GetMyProfileResponse, err error) {
	/*
		// Get the user ID from the access token
		subject, err := commongrpcserverctx.GetCtxTokenClaimsSubject(ctx)
		if err != nil {
			s.logger.MissingTokenClaimsSubject()
			return nil, InternalServerError
		}

		// Get the full profile by user ID
		username, err := s.userDatabase.GetFullProfile(subject)
		if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
			s.logger.GetProfileFailed(err)
			return nil, InternalServerError
		}

		// Check if the user ID doesn't exist
		if username == "" {
			// User ID does not exist
			s.logger.UserNotFoundByUserId(request.GetUserId())

			return nil, status.Error(codes.NotFound, FoundByUserIdFailure)
		}

		// User found by user ID
		s.logger.UserFoundByUsername(request.GetUserId())

		return &protobuf.GetUsernameByUserIdResponse{
			Message:  FoundByUserIdSuccess,
			Username: username,
		}, nil
	*/
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

// GetActiveEmails gets the user's active emails
func (s Server) GetActiveEmails(
	ctx context.Context,
	request *protobuf.GetActiveEmailsRequest,
) (*protobuf.GetActiveEmailsResponse, error) {
	return nil, InDevelopmentError
}

// DeleteUser deletes the user's account
func (s Server) DeleteUser(
	ctx context.Context,
	request *protobuf.DeleteUserRequest,
) (*protobuf.DeleteUserResponse, error) {
	return nil, InDevelopmentError
}

// --- Requires more development ---

// VerifyEmail verifies the user's email
func (s Server) VerifyEmail(
	ctx context.Context,
	request *protobuf.VerifyEmailRequest,
) (*protobuf.VerifyEmailResponse, error) {
	return nil, InDevelopmentError
}

// SendVerificationEmail sends a verification email to the user
func (s Server) SendVerificationEmail(
	ctx context.Context,
	request *protobuf.SendVerificationEmailRequest,
) (*protobuf.SendVerificationEmailResponse, error) {
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
