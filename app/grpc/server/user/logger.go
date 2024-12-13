package user

import commonlogger "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/utils/logger"

type Logger struct {
	logger commonlogger.Logger
}

// NewLogger is the logger for the user database
func NewLogger(logger commonlogger.Logger) (*Logger, error) {
	// Check if the logger is nil
	if logger == nil {
		return nil, commonlogger.NilLoggerError
	}

	return &Logger{logger: logger}, nil
}

// SignedUp logs that the user signed up
func (l *Logger) SignedUp(userId string, username string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User signed up",
			commonlogger.StatusSuccess,
			userId,
			username,
		),
	)
}

// FailedToSignUp logs the user sign up failure
func (l *Logger) FailedToSignUp(err error) {
	l.logger.LogError(commonlogger.NewLogError("User sign up failed", err))
}

// PasswordIsCorrect logs the password check success
func (l *Logger) PasswordIsCorrect(userId string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"Password is correct",
			commonlogger.StatusSuccess,
			userId,
		),
	)
}

// PasswordIsIncorrect logs the password check failure
func (l *Logger) PasswordIsIncorrect(userId string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"Password is incorrect",
			commonlogger.StatusFailed,
			userId,
		),
	)
}

// FailedToComparePassword logs the password check failure
func (l *Logger) FailedToComparePassword(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"Failed to compare password",
			err,
		),
	)
}

// UserFoundByUsername logs the user retrieval success
func (l *Logger) UserFoundByUsername(username string, userId string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User found by username",
			commonlogger.StatusSuccess,
			username,
			userId,
		),
	)
}

// UserNotFoundByUsername logs the user retrieval failure
func (l *Logger) UserNotFoundByUsername(username string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User not found by username",
			commonlogger.StatusFailed,
			username,
		),
	)
}

// UserFoundByUserId logs the user retrieval success
func (l *Logger) UserFoundByUserId(userId string, username string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User found by user ID",
			commonlogger.StatusSuccess,
			userId,
			username,
		),
	)
}

// UserNotFoundByUserId logs the user retrieval failure
func (l *Logger) UserNotFoundByUserId(userId string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User not found by user ID",
			commonlogger.StatusFailed,
			userId,
		),
	)
}

// UsernameExists logs the username check success
func (l *Logger) UsernameExists(username string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"Username exists",
			commonlogger.StatusSuccess,
			username,
		),
	)
}

// FailedToCheckIfUsernameExists logs the username check failure
func (l *Logger) FailedToCheckIfUsernameExists(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"Username exists check failed",
			err,
		),
	)
}

// FailedToGetUsernameByUserId logs the username retrieval failure
func (l *Logger) FailedToGetUsernameByUserId(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"Failed to get username by user ID",
			err,
		),
	)
}

// FailedToGetUserIdByUsername logs the user ID retrieval failure
func (l *Logger) FailedToGetUserIdByUsername(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"Failed to get user ID by username",
			err,
		),
	)
}

// UpdatedUser logs the user update
func (l *Logger) UpdatedUser(userId string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User updated",
			commonlogger.StatusSuccess,
			userId,
		),
	)
}

// FailedToUpdateUser logs the user update failure
func (l *Logger) FailedToUpdateUser(err error) {
	l.logger.LogError(commonlogger.NewLogError("User update failed", err))
}

// GetUserPhoneNumber logs the user phone number retrieval
func (l *Logger) GetUserPhoneNumber(userId string, phoneNumber string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"Fetched user phone number",
			commonlogger.StatusSuccess,
			userId,
			phoneNumber,
		),
	)
}

// FailedToGetUserPhoneNumber logs the user phone number retrieval failure
func (l *Logger) FailedToGetUserPhoneNumber(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"Failed to get user phone number",
			err,
		),
	)
}

// GetUserProfile logs the user profile update
func (l *Logger) GetUserProfile(username string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"Fetched user profile",
			commonlogger.StatusSuccess,
			username,
		),
	)
}

// FailedToGetUserProfile logs the user profile update failure
func (l *Logger) FailedToGetUserProfile(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"Failed to fetch user profile",
			err,
		),
	)
}

// GetUserOwnProfile logs the user own profile retrieval
func (l *Logger) GetUserOwnProfile(userId string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"Fetched user own profile",
			commonlogger.StatusSuccess,
			userId,
		),
	)
}

// FailedToGetUserOwnProfile logs the user own profile retrieval failure
func (l *Logger) FailedToGetUserOwnProfile(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"Failed to fetch user own profile",
			err,
		),
	)
}

