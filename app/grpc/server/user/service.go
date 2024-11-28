package user

import (
	"errors"
	"github.com/google/uuid"
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

	// Create the user shared identifier
	newUserSharedId := commonuser.UserSharedIdentifier{
		ID:     primitive.NewObjectID(),
		UserID: userId,
		UUID:   uuid.New().String(),
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
	if err = s.userDatabase.CreateUser(
		&newUser,
		&newUserSharedId,
		&newUserEmail,
		&newUserPhoneNumber,
	); err != nil {
		s.logger.SignUpFailed(err)
		return nil, InternalServerError
	}

	// User signed up successfully
	s.logger.SignedUp(userId.Hex(), request.GetUsername())

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

		return nil, status.Error(codes.OK, IsPasswordCorrectFailed)
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

		return nil, status.Error(codes.NotFound, FoundByUsernameFailed)
	}

	// User found by username
	s.logger.UsernameExists(request.GetUsername())

	return &protobuf.UsernameExistsResponse{
		Message: FoundByUsernameSuccess,
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
	if err != nil {
		// Username does not exist
		s.logger.UserNotFoundByUsername(request.GetUsername())

		return nil, status.Error(codes.NotFound, FoundByUsernameFailed)
	}

	// User found by username
	s.logger.UserFoundByUsername(request.GetUsername(), userId)

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
	if err != nil {
		// User ID does not exist
		s.logger.UserNotFoundByUserId(request.GetUserId())

		return nil, status.Error(codes.NotFound, FoundByUserIdFailed)
	}

	// User found by user ID
	s.logger.UserFoundByUsername(request.GetUserId(), username)

	return &protobuf.GetUsernameByUserIdResponse{
		Message:  FoundByUserIdSuccess,
		Username: username,
	}, nil
}

// GetUserSharedIdByUserId gets the user's shared ID by user ID
func (s Server) GetUserSharedIdByUserId(ctx context.Context, request *protobuf.GetUserSharedIdByUserIdRequest) (response *protobuf.GetUserSharedIdByUserIdResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateGetUserSharedIdByUserIdRequest(request); err != nil {
		s.logger.GetUserSharedIdByUserIdFailed(err)
		return nil, err
	}

	// Get the user shared ID by user ID
	userSharedId, err := s.userDatabase.GetUserSharedIdByUserId(context.Background(), request.GetUserId())
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.GetUserSharedIdByUserIdFailed(err)
		return nil, InternalServerError
	}

	// Check if the user shared ID doesn't exist
	if err != nil {
		// User ID does not exist
		s.logger.UserNotFoundByUserId(request.GetUserId())

		return nil, status.Error(codes.NotFound, FoundByUserIdFailed)
	}

	// User shared ID found by user ID
	s.logger.UserSharedIdFoundByUserId(request.GetUserId(), userSharedId)

	return &protobuf.GetUserSharedIdByUserIdResponse{
		Message:      FoundByUserIdSuccess,
		UserSharedId: userSharedId,
	}, nil
}

