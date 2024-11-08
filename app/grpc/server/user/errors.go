package user

import (
	commonmessage "github.com/pixel-plaza-dev/uru-databases-2-go-api-common/message"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InternalServerError = status.Error(codes.Internal, commonmessage.InternalServerError)
	InDevelopmentError  = status.Error(codes.Internal, "in development")
)
