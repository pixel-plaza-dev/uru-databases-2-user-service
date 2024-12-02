package validator

import (
	"context"
	commongrpcvalidator "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/grpc/server/validator"
	pbuser "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/compiled/pixel_plaza/user"
	mongodbuser "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/database/mongodb/user"
	"google.golang.org/grpc/codes"
)

type (
	// Validator is the default validator for the user service gRPC methods
	Validator struct {
		userDatabase *mongodbuser.Database
		validator    commongrpcvalidator.Validator
	}
)

// NewValidator creates a new validator
func NewValidator(
	userDatabase *mongodbuser.Database,
	validator commongrpcvalidator.Validator,
) *Validator {
	return &Validator{userDatabase: userDatabase, validator: validator}
}

// UsernameExists checks if the username exists
func (v Validator) UsernameExists(
	usernameField string,
	username string,
	validations *map[string][]error,
) bool {
	if exists, _ := v.userDatabase.UsernameExists(
		context.Background(),
		username,
	); exists {
		(*validations)[usernameField] = append(
			(*validations)[usernameField],
			UsernameTakenError,
		)
		return true
	}
	return false
}

// ValidateSignUpRequest validates the sign up request
func (v Validator) ValidateSignUpRequest(request *pbuser.SignUpRequest) error {
	// Get the request fields
	usernameField := "username"
	emailField := "email"
	birthdateField := "birthdate"

	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Password":    "password",
			"Username":    usernameField,
			"FirstName":   "first_name",
			"LastName":    "last_name",
			"Email":       emailField,
			"PhoneNumber": "phone_number",
		},
	)

	// Check if the user already exists
	usernameExists := v.UsernameExists(
		usernameField,
		request.GetUsername(),
		validations,
	)

	// Check if the email is valid
	v.validator.ValidateEmail(emailField, request.GetEmail(), validations)

	// Check if the birthdate is valid
	if birthdate := request.GetBirthdate(); birthdate != nil {
		v.validator.ValidateBirthdate(birthdateField, birthdate, validations)
	}

	// Get the code
	code := codes.InvalidArgument
	if usernameExists {
		code = codes.AlreadyExists
	}

	return v.validator.CheckValidations(validations, code)
}

// ValidateIsPasswordCorrectRequest validates the is password correct request
func (v Validator) ValidateIsPasswordCorrectRequest(request *pbuser.IsPasswordCorrectRequest) error {
	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Username": "username",
			"Password": "password",
		},
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateUsernameExistsRequest validates the username exists request
func (v Validator) ValidateUsernameExistsRequest(request *pbuser.UsernameExistsRequest) error {
	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Username": "username",
		},
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateGetUserIdByUsernameRequest validates the get user ID by username request
func (v Validator) ValidateGetUserIdByUsernameRequest(request *pbuser.GetUserIdByUsernameRequest) error {
	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Username": "username",
		},
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateGetUsernameByUserIdRequest validates the get username by user ID request
func (v Validator) ValidateGetUsernameByUserIdRequest(request *pbuser.GetUsernameByUserIdRequest) error {
	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"UserId": "user_id",
		},
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateGetUserSharedIdByUserIdRequest validates the get user shared ID by user ID request
func (v Validator) ValidateGetUserSharedIdByUserIdRequest(request *pbuser.GetUserSharedIdByUserIdRequest) error {
	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"UserId": "user_id",
		},
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateGetUserIdByUserSharedIdRequest validates the get user ID by user shared ID request
func (v Validator) ValidateGetUserIdByUserSharedIdRequest(request *pbuser.GetUserIdByUserSharedIdRequest) error {
	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"UserSharedId": "user_shared_id",
		},
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateGetProfileRequest validates the get profile request
func (v Validator) ValidateGetProfileRequest(request *pbuser.GetProfileRequest) error {
	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Username": "username",
		},
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateChangeUsernameRequest validates the change username request
func (v Validator) ValidateChangeUsernameRequest(request *pbuser.ChangeUsernameRequest) error {
	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Username": "username",
		},
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateChangePasswordRequest validates the change password request
func (v Validator) ValidateChangePasswordRequest(request *pbuser.ChangePasswordRequest) error {
	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"OldPassword": "old_password",
			"NewPassword": "new_password",
		},
	)

	// Check if the new password is different from the old password
	if request.GetOldPassword() == request.GetNewPassword() {
		(*validations)["new_password"] = append(
			(*validations)["new_password"],
			NewPasswordSameAsOldError,
		)
	}

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateChangePhoneNumberRequest validates the change phone number request
func (v Validator) ValidateChangePhoneNumberRequest(request *pbuser.ChangePhoneNumberRequest) error {
	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"PhoneNumber": "phone_number",
		},
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateDeleteUserRequest validates the delete user request
func (v Validator) ValidateDeleteUserRequest(request *pbuser.DeleteUserRequest) error {
	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Password": "password",
		},
	)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateAddEmailRequest validates the add email request
func (v Validator) ValidateAddEmailRequest(request *pbuser.AddEmailRequest) error {
	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Email": "email",
		},
	)

	// Check if the email is valid
	v.validator.ValidateEmail("email", request.GetEmail(), validations)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateChangePrimaryEmailRequest validates the change primary email request
func (v Validator) ValidateChangePrimaryEmailRequest(request *pbuser.ChangePrimaryEmailRequest) error {
	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Email": "email",
		},
	)

	// Check if the email is valid
	v.validator.ValidateEmail("email", request.GetEmail(), validations)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateDeleteEmailRequest validates the delete email request
func (v Validator) ValidateDeleteEmailRequest(request *pbuser.DeleteEmailRequest) error {
	// Get validations from fields to validate
	validations := v.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Email": "email",
		},
	)

	// Check if the email is valid
	v.validator.ValidateEmail("email", request.GetEmail(), validations)

	return v.validator.CheckValidations(validations, codes.InvalidArgument)
}
