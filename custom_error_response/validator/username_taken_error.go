package validator

import (
	"github.com/pixel-plaza-dev/uru-databases-2-go-service-common/utils"
	"strings"
)

type UsernameTakenError struct {
	Username string
}

// Error returns a formatted error message for UsernameTakenError
func (u UsernameTakenError) Error() (message string) {
	formattedUsername := utils.AddBrackets(u.Username)
	return strings.Join([]string{"Username is already taken", formattedUsername}, " ")
}
