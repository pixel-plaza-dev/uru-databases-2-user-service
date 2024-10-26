package mongodb

import (
	"context"
	commonMongodb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"time"
)

type UserDatabase struct {
	mongodbConnection *commonMongodb.Connection
	database          *mongo.Database
	collections       *map[string]*commonMongodb.Collection
}

// NewUserDatabase creates a new MongoDB user database handler
func NewUserDatabase(connection *commonMongodb.Connection, database string) (userDatabase *UserDatabase, err error) {
	// Connect to the MongoDB required database
	usersServiceDb := connection.Client.Database(database)

	// Create map of collections
	collections := make(map[string]*commonMongodb.Collection)
	for _, collection := range []*commonMongodb.Collection{
		UserCollection, UserEmailCollection, UserPhoneNumberCollection,
		UserUsernameLogCollection, UserHashedPasswordLogCollection} {
		collections[collection.Name] = collection
	}

	// Create the user database instance
	instance := &UserDatabase{mongodbConnection: connection, database: usersServiceDb, collections: &collections}

	return instance, nil
}

// Client returns the MongoDB client
func (u *UserDatabase) Client() *mongo.Client {
	return u.mongodbConnection.Client
}

// Database returns the MongoDB users database
func (u *UserDatabase) Database() *mongo.Database {
	return u.database
}

// GetQueryContext returns a new query context
func (u *UserDatabase) GetQueryContext() (ctx context.Context, cancelFunc context.CancelFunc) {
	return context.WithTimeout(context.Background(), QueryCtxTimeout)
}

// GetTransactionContext returns a new transaction context
func (u *UserDatabase) GetTransactionContext() (ctx context.Context, cancelFunc context.CancelFunc) {
	return context.WithTimeout(context.Background(), TransactionCtxTimeout)
}

// GetCollection returns a collection
func (u *UserDatabase) GetCollection(collection *commonMongodb.Collection) *mongo.Collection {
	return u.database.Collection(collection.Name)
}

// FindUser finds a user
func (u *UserDatabase) FindUser(filter bson.M) (user *User, err error) {
	// Create the context
	ctx, cancelFunc := u.GetQueryContext()
	defer cancelFunc()

	// Find the user
	err = u.GetCollection(UserCollection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindUserByUsername finds a user by username
func (u *UserDatabase) FindUserByUsername(username string) (user *User, err error) {
	// Create the filter
	filter := bson.M{"username": username}
	return u.FindUser(filter)
}

// InsertOne inserts a document into a collection
func (u *UserDatabase) InsertOne(collection *commonMongodb.Collection, document interface{}) (result *mongo.InsertOneResult, err error) {
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
func (u *UserDatabase) CreateUserUsernameLog(userUsernameLog UserUsernameLog) (result *mongo.InsertOneResult, err error) {
	return u.InsertOne(UserUsernameLogCollection, userUsernameLog)
}

// CreateUserHashedPasswordLog creates a new user hashed password log
func (u *UserDatabase) CreateUserHashedPasswordLog(userHashedPasswordLog *UserHashedPasswordLog) (result *mongo.InsertOneResult, err error) {
	return u.InsertOne(UserHashedPasswordLogCollection, userHashedPasswordLog)
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

	// Create the UserUsernameLog
	currentTime := time.Now()
	userUsernameLog := UserUsernameLog{
		ID:         primitive.NewObjectID(),
		UserID:     user.ID,
		Username:   user.Username,
		AssignedAt: currentTime,
	}

	// Create the UserHashedPasswordLog
	userHashedPasswordLog := UserHashedPasswordLog{
		ID:             primitive.NewObjectID(),
		UserID:         user.ID,
		HashedPassword: user.HashedPassword,
		AssignedAt:     currentTime,
	}

	// Run the transaction
	result, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		// Create a new email for the user
		if _, err = u.GetCollection(UserEmailCollection).InsertOne(ctx, email); err != nil {
			return nil, err
		}

		// Create a new phone number for the user
		if _, err = u.GetCollection(UserPhoneNumberCollection).InsertOne(ctx, phoneNumber); err != nil {
			return nil, err
		}

		// Create a new user
		userResult, err := u.GetCollection(UserCollection).InsertOne(ctx, user)
		if err != nil {
			return nil, err
		}

		// Create a new user username log
		if _, err = u.GetCollection(UserUsernameLogCollection).InsertOne(ctx, userUsernameLog); err != nil {
			return nil, err
		}

		// Create a new user hashed password log
		if _, err = u.GetCollection(UserHashedPasswordLogCollection).InsertOne(ctx, userHashedPasswordLog); err != nil {
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
	return u.InsertOne(UserEmailCollection, userEmail)
}

// CreateUserPhoneNumber creates a new user phone number
func (u *UserDatabase) CreateUserPhoneNumber(userPhoneNumber UserPhoneNumber) (result *mongo.InsertOneResult, err error) {
	return u.InsertOne(UserPhoneNumberCollection, userPhoneNumber)
}
