package mongodb

import (
	commonMongodb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"golang.org/x/net/context"
	"pixel_plaza/users_service/config"
)

const (
	// UserCollection is the name of the users collection in MongoDB
	UserCollection = "User"

	// UserEmailCollection is the name of the user emails collection in MongoDB
	UserEmailCollection = "UserEmail"

	// UserPhoneNumberCollection is the name of the user phone numbers collection in MongoDB
	UserPhoneNumberCollection = "UserPhoneNumber"

	// UserEmailVerificationCollection is the name of the user email verifications collection in MongoDB
	UserEmailVerificationCollection = "UserEmailVerification"

	// UserPhoneNumberVerificationCollection is the name of the user phone number verifications collection in MongoDB
	UserPhoneNumberVerificationCollection = "UserPhoneNumberVerification"

	// UserResetPasswordCollection is the name of the user reset password collection in MongoDB
	UserResetPasswordCollection = "UserResetPassword"
)

type UserDatabase struct {
	mongodbConnection *commonMongodb.Connection
	database          *mongo.Database
}

// NewUserDatabase creates a new MongoDB user database handler
func NewUserDatabase(connection *commonMongodb.Connection, database string) (userDatabase *UserDatabase) {
	// Connect to the MongoDB required database
	usersServiceDb := connection.Client.Database(database)

	return &UserDatabase{mongodbConnection: connection, database: usersServiceDb}
}

// Client returns the MongoDB client
func (u *UserDatabase) Client() *mongo.Client {
	return u.mongodbConnection.Client
}

// Database returns the MongoDB users database
func (u *UserDatabase) Database() *mongo.Database {
	return u.database
}

// UserCollection returns the MongoDB users collection
func (u *UserDatabase) UserCollection() *mongo.Collection {
	return u.database.Collection(UserCollection)
}

// UserEmailCollection returns the MongoDB user emails collection
func (u *UserDatabase) UserEmailCollection() *mongo.Collection {
	return u.database.Collection(UserEmailCollection)
}

// UserPhoneNumberCollection returns the MongoDB user phone numbers collection
func (u *UserDatabase) UserPhoneNumberCollection() *mongo.Collection {
	return u.database.Collection(UserPhoneNumberCollection)
}

// GetQueryContext returns a new query context
func (u *UserDatabase) GetQueryContext() (ctx context.Context, cancelFunc context.CancelFunc) {
	return context.WithTimeout(context.Background(), config.MongoDbQueryCtxTimeout)
}

// GetTransactionContext returns a new transaction context
func (u *UserDatabase) GetTransactionContext() (ctx context.Context, cancelFunc context.CancelFunc) {
	return context.WithTimeout(context.Background(), config.MongoDbTransactionCtxTimeout)
}

// FindUserByUsername finds a user by username
func (u *UserDatabase) FindUserByUsername(username string) (user User, err error) {
	// Create the context
	ctx, cancelFunc := u.GetQueryContext()
	defer cancelFunc()

	// Find the user by username
	err = u.UserCollection().FindOne(ctx, User{Username: username}).Decode(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// CreateUser creates a new user
func (u *UserDatabase) CreateUser(user *User, email *UserEmail, phoneNumber *UserPhoneNumber) (result interface{}, err error) {
	// Create the transaction options
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	// Create the context
	ctx, cancelFunc := u.GetTransactionContext()
	defer cancelFunc()

	// Starts a session on the client
	session, err := u.Client().StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	// Run the transaction
	result, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		// Create a new email for the user
		if _, err = u.UserEmailCollection().InsertOne(ctx, email); err != nil {
			return nil, err
		}

		// Create a new phone number for the user
		if _, err = u.UserPhoneNumberCollection().InsertOne(ctx, phoneNumber); err != nil {
			return nil, err
		}

		// Create a new user
		userResult, err := u.UserCollection().InsertOne(ctx, user)
		if err != nil {
			return nil, err
		}

		return userResult, nil
	}, txnOptions)

	// Check if there are any errors
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateUserEmail creates a new user email
func (u *UserDatabase) CreateUserEmail(userEmail UserEmail) (result *mongo.InsertOneResult, err error) {
	// Create the context
	ctx, cancelFunc := u.GetQueryContext()
	defer cancelFunc()

	// Create a new user email
	result, err = u.UserEmailCollection().InsertOne(ctx, userEmail)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateUserPhoneNumber creates a new user phone number
func (u *UserDatabase) CreateUserPhoneNumber(userPhoneNumber UserPhoneNumber) (result *mongo.InsertOneResult, err error) {
	// Create the context
	ctx, cancelFunc := u.GetQueryContext()
	defer cancelFunc()

	// Create a new user phone number
	result, err = u.UserPhoneNumberCollection().InsertOne(ctx, userPhoneNumber)
	if err != nil {
		return nil, err
	}

	return result, nil
}
