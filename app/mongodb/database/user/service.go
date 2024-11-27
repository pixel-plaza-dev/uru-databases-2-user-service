package user

import (
	"context"
	"errors"
	commonmongodb "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/database/mongodb"
	commonuser "github.com/pixel-plaza-dev/uru-databases-2-go-service-common/database/mongodb/model/user"
	pbauth "github.com/pixel-plaza-dev/uru-databases-2-protobuf-common/protobuf/compiled/auth"
	"github.com/pixel-plaza-dev/uru-databases-2-user-service/app/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Database struct {
	database    *mongo.Database
	collections *map[string]*commonmongodb.Collection
	client      *mongo.Client
	logger      Logger
	authClient  pbauth.AuthClient
}

// NewDatabase creates a new MongoDB user database handler
func NewDatabase(client *mongo.Client, databaseName string, logger Logger, authClient pbauth.AuthClient) (database *Database, err error) {
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
	instance := &Database{client: client, database: userServiceDb, collections: &collections, logger: logger,
		authClient: authClient}

	return instance, nil
}

// Database returns the MongoDB users database
func (d *Database) Database() *mongo.Database {
	return d.database
}

// GetCollection returns a collection
func (d *Database) GetCollection(collection *commonmongodb.Collection) *mongo.Collection {
	return d.database.Collection(collection.Name)
}

// CreateUserUsernameLogObject creates a new user username log object
func (d *Database) CreateUserUsernameLogObject(userId primitive.ObjectID, username string) commonuser.UserUsernameLog {
	return commonuser.UserUsernameLog{
		ID:         primitive.NewObjectID(),
		UserID:     userId,
		Username:   username,
		AssignedAt: time.Now(),
	}
}

// CreateUserHashedPasswordLogObject creates a new user hashed password log object
func (d *Database) CreateUserHashedPasswordLogObject(userId primitive.ObjectID, hashedPassword string) commonuser.UserHashedPasswordLog {
	return commonuser.UserHashedPasswordLog{
		ID:             primitive.NewObjectID(),
		UserID:         userId,
		HashedPassword: hashedPassword,
		AssignedAt:     time.Now(),
	}
}

// CreateUser creates a new user
func (d *Database) CreateUser(user *commonuser.User, email *commonuser.UserEmail, phoneNumber *commonuser.UserPhoneNumber) error {
	// Create the UserHashedPasswordLog and UserUsernameLog objects
	userHashedPasswordLog := d.CreateUserHashedPasswordLogObject(user.ID, user.HashedPassword)
	userUsernameLog := d.CreateUserUsernameLogObject(user.ID, user.Username)

	// Run the transaction
	err := commonmongodb.CreateTransaction(d.client, func(sc mongo.SessionContext) error {
		// Create a new email for the user
		if _, err := d.GetCollection(mongodb.UserEmailCollection).InsertOne(sc, email); err != nil {
			return err
		}

		// Create a new phone number for the user
		if _, err := d.GetCollection(mongodb.UserPhoneNumberCollection).InsertOne(sc, phoneNumber); err != nil {
			return err
		}

		// Create a new user
		if _, err := d.GetCollection(mongodb.UserCollection).InsertOne(sc, user); err != nil {
			return err
		}

		// Create a new user hashed password log
		if _, err := d.GetCollection(mongodb.UserHashedPasswordLogCollection).InsertOne(sc, userHashedPasswordLog); err != nil {
			return err
		}

		// Create a new user username log
		if _, err := d.GetCollection(mongodb.UserUsernameLogCollection).InsertOne(sc, userUsernameLog); err != nil {
			return err
		}

		return nil
	})

	// Check if there are any errors
	if err != nil {
		d.logger.FailedToCreateDocument(err)
		return err
	}

	return nil
}

