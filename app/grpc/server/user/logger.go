package user

import commonlogger "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/utils/logger"

type Logger struct {
	logger commonlogger.Logger
}

// NewLogger is the logger for the user database
func NewLogger(logger commonlogger.Logger) Logger {
	return Logger{logger: logger}
}

// SignedUp logs the user sign up
func (l Logger) SignedUp(userId string, username string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User signed up", commonlogger.StatusSuccess, userId, username))
}

// SignUpFailed logs the user sign up failure
func (l Logger) SignUpFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("User sign up failed", err))
}

// PasswordIsCorrect logs the password check success
func (l Logger) PasswordIsCorrect(userId string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("Password is correct", commonlogger.StatusSuccess, userId))
}

// PasswordIsIncorrect logs the password check failure
func (l Logger) PasswordIsIncorrect(userId string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("Password is incorrect", commonlogger.StatusFailed, userId))
}

// PasswordIsCorrectFailed logs the password check failure
func (l Logger) PasswordIsCorrectFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("Password check failed", err))
}

// UserFoundByUsername logs the user retrieval success
func (l Logger) UserFoundByUsername(username string, userId string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User found by username", commonlogger.StatusSuccess, username, userId))
}

// UserNotFoundByUsername logs the user retrieval failure
func (l Logger) UserNotFoundByUsername(username string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User not found by username", commonlogger.StatusFailed, username))
}

// UserFoundByUserId logs the user retrieval success
func (l Logger) UserFoundByUserId(userId string, username string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User found by user ID", commonlogger.StatusSuccess, userId, username))
}

// UserNotFoundByUserId logs the user retrieval failure
func (l Logger) UserNotFoundByUserId(userId string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User not found by user ID", commonlogger.StatusFailed, userId))
}

// UserFoundBySharedId logs the user retrieval success
func (l Logger) UserFoundBySharedId(sharedId string, userId string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User found by shared ID", commonlogger.StatusSuccess, sharedId, userId))
}

// UserNotFoundBySharedId logs the user retrieval failure
func (l Logger) UserNotFoundBySharedId(sharedId string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User not found by shared ID", commonlogger.StatusFailed, sharedId))
}

// UserSharedIdFoundByUserId logs the user shared ID retrieval success
func (l Logger) UserSharedIdFoundByUserId(userId string, userSharedId string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User shared ID found by user ID", commonlogger.StatusSuccess, userId, userSharedId))
}

// UsernameExists logs the username check success
func (l Logger) UsernameExists(username string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("Username exists", commonlogger.StatusSuccess, username))
}

// UsernameExistsFailed logs the username check failure
func (l Logger) UsernameExistsFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("Username exists check failed", err))
}

// GetUsernameByUserIdFailed logs the username retrieval failure
func (l Logger) GetUsernameByUserIdFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("Failed to fetch username by user ID", err))
}

// GetUserIdByUsernameFailed logs the user ID retrieval failure
func (l Logger) GetUserIdByUsernameFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("Failed to fetch user ID by username", err))
}

// GetUserSharedIdByUserIdFailed logs the user shared ID retrieval failure
func (l Logger) GetUserSharedIdByUserIdFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("Failed to fetch user shared ID by user ID", err))
}

// GetUserIdByUserSharedIdFailed logs the user ID retrieval failure
func (l Logger) GetUserIdByUserSharedIdFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("Failed to fetch user ID by user shared ID", err))
}

// MissingTokenClaimsSubject logs the missing token claims subject
func (l Logger) MissingTokenClaimsSubject() {
	l.logger.LogMessage(commonlogger.NewLogMessage("Missing token claims subject", commonlogger.StatusFailed))
}

// UpdateUser logs the user update
func (l Logger) UpdateUser(userId string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User updated", commonlogger.StatusSuccess, userId))
}

// UpdateUserFailed logs the user update failure
func (l Logger) UpdateUserFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("User update failed", err))
}

// GetPhoneNumber logs the user phone number retrieval
func (l Logger) GetPhoneNumber(userId string, phoneNumber string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("Fetched user phone number", commonlogger.StatusSuccess, userId, phoneNumber))
}

// GetPhoneNumberFailed logs the user phone number retrieval failure
func (l Logger) GetPhoneNumberFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("Failed to fetch user phone number", err))
}

// GetProfile logs the user profile update
func (l Logger) GetProfile(userId string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("Fetched user profile", commonlogger.StatusSuccess, userId))
}

// GetProfileFailed logs the user profile update failure
func (l Logger) GetProfileFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("Failed to fetch user profile", err))
}

// UpdateUsername logs the user username update
func (l Logger) UpdateUsername(userId string, newUsername string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User username updated", commonlogger.StatusSuccess, userId, newUsername))
}

// UpdateUsernameFailed logs the user username update failure
func (l Logger) UpdateUsernameFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("User username update failed", err))
}

// UpdatePassword logs the user password update
func (l Logger) UpdatePassword(userId string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User password updated", commonlogger.StatusSuccess, userId))
}

// UpdatePasswordFailed logs the user password update failure
func (l Logger) UpdatePasswordFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("User password update failed", err))
}

// HashPasswordFailed logs a failed password hash attempt
func (l Logger) HashPasswordFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("Failed to hash password", err))
}

// UpdatePhoneNumber logs the user phone number update
func (l Logger) UpdatePhoneNumber(userId string, newPhoneNumber string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User phone number updated", commonlogger.StatusSuccess, userId, newPhoneNumber))
}

// UpdatePhoneNumberFailed logs the user phone number update failure
func (l Logger) UpdatePhoneNumberFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("User phone number update failed", err))
}

// DeleteUser logs the user deletion
func (l Logger) DeleteUser(userId string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User deleted", commonlogger.StatusSuccess, userId))
}

// DeleteUserFailed logs the user deletion failure
func (l Logger) DeleteUserFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("User deletion failed", err))
}

// UserEmailAlreadyExists logs the user email existence check success
func (l Logger) UserEmailAlreadyExists(userId string, email string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User email already exists", commonlogger.StatusFailed, userId, email))
}

// UserEmailNotFound logs the user email retrieval failure
func (l Logger) UserEmailNotFound(userId string, email string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User email not found", commonlogger.StatusFailed, userId, email))
}

// AddEmail logs the user email addition
func (l Logger) AddEmail(userId string, email string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User email added", commonlogger.StatusSuccess, userId, email))
}

// AddEmailFailed logs the user email addition failure
func (l Logger) AddEmailFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("User email addition failed", err))
}

// DeleteEmail logs the user email deletion
func (l Logger) DeleteEmail(userId string, email string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("User email deleted", commonlogger.StatusSuccess, userId, email))
}

// DeleteEmailFailed logs the user email deletion failure
func (l Logger) DeleteEmailFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("User email deletion failed", err))
}

// GetPrimaryEmail logs the user primary email retrieval
func (l Logger) GetPrimaryEmail(userId string, email string) {
	l.logger.LogMessage(commonlogger.NewLogMessage("Fetched user primary email", commonlogger.StatusSuccess, userId, email))
}

// GetPrimaryEmailFailed logs the user primary email retrieval failure
func (l Logger) GetPrimaryEmailFailed(err error) {
	l.logger.LogError(commonlogger.NewLogError("Failed to fetch user primary email", err))
}
