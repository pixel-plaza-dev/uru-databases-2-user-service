package user

import (
	"errors"
	commonbcrypt "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto/bcrypt"
	commonjwtvalidator "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto/jwt/validator"
	commonuser "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/database/mongodb/model/user"
	commongrpcclientctx "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/grpc/client/context"
	commongrpcserverctx "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/grpc/server/context"
	pbauth "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/compiled/pixel_plaza/auth"
	pbuser "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/compiled/pixel_plaza/user"
	appmongodbuser "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/database/mongodb/user"
	userservervalidator "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/grpc/server/user/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

// Server is the gRPC user server
type Server struct {
	userDatabase       *appmongodbuser.Database
	authClient         pbauth.AuthClient
	logger             *Logger
	validator          *userservervalidator.Validator
	jwtValidatorLogger *commonjwtvalidator.Logger
	pbuser.UnimplementedUserServer
}

// NewServer creates a new gRPC user server
func NewServer(
	userDatabase *appmongodbuser.Database,
	authClient pbauth.AuthClient,
	logger *Logger,
	validator *userservervalidator.Validator,
	jwtValidatorLogger *commonjwtvalidator.Logger,
) *Server {
	return &Server{
		userDatabase:       userDatabase,
		authClient:         authClient,
		logger:             logger,
		validator:          validator,
		jwtValidatorLogger: jwtValidatorLogger,
	}
}

// SignUp creates a new user
func (s *Server) SignUp(
	ctx context.Context,
	request *pbuser.SignUpRequest,
) (response *pbuser.SignUpResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateSignUpRequest(request); err != nil {
		s.logger.FailedToSignUp(err)
		return nil, err
	}

	// Hash the password
	hashedPassword, err := commonbcrypt.HashPassword(request.GetPassword())
	if err != nil {
		s.logger.FailedToHashPassword(err)
		return nil, InternalServerError
	}

	// Create a new user
	currentTime := time.Now()
	userId := primitive.NewObjectID()
	newUser := commonuser.User{
		ID:             userId,
		Username:       request.GetUsername(),
		FirstName:      request.GetFirstName(),
		LastName:       request.GetLastName(),
		HashedPassword: hashedPassword,
		JoinedAt:       currentTime,
	}

	// Add the birthdate if it exists
	if request.GetBirthdate() != nil {
		newUser.Birthdate = request.GetBirthdate().AsTime()
	}

	// Create the user email
	newUserEmail := commonuser.UserEmail{
		ID:         primitive.NewObjectID(),
		UserID:     userId,
		Email:      request.GetEmail(),
		AssignedAt: currentTime,
		IsPrimary:  true,
	}

	// Create the user phone number
	newUserPhoneNumber := commonuser.UserPhoneNumber{
		ID:          primitive.NewObjectID(),
		UserID:      userId,
		PhoneNumber: request.GetPhoneNumber(),
		AssignedAt:  currentTime,
	}

	// Insert the user into the user
	if err = s.userDatabase.InsertUser(
		&newUser,
		&newUserEmail,
		&newUserPhoneNumber,
	); err != nil {
		s.logger.FailedToSignUp(err)
		return nil, InternalServerError
	}

	// User signed up successfully
	s.logger.SignedUp(userId.Hex(), request.GetUsername())

	return &pbuser.SignUpResponse{
		Message: SignedUp,
	}, nil
}

func (s *Server) IsPasswordCorrect(
	ctx context.Context,
	request *pbuser.IsPasswordCorrectRequest,
) (response *pbuser.IsPasswordCorrectResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateIsPasswordCorrectRequest(request); err != nil {
		s.logger.FailedToComparePassword(err)
		return nil, err
	}

	// Get the user ID and hashed password by username
	user, err := s.userDatabase.GetUserHashedPassword(
		context.Background(), request.GetUsername(),
	)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		s.logger.FailedToComparePassword(err)
		return nil, InternalServerError
	}

	// Check if the user doesn't exist
	if err != nil {
		// User not found by username
		s.logger.UserNotFoundByUsername(request.GetUsername())

		return nil, status.Error(codes.NotFound, FailedToComparePassword)
	}

	// Check if the password matches
	matches := commonbcrypt.CheckPasswordHash(
		request.GetPassword(),
		user.HashedPassword,
	)

	// Get the user ID
	userId := user.ID.Hex()

	// Check if the password doesn't match or the user doesn't exist
	if !matches {
		// User checked password unsuccessfully
		s.logger.PasswordIsIncorrect(userId)

		return nil, status.Error(codes.InvalidArgument, FailedToComparePassword)
	}

	// User checked password successfully
	s.logger.PasswordIsCorrect(userId)

	return &pbuser.IsPasswordCorrectResponse{
		Message: PasswordIsCorrect,
		UserId:  userId,
	}, nil
}