// FindUser finds a user
func (d *Database) FindUser(ctx context.Context, filter interface{}, projection interface{}, sort interface{}) (user *commonuser.User, err error) {
	// Set the default projection
	if projection == nil {
		projection = bson.M{"_id": 1}
	}

	// Create the find options
	findOptions := commonmongodb.PrepareFindOneOptions(projection, sort)

	// Find the user
	err = d.GetCollection(mongodb.UserCollection).FindOne(ctx, filter, findOptions).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// FindUserByUsername finds a user by username
func (d *Database) FindUserByUsername(ctx context.Context, username string, projection interface{}, sort interface{}) (user *commonuser.User, err error) {
	// Check if the username is empty
	if username == "" {
		return nil, mongo.ErrNoDocuments
	}

	// Find the user
	return d.FindUser(ctx, bson.M{"username": username}, projection, sort)
}

// FindUserByUserId finds a user by the user ID
func (d *Database) FindUserByUserId(ctx context.Context, userId string, projection interface{}, sort interface{}) (user *commonuser.User, err error) {
	// Check if the user ID is empty
	if userId == "" {
		return nil, mongo.ErrNoDocuments
	}

	// Convert the user ID to an object ID
	objectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return nil, err
	}

	// Find the user
	return d.FindUser(ctx, bson.M{"_id": *objectId}, projection, sort)
}

// GetUserHashedPassword gets the user's hashed password
func (d *Database) GetUserHashedPassword(ctx context.Context, username string) (user *commonuser.User, err error) {
	// Check if the username is empty
	if username == "" {
		return nil, mongo.ErrNoDocuments
	}

	// Find the user
	return d.FindUserByUsername(ctx, username, bson.M{"_id": 1, "hashed_password": 1}, nil)
}

// GetUsernameByUserId gets the username by the user ID
func (d *Database) GetUsernameByUserId(ctx context.Context, userId string) (username string, err error) {
	// Find the user
	user, err := d.FindUserByUserId(ctx, userId, bson.M{"username": 1}, nil)
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

// GetUserIdByUsername gets the user ID by the username
func (d *Database) GetUserIdByUsername(ctx context.Context, username string) (userId string, err error) {
	// Find the user
	user, err := d.FindUserByUsername(ctx, username, bson.M{"_id": 1}, nil)
	if err != nil {
		return "", err
	}
	return user.ID.Hex(), nil
}

// UsernameExists checks if the username exists
func (d *Database) UsernameExists(ctx context.Context, username string) (exists bool, err error) {
	// Find the user
	user, err := d.FindUserByUsername(ctx, username, bson.M{"_id": 1}, nil)
	if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
		return false, err
	}
	return user != nil, nil
}

// UpdateUser updates a user
func (d *Database) UpdateUser(ctx context.Context, userId string, update interface{}) (result *mongo.UpdateResult, err error) {
	// Convert the user ID to an object ID
	objectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return nil, err
	}

	// Create the filter
	filter := bson.M{"_id": *objectId}

	// Update the user
	result, err = d.GetCollection(mongodb.UserCollection).UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateUserUsername updates the user username
func (d *Database) UpdateUserUsername(userId string, username string) error {
	// Convert the user ID to an object ID
	objectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return err
	}

	// Create the UserUsernameLog object
	userUsernameLog := d.CreateUserUsernameLogObject(*objectId, username)

	// Run the transaction
	err = commonmongodb.CreateTransaction(d.client, func(sc mongo.SessionContext) error {
		// Update the user username
		if _, err = d.GetCollection(mongodb.UserCollection).UpdateOne(sc, bson.M{"_id": *objectId}, bson.M{"username": username}); err != nil {
			return err
		}

		// Create a new user username log
		if _, err = d.GetCollection(mongodb.UserUsernameLogCollection).InsertOne(sc, userUsernameLog); err != nil {
			return err
		}

		return nil
	})
	return err
}

