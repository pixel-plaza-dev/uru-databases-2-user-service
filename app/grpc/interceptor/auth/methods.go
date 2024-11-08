package auth

// MethodsToIntercept is a map of methods to not intercept
var MethodsToIntercept = map[string]bool{
	"/user.User/SignUp":                false,
	"/user.User/IsPasswordCorrect":     false,
	"/user.User/UpdateProfile":         true,
	"/user.User/GetProfile":            true,
	"/user.User/GetFullProfile":        true,
	"/user.User/ChangePassword":        true,
	"/user.User/ChangeUsername":        true,
	"/user.User/AddEmail":              true,
	"/user.User/DeleteEmail":           true,
	"/user.User/SendVerificationEmail": true,
	"/user.User/VerifyEmail":           true,
	"/user.User/GetPrimaryEmail":       true,
	"/user.User/ChangePrimaryEmail":    true,
	"/user.User/GetActiveEmails":       true,
	"/user.User/ChangePhoneNumber":     true,
	"/user.User/GetPhoneNumber":        true,
	"/user.User/VerifyPhoneNumber":     true,
	"/user.User/ForgotPassword":        false,
	"/user.User/ResetPassword":         false,
	"/user.User/DeleteUser":            true,
}
