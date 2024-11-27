package validator

import (
	"context"
	commongrpcvalidator "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/grpc/server/validator"
	pbuser "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/protobuf/compiled/user"
	mongodbuser "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb/database/user"
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
func NewValidator(userDatabase *mongodbuser.Database, validator commongrpcvalidator.Validator) *Validator {
	return &Validator{userDatabase: userDatabase, validator: validator}
}

// UsernameExists checks if the username exists
func (d Validator) UsernameExists(usernameField string, username string, validations *map[string][]error) bool {
	if exists, _ := d.userDatabase.UsernameExists(context.Background(), username); exists {
		(*validations)[usernameField] = append(
			(*validations)[usernameField],
			UsernameTakenError,
		)
		return true
	}
	return false
}

// ValidateSignUpRequest validates the sign up request
func (d Validator) ValidateSignUpRequest(request *pbuser.SignUpRequest) error {
	// Get the request fields
	usernameField := "username"
	emailField := "email"
	birthdateField := "birthdate"

	// Get validations from fields to validate
	validations := d.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Password":    "password",
			"Username":    usernameField,
			"FirstName":   "first_name",
			"LastName":    "last_name",
			"Email":       emailField,
			"PhoneNumber": "phone_number",
		})

	// Check if the user already exists
	usernameExists := d.UsernameExists(usernameField, request.GetUsername(), validations)

	// Check if the email is valid
	d.validator.ValidateEmail(emailField, request.GetEmail(), validations)

	// Check if the birthdate is valid
	if birthdate := request.GetBirthdate(); birthdate != nil {
		d.validator.ValidateBirthdate(birthdateField, birthdate, validations)
	}

	// Get the code
	code := codes.InvalidArgument
	if usernameExists {
		code = codes.AlreadyExists
	}

	return d.validator.CheckValidations(validations, code)
}

// ValidateIsPasswordCorrectRequest validates the is password correct request
func (d Validator) ValidateIsPasswordCorrectRequest(request *pbuser.IsPasswordCorrectRequest) error {
	// Get validations from fields to validate
	validations := d.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Username": "username",
			"Password": "password",
		})

	return d.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateUsernameExistsRequest validates the username exists request
func (d Validator) ValidateUsernameExistsRequest(request *pbuser.UsernameExistsRequest) error {
	// Get validations from fields to validate
	validations := d.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Username": "username",
		})

	return d.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateGetUserIdByUsernameRequest validates the get user ID by username request
func (d Validator) ValidateGetUserIdByUsernameRequest(request *pbuser.GetUserIdByUsernameRequest) error {
	// Get validations from fields to validate
	validations := d.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Username": "username",
		})

	return d.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateGetUsernameByUserIdRequest validates the get username by user ID request
func (d Validator) ValidateGetUsernameByUserIdRequest(request *pbuser.GetUsernameByUserIdRequest) error {
	// Get validations from fields to validate
	validations := d.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"UserId": "user_id",
		})

	return d.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateGetProfileRequest validates the get profile request
func (d Validator) ValidateGetProfileRequest(request *pbuser.GetProfileRequest) error {
	// Get validations from fields to validate
	validations := d.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"UserId": "user_id",
		})

	return d.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateChangeUsernameRequest validates the change username request
func (d Validator) ValidateChangeUsernameRequest(request *pbuser.ChangeUsernameRequest) error {
	// Get validations from fields to validate
	validations := d.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"Username": "username",
		})

	return d.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateChangePasswordRequest validates the change password request
func (d Validator) ValidateChangePasswordRequest(request *pbuser.ChangePasswordRequest) error {
	// Get validations from fields to validate
	validations := d.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"OldPassword": "old_password",
			"NewPassword": "new_password",
		})

	// Check if the new password is different from the old password
	if request.GetOldPassword() == request.GetNewPassword() {
		(*validations)["new_password"] = append((*validations)["new_password"], NewPasswordSameAsOldError)
	}

	return d.validator.CheckValidations(validations, codes.InvalidArgument)
}

// ValidateChangePhoneNumberRequest validates the change phone number request
func (d Validator) ValidateChangePhoneNumberRequest(request *pbuser.ChangePhoneNumberRequest) error {
	// Get validations from fields to validate
	validations := d.validator.ValidateNonEmptyStringFields(
		request,
		&map[string]string{
			"PhoneNumber": "phone_number",
		})

	return d.validator.CheckValidations(validations, codes.InvalidArgument)
}
