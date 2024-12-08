package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	commongcloud "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/cloud/gcloud"
	commonenv "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/config/env"
	commonflag "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/config/flag"
	commonjwtvalidator "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto/jwt/validator"
	commonjwtvalidatorgrpc "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto/jwt/validator/grpc"
	commonmongodb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/database/mongodb"
	clientauth "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/grpc/client/interceptor/auth"
	serverauth "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/grpc/server/interceptor/auth"
	commongrpcvalidator "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/grpc/server/validator"
	commonlistener "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/listener"
	commontls "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/http/tls"
	pbauth "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/compiled/pixel_plaza/auth"
	pbuser "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/compiled/pixel_plaza/user"
	pbconfigauth "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/config/grpc/auth"
	pbconfiguser "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/config/grpc/user"
	pbtypesgrpc "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/types/grpc"
	appmongodb "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/database/mongodb"
	userdatabase "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/database/mongodb/user"
	appgrpc "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/grpc"
	userserver "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/grpc/server/user"
	userservervalidator "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/grpc/server/user/validator"
	appjwt "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/jwt"
	applistener "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/listener"
	applogger "github.com/pixel-plaza-dev/uru-databases-2-user-service/app/logger"
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
	applogger.Flag.ModeFlagSet(commonflag.Mode)

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
		applistener.PortKey,
	)
	if err != nil {
		panic(err)
	}
	applogger.Environment.EnvironmentVariableLoaded(applistener.PortKey)

	// Get the MongoDB URI
	mongoDbUri, err := commonenv.LoadVariable(userdatabase.UriKey)
	if err != nil {
		panic(err)
	}
	applogger.Environment.EnvironmentVariableLoaded(userdatabase.UriKey)

	// Get the required MongoDB database name
	mongoDbName, err := commonenv.LoadVariable(userdatabase.DbNameKey)
	if err != nil {

		panic(err)
	}
	applogger.Environment.EnvironmentVariableLoaded(userdatabase.DbNameKey)

	// Get the gRPC services URI
	var uriKeys = []string{appgrpc.AuthServiceUriKey}
	var uris = make(map[string]string)
	for _, uriKey := range uriKeys {
		uri, err := commonenv.LoadVariable(uriKey)
		if err != nil {
			panic(err)
		}
		applogger.Environment.EnvironmentVariableLoaded(uriKey)
		uris[uriKey] = uri
	}

	// Get the JWT public key
	jwtPublicKey, err := commonenv.LoadVariable(appjwt.PublicKey)
	if err != nil {
		panic(err)
	}
	applogger.Environment.EnvironmentVariableLoaded(appjwt.PublicKey)

	// Load Google Cloud service account credentials
	googleCredentials, err := commongcloud.LoadGoogleCredentials(context.Background())
	if err != nil {
		panic(err)
	}

	// Get the service account token source for each gRPC server URI
	var tokenSources = make(map[string]*oauth.TokenSource)
	for _, uriKey := range uriKeys {
		tokenSource, err := commongcloud.LoadServiceAccountCredentials(
			context.Background(), "https://"+uris[uriKey], googleCredentials,
		)
		if err != nil {
			panic(err)
		}
		tokenSources[uriKey] = tokenSource
	}

	// Get the MongoDB configuration
	mongoDbConfig := &commonmongodb.Config{
		Uri:     mongoDbUri,
		Timeout: appmongodb.ConnectionCtxTimeout,
	}

	// Get the connection handler
	mongodbConnection, err := commonmongodb.NewDefaultConnectionHandler(mongoDbConfig)
	if err != nil {
		panic(err)
	}

	// Connect to MongoDB and get the client
	mongodbClient, err := mongodbConnection.Connect()
	if err != nil {
		panic(err)
	}

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

	// Create gRPC interceptions map
	var grpcInterceptions = map[string]*map[pbtypesgrpc.Method]pbtypesgrpc.Interception{
		appgrpc.AuthServiceUriKey: &pbconfigauth.Interceptions,
	}

	// Create client authentication interceptors
	var clientAuthInterceptors = make(map[string]clientauth.Authentication)
	for _, uriKey := range uriKeys {
		clientAuthInterceptor, err := clientauth.NewInterceptor(tokenSources[uriKey], grpcInterceptions[uriKey])
		if err != nil {
			panic(err)
		}
		clientAuthInterceptors[uriKey] = clientAuthInterceptor
	}

	// Create gRPC connections
	var conns = make(map[string]*grpc.ClientConn)
	for _, uriKey := range uriKeys {
		conn, err := grpc.NewClient(
			uris[uriKey], grpc.WithTransportCredentials(transportCredentials),
			grpc.WithChainUnaryInterceptor(clientAuthInterceptors[uriKey].Authenticate()),
		)
		if err != nil {
			panic(err)
		}
		conns[uriKey] = conn
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

	// Create user database handler
	userDatabase, err := userdatabase.NewDatabase(
		mongodbClient,
		mongoDbName,
		authClient,
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		// Disconnect from MongoDB
		mongodbConnection.Disconnect()
		applogger.MongoDb.DisconnectedFromDatabase()
	}()
	applogger.MongoDb.ConnectedToDatabase()

	// Create token validator
	tokenValidator, err := commonjwtvalidatorgrpc.NewDefaultTokenValidator(
		tokenSources[appgrpc.AuthServiceUriKey], authClient, nil,
	)
	if err != nil {
		panic(err)
	}

	// Create JWT validator with ED25519 public key
	jwtValidator, err := commonjwtvalidator.NewEd25519Validator(
		[]byte(jwtPublicKey),
		tokenValidator,
		commonflag.Mode,
	)
	if err != nil {
		panic(err)
	}

	// Create server authentication interceptor
	serverAuthInterceptor, err := serverauth.NewInterceptor(
		jwtValidator,
		&pbconfiguser.Interceptions,
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

	// Create the gRPC server validator
	serverValidator := commongrpcvalidator.NewDefaultValidator()

	// Create the gRPC user server validator
	userServerValidator, err := userservervalidator.NewValidator(
		userDatabase,
		serverValidator,
	)
	if err != nil {
		panic(err)
	}

	// Create the gRPC user server
	userServer := userserver.NewServer(
		userDatabase,
		authClient,
		applogger.UserServer,
		userServerValidator,
		applogger.JwtValidator,
	)

	// Register the user server with the gRPC server
	pbuser.RegisterUserServer(s, userServer)

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
	applogger.Listener.ServerStarted(servicePort.Port)
	if err = s.Serve(portListener); err != nil {
		panic(commonlistener.FailedToServeError)
	}
	defer s.Stop()
}