// UsernameExists checks if the username exists
func (s *Server) UsernameExists(
	ctx context.Context,
	request *pbuser.UsernameExistsRequest,
) (response *pbuser.UsernameExistsResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateUsernameExistsRequest(request); err != nil {
		s.logger.FailedToCheckIfUsernameExists(err)
		return nil, err
	}

	// Check if the username exists
	exists, err := s.userDatabase.UsernameExists(
		context.Background(),
		request.GetUsername(),
	)
	if err != nil {
		s.logger.FailedToCheckIfUsernameExists(err)
		return nil, InternalServerError
	}

	// Check if the username doesn't exist
	if !exists {
		// Username does not exist
		s.logger.UserNotFoundByUsername(request.GetUsername())

		return nil, status.Error(codes.NotFound, NotFoundByUsername)
	}

	// User found by username
	s.logger.UsernameExists(request.GetUsername())

	return &pbuser.UsernameExistsResponse{
		Message: FoundByUsername,
	}, nil
}

// GetUserIdByUsername gets the user's ID by username
func (s *Server) GetUserIdByUsername(
	ctx context.Context,
	request *pbuser.GetUserIdByUsernameRequest,
) (response *pbuser.GetUserIdByUsernameResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateGetUserIdByUsernameRequest(request); err != nil {
		s.logger.FailedToGetUserIdByUsername(err)
		return nil, err
	}

	// Get the user ID by username
	userId, err := s.userDatabase.GetUserIdByUsername(
		context.Background(),
		request.GetUsername(),
	)
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.FailedToGetUserIdByUsername(err)
		return nil, InternalServerError
	}

	// Check if the username doesn't exist
	if err != nil {
		// Username does not exist
		s.logger.UserNotFoundByUsername(request.GetUsername())

		return nil, status.Error(codes.NotFound, NotFoundByUsername)
	}

	// User found by username
	s.logger.UserFoundByUsername(request.GetUsername(), userId)

	return &pbuser.GetUserIdByUsernameResponse{
		Message: FoundByUsername,
		UserId:  userId,
	}, nil
}

// GetUsernameByUserId gets the user's username by ID
func (s *Server) GetUsernameByUserId(
	ctx context.Context,
	request *pbuser.GetUsernameByUserIdRequest,
) (response *pbuser.GetUsernameByUserIdResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateGetUsernameByUserIdRequest(request); err != nil {
		s.logger.FailedToGetUsernameByUserId(err)
		return nil, err
	}

	// Get the username by user ID
	username, err := s.userDatabase.GetUsernameByUserId(
		context.Background(),
		request.GetUserId(),
	)
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.FailedToGetUsernameByUserId(err)
		return nil, InternalServerError
	}

	// Check if the user ID doesn't exist
	if err != nil {
		// User ID does not exist
		s.logger.UserNotFoundByUserId(request.GetUserId())

		return nil, status.Error(codes.NotFound, NotFoundByUserId)
	}

	// User found by user ID
	s.logger.UserFoundByUsername(request.GetUserId(), username)

	return &pbuser.GetUsernameByUserIdResponse{
		Message:  FoundByUserId,
		Username: username,
	}, nil
}