// GetUserIdByUserSharedId gets the user's ID by shared ID
func (s Server) GetUserIdByUserSharedId(ctx context.Context, request *protobuf.GetUserIdByUserSharedIdRequest) (response *protobuf.GetUserIdByUserSharedIdResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateGetUserIdByUserSharedIdRequest(request); err != nil {
		s.logger.GetUserIdByUserSharedIdFailed(err)
		return nil, err
	}

	// Get the user ID by shared ID
	userId, err := s.userDatabase.GetUserIdByUserSharedId(context.Background(), request.GetUserSharedId())
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.GetUserIdByUserSharedIdFailed(err)
		return nil, InternalServerError
	}

	// Check if the user shared ID doesn't exist
	if err != nil {
		// User shared ID does not exist
		s.logger.UserNotFoundBySharedId(request.GetUserSharedId())

		return nil, status.Error(codes.NotFound, FoundByUserSharedIdFailed)
	}

	// User found by user shared ID
	s.logger.UserFoundBySharedId(request.GetUserSharedId(), userId)

	return &protobuf.GetUserIdByUserSharedIdResponse{
		Message: FoundByUserSharedIdSuccess,
		UserId:  userId,
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

		return nil, status.Error(codes.NotFound, FoundByUserIdFailed)
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
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
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
	_, err = s.userDatabase.UpdateUser(context.Background(), userId, &update)
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.UpdateUserFailed(err)
		return nil, InternalServerError
	}

	// User found by user ID
	s.logger.UpdateUser(userId)

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
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.logger.MissingTokenClaimsSubject()
		return nil, InternalServerError
	}

	// Update the user's username
	err = s.userDatabase.UpdateUserUsername(userId, request.GetUsername())
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
	s.logger.UpdateUsername(userId, request.GetUsername())

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
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.logger.MissingTokenClaimsSubject()
		return nil, InternalServerError
	}

	// Check if the old password is correct
	userHashedPassword, err := s.userDatabase.GetUserHashedPassword(context.Background(), userId)
	if err != nil {
		s.logger.PasswordIsCorrectFailed(err)
		return nil, InternalServerError
	}

	// Check if the password matches
	matches := commonbcrypt.CheckPasswordHash(userHashedPassword.HashedPassword, request.GetOldPassword())
	if !matches {
		s.logger.PasswordIsIncorrect(userId)
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
	err = s.userDatabase.UpdateUserPassword(grpcCtx, userId, hashedNewPassword)
	if err != nil {
		s.logger.UpdatePasswordFailed(err)
		return nil, InternalServerError
	}

	// Updated the user's password
	s.logger.UpdatePassword(userId)

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
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.logger.MissingTokenClaimsSubject()
		return nil, InternalServerError
	}

	// Get the current phone number by user ID
	phoneNumber, err := s.userDatabase.GetUserPhoneNumber(context.Background(), userId)
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.GetPhoneNumberFailed(err)
		return nil, InternalServerError
	}

	// User found by user ID
	s.logger.GetPhoneNumber(userId, phoneNumber)

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
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.logger.MissingTokenClaimsSubject()
		return nil, InternalServerError
	}

	// Update the user's phone number
	err = s.userDatabase.UpdateUserPhoneNumber(userId, request.GetPhoneNumber())
	if err != nil {
		s.logger.UpdatePhoneNumberFailed(err)
		return nil, InternalServerError
	}

	// Updated the user's phone number
	s.logger.UpdatePhoneNumber(userId, request.GetPhoneNumber())

	return &protobuf.ChangePhoneNumberResponse{
		Message: ChangePhoneNumberSuccess,
	}, nil
}

// AddEmail adds an email to the user's account
func (s Server) AddEmail(
	ctx context.Context,
	request *protobuf.AddEmailRequest,
) (response *protobuf.AddEmailResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateAddEmailRequest(request); err != nil {
		s.logger.AddEmailFailed(err)
		return nil, err
	}

	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.logger.MissingTokenClaimsSubject()
		return nil, InternalServerError
	}

	// Add the email to the user's account
	err = s.userDatabase.UpdateUserPhoneNumber(userId, request.GetEmail())
	if err != nil && !errors.Is(mongodbuser.EmailAlreadyExistsError, err) {
		s.logger.AddEmailFailed(err)
		return nil, InternalServerError
	}

	// Check if the email already exists
	if err != nil {
		s.logger.UserEmailAlreadyExists(userId, request.GetEmail())

		return nil, status.Error(codes.AlreadyExists, AddEmailFailed)
	}

	// Added email to the user's account
	s.logger.AddEmail(userId, request.GetEmail())

	return &protobuf.AddEmailResponse{
		Message: AddEmailSuccess,
	}, nil
}

