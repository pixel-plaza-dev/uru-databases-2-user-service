package user

import (
	commongin "github.com/pixel-plaza-dev/uru-databases-2-go-api-common/server/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InternalServerError = status.Error(codes.Internal, commongin.InternalServerError)
	InDevelopmentError  = status.Error(codes.Internal, "in development")
)