// GetProfile gets the user's profile
func (s *Server) GetProfile(
	ctx context.Context,
	request *pbuser.GetProfileRequest,
) (response *pbuser.GetProfileResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateGetProfileRequest(request); err != nil {
		s.logger.FailedToGetUserProfile(err)
		return nil, err
	}

	// Get the profile by username
	profile, err := s.userDatabase.GetUserProfile(
		context.Background(),
		request.GetUsername(),
	)
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.FailedToGetUserProfile(err)
		return nil, InternalServerError
	}

	// Check if the username doesn't exist
	if err != nil {
		// Username does not exist
		s.logger.UserNotFoundByUsername(request.GetUsername())

		return nil, status.Error(codes.NotFound, NotFoundByUserId)
	}

	// User profile found by username
	s.logger.GetUserProfile(request.GetUsername())

	return &pbuser.GetProfileResponse{
		Message:   FetchedUserProfile,
		FirstName: profile.FirstName,
		LastName:  profile.LastName,
		JoinedAt:  timestamppb.New(profile.JoinedAt),
	}, nil
}

// UpdateUser updates the user
func (s *Server) UpdateUser(
	ctx context.Context,
	request *pbuser.UpdateUserRequest,
) (response *pbuser.UpdateUserResponse, err error) {
	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.jwtValidatorLogger.MissingTokenClaimsUserId()
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
	_, err = s.userDatabase.UpdateUserByUserId(
		context.Background(),
		userId,
		&update,
	)
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.FailedToUpdateUser(err)
		return nil, InternalServerError
	}

	// User found by user ID
	s.logger.UpdatedUser(userId)

	return &pbuser.UpdateUserResponse{
		Message: Updated,
	}, nil
}

// ChangeUsername changes the user's username
func (s *Server) ChangeUsername(
	ctx context.Context,
	request *pbuser.ChangeUsernameRequest,
) (response *pbuser.ChangeUsernameResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateChangeUsernameRequest(request); err != nil {
		s.logger.FailedToGetUserProfile(err)
		return nil, err
	}

	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.jwtValidatorLogger.MissingTokenClaimsUserId()
		return nil, InternalServerError
	}

	// Update the user's username
	err = s.userDatabase.UpdateUserUsername(userId, request.GetUsername())
	if err != nil && !mongo.IsDuplicateKeyError(err) {
		s.logger.FailedToUpdateUsername(err)
		return nil, InternalServerError
	}

	// Check if the username already exists
	if err != nil {
		// Username exists
		s.logger.UsernameExists(request.GetUsername())

		return nil, status.Error(codes.AlreadyExists, UsernameExists)
	}

	// Updated the user's username
	s.logger.UpdatedUsername(userId, request.GetUsername())

	return &pbuser.ChangeUsernameResponse{
		Message: UpdatedUsername,
	}, nil
}

// ChangePassword changes the user's password
func (s *Server) ChangePassword(
	ctx context.Context,
	request *pbuser.ChangePasswordRequest,
) (response *pbuser.ChangePasswordResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateChangePasswordRequest(request); err != nil {
		s.logger.FailedToUpdatePassword(err)
		return nil, err
	}

	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.jwtValidatorLogger.MissingTokenClaimsUserId()
		return nil, InternalServerError
	}

	// Check if the old password is correct
	userHashedPassword, err := s.userDatabase.GetUserHashedPassword(
		context.Background(),
		userId,
	)
	if err != nil {
		s.logger.FailedToComparePassword(err)
		return nil, InternalServerError
	}

	// Check if the password matches
	matches := commonbcrypt.CheckPasswordHash(
		userHashedPassword.HashedPassword,
		request.GetOldPassword(),
	)
	if !matches {
		s.logger.PasswordIsIncorrect(userId)
		return nil, status.Error(codes.InvalidArgument, FailedToComparePassword)
	}

	// Get the user's hashed password
	hashedNewPassword, err := commonbcrypt.HashPassword(request.GetNewPassword())
	if err != nil {
		s.logger.FailedToHashPassword(err)
		return nil, InternalServerError
	}

	// Get outgoing gRPC context
	grpcCtx, err := commongrpcclientctx.GetOutgoingCtx(ctx)
	if err != nil {
		return nil, InternalServerError
	}

	// Update the user's password
	err = s.userDatabase.UpdateUserPassword(grpcCtx, userId, hashedNewPassword)
	if err != nil {
		s.logger.FailedToUpdatePassword(err)
		return nil, InternalServerError
	}

	// Updated the user's password
	s.logger.UpdatedPassword(userId)

	return &pbuser.ChangePasswordResponse{
		Message: UpdatedPassword,
	}, nil
}