// UpdatedUsername logs the user username update
func (l *Logger) UpdatedUsername(userId string, newUsername string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User username updated",
			commonlogger.StatusSuccess,
			userId,
			newUsername,
		),
	)
}

// FailedToUpdateUsername logs the user username update failure
func (l *Logger) FailedToUpdateUsername(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"User username update failed",
			err,
		),
	)
}

// UpdatedPassword logs the user password update
func (l *Logger) UpdatedPassword(userId string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User password updated",
			commonlogger.StatusSuccess,
			userId,
		),
	)
}

// FailedToUpdatePassword logs the user password update failure
func (l *Logger) FailedToUpdatePassword(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"User password update failed",
			err,
		),
	)
}

// FailedToHashPassword logs a failed password hash attempt
func (l *Logger) FailedToHashPassword(err error) {
	l.logger.LogError(commonlogger.NewLogError("Failed to hash password", err))
}

// UpdatedUserPhoneNumber logs the user phone number update
func (l *Logger) UpdatedUserPhoneNumber(userId string, newPhoneNumber string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User phone number updated",
			commonlogger.StatusSuccess,
			userId,
			newPhoneNumber,
		),
	)
}

// FailedToUpdatePhoneNumber logs the user phone number update failure
func (l *Logger) FailedToUpdatePhoneNumber(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"User phone number update failed",
			err,
		),
	)
}

// DeletedUser logs the user deletion
func (l *Logger) DeletedUser(userId string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User deleted",
			commonlogger.StatusSuccess,
			userId,
		),
	)
}

// FailedToDeleteUser logs the user deletion failure
func (l *Logger) FailedToDeleteUser(err error) {
	l.logger.LogError(commonlogger.NewLogError("User deletion failed", err))
}

// UserEmailAlreadyExists logs the user email existence check success
func (l *Logger) UserEmailAlreadyExists(userId string, email string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User email already exists",
			commonlogger.StatusFailed,
			userId,
			email,
		),
	)
}

// UserEmailNotFound logs the user email retrieval failure
func (l *Logger) UserEmailNotFound(userId string, email string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User email not found",
			commonlogger.StatusFailed,
			userId,
			email,
		),
	)
}

// AddedUserEmail logs the user email addition
func (l *Logger) AddedUserEmail(userId string, email string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User email added",
			commonlogger.StatusSuccess,
			userId,
			email,
		),
	)
}

// FailedToAddUserEmail logs the user email addition failure
func (l *Logger) FailedToAddUserEmail(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"User email addition failed",
			err,
		),
	)
}

// UpdatedUserPrimaryEmail logs the user primary email change
func (l *Logger) UpdatedUserPrimaryEmail(userId string, email string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User primary email changed",
			commonlogger.StatusSuccess,
			userId,
			email,
		),
	)
}

// FailedToUpdateUserPrimaryEmail logs the user primary email change failure
func (l *Logger) FailedToUpdateUserPrimaryEmail(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"User primary email change failed",
			err,
		),
	)
}

// DeletedUserEmail logs the user email deletion
func (l *Logger) DeletedUserEmail(userId string, email string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"User email deleted",
			commonlogger.StatusSuccess,
			userId,
			email,
		),
	)
}

// FailedToDeleteUserEmail logs the user email deletion failure
func (l *Logger) FailedToDeleteUserEmail(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"User email deletion failed",
			err,
		),
	)
}

// GetUserPrimaryEmail logs the user primary email retrieval
func (l *Logger) GetUserPrimaryEmail(userId string, email string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"Fetched user primary email",
			commonlogger.StatusSuccess,
			userId,
			email,
		),
	)
}

// FailedToGetPrimaryEmail logs the user primary email retrieval failure
func (l *Logger) FailedToGetPrimaryEmail(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"Failed to fetch user primary email",
			err,
		),
	)
}

// GetUserActiveEmails logs the user active emails retrieval
func (l *Logger) GetUserActiveEmails(userId string) {
	l.logger.LogMessage(
		commonlogger.NewLogMessage(
			"Fetched user active emails",
			commonlogger.StatusSuccess,
			userId,
		),
	)
}

// FailedToGetActiveEmails logs the user active emails retrieval failure
func (l *Logger) FailedToGetActiveEmails(err error) {
	l.logger.LogError(
		commonlogger.NewLogError(
			"Failed to fetch user active emails",
			err,
		),
	)
}
