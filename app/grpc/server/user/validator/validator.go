package validator

import (
	"context"
	commonflag "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/config/flag"
	commongrpcvalidator "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/grpc/server/validator"
	commonvalidatorfields "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/utils/validator/fields"
	pbuser "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/compiled/pixel_plaza/user"
	appmongodbuser "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/database/mongodb/user"
	"google.golang.org/grpc/codes"
)

type (
	// Validator is the default validator for the user service gRPC methods
	Validator struct {
		userDatabase *appmongodbuser.Database
		validator    commongrpcvalidator.Validator
	}
)

// StructFieldsToValidate
var (
	SignUpRequestFieldsToValidate, _ = commonvalidatorfields.CreateGRPCStructFieldsToValidate(
		&pbuser.SignUpRequest{},
		commonflag.Mode,
	)
	IsPasswordCorrectRequestFieldsToValidate, _ = commonvalidatorfields.CreateGRPCStructFieldsToValidate(
		&pbuser.IsPasswordCorrectRequest{},
		commonflag.Mode,
	)
	UsernameExistsRequestFieldsToValidate, _ = commonvalidatorfields.CreateGRPCStructFieldsToValidate(
		&pbuser.UsernameExistsRequest{},
		commonflag.Mode,
	)
	GetUserIdByUsernameRequestFieldsToValidate, _ = commonvalidatorfields.CreateGRPCStructFieldsToValidate(
		&pbuser.GetUserIdByUsernameRequest{},
		commonflag.Mode,
	)
	GetUsernameByUserIdRequestFieldsToValidate, _ = commonvalidatorfields.CreateGRPCStructFieldsToValidate(
		&pbuser.GetUsernameByUserIdRequest{},
		commonflag.Mode,
	)
	GetProfileRequestFieldsToValidate, _ = commonvalidatorfields.CreateGRPCStructFieldsToValidate(
		&pbuser.GetProfileRequest{},
		commonflag.Mode,
	)
	ChangeUsernameRequestFieldsToValidate, _ = commonvalidatorfields.CreateGRPCStructFieldsToValidate(
		&pbuser.ChangeUsernameRequest{},
		commonflag.Mode,
	)
	ChangePasswordRequestFieldsToValidate, _ = commonvalidatorfields.CreateGRPCStructFieldsToValidate(
		&pbuser.ChangePasswordRequest{},
		commonflag.Mode,
	)
	ChangePhoneNumberRequestFieldsToValidate, _ = commonvalidatorfields.CreateGRPCStructFieldsToValidate(
		&pbuser.ChangePhoneNumberRequest{},
		commonflag.Mode,
	)
	DeleteUserRequestFieldsToValidate, _ = commonvalidatorfields.CreateGRPCStructFieldsToValidate(
		&pbuser.DeleteUserRequest{},
		commonflag.Mode,
	)
	AddEmailRequestFieldsToValidate, _ = commonvalidatorfields.CreateGRPCStructFieldsToValidate(
		&pbuser.AddEmailRequest{},
		commonflag.Mode,
	)
	ChangePrimaryEmailRequestFieldsToValidate, _ = commonvalidatorfields.CreateGRPCStructFieldsToValidate(
		&pbuser.ChangePrimaryEmailRequest{},
		commonflag.Mode,
	)
	DeleteEmailRequestFieldsToValidate, _ = commonvalidatorfields.CreateGRPCStructFieldsToValidate(
		&pbuser.DeleteEmailRequest{},
		commonflag.Mode,
	)
)

// NewValidator creates a new validator
func NewValidator(
	userDatabase *appmongodbuser.Database,
	validator commongrpcvalidator.Validator,
) (*Validator, error) {
	// Check if either the user database or the validator is nil
	if userDatabase == nil {
		return nil, appmongodbuser.NilDatabaseError
	}
	if validator == nil {
		return nil, commongrpcvalidator.NilValidatorError
	}

	return &Validator{userDatabase: userDatabase, validator: validator}, nil
}

// UsernameExists checks if the username exists
func (v *Validator) UsernameExists(
	usernameField string,
	username string,
	structFieldsValidations *commonvalidatorfields.StructFieldsValidations,
) bool {
	if exists, _ := v.userDatabase.UsernameExists(
		context.Background(),
		username,
	); exists {
		structFieldsValidations.AddFailedFieldValidationError(usernameField, UsernameTakenError)
		return true
	}
	return false
}

// ValidateSignUpRequest validates the sign up request
func (v *Validator) ValidateSignUpRequest(request *pbuser.SignUpRequest) error {
	// Get validations from fields to validate
	validations, _ := v.validator.ValidateNilFields(
		request,
		SignUpRequestFieldsToValidate,
	)

	// Check if the user already exists
	usernameExists := v.UsernameExists(
		"username",
		request.GetUsername(),
		validations,
	)

	// Check if the email is valid
	v.validator.ValidateEmail("email", request.GetEmail(), validations)

	// Check if the birthdate is valid
	if birthdate := request.GetBirthdate(); birthdate != nil {
		v.validator.ValidateBirthdate("birthdate", birthdate, validations)
	}

	// Get the code
	code := codes.InvalidArgument
	if usernameExists {
		code = codes.AlreadyExists
	}

	return v.validator.CheckValidations(validations, code)
}