// UpdateUserPassword updates the user password
func (d *Database) UpdateUserPassword(grpcCtx context.Context, userId string, hashedPassword string) error {
	// Convert the user ID to an object ID
	objectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return err
	}

	// Create the UserHashedPasswordLog object
	userHashedPasswordLog := d.CreateUserHashedPasswordLogObject(*objectId, hashedPassword)

	// Run the transaction
	err = commonmongodb.CreateTransaction(d.client, func(sc mongo.SessionContext) error {
		// Update the user password
		if _, err = d.GetCollection(mongodb.UserCollection).UpdateOne(sc, bson.M{"_id": *objectId}, bson.M{"hashed_password": hashedPassword}); err != nil {
			return err
		}

		// Create a new user hashed password log
		if _, err = d.GetCollection(mongodb.UserHashedPasswordLogCollection).InsertOne(sc, userHashedPasswordLog); err != nil {
			return err
		}

		// Close all user sessions
		_, err = d.authClient.CloseSessions(grpcCtx, &pbauth.CloseSessionsRequest{})

		return nil
	})
	return err
}

// UpdateProfile updates a user
func (d *Database) UpdateProfile(ctx context.Context, userId string, update interface{}) (result *mongo.UpdateResult, err error) {
	return d.UpdateUser(ctx, userId, update)
}

// GetProfile gets the user's profile
func (d *Database) GetProfile(ctx context.Context, userId string) (user *commonuser.User, err error) {
	return d.FindUserByUserId(ctx, userId, bson.M{"username": 1, "first_name": 1, "last_name": 1}, nil)
}

// FindUserPhoneNumber finds a user's phone number
func (d *Database) FindUserPhoneNumber(ctx context.Context, filter interface{}, projection interface{}, sort interface{}) (phoneNumber *commonuser.UserPhoneNumber, err error) {
	// Create the find options
	findOptions := commonmongodb.PrepareFindOneOptions(projection, sort)

	// Find the user's phone number
	err = d.GetCollection(mongodb.UserPhoneNumberCollection).FindOne(ctx, filter, findOptions).Decode(phoneNumber)
	if err != nil {
		return nil, err
	}
	return phoneNumber, nil
}

// GetUserPhoneNumber gets the user's phone number
func (d *Database) GetUserPhoneNumber(ctx context.Context, userId string) (phoneNumber string, err error) {
	// Check if the user ID is empty
	if userId == "" {
		return "", mongo.ErrNoDocuments
	}

	// Convert the user ID to an object ID
	objectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return "", err
	}

	// Create the find options with the most recent document based on the given field
	sort := bson.M{"assigned_at": -1}

	// Find the user's phone number
	var userPhoneNumber *commonuser.UserPhoneNumber
	userPhoneNumber, err = d.FindUserPhoneNumber(ctx, bson.M{"user_id": objectId}, bson.M{"phone_number": 1}, sort)
	if err != nil {
		return "", err
	}
	return userPhoneNumber.PhoneNumber, nil
}

// UpdateUserPhoneNumber updates the user's phone number
func (d *Database) UpdateUserPhoneNumber(userId string, phoneNumber string) error {
	// Convert the user ID to an object ID
	objectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return err
	}

	// Run the transaction
	err = commonmongodb.CreateTransaction(d.client, func(sc mongo.SessionContext) error {
		// Revoke the user's phone number
		if _, err = d.GetCollection(mongodb.UserPhoneNumberCollection).UpdateOne(sc, bson.M{"user_id": *objectId, "revoked_at": bson.M{"$exists": false}}, bson.M{"revoked_at": time.Now()}); err != nil {
			return err
		}

		// Update the user's phone number
		if _, err = d.GetCollection(mongodb.UserPhoneNumberCollection).InsertOne(sc, commonuser.UserPhoneNumber{
			ID:          primitive.NewObjectID(),
			UserID:      *objectId,
			PhoneNumber: phoneNumber,
			AssignedAt:  time.Now(),
		}); err != nil {
			return err
		}

		return nil
	})
	return err
}
