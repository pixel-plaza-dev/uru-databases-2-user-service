package user

import (
	"context"
	commonmongodb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb"
	commonuser "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb/database/user"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

type Database struct {
	mongodbConnection *commonmongodb.Connection
	database          *mongo.Database
	collections       *map[string]*commonmongodb.Collection
}

// NewDatabase creates a new MongoDB user database handler
func NewDatabase(connection *commonmongodb.Connection, databaseName string) (database *Database, err error) {
	// Connect to the MongoDB required database
	usersServiceDb := connection.Client.Database(databaseName)

	// Create map of collections
	collections := make(map[string]*commonmongodb.Collection)
	for _, collection := range []*commonmongodb.Collection{
		mongodb.UserCollection, mongodb.UserEmailCollection, mongodb.UserPhoneNumberCollection,
		mongodb.UserUsernameLogCollection, mongodb.UserHashedPasswordLogCollection} {
		collections[collection.Name] = collection
	}

	// Create the user database instance
	instance := &Database{mongodbConnection: connection, database: usersServiceDb, collections: &collections}

	return instance, nil
}

// Client returns the MongoDB client
func (u *Database) Client() *mongo.Client {
	return u.mongodbConnection.Client
}

// Database returns the MongoDB users database
func (u *Database) Database() *mongo.Database {
	return u.database
}

// GetQueryContext returns a new query context
func (u *Database) GetQueryContext() (ctx context.Context, cancelFunc context.CancelFunc) {
	return context.WithTimeout(context.Background(), mongodb.QueryCtxTimeout)
}

// GetTransactionContext returns a new transaction context
func (u *Database) GetTransactionContext() (ctx context.Context, cancelFunc context.CancelFunc) {
	return context.WithTimeout(context.Background(), mongodb.TransactionCtxTimeout)
}

// GetCollection returns a collection
func (u *Database) GetCollection(collection *commonmongodb.Collection) *mongo.Collection {
	return u.database.Collection(collection.Name)
}

// FindUser finds a user
func (u *Database) FindUser(filter bson.M) (user *commonuser.User, err error) {
	// Create the context
	ctx, cancelFunc := u.GetQueryContext()
	defer cancelFunc()

	// Find the user
	err = u.GetCollection(mongodb.UserCollection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindUserByUsername finds a user by username
func (u *Database) FindUserByUsername(username string) (user *commonuser.User, err error) {
	// Create the filter
	filter := bson.M{"username": username}
	return u.FindUser(filter)
}

// InsertOne inserts a document into a collection
func (u *Database) InsertOne(collection *commonmongodb.Collection, document interface{}) (result *mongo.InsertOneResult, err error) {
	// Create the context
	ctx, cancelFunc := u.GetQueryContext()
	defer cancelFunc()

	// Insert the document
	result, err = u.GetCollection(collection).InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// CreateUserUsernameLog creates a new user username log
func (u *Database) CreateUserUsernameLog(userUsernameLog commonuser.UserUsernameLog) (result *mongo.InsertOneResult, err error) {
	return u.InsertOne(mongodb.UserUsernameLogCollection, userUsernameLog)
}

// CreateUserHashedPasswordLog creates a new user hashed password log
func (u *Database) CreateUserHashedPasswordLog(userHashedPasswordLog *commonuser.UserHashedPasswordLog) (result *mongo.InsertOneResult, err error) {
	return u.InsertOne(mongodb.UserHashedPasswordLogCollection, userHashedPasswordLog)
}

// CreateUser creates a new user
func (u *Database) CreateUser(user *commonuser.User, email *commonuser.UserEmail, phoneNumber *commonuser.UserPhoneNumber) (result interface{}, err error) {
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

	// Create the UserUsernameLog
	currentTime := time.Now()
	userUsernameLog := commonuser.UserUsernameLog{
		ID:         primitive.NewObjectID(),
		UserID:     user.ID,
		Username:   user.Username,
		AssignedAt: currentTime,
	}

	// Create the UserHashedPasswordLog
	userHashedPasswordLog := commonuser.UserHashedPasswordLog{
		ID:             primitive.NewObjectID(),
		UserID:         user.ID,
		HashedPassword: user.HashedPassword,
		AssignedAt:     currentTime,
	}

	// Run the transaction
	result, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		// Create a new email for the user
		if _, err = u.GetCollection(mongodb.UserEmailCollection).InsertOne(ctx, email); err != nil {
			return nil, err
		}

		// Create a new phone number for the user
		if _, err = u.GetCollection(mongodb.UserPhoneNumberCollection).InsertOne(ctx, phoneNumber); err != nil {
			return nil, err
		}

		// Create a new user
		userResult, err := u.GetCollection(mongodb.UserCollection).InsertOne(ctx, user)
		if err != nil {
			return nil, err
		}

		// Create a new user username log
		if _, err = u.GetCollection(mongodb.UserUsernameLogCollection).InsertOne(ctx, userUsernameLog); err != nil {
			return nil, err
		}

		// Create a new user hashed password log
		if _, err = u.GetCollection(mongodb.UserHashedPasswordLogCollection).InsertOne(ctx, userHashedPasswordLog); err != nil {
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
func (u *Database) CreateUserEmail(userEmail commonuser.UserEmail) (result *mongo.InsertOneResult, err error) {
	return u.InsertOne(mongodb.UserEmailCollection, userEmail)
}

// CreateUserPhoneNumber creates a new user phone number
func (u *Database) CreateUserPhoneNumber(userPhoneNumber commonuser.UserPhoneNumber) (result *mongo.InsertOneResult, err error) {
	return u.InsertOne(mongodb.UserPhoneNumberCollection, userPhoneNumber)
}
