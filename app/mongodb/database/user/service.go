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
	authClient  pbauth.AuthClient
}

// NewDatabase creates a new MongoDB user database handler
func NewDatabase(client *mongo.Client, databaseName string, authClient pbauth.AuthClient) (database *Database, err error) {
	// Get the user service database
	userServiceDb := client.Database(databaseName)

	// Create map of collections
	collections := make(map[string]*commonmongodb.Collection)

	for _, collection := range []*commonmongodb.Collection{
		mongodb.UserCollection, mongodb.UserSharedIdentifierCollection, mongodb.UserEmailCollection, mongodb.UserPhoneNumberCollection,
		mongodb.UserUsernameLogCollection, mongodb.UserHashedPasswordLogCollection} {
		// Create the collection
		collections[collection.Name] = collection
		if _, err = collection.CreateCollection(userServiceDb); err != nil {
			return nil, err
		}
	}

	// Create the user database instance
	instance := &Database{client: client, database: userServiceDb, collections: &collections,
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
func (d *Database) CreateUser(user *commonuser.User, userSharedIdentifier *commonuser.UserSharedIdentifier, userEmail *commonuser.UserEmail, userPhoneNumber *commonuser.UserPhoneNumber) error {
	// Create the UserHashedPasswordLog and UserUsernameLog objects
	userHashedPasswordLog := d.CreateUserHashedPasswordLogObject(user.ID, user.HashedPassword)
	userUsernameLog := d.CreateUserUsernameLogObject(user.ID, user.Username)

	// Run the transaction
	err := commonmongodb.CreateTransaction(d.client, func(sc mongo.SessionContext) error {
		// Create a new email for the user
		if _, err := d.GetCollection(mongodb.UserEmailCollection).InsertOne(sc, userEmail); err != nil {
			return err
		}

		// Create a new phone number for the user
		if _, err := d.GetCollection(mongodb.UserPhoneNumberCollection).InsertOne(sc, userPhoneNumber); err != nil {
			return err
		}

		// Create a new shared identifier for the user
		if _, err := d.GetCollection(mongodb.UserSharedIdentifierCollection).InsertOne(sc, userSharedIdentifier); err != nil {
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
		_, err := d.GetCollection(mongodb.UserUsernameLogCollection).InsertOne(sc, userUsernameLog)
		return err
	})
	return err
}

// FindUser finds a user
func (d *Database) FindUser(ctx context.Context, filter interface{}, projection interface{}, sort interface{}) (user *commonuser.User, err error) {
	// Set the default projection
	if projection == nil {
		projection = bson.M{"_id": 1}
	}

	// Create the find options
	findOptions := commonmongodb.PrepareFindOneOptions(projection, sort)

	// Add not deleted filter
	filter = bson.M{"$and": []interface{}{filter, bson.M{"deleted_at": bson.M{"$exists": false}}}}

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
		_, err = d.GetCollection(mongodb.UserUsernameLogCollection).InsertOne(sc, userUsernameLog)
		return err
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

// FindUserSharedIdentifier finds a user's shared identifier
func (d *Database) FindUserSharedIdentifier(ctx context.Context, filter interface{}, projection interface{}, sort interface{}) (userSharedIdentifier *commonuser.UserSharedIdentifier, err error) {
	// Set the default projection
	if projection == nil {
		projection = bson.M{"_id": 1}
	}

	// Create the find options
	findOptions := commonmongodb.PrepareFindOneOptions(projection, sort)

	// Find the user's shared identifier
	err = d.GetCollection(mongodb.UserSharedIdentifierCollection).FindOne(ctx, filter, findOptions).Decode(userSharedIdentifier)
	if err != nil {
		return nil, err
	}
	return userSharedIdentifier, nil
}

// GetUserSharedIdByUserId gets the user's shared identifier by the user ID
func (d *Database) GetUserSharedIdByUserId(ctx context.Context, userId string) (userSharedId string, err error) {
	// Convert the user ID to an object ID
	objectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return "", err
	}

	// Find the user's shared identifier
	var userSharedIdentifier *commonuser.UserSharedIdentifier
	userSharedIdentifier, err = d.FindUserSharedIdentifier(ctx, bson.M{"user_id": *objectId}, bson.M{"uuid": 1}, nil)
	if err != nil {
		return "", err
	}
	return userSharedIdentifier.UUID, nil
}

// GetUserIdByUserSharedId gets the user ID by the user shared ID
func (d *Database) GetUserIdByUserSharedId(ctx context.Context, userSharedId string) (userId string, err error) {
	// Find the user's shared identifier
	userSharedIdentifier, err := d.FindUserSharedIdentifier(ctx, bson.M{"uuid": userSharedId}, bson.M{"user_id": 1}, nil)
	if err != nil {
		return "", err
	}
	return userSharedIdentifier.UserID.Hex(), nil
}

// FindUserPhoneNumber finds a user's phone number
func (d *Database) FindUserPhoneNumber(ctx context.Context, filter interface{}, projection interface{}, sort interface{}) (userPhoneNumber *commonuser.UserPhoneNumber, err error) {
	// Set the default projection
	if projection == nil {
		projection = bson.M{"_id": 1}
	}

	// Create the find options
	findOptions := commonmongodb.PrepareFindOneOptions(projection, sort)

	// Find the user's phone number
	err = d.GetCollection(mongodb.UserPhoneNumberCollection).FindOne(ctx, filter, findOptions).Decode(userPhoneNumber)
	if err != nil {
		return nil, err
	}
	return userPhoneNumber, nil
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

// DeleteUser deletes a user
func (d *Database) DeleteUser(grpcCtx context.Context, userId string) error {
	// Convert the user ID to an object ID
	objectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return err
	}

	// Run the transaction
	err = commonmongodb.CreateTransaction(d.client, func(sc mongo.SessionContext) error {
		// Update the user deleted at field
		if _, err = d.GetCollection(mongodb.UserCollection).UpdateOne(sc, bson.M{"_id": *objectId}, bson.M{"deleted_at": time.Now()}); err != nil {
			return err
		}

		// Close all user sessions
		_, err = d.authClient.CloseSessions(grpcCtx, &pbauth.CloseSessionsRequest{})

		return nil
	})
	return err
}

// AddEmail adds an email to a user
func (d *Database) AddEmail(ctx context.Context, userId string, email string) error {
	// Convert the user ID to an object ID
	objectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return err
	}

	// Run the transaction
	err = commonmongodb.CreateTransaction(d.client, func(sc mongo.SessionContext) error {
		// Check if the user email already exists
		_, err = d.FindUserEmail(ctx, bson.M{"user_id": *objectId, "email": email, "revoked_at": bson.M{"$exists": false}}, bson.M{"_id": 1}, nil)
		if err == nil || !errors.Is(mongo.ErrNoDocuments, err) {
			return EmailAlreadyExistsError
		}

		// Add the email to the user
		_, err = d.GetCollection(mongodb.UserEmailCollection).InsertOne(ctx, commonuser.UserEmail{
			ID:         primitive.NewObjectID(),
			UserID:     *objectId,
			Email:      email,
			AssignedAt: time.Now(),
		})

		return err
	})
	return err
}

// DeleteEmail deletes an email from a user
func (d *Database) DeleteEmail(ctx context.Context, userId string, email string) error {
	// Convert the user ID to an object ID
	objectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return err
	}

	// Revoke the user's email
	_, err = d.GetCollection(mongodb.UserEmailCollection).UpdateOne(ctx, bson.M{"user_id": *objectId, "email": email, "is_primary": false, "revoked_at": bson.M{"$exists": false}}, bson.M{"revoked_at": time.Now()})
	return err
}

// FindUserEmail finds a user's email
func (d *Database) FindUserEmail(ctx context.Context, filter interface{}, projection interface{}, sort interface{}) (userEmail *commonuser.UserEmail, err error) {
	// Set the default projection
	if projection == nil {
		projection = bson.M{"_id": 1}
	}

	// Create the find options
	findOptions := commonmongodb.PrepareFindOneOptions(projection, sort)

	// Find the user's email
	err = d.GetCollection(mongodb.UserEmailCollection).FindOne(ctx, filter, findOptions).Decode(userEmail)
	if err != nil {
		return nil, err
	}
	return userEmail, nil
}

// FindUserEmailByEmail finds a user's email by email
func (d *Database) FindUserEmailByEmail(ctx context.Context, userId primitive.ObjectID, email string, projection interface{}, sort interface{}) (userEmail *commonuser.UserEmail, err error) {
	// Check if the user email already exists
	userEmail, err = d.FindUserEmail(ctx, bson.M{"user_id": userId, "email": email, "revoked_at": bson.M{"$exists": false}}, projection, sort)
	if err != nil {
		return nil, err
	}
	return userEmail, nil
}

// FindUserEmailPrimaryEmail finds a user's primary email
func (d *Database) FindUserEmailPrimaryEmail(ctx context.Context, userId primitive.ObjectID, projection interface{}, sort interface{}) (userEmail *commonuser.UserEmail, err error) {
	// Find the user's primary email
	userEmail, err = d.FindUserEmail(ctx, bson.M{"user_id": userId, "is_primary": true, "revoked_at": bson.M{"$exists": false}}, projection, sort)
	if err != nil {
		return nil, err
	}
	return userEmail, nil
}

// GetPrimaryEmail gets the user's primary email
func (d *Database) GetPrimaryEmail(ctx context.Context, userId string) (email string, err error) {
	// Convert the user ID to an object ID
	objectId, err := commonmongodb.GetObjectIdFromString(userId)
	if err != nil {
		return "", err
	}

	// Find the user's primary email
	var userEmail *commonuser.UserEmail
	userEmail, err = d.FindUserEmailPrimaryEmail(ctx, *objectId, bson.M{"email": 1}, nil)
	if err != nil {
		return "", err
	}
	return userEmail.Email, nil
}

// UserEmailExists checks if the user's email exists
func (d *Database) UserEmailExists(ctx context.Context, userId primitive.ObjectID, email string) (userEmailId string, err error) {
	// Check if the user email already exists
	userEmail, err := d.FindUserEmail(ctx, bson.M{"user_id": userId, "email": email, "revoked_at": bson.M{"$exists": false}}, bson.M{"_id": 1}, nil)
	if err != nil {
		return "", err
	}
	return userEmail.ID.Hex(), nil
}
