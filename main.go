package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	commongcloud "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/cloud/gcloud"
	commonenv "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/config/env"
	commonflag "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/config/flag"
	commonjwt "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto/jwt"
	commonjwtvalidator "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto/jwt/validator"
	commonjwtvalidatorgrpc "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto/jwt/validator/grpc"
	commonmongodb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/database/mongodb"
	commongrpc "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/grpc"
	clientauth "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/grpc/client/interceptor/auth"
	serverauth "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/grpc/server/interceptor/auth"
	commonlistener "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/listener"
	commontls "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/tls"
	pbauth "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/protobuf/compiled/auth"
	protobuf "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/protobuf/compiled/user"
	detailsuser "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/protobuf/details/user"
	appgrpc "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/grpc"
	userserver "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/grpc/server/user"
	appjwt "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/jwt"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/listener"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/logger"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb"
	userdatabase "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb/database/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/credentials/oauth"
	"net"
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
	servicePort, err := commonlistener.LoadServicePort(
		"0.0.0.0",
		listener.PortKey,
	)
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

	// Get the gRPC services URI
	var uris = make(map[string]string)
	for _, key := range []string{appgrpc.AuthServiceUriKey} {
		uri, err := commongrpc.LoadServiceURI(key)
		if err != nil {
			panic(err)
		}
		logger.EnvironmentLogger.EnvironmentVariableLoaded(key)
		uris[key] = uri
	}

	// Get the JWT public key
	jwtPublicKey, err := commonjwt.LoadJwtKey(appjwt.PublicKey)
	if err != nil {
		panic(err)
	}
	logger.EnvironmentLogger.EnvironmentVariableLoaded(appjwt.PublicKey)

	// Load Google Cloud service account credentials
	googleCredentials, err := commongcloud.LoadGoogleCredentials(context.Background())
	if err != nil {
		panic(err)
	}

	// Get the service account token source for each gRPC server URI
	var tokenSources = make(map[string]*oauth.TokenSource)
	for key, uri := range uris {
		tokenSource, err := commongcloud.LoadServiceAccountCredentials(
			context.Background(), "https://"+uri, googleCredentials,
		)
		if err != nil {
			panic(err)
		}
		tokenSources[key] = tokenSource
	}

	// Get the MongoDB configuration
	mongoDbConfig := &commonmongodb.Config{
		Uri:     mongoDbUri,
		Timeout: mongodb.ConnectionCtxTimeout,
	}

	// Get the connection handler
	mongodbConnection := commonmongodb.NewDefaultConnectionHandler(mongoDbConfig)

	// Connect to MongoDB and get the client
	mongodbClient, err := mongodbConnection.Connect()
	if err != nil {
		panic(err)
	}

	// Create user database handler
	userDatabase, err := userdatabase.NewDatabase(
		mongodbClient,
		mongoDbName,
		logger.UserDatabaseLogger,
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		// Disconnect from MongoDB
		mongodbConnection.Disconnect()
		logger.MongoDbLogger.DisconnectedFromMongoDB()
	}()
	logger.MongoDbLogger.ConnectedToMongoDB()

	// Load transport credentials
	var transportCredentials credentials.TransportCredentials

	if commonflag.Mode.IsDev() {
		// Load server TLS credentials
		transportCredentials, err = credentials.NewServerTLSFromFile(
			appgrpc.ServerCertPath, appgrpc.ServerKeyPath,
		)
		if err != nil {
			panic(err)
		}
	} else {
		// Load system certificates pool
		transportCredentials, err = commontls.LoadSystemCredentials()
		if err != nil {
			panic(err)
		}
	}

	// Create client authentication interceptors
	var clientAuthInterceptors = make(map[string]*clientauth.Interceptor)
	for key, tokenSource := range tokenSources {
		clientAuthInterceptor, err := clientauth.NewInterceptor(tokenSource)
		if err != nil {
			panic(err)
		}
		clientAuthInterceptors[key] = clientAuthInterceptor
	}

	// Create gRPC connections
	var conns = make(map[string]*grpc.ClientConn)
	for key, uri := range uris {
		conn, err := grpc.NewClient(
			uri, grpc.WithTransportCredentials(transportCredentials),
			grpc.WithChainUnaryInterceptor(clientAuthInterceptors[key].Authenticate()),
		)
		if err != nil {
			panic(err)
		}
		conns[key] = conn
	}
	defer func(conns map[string]*grpc.ClientConn) {
		for _, conn := range conns {
			err = conn.Close()
			if err != nil {
				panic(err)
			}
		}
	}(conns)

	// Create gRPC server clients
	authClient := pbauth.NewAuthClient(conns[appgrpc.AuthServiceUriKey])

	// Create token validator
	tokenValidator, err := commonjwtvalidatorgrpc.NewDefaultTokenValidator(
		tokenSources[appgrpc.AuthServiceUriKey], &authClient,
	)
	if err != nil {
		panic(err)
	}

	// Create JWT validator
	jwtValidator, err := commonjwtvalidator.NewDefaultValidator(
		[]byte(jwtPublicKey),
		tokenValidator,
	)
	if err != nil {
		panic(err)
	}

	// Create server authentication interceptor
	serverAuthInterceptor, err := serverauth.NewInterceptor(
		jwtValidator,
		&detailsuser.GRPCInterceptions,
	)
	if err != nil {
		panic(err)
	}

	// Create the gRPC server
	s := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			serverAuthInterceptor.Authenticate(),
		),
	)

	// Create the gRPC user server
	userServer := userserver.NewServer(
		userDatabase,
		authClient,
		logger.UserServerLogger,
	)

	// Register the user server with the gRPC server
	protobuf.RegisterUserServer(s, userServer)
	logger.ListenerLogger.ServerStarted(servicePort.Port)

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

	// Serve the gRPC server
	if err = s.Serve(portListener); err != nil {
		panic(commonlistener.FailedToServeError)
	}
	logger.ListenerLogger.ServerStarted(servicePort.Port)
	defer s.Stop()
}