// DeleteEmail deletes an email from the user's account
func (s Server) DeleteEmail(
	ctx context.Context,
	request *protobuf.DeleteEmailRequest,
) (response *protobuf.DeleteEmailResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateDeleteEmailRequest(request); err != nil {
		s.logger.DeleteEmailFailed(err)
		return nil, err
	}

	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.logger.MissingTokenClaimsSubject()
		return nil, InternalServerError
	}

	// Delete the email from the user's account
	err = s.userDatabase.DeleteEmail(context.Background(), userId, request.GetEmail())
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		s.logger.DeleteEmailFailed(err)
		return nil, InternalServerError
	}

	// Check if the email doesn't exist, or it's the primary email
	if err != nil {
		s.logger.DeleteEmailFailed(err)

		return nil, status.Error(codes.NotFound, DeleteEmailFailed)
	}

	// Deleted email from the user's account
	s.logger.DeleteEmail(userId, request.GetEmail())

	return &protobuf.DeleteEmailResponse{
		Message: DeleteEmailSuccess,
	}, nil
}

// GetPrimaryEmail gets the user's primary email
func (s Server) GetPrimaryEmail(
	ctx context.Context,
	request *protobuf.GetPrimaryEmailRequest,
) (*protobuf.GetPrimaryEmailResponse, error) {
	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.logger.MissingTokenClaimsSubject()
		return nil, InternalServerError
	}

	// Get the current primary email by user ID
	primaryEmail, err := s.userDatabase.GetPrimaryEmail(context.Background(), userId)
	if err != nil {
		s.logger.GetPrimaryEmailFailed(err)
		return nil, InternalServerError
	}

	// User primary email found by user ID
	s.logger.GetPrimaryEmail(userId, primaryEmail)

	return &protobuf.GetPrimaryEmailResponse{
		Message: GetPrimaryEmailSuccess,
		Email:   primaryEmail,
	}, nil
}

// DeleteUser deletes the user's account
func (s Server) DeleteUser(
	ctx context.Context,
	request *protobuf.DeleteUserRequest,
) (response *protobuf.DeleteUserResponse, err error) {
	// Validate the request
	if err = s.validator.ValidateDeleteUserRequest(request); err != nil {
		s.logger.DeleteUserFailed(err)
		return nil, err
	}

	// Get the user ID from the access token
	userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
	if err != nil {
		s.logger.MissingTokenClaimsSubject()
		return nil, InternalServerError
	}

	// Check if the password is correct
	userHashedPassword, err := s.userDatabase.GetUserHashedPassword(context.Background(), userId)
	if err != nil {
		s.logger.PasswordIsCorrectFailed(err)
		return nil, InternalServerError
	}

	// Check if the password matches
	matches := commonbcrypt.CheckPasswordHash(userHashedPassword.HashedPassword, request.GetPassword())
	if !matches {
		s.logger.PasswordIsIncorrect(userId)
		return nil, status.Error(codes.InvalidArgument, PasswordIsIncorrect)
	}

	// Get outgoing gRPC context
	grpcCtx, err := commongrpcclientctx.GetOutgoingCtx(ctx)
	if err != nil {
		return nil, InternalServerError
	}

	// Delete user
	err = s.userDatabase.DeleteUser(grpcCtx, userId)
	if err != nil {
		s.logger.DeleteUserFailed(err)
		return nil, InternalServerError
	}

	// User deleted successfully
	s.logger.DeleteUser(userId)

	return &protobuf.DeleteUserResponse{
		Message: DeleteUserSuccess,
	}, nil
}

// GetMyProfile gets the user's profile
func (s Server) GetMyProfile(
	ctx context.Context,
	request *protobuf.GetMyProfileRequest,
) (response *protobuf.GetMyProfileResponse, err error) {
	/*
		// Get the user ID from the access token
		userId, err := commongrpcserverctx.GetCtxTokenClaimsUserId(ctx)
		if err != nil {
			s.logger.MissingTokenClaimsSubject()
			return nil, InternalServerError
		}

		// Get the full profile by user ID
		username, err := s.userDatabase.GetFullProfile(userId)
		if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
			s.logger.GetProfileFailed(err)
			return nil, InternalServerError
		}

		// Check if the user ID doesn't exist
		if username == "" {
			// User ID does not exist
			s.logger.UserNotFoundByUserId(request.GetUserId())

			return nil, status.Error(codes.NotFound, FoundByUserIdFailed)
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