// GetPhoneNumber gets the user's phone number
func (s *Server) GetPhoneNumber(
	ctx context.Context,
	request *emptypb.Empty,
) (*pbuser.GetPhoneNumberResponse, error) {
	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.jwtValidatorLogger.MissingTokenClaimsUserId()
		return nil, InternalServerError
	}

	// Get the current phone number by user ID
	phoneNumber, err := s.userDatabase.GetUserPhoneNumber(
		context.Background(),
		userId,
	)
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.FailedToGetUserPhoneNumber(err)
		return nil, InternalServerError
	}

	// User found by user ID
	s.logger.GetUserPhoneNumber(userId, phoneNumber)

	return &pbuser.GetPhoneNumberResponse{
		Message:     FetchedPhoneNumber,
		PhoneNumber: phoneNumber,
	}, nil
}

// ChangePhoneNumber changes the user's phone number
func (s *Server) ChangePhoneNumber(
	ctx context.Context,
	request *pbuser.ChangePhoneNumberRequest,
) (response *pbuser.ChangePhoneNumberResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateChangePhoneNumberRequest(request); err != nil {
		s.logger.FailedToUpdatePhoneNumber(err)
		return nil, err
	}

	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.jwtValidatorLogger.MissingTokenClaimsUserId()
		return nil, InternalServerError
	}

	// Update the user's phone number
	err = s.userDatabase.UpdateUserPhoneNumber(userId, request.GetPhoneNumber())
	if err != nil {
		s.logger.FailedToUpdatePhoneNumber(err)
		return nil, InternalServerError
	}

	// Updated the user's phone number
	s.logger.UpdatedUserPhoneNumber(userId, request.GetPhoneNumber())

	return &pbuser.ChangePhoneNumberResponse{
		Message: UpdatedPhoneNumber,
	}, nil
}

// AddEmail adds an email to the user's account
func (s *Server) AddEmail(
	ctx context.Context,
	request *pbuser.AddEmailRequest,
) (response *pbuser.AddEmailResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateAddEmailRequest(request); err != nil {
		s.logger.FailedToAddUserEmail(err)
		return nil, err
	}

	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.jwtValidatorLogger.MissingTokenClaimsUserId()
		return nil, InternalServerError
	}

	// Add the email to the user's account
	err = s.userDatabase.UpdateUserPhoneNumber(userId, request.GetEmail())
	if err != nil && !errors.Is(appmongodbuser.EmailAlreadyExistsError, err) {
		s.logger.FailedToAddUserEmail(err)
		return nil, InternalServerError
	}

	// Check if the email already exists
	if err != nil {
		s.logger.UserEmailAlreadyExists(userId, request.GetEmail())

		return nil, status.Error(codes.AlreadyExists, FailedToAddUserEmail)
	}

	// Added email to the user's account
	s.logger.AddedUserEmail(userId, request.GetEmail())

	return &pbuser.AddEmailResponse{
		Message: AddedUserEmail,
	}, nil
}

