package mongodb

import (
	"github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb"
)

type UsersService struct {
	mongodbConnection *mongodb.Connection
}

// NewUsersService creates a new MongoDB users service
func NewUsersService(connection *mongodb.Connection) *UsersService {
	return &UsersService{mongodbConnection: connection}
}
