package main

import (
	"flag"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	commonenv "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/env"
	commonflag "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/flag"
	commonauthinterceptor "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/grpc/server/interceptor/auth"
	commonjwt "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/jwt"
	commonjwtvalidator "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/jwt/validator"
	commonlistener "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/listener"
	commonmongodb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb"
	protobuf "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/compiled-protobuf/user"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/grpc/interceptor/auth"
	userserver "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/grpc/server/user"
	appjwt "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/jwt"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/listener"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/logger"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb"
	userdatabase "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb/database/user"
	"google.golang.org/grpc"
	"net"
	"time"
)

// Load environment variables
func init() {
	if err := godotenv.Load(); err != nil {
		panic(commonenv.FailedToLoadEnvironmentVariablesError)
	}
}

func main() {
	// Declare flags and parse them
	commonflag.SetModeFlag()
	flag.Parse()
	logger.FlagLogger.ModeFlagSet(commonflag.Mode)

	// Get the listener port
	servicePort, err := commonlistener.LoadServicePort(listener.UserServicePortKey)
	if err != nil {
		panic(err)
	}
	logger.EnvironmentLogger.EnvironmentVariableLoaded(listener.UserServicePortKey)

	// Get the MongoDB URI
	mongoDbUri, err := commonmongodb.LoadMongoDBURI(mongodb.UriKey)
	if err != nil {
		panic(err)
	}
	logger.EnvironmentLogger.EnvironmentVariableLoaded(mongodb.UriKey)

	// Get the required MongoDB database name
	mongoDbName, err := commonmongodb.LoadMongoDBName(mongodb.DbNameKey)
	if err != nil {

		panic(err)
	}
	logger.EnvironmentLogger.EnvironmentVariableLoaded(mongodb.DbNameKey)

	// Get the JWT public key
	jwtPublicKey, err := commonjwt.LoadJwtKey(appjwt.PublicKey)
	if err != nil {
		panic(err)
	}

	// Get the MongoDB configuration
	mongoDbConfig := &commonmongodb.Config{Uri: mongoDbUri, Timeout: mongodb.ConnectionCtxTimeout}

	// Get the connection handler
	mongodbConnection := commonmongodb.NewDefaultConnectionHandler(mongoDbConfig)

	// Connect to MongoDB and get the client
	mongodbClient, err := mongodbConnection.Connect()
	if err != nil {
		panic(err)
	}

	// Create user database handler
	userDatabase, err := userdatabase.NewDatabase(mongodbClient, mongoDbName, logger.UserDatabaseLogger)
	if err != nil {
		panic(err)
	}
	defer func() {
		// Disconnect from MongoDB
		mongodbConnection.Disconnect()
		logger.MongoDbLogger.DisconnectedFromMongoDB()
	}()
	logger.MongoDbLogger.ConnectedToMongoDB()

	// Listen on the given port
	portListener, err := net.Listen("tcp", servicePort.FormattedPort)
	if err != nil {
		panic(commonlistener.FailedToListenError)
	}
	defer func() {
		if err := portListener.Close(); err != nil {
			panic(commonlistener.FailedToCloseError)
		}
	}()

	// Create JWT validator
	jwtValidator, err := commonjwtvalidator.NewDefaultValidator(jwtPublicKey, func(claims *jwt.MapClaims) (*jwt.MapClaims, error) {
		// Get the expiration time
		exp, err := claims.GetExpirationTime()
		if err != nil {
			return nil, commonjwt.InvalidClaimsError
		}

		// Check if the token is expired
		if exp.Before(time.Now()) {
			return nil, commonjwt.TokenExpiredError
		}
		return claims, nil
	})
	if err != nil {
		panic(err)
	}

	// Create gRPC Authentication interceptor
	authInterceptor := commonauthinterceptor.NewInterceptor(jwtValidator, auth.MethodsToIntercept)

	// Create a new gRPC server
	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		authInterceptor.UnaryServerInterceptor()))

	// Create a new gRPC user server
	userServer := userserver.NewServer(userDatabase, logger.UserServerLogger)

	// Register the user server with the gRPC server
	protobuf.RegisterUserServer(s, userServer)
	logger.ListenerLogger.ServerStarted(servicePort.Port)

	// Serve the gRPC server
	if err = s.Serve(portListener); err != nil {
		panic(commonlistener.FailedToServeError)
	}
	logger.ListenerLogger.ServerStarted(servicePort.Port)
	defer s.Stop()
}