// DeleteEmail deletes an email from the user's account
func (s *Server) DeleteEmail(
	ctx context.Context,
	request *pbuser.DeleteEmailRequest,
) (response *pbuser.DeleteEmailResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateDeleteEmailRequest(request); err != nil {
		s.logger.FailedToDeleteUserEmail(err)
		return nil, err
	}

	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.jwtValidatorLogger.MissingTokenClaimsUserId()
		return nil, InternalServerError
	}

	// Delete the email from the user's account
	err = s.userDatabase.DeleteUserEmail(
		context.Background(),
		userId,
		request.GetEmail(),
	)
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.FailedToDeleteUserEmail(err)
		return nil, InternalServerError
	}

	// Check if the email doesn't exist, or it's the primary email
	if err != nil {
		s.logger.FailedToDeleteUserEmail(err)

		return nil, status.Error(codes.NotFound, FailedToDeleteUserEmail)
	}

	// Deleted email from the user's account
	s.logger.DeletedUserEmail(userId, request.GetEmail())

	return &pbuser.DeleteEmailResponse{
		Message: DeletedUserEmail,
	}, nil
}

// GetPrimaryEmail gets the user's primary email
func (s *Server) GetPrimaryEmail(
	ctx context.Context,
	request *emptypb.Empty,
) (*pbuser.GetPrimaryEmailResponse, error) {
	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.jwtValidatorLogger.MissingTokenClaimsUserId()
		return nil, InternalServerError
	}

	// Get the current primary email by user ID
	primaryEmail, err := s.userDatabase.GetUserPrimaryEmail(
		context.Background(),
		userId,
	)
	if err != nil {
		s.logger.FailedToGetPrimaryEmail(err)
		return nil, InternalServerError
	}

	// User primary email found by user ID
	s.logger.GetUserPrimaryEmail(userId, primaryEmail)

	return &pbuser.GetPrimaryEmailResponse{
		Message: FetchedUserPrimaryEmail,
		Email:   primaryEmail,
	}, nil
}

// ChangePrimaryEmail changes the user's primary email
func (s *Server) ChangePrimaryEmail(
	ctx context.Context,
	request *pbuser.ChangePrimaryEmailRequest,
) (response *pbuser.ChangePrimaryEmailResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateChangePrimaryEmailRequest(request); err != nil {
		s.logger.FailedToUpdateUserPrimaryEmail(err)
		return nil, err
	}

	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.jwtValidatorLogger.MissingTokenClaimsUserId()
		return nil, InternalServerError
	}

	// Update the user's primary email
	err = s.userDatabase.UpdateUserPrimaryEmail(userId, request.GetEmail())
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.FailedToUpdateUserPrimaryEmail(err)
		return nil, InternalServerError
	}

	// Check if the user email doesn't exist
	if err != nil {
		s.logger.UserEmailNotFound(userId, request.GetEmail())

		return nil, status.Error(codes.NotFound, NotFoundUserEmail)
	}

	// Change user primary email
	s.logger.UpdatedUserPrimaryEmail(userId, request.GetEmail())

	return &pbuser.ChangePrimaryEmailResponse{
		Message: UpdatedUserPrimaryEmail,
	}, nil
}

// GetActiveEmails gets the user's active emails
func (s *Server) GetActiveEmails(
	ctx context.Context,
	request *emptypb.Empty,
) (*pbuser.GetActiveEmailsResponse, error) {
	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.jwtValidatorLogger.MissingTokenClaimsUserId()
		return nil, InternalServerError
	}

	// Get the active emails by user ID
	activeEmails, err := s.userDatabase.GetUserActiveEmails(
		context.Background(),
		userId,
	)
	if err != nil {
		s.logger.FailedToGetActiveEmails(err)
		return nil, InternalServerError
	}

	// User active emails found by user ID
	s.logger.GetUserActiveEmails(userId)

	return &pbuser.GetActiveEmailsResponse{
		Message: FetchedUserActiveEmails,
		Emails:  activeEmails,
	}, nil
}

