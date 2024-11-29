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
	pbauth "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/protobuf/compiled/auth"
	protobuf "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/protobuf/compiled/user"
	detailsuser "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/protobuf/details/user"
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
	var uris = make(map[string]string)
	for _, key := range []string{appgrpc.AuthServiceUriKey} {
		uri, err := commonenv.LoadVariable(key)
		if err != nil {
			panic(err)
		}
		applogger.Environment.EnvironmentVariableLoaded(key)
		uris[key] = uri
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
		Timeout: appmongodb.ConnectionCtxTimeout,
	}

	// Get the connection handler
	mongodbConnection := commonmongodb.NewDefaultConnectionHandler(mongoDbConfig)

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
		tokenSources[appgrpc.AuthServiceUriKey], &authClient, nil,
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

	// Create the gRPC server validator
	serverValidator := commongrpcvalidator.NewDefaultValidator()

	// Create the gRPC user server validator
	userServerValidator := userservervalidator.NewValidator(
		userDatabase,
		serverValidator,
	)

	// Create the gRPC user server
	userServer := userserver.NewServer(
		userDatabase,
		authClient,
		applogger.UserServer,
		userServerValidator,
	)

	// Register the user server with the gRPC server
	protobuf.RegisterUserServer(s, userServer)

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
