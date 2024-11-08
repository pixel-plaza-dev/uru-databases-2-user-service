package user

import (
	"context"
	commonbcrypt "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/crypto/bcrypt"
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
	database    *mongo.Database
	collections *map[string]*commonmongodb.Collection
	client      *mongo.Client
	logger      Logger
}

// NewDatabase creates a new MongoDB user database handler
func NewDatabase(client *mongo.Client, databaseName string, logger Logger) (database *Database, err error) {
	// Get the user service database
	userServiceDb := client.Database(databaseName)

	// Create map of collections
	collections := make(map[string]*commonmongodb.Collection)

	for _, collection := range []*commonmongodb.Collection{
		mongodb.UserCollection, mongodb.UserEmailCollection, mongodb.UserPhoneNumberCollection,
		mongodb.UserUsernameLogCollection, mongodb.UserHashedPasswordLogCollection} {
		// Create the collection
		collections[collection.Name] = collection
		if _, err = collection.CreateCollection(userServiceDb); err != nil {
			return nil, err
		}
	}

	// Create the user database instance
	instance := &Database{client: client, database: userServiceDb, collections: &collections, logger: logger}

	return instance, nil
}

// Database returns the MongoDB users database
func (d *Database) Database() *mongo.Database {
	return d.database
}

// GetQueryContext returns a new query context
func (d *Database) GetQueryContext() (ctx context.Context, cancelFunc context.CancelFunc) {
	return context.WithTimeout(context.Background(), mongodb.QueryCtxTimeout)
}

// GetTransactionContext returns a new transaction context
func (d *Database) GetTransactionContext() (ctx context.Context, cancelFunc context.CancelFunc) {
	return context.WithTimeout(context.Background(), mongodb.TransactionCtxTimeout)
}

// GetCollection returns a collection
func (d *Database) GetCollection(collection *commonmongodb.Collection) *mongo.Collection {
	return d.database.Collection(collection.Name)
}

// FindUser finds a user
func (d *Database) FindUser(filter bson.M, projection interface{}) (user *commonuser.User, err error) {
	// Create the context
	ctx, cancelFunc := d.GetQueryContext()
	defer cancelFunc()

	// Create the find options
	if projection == nil {
		projection = bson.M{"_id": 1}
	}
	findOptions := options.FindOne().SetProjection(projection)

	// Find the user
	err = d.GetCollection(mongodb.UserCollection).FindOne(ctx, filter, findOptions).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindUserByUsername finds a user by username
func (d *Database) FindUserByUsername(username string, projection interface{}) (user *commonuser.User, err error) {
	// Create the filter
	filter := bson.M{"username": username}
	return d.FindUser(filter, projection)
}

// CreateUser creates a new user
func (d *Database) CreateUser(user *commonuser.User, email *commonuser.UserEmail, phoneNumber *commonuser.UserPhoneNumber) (result interface{}, err error) {
	// Create the transaction options
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	// Create the context
	ctx, cancelFunc := d.GetTransactionContext()
	defer cancelFunc()

	// Starts a session on the client
	session, err := d.client.StartSession()
	if err != nil {
		panic(err)
	}
	defer session.EndSession(ctx)

	// Create the UserHashedPasswordLog
	currentTime := time.Now()
	userHashedPasswordLog := commonuser.UserHashedPasswordLog{
		ID:             primitive.NewObjectID(),
		UserID:         user.ID,
		HashedPassword: user.HashedPassword,
		AssignedAt:     currentTime,
	}

	// Create the UserUsernameLog
	userUsernameLog := commonuser.UserUsernameLog{
		ID:         primitive.NewObjectID(),
		UserID:     user.ID,
		Username:   user.Username,
		AssignedAt: currentTime,
	}

	// Run the transaction
	result, err = session.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		// Create a new email for the user
		if _, err = d.GetCollection(mongodb.UserEmailCollection).InsertOne(ctx, email); err != nil {
			return nil, err
		}

		// Create a new phone number for the user
		if _, err = d.GetCollection(mongodb.UserPhoneNumberCollection).InsertOne(ctx, phoneNumber); err != nil {
			return nil, err
		}

		// Create a new user
		userResult, err := d.GetCollection(mongodb.UserCollection).InsertOne(ctx, user)
		if err != nil {
			return nil, err
		}

		// Create a new user hashed password log
		if _, err = d.GetCollection(mongodb.UserHashedPasswordLogCollection).InsertOne(ctx, userHashedPasswordLog); err != nil {
			return nil, err
		}

		// Create a new user username log
		if _, err = d.GetCollection(mongodb.UserUsernameLogCollection).InsertOne(ctx, userUsernameLog); err != nil {
			return nil, err
		}

		return userResult, nil
	}, txnOptions)

	// Check if there are any errors
	if err != nil {
		d.logger.FailedToCreateDocument(err)
		return nil, err
	}

	return result, nil
}

// IsPasswordCorrect checks if the password is correct
func (d *Database) IsPasswordCorrect(username string, hashedPassword string) (userId string, err error) {
	// Create the projection
	projection := bson.M{"_id": 1, "hashed_password": 1}

	// Find the user
	user, err := d.FindUserByUsername(username, projection)
	if err != nil {
		return "", PasswordDoesNotMatchError
	}

	// Check if the password is correct
	if commonbcrypt.CheckPasswordHash(hashedPassword, user.HashedPassword) {
		return user.ID.Hex(), nil
	}
	return "", PasswordDoesNotMatchError
}

// UpdateUserField updates a user field
func (d *Database) UpdateUserField(userId string, field string, value interface{}) (result *mongo.UpdateResult, err error) {
	// Create the filter
	filter := bson.M{"_id": userId}

	// Create the update
	update := bson.M{"$set": bson.M{field: value}}

	// Create the context
	ctx, cancelFunc := d.GetQueryContext()
	defer cancelFunc()

	// Update the user
	result, err = d.GetCollection(mongodb.UserCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateUserPassword updates a user password
func (d *Database) UpdateUserPassword(userId string, hashedPassword string) (result *mongo.UpdateResult, err error) {
	return d.UpdateUserField(userId, "hashed_password", hashedPassword)
}

// UpdateUserUsername updates a user username
func (d *Database) UpdateUserUsername(userId string, username string) (result *mongo.UpdateResult, err error) {
	return d.UpdateUserField(userId, "username", username)
}

// UpdateUser updates a user
func (d *Database) UpdateUser(user *commonuser.User) (result *mongo.UpdateResult, err error) {
	// Create the filter
	filter := bson.M{"_id": user.ID}

	// Create the update
	update := bson.M{"$set": user}

	// Create the context
	ctx, cancelFunc := d.GetQueryContext()
	defer cancelFunc()

	// Update the user
	result, err = d.GetCollection(mongodb.UserCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return result, nil
}