// DeleteUser deletes the user's account
func (s *Server) DeleteUser(
	ctx context.Context,
	request *pbuser.DeleteUserRequest,
) (response *pbuser.DeleteUserResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateDeleteUserRequest(request); err != nil {
		s.logger.FailedToDeleteUser(err)
		return nil, err
	}

	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.jwtValidatorLogger.MissingTokenClaimsUserId()
		return nil, InternalServerError
	}

	// Check if the password is correct
	userHashedPassword, err := s.userDatabase.GetUserHashedPassword(
		context.Background(),
		userId,
	)
	if err != nil {
		s.logger.FailedToComparePassword(err)
		return nil, InternalServerError
	}

	// Check if the password matches
	matches := commonbcrypt.CheckPasswordHash(
		userHashedPassword.HashedPassword,
		request.GetPassword(),
	)
	if !matches {
		s.logger.PasswordIsIncorrect(userId)
		return nil, status.Error(codes.InvalidArgument, FailedToComparePassword)
	}

	// Get outgoing gRPC context
	grpcCtx, err := commongrpcclientctx.GetOutgoingCtx(ctx)
	if err != nil {
		return nil, InternalServerError
	}

	// Delete user
	err = s.userDatabase.DeleteUser(grpcCtx, userId)
	if err != nil {
		s.logger.FailedToDeleteUser(err)
		return nil, InternalServerError
	}

	// User deleted successfully
	s.logger.DeletedUser(userId)

	return &pbuser.DeleteUserResponse{
		Message: DeletedUser,
	}, nil
}

// GetMyProfile gets the user's profile
func (s *Server) GetMyProfile(
	ctx context.Context,
	request *emptypb.Empty,
) (response *pbuser.GetMyProfileResponse, err error) {
	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.jwtValidatorLogger.MissingTokenClaimsUserId()
		return nil, InternalServerError
	}

	// Get the user own profile by user ID
	fullProfile, emails, phoneNumber, err := s.userDatabase.GetMyProfile(userId)
	if err != nil {
		s.logger.FailedToGetUserOwnProfile(err)
		return nil, InternalServerError
	}

	// User own profile found by user ID
	s.logger.GetUserOwnProfile(userId)

	return &pbuser.GetMyProfileResponse{
		Message:     FetchedUserOwnProfile,
		Username:    fullProfile.Username,
		FirstName:   fullProfile.FirstName,
		LastName:    fullProfile.LastName,
		Birthdate:   timestamppb.New(fullProfile.Birthdate),
		JoinedAt:    timestamppb.New(fullProfile.JoinedAt),
		Emails:      *emails,
		PhoneNumber: phoneNumber,
	}, nil
}

// --- Requires more development ---

func (s *Server) SetProfilePicture(
	ctx context.Context,
	request *pbuser.SetProfilePictureRequest,
) (*pbuser.SetProfilePictureResponse, error) {
	return nil, InDevelopmentError
}

// SendVerificationEmail sends a verification email to the user
func (s *Server) SendVerificationEmail(
	ctx context.Context,
	request *pbuser.SendVerificationEmailRequest,
) (*pbuser.SendVerificationEmailResponse, error) {
	return nil, InDevelopmentError
}

// VerifyEmail verifies the user's email
func (s *Server) VerifyEmail(
	ctx context.Context,
	request *pbuser.VerifyEmailRequest,
) (*pbuser.VerifyEmailResponse, error) {
	return nil, InDevelopmentError
}

// VerifyPhoneNumber verifies the user's phone number
func (s *Server) VerifyPhoneNumber(
	ctx context.Context,
	request *pbuser.VerifyPhoneNumberRequest,
) (*pbuser.VerifyPhoneNumberResponse, error) {
	return nil, InDevelopmentError
}

// SendVerificationSMS sends a verification SMS to the user
func (s *Server) SendVerificationSMS(
	ctx context.Context,
	request *pbuser.SendVerificationSMSRequest,
) (*pbuser.SendVerificationSMSResponse, error) {
	return nil, InDevelopmentError
}

// ForgotPassword sends a password reset link to the user's email
func (s *Server) ForgotPassword(
	ctx context.Context,
	request *pbuser.ForgotPasswordRequest,
) (*pbuser.ForgotPasswordResponse, error) {
	return nil, InDevelopmentError
}

// ResetPassword resets the user's password
func (s *Server) ResetPassword(
	ctx context.Context,
	request *pbuser.ResetPasswordRequest,
) (*pbuser.ResetPasswordResponse, error) {
	return nil, InDevelopmentError
}