// ValidateIsPasswordCorrectRequest validates the is password correct request
func (v *Validator) ValidateIsPasswordCorrectRequest(request *pbuser.IsPasswordCorrectRequest) error {
	// Get validations from fields to validate
	validations, _ := v.validator.ValidateNilFields(
		request,
		IsPasswordCorrectRequestFieldsToValidate,
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateUsernameExistsRequest validates the username exists request
func (v *Validator) ValidateUsernameExistsRequest(request *pbuser.UsernameExistsRequest) error {
	// Get validations from fields to validate
	validations, _ := v.validator.ValidateNilFields(
		request,
		UsernameExistsRequestFieldsToValidate,
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateGetUserIdByUsernameRequest validates the get user ID by username request
func (v *Validator) ValidateGetUserIdByUsernameRequest(request *pbuser.GetUserIdByUsernameRequest) error {
	// Get validations from fields to validate
	validations, _ := v.validator.ValidateNilFields(
		request,
		GetUserIdByUsernameRequestFieldsToValidate,
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateGetUsernameByUserIdRequest validates the get username by user ID request
func (v *Validator) ValidateGetUsernameByUserIdRequest(request *pbuser.GetUsernameByUserIdRequest) error {
	// Get validations from fields to validate
	validations, _ := v.validator.ValidateNilFields(
		request,
		GetUsernameByUserIdRequestFieldsToValidate,
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateGetProfileRequest validates the get profile request
func (v *Validator) ValidateGetProfileRequest(request *pbuser.GetProfileRequest) error {
	// Get validations from fields to validate
	validations, _ := v.validator.ValidateNilFields(
		request,
		GetProfileRequestFieldsToValidate,
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateChangeUsernameRequest validates the change username request
func (v *Validator) ValidateChangeUsernameRequest(request *pbuser.ChangeUsernameRequest) error {
	// Get validations from fields to validate
	validations, _ := v.validator.ValidateNilFields(
		request,
		ChangeUsernameRequestFieldsToValidate,
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateChangePasswordRequest validates the change password request
func (v *Validator) ValidateChangePasswordRequest(request *pbuser.ChangePasswordRequest) error {
	// Get validations from fields to validate
	validations, _ := v.validator.ValidateNilFields(
		request,
		ChangePasswordRequestFieldsToValidate,
	)

	// Check if the new password is different from the old password
	if request.GetOldPassword() == request.GetNewPassword() {
		validations.AddFailedFieldValidationError("new_password", NewPasswordSameAsOldError)
	}

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateChangePhoneNumberRequest validates the change phone number request
func (v *Validator) ValidateChangePhoneNumberRequest(request *pbuser.ChangePhoneNumberRequest) error {
	// Get validations from fields to validate
	validations, _ := v.validator.ValidateNilFields(
		request,
		ChangePhoneNumberRequestFieldsToValidate,
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateDeleteUserRequest validates the delete user request
func (v *Validator) ValidateDeleteUserRequest(request *pbuser.DeleteUserRequest) error {
	// Get validations from fields to validate
	validations, _ := v.validator.ValidateNilFields(
		request,
		DeleteUserRequestFieldsToValidate,
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateAddEmailRequest validates the add email request
func (v *Validator) ValidateAddEmailRequest(request *pbuser.AddEmailRequest) error {
	// Get validations from fields to validate
	validations, _ := v.validator.ValidateNilFields(
		request,
		AddEmailRequestFieldsToValidate,
	)

	// Check if the email is valid
	v.validator.ValidateEmail("email", request.GetEmail(), validations)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateChangePrimaryEmailRequest validates the change primary email request
func (v *Validator) ValidateChangePrimaryEmailRequest(request *pbuser.ChangePrimaryEmailRequest) error {
	// Get validations from fields to validate
	validations, _ := v.validator.ValidateNilFields(
		request,
		ChangePrimaryEmailRequestFieldsToValidate,
	)

	// Check if the email is valid
	v.validator.ValidateEmail("email", request.GetEmail(), validations)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateDeleteEmailRequest validates the delete email request
func (v *Validator) ValidateDeleteEmailRequest(request *pbuser.DeleteEmailRequest) error {
	// Get validations from fields to validate
	validations, _ := v.validator.ValidateNilFields(
		request,
		DeleteEmailRequestFieldsToValidate,
	)

	// Check if the email is valid
	v.validator.ValidateEmail("email", request.GetEmail(), validations)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}
