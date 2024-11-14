package main

import (
	"flag"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	commonenv "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/env"
	commonflag "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/flag"
	commongrpc "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/grpc"
	commonauthinterceptor "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/grpc/server/interceptor/auth"
	commonjwt "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/jwt"
	commonjwtvalidator "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/jwt/validator"
	commonlistener "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/listener"
	commonmongodb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb"
	pbauth "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/compiled-protobuf/auth"
	protobuf "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/compiled-protobuf/user"
	appgrpc "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/grpc"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/grpc/interceptor/auth"
	userserver "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/grpc/server/user"
	appjwt "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/jwt"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/listener"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/logger"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb"
	userdatabase "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb/database/user"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"net"
	"time"
)

// Load environment variables
func init() {
	// Declare flags and parse them
	commonflag.SetModeFlag()
	flag.Parse()
	logger.FlagLogger.ModeFlagSet(commonflag.Mode)

	// Check if the environment is production
	if commonflag.Mode.IsProd() {
		return
	}

	if err := godotenv.Load(); err != nil {
		panic(commonenv.FailedToLoadEnvironmentVariablesError)
	}
}

func main() {
	// Get the listener port
	servicePort, err := commonlistener.LoadServicePort("0.0.0.0", listener.PortKey)
	if err != nil {
		panic(err)
	}
	logger.EnvironmentLogger.EnvironmentVariableLoaded(listener.PortKey)

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

	// Get the auth service URI
	authUri, err := commongrpc.LoadServiceURI(appgrpc.AuthServiceUriKey)
	if err != nil {
		panic(err)
	}
	logger.EnvironmentLogger.EnvironmentVariableLoaded(appgrpc.AuthServiceUriKey)

	// Get the user service URI
	userUri, err := commongrpc.LoadServiceURI(appgrpc.UserServiceUriKey)
	if err != nil {
		panic(err)
	}
	logger.EnvironmentLogger.EnvironmentVariableLoaded(appgrpc.UserServiceUriKey)

	// Get the JWT public key
	jwtPublicKey, err := commonjwt.LoadJwtKey(appjwt.PublicKey)
	if err != nil {
		panic(err)
	}
	logger.EnvironmentLogger.EnvironmentVariableLoaded(appjwt.PublicKey)

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

	// Connect to gRPC servers
	var authConn *grpc.ClientConn

	if commonflag.Mode.IsDev() {
		// Load the self-signed CA certificates for the Pixel Plaza's services
		CACredentials, err := commongrpc.LoadTLSCredentials(appgrpc.CACertificatePath)
		if err != nil {
			panic(err)
		}

		authConn, err = grpc.NewClient(authUri, grpc.WithTransportCredentials(CACredentials))
		if err != nil {
			panic(err)
		}
	} else {
		// Load default account credentials
		tokenSource, err := commongrpc.LoadServiceAccountCredentials(context.Background(), userUri)
		if err != nil {
			panic(err)
		}

		authConn, err = grpc.NewClient(authUri, grpc.WithPerRPCCredentials(tokenSource))
		if err != nil {
			panic(err)
		}
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			panic(err)
		}
	}(authConn)

	// Create auth client
	authClient := pbauth.NewAuthClient(authConn)

	// Create JWT validator
	jwtValidator, err := commonjwtvalidator.NewDefaultValidator([]byte(jwtPublicKey), func(claims *jwt.MapClaims) (*jwt.MapClaims, error) {
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
	var s *grpc.Server

	if commonflag.Mode.IsDev() {
		// Load the TLS certificate and key
		grpcTransportCredentials, err := credentials.NewServerTLSFromFile(appgrpc.CertificateFilePath, appgrpc.KeyFilePath)
		if err != nil {
			panic(err)
		}

		s = grpc.NewServer(grpc.Creds(grpcTransportCredentials),
			grpc.ChainUnaryInterceptor(
				authInterceptor.UnaryServerInterceptor()))
	} else {
		s = grpc.NewServer(grpc.ChainUnaryInterceptor(
			authInterceptor.UnaryServerInterceptor()))
	}

	// Create a new gRPC user server
	userServer := userserver.NewServer(userDatabase, authClient, logger.UserServerLogger)

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
